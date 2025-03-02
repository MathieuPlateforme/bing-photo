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
    media.IsPrivate = true
    media.Path = fmt.Sprintf("%s/%s", targetBucket, sourceKey)

    // Sauvegarder les modifications
    if err := s.DBManager.DB.Save(&media).Error; err != nil {
        return fmt.Errorf("échec de la mise à jour du média : %v", err)
    }

    return nil
}

func (s *MediaService) DownloadMedia(mediaID uint, w io.Writer) error {
	// Récupérer le média à partir de la base de données
	var media models.Media
	if err := s.DBManager.DB.First(&media, mediaID).Error; err != nil {
		return fmt.Errorf("média non trouvé pour l'ID %d : %v", mediaID, err)
	}

	// Vérifier si l'album existe
	var album models.Album
	if err := s.DBManager.DB.First(&album, media.AlbumID).Error; err != nil {
		return fmt.Errorf("album non trouvé pour l'ID %d : %v", media.AlbumID, err)
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
    log.Printf("Starting detection of similar media for userID=%d in albumID=%d", userID, albumID)

    var album models.Album
    if err := s.DBManager.DB.First(&album, albumID).Error; err != nil {
        return nil, fmt.Errorf("album not found")
    }
    if album.UserID != userID {
        return nil, fmt.Errorf("access denied")
    }

    fileNames, err := s.S3Service.GetFilesInAlbum(album.BucketName)
    if err != nil {
        return nil, fmt.Errorf("error retrieving files")
    }
    if len(fileNames) < 2 {
        return nil, errors.New("not enough files to compare")
    }

    log.Printf("Number of files retrieved for albumID=%d: %d", albumID, len(fileNames))

    newGroup := models.SimilarGroup{
        UserID:    userID,
        CreatedAt: time.Now(),
    }
    if err := s.DBManager.DB.Create(&newGroup).Error; err != nil {
        return nil, fmt.Errorf("failed to create similarity group: %v", err)
    }

    hashes := make(map[string]uint64)
    tempFiles := make(map[string]string)

    for _, fileName := range fileNames {
        tempFilePath, err := s.S3Service.DownloadTempFile(album.BucketName, fileName)
        if err != nil {
            log.Printf("Download error for %s: %v", fileName, err)
            continue
        }
        tempFiles[tempFilePath] = fileName

        hash, err := utils.ComputePHash(tempFilePath)
        if err != nil {
            log.Printf("Error computing pHash for %s: %v", tempFilePath, err)
            continue
        }
        hashes[tempFilePath] = hash
    }

    if len(hashes) < 2 {
        return nil, errors.New("not enough valid media for comparison")
    }

    similarEntries := []models.SimilarMedia{}
    threshold := 20
    files := make([]string, 0, len(hashes))
    for file := range hashes {
        files = append(files, file)
    }

    for i := 0; i < len(files); i++ {
        for j := i + 1; j < len(files); j++ {
            dist := utils.HammingDistance(hashes[files[i]], hashes[files[j]])
            log.Printf("Hamming distance between %s and %s: %d", files[i], files[j], dist)
            if dist < threshold {
                log.Printf("Similar media found: %s and %s (distance: %d)", files[i], files[j], dist)

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
                })
                similarEntries = append(similarEntries, models.SimilarMedia{
                    SimilarGroupID: newGroup.ID,
                    MediaID:        media2.ID,
                    SimilarityScore: float64(100 - dist),
                })
            }
        }
    }

    if len(similarEntries) > 0 {
        if err := s.DBManager.DB.Create(&similarEntries).Error; err != nil {
            return nil, fmt.Errorf("failed to save similarity data: %v", err)
        }
    }

    for tempPath := range tempFiles {
        _ = os.Remove(tempPath)
    }

    var similarMedia []models.Media
    if len(similarEntries) > 0 {
        var mediaIDs []uint
        for _, entry := range similarEntries {
            mediaIDs = append(mediaIDs, entry.MediaID)
        }
        s.DBManager.DB.Where("id IN (?)", mediaIDs).Find(&similarMedia)
    }

    log.Printf("Similar media detection completed for albumID=%d: %d media found", albumID, len(similarMedia))

    return similarMedia, nil
}

