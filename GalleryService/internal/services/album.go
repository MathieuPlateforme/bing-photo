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
	// Récupérer l'utilisateur pour accéder à ses albums spéciaux
	var user models.User
	if err := s.DBManager.DB.First(&user, userID).Error; err != nil {
		log.Printf("Utilisateur introuvable pour userID %d : %v", userID, err)
		return nil, fmt.Errorf("utilisateur non trouvé")
	}

	// Récupérer tous les albums sauf le main et le privé, et précharger les médias associés
	var albums []models.Album
	err := s.DBManager.DB.
		Preload("Media"). // <- préchargement des médias liés
		Where("user_id = ? AND id NOT IN (?, ?)", userID, user.PrivateAlbumID, user.MainAlbumID).
		Find(&albums).Error

	if err != nil {
		log.Printf("Erreur lors de la récupération des albums : %v", err)
		return nil, fmt.Errorf("échec de la récupération des albums")
	}

	// Vérifier l'existence des buckets S3
	s3Buckets, err := s.S3Service.ListBuckets()
	if err != nil {
		log.Printf("Erreur lors de la récupération des buckets S3 : %v", err)
	}

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


// Méthode publique appelée par le gRPC
func (s *AlbumService) GetPrivateAlbum(userID uint, albumType string) (*models.Album, error) {
	if albumType == "main" {
		return s.getMainAlbum(userID)
	}
	return s.getPrivateAlbum(userID)
}

// Méthode interne pour récupérer l'album privé
func (s *AlbumService) getPrivateAlbum(userID uint) (*models.Album, error) {
	var album models.Album
	if err := s.DBManager.DB.Where("user_id = ? AND is_private = true", userID).First(&album).Error; err != nil {
		return nil, fmt.Errorf("album privé non trouvé : %v", err)
	}
	return &album, nil
}

// Méthode interne pour récupérer l'album principal
func (s *AlbumService) getMainAlbum(userID uint) (*models.Album, error) {
	var album models.Album
	if err := s.DBManager.DB.Where("user_id = ? AND is_main = true", userID).First(&album).Error; err != nil {
		return nil, fmt.Errorf("album principal non trouvé : %v", err)
	}
	return &album, nil
}


