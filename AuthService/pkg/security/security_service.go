package security

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// SecurityService structure
type SecurityService struct {
}

// NewSecurityService initialise et retourne une nouvelle instance de SecurityService
func NewSecurityService() (*SecurityService, error) {
	fmt.Println("Initializing SecurityService...")
	return &SecurityService{}, nil
}

func HashPassword(password string) string {
	// Logique de hachage de mot de passe
	const cost = 10

	// Hacher le mot de passe avec un coût de 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		fmt.Println("Error hashing password: ", err)
	}
	return string(hashedPassword)
}

func ComparePasswords(hashedPassword string, password string) bool {
	// Logique de comparaison de mot de passe
	return true
}

func GenerateSecureToken() string {
	// Logique de génération de jeton JWT
	return "Token"
}