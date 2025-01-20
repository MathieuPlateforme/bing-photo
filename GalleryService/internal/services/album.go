package services

import (
	"GalleryService/internal/db"
	"GalleryService/internal/models"
	"fmt"
	"log"
)

type AlbumService struct {
	DBManager *db.DBManagerService
	S3Service  *S3Service
}

// NewAlbumService initialise le service AlbumService
func NewAlbumService(dbManager *db.DBManagerService, S3Service *S3Service) *AlbumService {
	return &AlbumService{
		DBManager: dbManager,
		S3Service:  S3Service,
	}
}

func (s *AlbumService) CreateAlbum(album models.Album) error {
	// Étape 1 : Créer un bucket pour l'album
	err := s.S3Service.CreateBucket(album.Name)
	if err != nil {
		return fmt.Errorf("failed to create bucket: %v", err)
	}
	log.Printf("Bucket '%s' created successfully", album.Name)

	// Étape 2 : Sauvegarder l'album dans la base de données
	log.Printf("Attempting to save album: %+v", album)
	if err := s.DBManager.DB.Create(&album).Error; err != nil {
		log.Printf("Failed to save album in database: %v", err)
		return fmt.Errorf("failed to save album: %v", err)
	}
	log.Printf("Album '%s' saved successfully in database", album.Name)

	return nil
}

