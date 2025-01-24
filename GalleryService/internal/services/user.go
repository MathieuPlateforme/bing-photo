package services

import (
	"GalleryService/internal/db"
	"GalleryService/internal/models"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// UserService gère les opérations liées aux utilisateurs
type UserService struct {
	DBManager *db.DBManagerService
	S3Service *S3Service
}

// NewUserService initialise un nouveau UserService
func NewUserService(dbManager *db.DBManagerService, S3Service *S3Service) *UserService {
	return &UserService{
		DBManager: dbManager,
		S3Service: S3Service,
	}
}

func (s *UserService) CreateUser(email string, username string) error {
	// Commencer une transaction
	tx := s.DBManager.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Étape 1 : Créer un utilisateur sans album privé pour récupérer son ID
	user := models.User{
		Email:    email,
		Username: username,
	}
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("échec de la création de l'utilisateur : %v", err)
	}

	// Étape 2 : Générer le nom du bucket
	bucketName := fmt.Sprintf("private-album-%d", user.ID)

	// Étape 3 : Créer un album privé avec l'ID utilisateur et le nom du bucket
	privateAlbum := models.Album{
		Name:       fmt.Sprintf("Private Album - %d", user.ID),
		UserID:     user.ID,
		BucketName: bucketName,
		IsPrivate:  true,
		ExistsInS3: false, // Par défaut à false, sera mis à jour après la création du bucket S3
	}
	if err := tx.Create(&privateAlbum).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("échec de la création de l'album privé : %v", err)
	}

	// Étape 4 : Créer le bucket dans S3
	err := s.S3Service.CreateBucket(bucketName)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("échec de la création du bucket S3 : %v", err)
	}

	// Étape 5 : Mettre à jour l'album pour indiquer qu'il existe dans S3
	privateAlbum.ExistsInS3 = true
	if err := tx.Save(&privateAlbum).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("échec de la mise à jour de l'album privé : %v", err)
	}

	// Étape 6 : Associer l'album privé à l'utilisateur
	user.PrivateAlbumID = privateAlbum.ID
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("échec de la mise à jour de l'utilisateur avec l'album privé : %v", err)
	}

	// Commit de la transaction
	return tx.Commit().Error
}


func (s *UserService) VerifyPrivateAlbumPin(userID uint, pin string) error {
	var user models.User
	err := s.DBManager.DB.First(&user, userID).Error
	if err != nil {
		return fmt.Errorf("utilisateur introuvable : %v", err)
	}

	if !compareHashAndPin(user.PrivateAlbumPin, pin) {
		return fmt.Errorf("PIN incorrect")
	}
	return nil
}

func (s *UserService) SetPrivateAlbumPin(userID uint, pin string) error {
	if len(pin) != 6 {
		return fmt.Errorf("le PIN doit contenir exactement 6 chiffres")
	}

	// Hacher le PIN
	hashedPin, err := hashPin(pin)
	if err != nil {
		return fmt.Errorf("échec du hachage du PIN : %v", err)
	}

	// Récupérer l'utilisateur
	var user models.User
	err = s.DBManager.DB.First(&user, userID).Error
	if err != nil {
		return fmt.Errorf("utilisateur introuvable : %v", err)
	}

	// Mettre à jour le PIN de l'utilisateur
	user.PrivateAlbumPin = hashedPin
	if err := s.DBManager.DB.Save(&user).Error; err != nil {
		return fmt.Errorf("échec de la mise à jour du PIN : %v", err)
	}

	return nil
}

// hashPin génère un hash pour un PIN donné
func hashPin(pin string) (string, error) {
	hashedPin, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPin), nil
}

// compareHashAndPin compare un PIN brut avec son hash
func compareHashAndPin(hashedPin, plainPin string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPin), []byte(plainPin))
	return err == nil
}