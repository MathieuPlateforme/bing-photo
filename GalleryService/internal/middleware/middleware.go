package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type AuthServiceClient interface {
	ValidateToken(ctx context.Context, token string) (bool, error)
}

// AuthMiddleware protège les routes avec un token JWT
func AuthMiddleware(authClient AuthServiceClient, jwtSecret string) func(http.Handler) http.Handler {
	// Déclaration de la fonction ParseToken
	parseToken := func(token string) (map[string]interface{}, error) {
		parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil || !parsedToken.Valid {
			return nil, err
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			return nil, err
		}
		return claims, nil
	}

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

			// Utiliser ParseToken pour extraire les claims
			claims, err := parseToken(token)
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
