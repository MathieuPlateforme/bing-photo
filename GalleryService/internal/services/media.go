package services

import (
	"GalleryService/internal/db"
	"GalleryService/internal/models"
	"fmt"
	"io"
	"log"
	"strings"
    "errors"
    "GalleryService/internal/utils"
    "os"
    "time"
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

// ajouter un fichier dans un album
func (s *MediaService) AddMedia(media *models.Media, file io.Reader, fileSize int64) error {
	// Vérifier si l'album existe
	var album models.Album
	if err := s.DBManager.DB.First(&album, media.AlbumID).Error; err != nil {
		return fmt.Errorf("album non trouvé : %v", err)
	}
	media.Path = fmt.Sprintf("%s/%s", album.BucketName, media.Name)

	// Téléverser le fichier dans S3
	err := s.S3Service.UploadFile(media.Path, file, fileSize)
	if err != nil {
		return fmt.Errorf("échec du téléversement du fichier : %v", err)
	}

	// Enregistrer les métadonnées dans la base de données
	if err := s.DBManager.DB.Create(media).Error; err != nil {
		return fmt.Errorf("échec de l'enregistrement des métadonnées : %v", err)
	}

	log.Printf("Média ajouté avec succès : %+v", media)
	return nil
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

func (s *MediaService) DetectSimilarMedia(userID uint, albumID uint) ([]models.Media, error) {
	log.Printf("Détection de médias similaires pour userID=%d dans albumID=%d", userID, albumID)

	// Étape 1 : Vérifier que l'album appartient bien à l'utilisateur
	var album models.Album
	if err := s.DBManager.DB.First(&album, albumID).Error; err != nil {
		return nil, fmt.Errorf("album introuvable pour albumID : %d", albumID)
	}
	if album.UserID != userID {
		return nil, fmt.Errorf("l'utilisateur %d n'a pas accès à cet album", userID)
	}

	// Étape 2 : Récupérer les fichiers dans le bucket S3
	fileNames, err := s.S3Service.GetFilesInAlbum(album.BucketName)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des fichiers : %v", err)
	}
	if len(fileNames) < 2 {
		return nil, errors.New("pas assez de fichiers pour effectuer une comparaison")
	}

	log.Printf("%d fichiers récupérés depuis le bucket %s", len(fileNames), album.BucketName)

	// Étape 3 : Créer un groupe de similarité
	newGroup := models.SimilarGroup{
		UserID:    userID,
		CreatedAt: time.Now(),
	}
	if err := s.DBManager.DB.Create(&newGroup).Error; err != nil {
		return nil, fmt.Errorf("échec de la création du groupe de similarité : %v", err)
	}

	// Étape 4 : Générer les pHash
	hashes := make(map[string]uint64)
	tempFiles := make(map[string]string)

	for _, fileName := range fileNames {
		tempPath, err := s.S3Service.DownloadTempFile(album.BucketName, fileName)
		if err != nil {
			log.Printf("Échec du téléchargement du fichier %s : %v", fileName, err)
			continue
		}
		defer os.Remove(tempPath)

		tempFiles[tempPath] = fileName

		hash, err := utils.ComputePHash(tempPath)
		if err != nil {
			log.Printf("Erreur de calcul de pHash pour %s : %v", tempPath, err)
			continue
		}
		hashes[tempPath] = hash
	}

	if len(hashes) < 2 {
		return nil, errors.New("pas assez de médias valides pour la comparaison")
	}

	// Étape 5 : Comparer les fichiers
	similarEntries := []models.SimilarMedia{}
	threshold := 20
	files := make([]string, 0, len(hashes))
	for file := range hashes {
		files = append(files, file)
	}

	for i := 0; i < len(files); i++ {
		for j := i + 1; j < len(files); j++ {
			dist := utils.HammingDistance(hashes[files[i]], hashes[files[j]])
			log.Printf("📏 Distance entre %s et %s : %d", tempFiles[files[i]], tempFiles[files[j]], dist)
			if dist < threshold {
				log.Printf("Médias similaires détectés : %s et %s", tempFiles[files[i]], tempFiles[files[j]])

				var media1, media2 models.Media
				if err := s.DBManager.DB.Where("album_id = ? AND name = ?", albumID, tempFiles[files[i]]).First(&media1).Error; err != nil {
					continue
				}
				if err := s.DBManager.DB.Where("album_id = ? AND name = ?", albumID, tempFiles[files[j]]).First(&media2).Error; err != nil {
					continue
				}

				similarEntries = append(similarEntries, models.SimilarMedia{
					SimilarGroupID: newGroup.ID,
					MediaID:        media1.ID,
					SimilarityScore: float64(100 - dist),
				}, models.SimilarMedia{
					SimilarGroupID: newGroup.ID,
					MediaID:        media2.ID,
					SimilarityScore: float64(100 - dist),
				})
			}
		}
	}

	// Étape 6 : Sauvegarde des similarités
	if len(similarEntries) > 0 {
		if err := s.DBManager.DB.Create(&similarEntries).Error; err != nil {
			return nil, fmt.Errorf("échec de l'enregistrement des données de similarité : %v", err)
		}
	}

	// Étape 7 : Récupération des médias similaires
	var similarMedia []models.Media
	if len(similarEntries) > 0 {
		var mediaIDs []uint
		for _, entry := range similarEntries {
			mediaIDs = append(mediaIDs, entry.MediaID)
		}
		s.DBManager.DB.Where("id IN ?", mediaIDs).Find(&similarMedia)
	}

	log.Printf("Fin de la détection - %d médias similaires trouvés", len(similarMedia))
	return similarMedia, nil
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

