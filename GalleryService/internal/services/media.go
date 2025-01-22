package services

import (
	"GalleryService/internal/db"
	"GalleryService/internal/models"
	"fmt"
	"io"
	"log"
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