package utils

import (
	"context"
	"net/http"
	"strings"
	"google.golang.org/grpc/metadata"

)
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
