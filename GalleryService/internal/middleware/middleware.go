package middleware

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// --- INTERFACES ---

type JWTService interface {
	VerifyToken(token string) (map[string]interface{}, error)
}

type AuthServiceClient interface {
	ValidateToken(ctx context.Context, token string) (bool, error)
}

// --- HTTP MIDDLEWARE ---

// AuthMiddleware protège les routes HTTP avec un token JWT
func AuthMiddleware(authClient AuthServiceClient, jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Routes ignorées par le middleware
			if r.Method == http.MethodPost && (r.URL.Path == "/albums" || r.URL.Path == "/users") {
				log.Printf("Middleware ignoré pour la route : %s %s", r.Method, r.URL.Path)
				next.ServeHTTP(w, r)
				return
			}
			// Extraire le token de l'en-tête Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				log.Println("Erreur : en-tête Authorization manquant ou invalide")
				http.Error(w, "Token manquant ou invalide", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			log.Printf("Token extrait : %s", token)

			// Valider le token via AuthService
			valid, err := authClient.ValidateToken(r.Context(), token)
			if err != nil || !valid {
				log.Printf("Erreur lors de la validation du token : %v", err)
				http.Error(w, "Token invalide", http.StatusUnauthorized)
				return
			}

			// Parser localement pour récupérer les claims
			claims, err := parseToken(token, jwtSecret)
			if err != nil {
				log.Printf("Erreur lors de l'analyse du token : %v", err)
				http.Error(w, "Token invalide", http.StatusUnauthorized)
				return
			}

			// Extraire le userID des claims et l'ajouter au contexte
			userID, ok := claims["userID"].(float64)
			if !ok {
				log.Println("Erreur : userID introuvable ou type incorrect dans les claims")
				http.Error(w, "Token invalide", http.StatusUnauthorized)
				return
			}

			log.Printf("Contexte enrichi avec userID : %d", uint(userID))
			ctx := context.WithValue(r.Context(), UserIDKey, uint(userID))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// parseToken extrait les claims d’un token JWT
func parseToken(token string, secret string) (map[string]interface{}, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !parsedToken.Valid {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}

// --- gRPC INTERCEPTOR ---

// AuthInterceptor protège les méthodes gRPC listées dans methodsToIntercept
func AuthInterceptor(jwtService JWTService, methodsToIntercept map[string]bool) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		if !methodsToIntercept[info.FullMethod] {
			log.Printf("AuthInterceptor ignoré pour la méthode : %s", info.FullMethod)
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, fmt.Errorf("aucune métadonnée trouvée dans le contexte")
		}

		authHeader := md["authorization"]
		if len(authHeader) == 0 {
			return nil, fmt.Errorf("en-tête Authorization manquant")
		}

		token := strings.TrimPrefix(authHeader[0], "Bearer ")
		log.Printf("Token extrait : %s", token)

		// Vérifier le token
		claims, err := jwtService.VerifyToken(token)
		if err != nil {
			return nil, fmt.Errorf("token invalide : %v", err)
		}

		// Extraire et injecter le userID
		userID, ok := claims["userID"].(float64)
		if !ok {
			return nil, fmt.Errorf("userID introuvable dans le token")
		}

		ctx = context.WithValue(ctx, UserIDKey, uint(userID))
		log.Printf("userID ajouté au contexte : %d", uint(userID))

		return handler(ctx, req)
	}
}
