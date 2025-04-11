package services

import (
	"GalleryService/internal/db"
	"GalleryService/internal/models"
	"fmt"
	"io"
	"log"
	"strings"
    "GalleryService/internal/utils"
    "os"
	"strconv"
)

type MediaService struct {
	DBManager *db.DBManagerService
	S3Service *S3Service
}

// NewMediaService initialise un MediaService
func NewMediaService(dbManager *db.DBManagerService, s3Service *S3Service) *MediaService {
	return &MediaService{
		DBManager: dbManager,
		S3Service: s3Service,
	}
}
func (s *MediaService) AddMedia(media *models.Media, file io.Reader, fileSize int64) error {
	log.Printf(" Début d'ajout du média : %+v", media)

	// 1. Vérifier que l'album existe
	var album models.Album
	if err := s.DBManager.DB.First(&album, media.AlbumID).Error; err != nil {
		log.Printf(" Album non trouvé pour ID %d : %v", media.AlbumID, err)
		return fmt.Errorf("album non trouvé : %v", err)
	}
	log.Printf("Album trouvé : %s", album.Name)

	// 2. Préparer le chemin du média
	media.Path = fmt.Sprintf("%s/%s", album.BucketName, media.Name)
	log.Printf("Chemin du fichier : %s", media.Path)

	// 3. Sauvegarder le fichier temporairement
	tempFilePath := fmt.Sprintf("/tmp/%s", media.Name)
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		log.Printf("Erreur création fichier temporaire %s : %v", tempFilePath, err)
		return fmt.Errorf("échec de la création du fichier temporaire : %v", err)
	}
	defer tempFile.Close()
	defer os.Remove(tempFilePath)
	log.Printf("Fichier temporaire créé : %s", tempFilePath)

	// 4. Uploader le fichier en S3 et copier localement avec TeeReader
	tee := io.TeeReader(file, tempFile)
	if err := s.S3Service.UploadFile(media.Path, tee, fileSize); err != nil {
		log.Printf("Échec de l'upload S3 : %v", err)
		return fmt.Errorf("échec du téléversement du fichier : %v", err)
	}
	log.Printf("Upload S3 réussi")

	// 5. Calculer le pHash à partir du fichier temporaire
	hash, err := utils.ComputePHash(tempFilePath)
	if err != nil {
		log.Printf("Erreur calcul pHash : %v", err)
		return fmt.Errorf("échec du calcul du pHash : %v", err)
	}
	log.Printf("pHash calculé : %d", hash)

	// 6. Affecter le hash au média
	media.Hash = ptr(fmt.Sprintf("%d", hash))
	log.Printf("Hash converti en string et assigné : %s", *media.Hash)

	// 7. Enregistrer les métadonnées
	log.Printf("📥 Enregistrement du média en base : %+v", media)
	if err := s.DBManager.DB.Create(media).Error; err != nil {
		log.Printf("Erreur lors de la création en base : %v", err)
		return fmt.Errorf("échec de l'enregistrement des métadonnées : %v", err)
	}
	log.Printf("Média enregistré avec succès")

	return nil
}

// Helper pour pointer une string
func ptr(s string) *string {
	return &s
}


func (s *MediaService) GetMediaByUser(userID uint) ([]models.Media, error) {
	var mediaList []models.Media

	// Charger les médias et les albums associés
	err := s.DBManager.DB.
		Preload("Album").
		Joins("JOIN albums ON albums.id = media.album_id").
		Where("albums.user_id = ?", userID).
		Find(&mediaList).Error
	if err != nil {
		return nil, fmt.Errorf("échec de la récupération des médias pour l'utilisateur %d : %v", userID, err)
	}

	// Vérification de l'existence des buckets
	s3Buckets, err := s.S3Service.ListBuckets()
	if err != nil {
		return nil, fmt.Errorf("échec de la récupération des buckets depuis l'API S3-like : %v", err)
	}

	bucketExists := make(map[string]bool)
	for _, bucket := range s3Buckets {
		bucketExists[strings.TrimSpace(bucket.Name)] = true
	}

	for i := range mediaList {
		// Vérifier si l'album est chargé et mettre à jour ExistsInS3
		if mediaList[i].Album != nil {
			mediaList[i].Album.ExistsInS3 = bucketExists[strings.TrimSpace(mediaList[i].Album.BucketName)]
		}
	}

	return mediaList, nil
}

func (s *MediaService) GetPrivateMedia(userID uint) ([]models.Media, error) {
    var mediaList []models.Media

    // Charger les médias associés aux albums de l'utilisateur
    err := s.DBManager.DB.
        Preload("Album"). // Charger les détails des albums associés
        Joins("JOIN albums ON albums.id = media.album_id").
        Where("albums.user_id = ?", userID).
		Where("media.is_private = ?", "true").
        Find(&mediaList).Error

    if err != nil {
        return nil, fmt.Errorf("échec de la récupération des médias pour l'utilisateur %d : %v", userID, err)
    }

    // Récupérer les buckets existants depuis l'API S3-like
    bucketExists := make(map[string]bool)
    s3Buckets, err := s.S3Service.ListBuckets()
    if err != nil {
        return nil, fmt.Errorf("échec de la récupération des buckets depuis S3 : %v", err)
    }

    for _, bucket := range s3Buckets {
        bucketExists[bucket.Name] = true
    }

    // Vérifier si les albums des médias existent dans S3
	for i := range mediaList {
		// Vérifier si l'album est chargé et mettre à jour ExistsInS3
		if mediaList[i].Album != nil {
			mediaList[i].Album.ExistsInS3 = bucketExists[strings.TrimSpace(mediaList[i].Album.BucketName)]
		}
	}

    return mediaList, nil
}

func (s *MediaService) MarkAsPrivate(mediaID uint, userID uint) error {
    // Récupérer le média à partir de son ID
    var media models.Media
    if err := s.DBManager.DB.First(&media, mediaID).Error; err != nil {
        return fmt.Errorf("média introuvable pour mediaID : %d", mediaID)
    }

    // Vérifier si l'utilisateur est propriétaire du média
    var album models.Album
    if err := s.DBManager.DB.First(&album, media.AlbumID).Error; err != nil {
        return fmt.Errorf("album introuvable pour albumID : %d", media.AlbumID)
    }
    if album.UserID != userID {
        return fmt.Errorf("l'utilisateur %d n'est pas propriétaire de ce média", userID)
    }

    // Récupérer l'album privé de l'utilisateur
    var user models.User
    if err := s.DBManager.DB.First(&user, userID).Error; err != nil {
        return fmt.Errorf("utilisateur introuvable pour userID : %d", userID)
    }
    var privateAlbum models.Album
    if err := s.DBManager.DB.First(&privateAlbum, user.PrivateAlbumID).Error; err != nil {
        return fmt.Errorf("album privé introuvable pour userID : %d", userID)
    }

    // Construire les paramètres pour le déplacement dans S3
    sourceBucket := album.BucketName
    sourceKey := media.Name
    targetBucket := privateAlbum.BucketName

    // Déplacer le fichier dans S3
    if err := s.S3Service.MoveObject(sourceBucket, sourceKey, targetBucket); err != nil {
        return fmt.Errorf("échec du déplacement du média dans S3 : %v", err)
    }

    // Mettre à jour le média pour qu'il soit associé à l'album privé
    media.AlbumID = privateAlbum.ID
    media.Path = fmt.Sprintf("%s/%s", targetBucket, sourceKey)

    // Sauvegarder les modifications
    if err := s.DBManager.DB.Save(&media).Error; err != nil {
        return fmt.Errorf("échec de la mise à jour du média : %v", err)
    }

    return nil
}

func (s *MediaService) DownloadMedia(mediaID uint, userID uint, w io.Writer) error {
	// Récupérer le média à partir de la base de données
	var media models.Media
	if err := s.DBManager.DB.First(&media, mediaID).Error; err != nil {
		return fmt.Errorf("média non trouvé pour l'ID %d : %v", mediaID, err)
	}

	// Récupérer l'album auquel appartient le média
	var album models.Album
	if err := s.DBManager.DB.First(&album, media.AlbumID).Error; err != nil {
		return fmt.Errorf("album non trouvé pour l'ID %d : %v", media.AlbumID, err)
	}

	// Vérifier que l'utilisateur est bien le propriétaire
	if album.UserID != userID {
		return fmt.Errorf("l'utilisateur %d n'est pas autorisé à accéder à ce média", userID)
	}

	// Télécharger le fichier depuis S3
	mediaPath := fmt.Sprintf("%s/%s", album.BucketName, media.Name)
	if err := s.S3Service.DownloadFile(mediaPath, w); err != nil {
		return fmt.Errorf("échec du téléchargement du fichier : %v", err)
	}

	log.Printf("Média téléchargé avec succès : %s", mediaPath)
	return nil
}

func (s *MediaService) DeleteMedia(mediaID uint, userID uint) error {
    // Récupérer le média à partir de son ID
    var media models.Media
    if err := s.DBManager.DB.First(&media, mediaID).Error; err != nil {
        return fmt.Errorf("média introuvable pour mediaID : %d", mediaID)
    }

    // Vérifier si l'utilisateur est propriétaire de l'album contenant le média
    var album models.Album
    if err := s.DBManager.DB.First(&album, media.AlbumID).Error; err != nil {
        return fmt.Errorf("album introuvable pour albumID : %d", media.AlbumID)
    }
    if album.UserID != userID {
        return fmt.Errorf("l'utilisateur %d n'est pas propriétaire de ce média", userID)
    }

    // Appeler la méthode du S3Service pour supprimer l'objet
    if err := s.S3Service.DeleteObject(album.BucketName, media.Name); err != nil {
        return fmt.Errorf("échec de la suppression du média dans S3 : %v", err)
    }

    // Supprimer le média de la base de données
    if err := s.DBManager.DB.Delete(&media).Error; err != nil {
        return fmt.Errorf("échec de la suppression du média de la base de données : %v", err)
    }

    log.Printf("Média supprimé avec succès : mediaID=%d, path=%s", mediaID, media.Path)
    return nil
}

func (s *MediaService) DetectSimilarMedia(userID uint, albumID uint) ([][]models.Media, error) {
	log.Printf("Début de la détection de médias similaires pour userID=%d, albumID=%d", userID, albumID)

	// Étape 1 : Vérification de l'accès à l'album
	var album models.Album
	if err := s.DBManager.DB.First(&album, albumID).Error; err != nil {
		log.Printf("Album introuvable : %v", err)
		return nil, fmt.Errorf("album introuvable pour albumID : %d", albumID)
	}
	if album.UserID != userID {
		log.Printf("Accès refusé à l'album %d pour l'utilisateur %d", albumID, userID)
		return nil, fmt.Errorf("l'utilisateur %d n'a pas accès à cet album", userID)
	}
	log.Printf("Accès à l'album confirmé")

	// Étape 2 : Récupération des médias avec pHash
	var medias []models.Media
	if err := s.DBManager.DB.Where("album_id = ?", albumID).Find(&medias).Error; err != nil {
		log.Printf("Erreur lors de la récupération des médias : %v", err)
		return nil, err
	}
	log.Printf(" %d médias récupérés depuis l'album", len(medias))

	// Étape 3 : Création de la map hash → []Media
	hashes := make(map[uint64][]models.Media)
	for _, m := range medias {
		if m.Hash == nil {
			log.Printf("Média %d (%s) sans hash, ignoré", m.ID, m.Name)
			continue
		}
		parsed, err := strconv.ParseUint(*m.Hash, 10, 64)
		if err != nil {
			log.Printf("Erreur de parsing du hash pour media %d : %v", m.ID, err)
			continue
		}
		log.Printf(" Média %d ajouté avec pHash %d", m.ID, parsed)
		hashes[parsed] = append(hashes[parsed], m)
	}
	log.Printf("%d pHash uniques analysés", len(hashes))

	// Étape 4 : Comparaison des pHash et regroupement par similarité
	
	var similarGroups [][]models.Media
	visited := make(map[uint]bool)
	threshold := 20
	hashesList := make([]uint64, 0, len(hashes))
	for hash, group := range hashes {
		if len(group) > 1 {
			log.Printf("Duplication exacte détectée pour hash %d : %d médias", hash, len(group))
			similarGroups = append(similarGroups, group)
			for _, m := range group {
				visited[m.ID] = true
			}
		}
	}
	for h := range hashes {
		hashesList = append(hashesList, h)
	}

	for i := 0; i < len(hashesList); i++ {
		for j := i + 1; j < len(hashesList); j++ {
			h1 := hashesList[i]
			h2 := hashesList[j]
			dist := utils.HammingDistance(h1, h2)
			log.Printf("Distance Hamming entre %d et %d : %d", h1, h2, dist)

			if dist < threshold {
				log.Printf("Groupe détecté : pHash %d et %d sont similaires", h1, h2)
				group := []models.Media{}

				for _, m := range hashes[h1] {
					if !visited[m.ID] {
						log.Printf("Ajout de media %d (%s) au groupe", m.ID, m.Name)
						visited[m.ID] = true
						group = append(group, m)
					}
				}
				for _, m := range hashes[h2] {
					if !visited[m.ID] {
						log.Printf("Ajout de media %d (%s) au groupe", m.ID, m.Name)
						visited[m.ID] = true
						group = append(group, m)
					}
				}
				if len(group) > 1 {
					log.Printf("Groupe finalisé avec %d médias", len(group))
					similarGroups = append(similarGroups, group)
				}
			}
		}
	}

	log.Printf("Détection terminée : %d groupes similaires trouvés", len(similarGroups))
	return similarGroups, nil
}




func (s *MediaService) GetMediaByAlbum(albumID uint) ([]models.Media, error) {
	var medias []models.Media

	// Récupérer tous les médias associés à l'album donné
	if err := s.DBManager.DB.Where("album_id = ?", albumID).Find(&medias).Error; err != nil {
		log.Printf("Erreur lors de la récupération des médias pour l'album %d : %v", albumID, err)
		return nil, fmt.Errorf("échec de la récupération des médias pour l'album %d", albumID)
	}

	return medias, nil
}

