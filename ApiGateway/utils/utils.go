package utils

import (
	"context"
	"errors"
	"log"
	"net/http"
	"google.golang.org/grpc/metadata"
)
var ErrMissingToken = errors.New("token manquant ou invalide")

func AttachTokenToContext(r *http.Request) (context.Context, error) {
	authHeader := r.Header.Get("Authorization")
	log.Printf("Tentative d'extraction du token depuis l'en-tête Authorization : %s\n", authHeader)


	token := authHeader
	return metadata.AppendToOutgoingContext(r.Context(), "authorization", "Bearer "+token), nil
}