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

		// Extraire les métadonnées pour obtenir le token
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

		// Valider le token via JWTService
		claims, err := jwtService.VerifyToken(token)
		if err != nil {
			return nil, fmt.Errorf("token invalide : %v", err)
		}

		// Extraire et ajouter userID au contexte
		userID, ok := claims["userID"].(float64)
		if !ok {
			return nil, fmt.Errorf("userID introuvable dans le token")
		}
		ctx = context.WithValue(ctx, "userID", uint(userID))
		log.Printf("userID ajouté au contexte : %d", uint(userID))

		// Continuer avec le handler
		return handler(ctx, req)
	}
}
