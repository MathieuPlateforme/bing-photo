package security

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"crypto/rand"
	"encoding/hex"
	"log"
	"gorm.io/gorm"
)

// SecurityService structure
type SecurityService struct {
	DB *gorm.DB
}

// NewSecurityService initialise et retourne une nouvelle instance de SecurityService
func NewSecurityService() (*SecurityService, error) {
	fmt.Println("Initializing SecurityService...")
	return &SecurityService{}, nil
}

func (s *SecurityService) HashPassword(password string) string {
	// Logique de hachage de mot de passe
	const cost = 10

	// Hacher le mot de passe avec un coût de 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		fmt.Println("Error hashing password: ", err)
	}
	return string(hashedPassword)
}

func (s *SecurityService) ComparePasswords(hashedPassword string, password string) bool {
	// Logique de comparaison de mot de passe
	hashedPasswordBytes := []byte(hashedPassword)
	passwordBytes := []byte(password)

	err := bcrypt.CompareHashAndPassword(hashedPasswordBytes, passwordBytes)
	if err != nil {
		fmt.Println("Passwords do not match")
		return false
	}

	return true
}

func (s *SecurityService) GenerateSecureToken() string {
	// Logique de génération de jeton 

	// Taille du jeton en octets (32 octets = 256 bits)
		tokenSize := 32
		token := make([]byte, tokenSize)
	
		// Génère des octets aléatoires sécurisés
		_, err := rand.Read(token)
		if err != nil {
			log.Fatalf("Erreur while generate token : %v", err)
		}
	
		// Convertit les octets en une chaîne hexadécimale
		return hex.EncodeToString(token)
}

func (s *SecurityService) GeneratePasswordResetLink(email string, updateTokenFunc func(email, token string) error) (string, error) {
	// Logique pour générer un lien de réinitialisation de mot de passe
    token := s.GenerateSecureToken()

    // Appelle la fonction de mise à jour du token, fournie comme argument
    if err := updateTokenFunc(email, token); err != nil {
        return "", fmt.Errorf("erreur lors de la mise à jour du token de réinitialisation : %v", err)
    }

	//Générer le lien de réinitialisation de mot de passe

	resetLink := "http://localhost:5050/reset-password?token=" + token + "&email=" + email
    return resetLink, nil
}