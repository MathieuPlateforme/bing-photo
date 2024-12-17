package jwt_test

import (
	"AuthService/pkg/jwt"
	"testing"
)

// TestGenerateToken teste la fonction GenerateToken
func TestGenerateToken(t *testing.T) {
	// Créer une instance de JWTService
	jwtService, err := jwt.NewJWTService()
	if err != nil {
		t.Fatalf("Erreur lors de l'initialisation de JWTService : %v", err)
	}

	// Simuler un nom d'utilisateur de test
	testUsername := "test"

	// Exécuter la fonction
	token, err := jwtService.GenerateToken(testUsername)
	if err != nil {
		t.Errorf("Erreur lors de la génération du token : %v", err)
		return
	}

	t.Log("Token généré avec succès :", token)
}

// TestVerifyToken teste la fonction VerifyToken
func TestVerifyToken(t *testing.T) {
	// Créer une instance de JWTService
	jwtService, err := jwt.NewJWTService()
	if err != nil {
		t.Fatalf("Erreur lors de l'initialisation de JWTService : %v", err)
	}

	// Simuler un nom d'utilisateur de test
	testUsername := "test"

	// Générer un token valide pour le test
	token, err := jwtService.GenerateToken(testUsername)
	if err != nil {
		t.Errorf("Erreur lors de la génération du token : %v", err)
		return
	}

	// Vérifier le token
	claims, err := jwtService.VerifyToken(token)
	if err != nil {
		t.Errorf("Erreur lors de la vérification du token : %v", err)
		return
	}

	// Utiliser les informations contenues dans les claims
	t.Logf("Token vérifié avec succès, claims: %v", claims)

}

// TestVerifyTokenInvalid teste la fonction VerifyToken avec un token invalide
func TestVerifyTokenInvalid(t *testing.T) {
	// Créer une instance de JWTService
	jwtService, err := jwt.NewJWTService()
	if err != nil {
		t.Fatalf("Erreur lors de l'initialisation de JWTService : %v", err)
	}

	// Simuler un token de test invalide
	token := "invalid"

	// Vérifier le token
	claims, err := jwtService.VerifyToken(token)
	if err != nil {
		t.Errorf("Erreur lors de la vérification du token : %v", err)
		return
	}

	// Utiliser les informations contenues dans les claims
	t.Logf("Token vérifié avec succès, claims: %v", claims)

}
