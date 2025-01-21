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

func (s *AlbumService) ListAlbums() ([]models.Album, error) {
	// Étape 1 : Récupérer les albums depuis la base de données
	var albums []models.Album
	err := s.DBManager.DB.Find(&albums).Error
	if err != nil {
		log.Printf("Erreur lors de la récupération des albums depuis la base de données : %v", err)
		return nil, fmt.Errorf("failed to fetch albums from database: %v", err)
	}

	// Étape 2 : Enrichir avec les données du S3 (si nécessaire)
	s3Buckets, err := s.S3Service.ListBuckets()
	if err != nil {
		log.Printf("Erreur lors de la récupération des buckets depuis l'API S3 : %v", err)
	}

	// ajouter une indication si le bucket existe
	bucketExists := make(map[string]bool)
	for _, bucket := range s3Buckets {
		bucketExists[bucket.Name] = true
	}
	for i := range albums {
		if _, exists := bucketExists[albums[i].Name]; exists {
			albums[i].ExistsInS3 = true
		} else {
			albums[i].ExistsInS3 = false
		}
	}

	return albums, nil
}



