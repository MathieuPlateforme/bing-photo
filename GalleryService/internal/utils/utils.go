package utils

import (
	"context"
	"errors"
	"log"
	"GalleryService/internal/middleware"
    "github.com/corona10/goimagehash"
    "image"
    _ "image/jpeg"
    _ "image/png"
    "os"

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

func ComputePHash(imagePath string) (uint64, error) {
    log.Printf("Début du calcul de pHash pour %s", imagePath)

    file, err := os.Open(imagePath)
    if err != nil {
        log.Printf("Erreur lors de l'ouverture du fichier %s : %v", imagePath, err)
        return 0, err
    }
    defer file.Close()

    img, _, err := image.Decode(file)
    if err != nil {
        log.Printf("Erreur lors du décodage de l'image %s : %v", imagePath, err)
        return 0, err
    }

	hash, err := goimagehash.AverageHash(img) 
    if err != nil {
        log.Printf("Erreur lors du calcul du pHash pour %s : %v", imagePath, err)
        return 0, err
    }

    log.Printf("pHash calculé pour %s : %d", imagePath, hash.GetHash())
    return hash.GetHash(), nil
}

func HammingDistance(hash1, hash2 uint64) int {
    dist := 0
    x := hash1 ^ hash2

    log.Printf("Calcul de la distance de Hamming entre hash1=%d et hash2=%d", hash1, hash2)

    for x > 0 {
        dist += int(x & 1)
        x >>= 1
    }

    log.Printf("Distance de Hamming calculée : %d", dist)
    return dist
}
