package middleware

import (
	"context"
	"fmt"
	"log"
	"strings"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type JWTService interface {
	VerifyToken(token string) (map[string]interface{}, error)
}

func AuthInterceptor(jwtService JWTService, methodsToIntercept map[string]bool) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Vérifier si la méthode nécessite l'authentification
		if !methodsToIntercept[info.FullMethod] {
			log.Printf("AuthInterceptor ignoré pour la méthode : %s", info.FullMethod)
			return handler(ctx, req)
		}

		// Extraire les métadonnées
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, fmt.Errorf("aucune métadonnée trouvée dans le contexte")
		}

		// Chercher l'en-tête Authorization (minuscule ou majuscule)
		var rawToken string
		if vals := md.Get("authorization"); len(vals) > 0 {
			rawToken = vals[0]
		} else if vals := md.Get("Authorization"); len(vals) > 0 {
			rawToken = vals[0]
		} else {
			return nil, fmt.Errorf("en-tête Authorization manquant")
		}

		token := strings.TrimPrefix(rawToken, "Bearer ")
		token = strings.TrimSpace(token)

		log.Printf("Token extrait : %s", token)

		// Valider le token via JWTService
		claims, err := jwtService.VerifyToken(token)
		if err != nil {
			return nil, fmt.Errorf("token invalide : %v", err)
		}

		// Extraire et injecter userID
		userID, ok := claims["userID"].(float64)
		if !ok {
			return nil, fmt.Errorf("userID introuvable dans le token")
		}
		ctx = context.WithValue(ctx, UserIDKey, uint(userID))
		log.Printf("userID ajouté au contexte : %d", uint(userID))

		return handler(ctx, req)
	}
}
