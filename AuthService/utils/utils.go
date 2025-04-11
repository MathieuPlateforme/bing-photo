package utils

import (
	"context"
	"errors"
	"log"
	"AuthService/middleware"
	"net/http"
	"strings"
	"google.golang.org/grpc/metadata"

)

func GetUserIDFromContext(ctx context.Context) (uint, error) {
	log.Println("Tentative de récupération du userID à partir du contexte...")

	userID, ok := ctx.Value(middleware.UserIDKey).(uint)
	if !ok {
		log.Println("Échec : userID introuvable ou type incorrect dans le contexte")
		return 0, errors.New("userID non trouvé dans le contexte")
	}

	log.Printf("Succès : userID récupéré depuis le contexte : %d", userID)
	return userID, nil
}



// AttachTokenToContext extrait le token JWT de l'en-tête Authorization
// et le transmet dans les métadonnées du contexte pour un appel gRPC
func AttachTokenToContext(r *http.Request) (context.Context, error) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, ErrMissingToken
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	return metadata.AppendToOutgoingContext(r.Context(), "authorization", "Bearer "+token), nil
}

var ErrMissingToken = http.ErrNoCookie
