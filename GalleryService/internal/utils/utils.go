package utils

import (
	"context"
	"errors"
	"log"
	"GalleryService/internal/middleware"

)

func GetUserIDFromContext(ctx context.Context) (uint, error) {
	log.Println("Tentative de récupération du userID à partir du contexte...")

	// Inspecter le contenu complet du contexte
	log.Printf("Contenu du contexte : %+v\n", ctx)

	// Vérifier si la clé userID existe dans le contexte
	userID, ok := ctx.Value(middleware.UserIDKey).(uint)
	if !ok {
		log.Println("Échec : userID introuvable ou type incorrect dans le contexte")
		log.Printf("Valeurs possibles dans le contexte : %+v\n", ctx)
		return 0, errors.New("userID non trouvé dans le contexte")
	}

	log.Printf("Succès : userID récupéré depuis le contexte : %d", userID)
	return userID, nil
}
