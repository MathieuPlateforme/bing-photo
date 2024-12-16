package reset_password_test

import (
	"testing"
	"AuthService/services/auth"
	"AuthService/models"
)

func TestResetPassword(t *testing.T) {
	// Initialiser le service d'authentification
	authService, err := auth.Initialize()
	if err != nil {
		t.Fatalf("Erreur lors de l'initialisation du service d'authentification : %v", err)
	}

	// Simuler un utilisateur existant
	existingUser := &models.User{
		Email:      "testuser@example.com",
		ResetToken: "valid_token",
		Password:   authService.Security.HashPassword("OldPassword123"),
	}

	// Ajouter l'utilisateur à la base de données
	err = authService.DBManager.DB.Create(existingUser).Error
	if err != nil {
		t.Fatalf("Erreur lors de la création de l'utilisateur de test : %v", err)
	}

	// Test du changement de mot de passe réussi
	err = authService.ResetPassword("testuser@example.com", "valid_token", "NewStrongPassword123")
	if err != nil {
		t.Errorf("Échec du changement de mot de passe : %v", err)
	}

	// Vérifier la mise à jour dans la base de données
	var updatedUser models.User
	err = authService.DBManager.DB.Where("email = ?", existingUser.Email).First(&updatedUser).Error
	if err != nil {
		t.Fatalf("Erreur lors de la récupération de l'utilisateur mis à jour : %v", err)
	}

	// Vérifier que le mot de passe a bien été modifié
	if !authService.Security.ComparePasswords(updatedUser.Password, "NewStrongPassword123") {
		t.Error("Le mot de passe n'a pas été mis à jour correctement")
	}

	// Vérifier que le token a été invalidé
	if updatedUser.ResetToken != "" {
		t.Error("Le token de réinitialisation n'a pas été invalidé")
	}

	// Test avec un mauvais token
	err = authService.ResetPassword("testuser@example.com", "invalid_token", "SomePassword")
	if err == nil || err.Error() != "token invalide ou expiré" {
		t.Error("test de token invalide réussi")
	}

	// Test avec un utilisateur inexistant
	err = authService.ResetPassword("unknownuser@examplom", "valid_token", "SomePassword")
	if err == nil || err.Error() != "utilisateur non trouvé" {
		t.Error("test d'utilisateur inexistant réussi")
	}
}
