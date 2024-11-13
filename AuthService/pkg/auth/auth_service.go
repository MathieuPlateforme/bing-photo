package auth 

import (
	"authservice/pkg/user"
    "fmt"
)

type AuthService struct{}

// Initialize démarre le service d'authentification
func Initialize() error {
	fmt.Println("Initializing AuthService...")
	return nil
}

func (s *AuthService) LoginWithEmail(u user.User) {
	// Logique de connexion
}

func (s *AuthService) RegisterWithEmail(u user.User) {
	// Logique d'inscription
}



