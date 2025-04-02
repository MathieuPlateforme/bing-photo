package services

import (
	"GalleryService/internal/db"
	"GalleryService/internal/models"
	"fmt"
	"log"
	"time"
	"strings"
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
	// Générer un nom unique pour le bucket
	album.BucketName = fmt.Sprintf("bucket-%d", time.Now().UnixNano())

	// Étape 1 : Créer un bucket pour l'album
	err := s.S3Service.CreateBucket(album.BucketName)
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

func (s *AlbumService) GetAlbumsByUser(userID uint) ([]models.Album, error) {
	// Étape 1 : Récupérer les albums depuis la base de données
	var albums []models.Album
	err := s.DBManager.DB.Where("user_id = ?", userID).Find(&albums).Error
	if err != nil {
		log.Printf("Erreur lors de la récupération des albums depuis la base de données : %v", err)
		return nil, fmt.Errorf("failed to fetch albums from database: %v", err)
	}

	// Obtenir la liste des buckets S3 existants
	s3Buckets, err := s.S3Service.ListBuckets()
	if err != nil {
		log.Printf("Erreur lors de la récupération des buckets depuis l'API S3 : %v", err)
	}

	// ajouter une indication si le bucket existe
	bucketExists := make(map[string]bool)
	for _, bucket := range s3Buckets {
		bucketExists[strings.TrimSpace(bucket.Name)] = true
	}
	for i := range albums {
		albums[i].ExistsInS3 = bucketExists[strings.TrimSpace(albums[i].BucketName)]
	}

	return albums, nil
}

func (s *AlbumService) UpdateAlbum(id uint, name string, description string) error {
	// Récupérer l'album
	var album models.Album
	if err := s.DBManager.DB.First(&album, id).Error; err != nil {
		return fmt.Errorf("album non trouvé : %v", err)
	}

	// Mettre à jour les champs modifiables
	album.Name = name
	album.Description = description

	// Sauvegarder les modifications
	if err := s.DBManager.DB.Save(&album).Error; err != nil {
		return fmt.Errorf("échec de la mise à jour de l'album : %v", err)
	}

	return nil
}

func (s *AlbumService) DeleteAlbum(albumID uint) error {
	// Récupérer l'album dans la base de données
	var album models.Album
	err := s.DBManager.DB.First(&album, albumID).Error
	if err != nil {
		return fmt.Errorf("album non trouvé : %v", err)
	}

	// Supprimer le bucket associé dans S3
	err = s.S3Service.DeleteBucket(album.BucketName)
	if err != nil {
		return fmt.Errorf("échec de la suppression du bucket S3 : %v", err)
	}

	// Supprimer l'album de la base de données
	err = s.DBManager.DB.Delete(&album).Error
	if err != nil {
		return fmt.Errorf("échec de la suppression de l'album : %v", err)
	}

	return nil
}

func (s *AlbumService) GetPrivateAlbum(userID uint) (*models.Album, error) {
    var album models.Album
    err := s.DBManager.DB.Where("user_id = ? AND is_private = ?", userID, true).First(&album).Error
    if err != nil {
        return nil, fmt.Errorf("échec de la récupération de l'album privé : %v", err)
    }
    return &album, nil
}

func (s *AlbumService) GetMainAlbum(userID uint) (*models.Album, error) {
    var album models.Album
    err := s.DBManager.DB.Where("user_id = ? AND is_main = ?", userID, true).First(&album).Error
    if err != nil {
        return nil, fmt.Errorf("échec de la récupération de l'album privé : %v", err)
    }
    return &album, nil
}




