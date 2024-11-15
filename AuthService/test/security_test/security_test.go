package security_test

import (
	"AuthService/pkg/security"
	"testing"
)

// TestHashPassword vérifie si le mot de passe est correctement haché
func TestHashPassword(t *testing.T) {
	securityService, _ := security.NewSecurityService()

	password := "mySecurePassword"
	hashedPassword := securityService.HashPassword(password)

	if hashedPassword == "" {
		t.Errorf("Le hachage du mot de passe ne doit pas être vide")
	}

	// Vérifier que le mot de passe haché ne soit pas égal au mot de passe en clair
	if hashedPassword == password {
		t.Errorf("Le mot de passe haché ne doit pas être le même que le mot de passe en clair")
	}
}

// TestComparePasswords vérifie si les mots de passe correspondent
func TestComparePasswords(t *testing.T) {
	securityService, _ := security.NewSecurityService()

	password := "mySecurePassword"
	hashedPassword := securityService.HashPassword(password)

	// Comparer le mot de passe haché avec le mot de passe correct
	match := securityService.ComparePasswords(hashedPassword, password)
	if !match {
		t.Errorf("Les mots de passe devraient correspondre")
	}

	// Comparer le mot de passe haché avec un mot de passe incorrect
	wrongPassword := "wrongPassword"
	match = securityService.ComparePasswords(hashedPassword, wrongPassword)
	if match {
		t.Errorf("Les mots de passe ne devraient pas correspondre")
	}
}

// TestGenerateSecureToken vérifie si le jeton sécurisé est généré correctement
func TestGenerateSecureToken(t *testing.T) {
	securityService, _ := security.NewSecurityService()

	token := securityService.GenerateSecureToken()
	if token == "" {
		t.Errorf("Le jeton sécurisé ne doit pas être vide")
	}

	// Vérifier la longueur du jeton (64 caractères pour un jeton de 256 bits en hexadécimal)
	expectedLength := 64
	if len(token) != expectedLength {
		t.Errorf("La longueur du jeton doit être de %d caractères, mais elle est de %d", expectedLength, len(token))
	}
}
