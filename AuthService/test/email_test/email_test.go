package email_test

import (
	"AuthService/pkg/email"
	"testing"
)

// TestSendEmailVerification teste la fonction SendEmailVerification
func TestSendEmailVerification(t *testing.T) {
	
	// Simuler une adresse email de test
	testEmail := "alizeamasse@gmail.com"

	// Initialiser le service d'email
	emailService, err := email.NewEmailService()
	if err != nil {
		t.Fatalf("Erreur lors de l'initialisation du service email : %v", err)
	}

	// Exécuter la fonction
	err = emailService.SendEmailVerification(testEmail)
	if err != nil {
		t.Errorf("Échec de l'envoi de l'email de vérification : %v", err)
	} else {
		t.Log("Email de vérification envoyé avec succès")
	}
}

// TestSendPasswordResetEmail teste la fonction SendPasswordResetEmail
func TestSendPasswordResetEmail(t *testing.T) {
	
	// Simuler une adresse email de test
	testEmail := "alizeamasse@gmail.com"
	resetLink := "http://localhost:8080/reset-password?token=test_token"

	// Initialiser le service d'email
	emailService, err := email.NewEmailService()
	if err != nil {
		t.Fatalf("Erreur lors de l'initialisation du service email : %v", err)
	}

	// Exécuter la fonction
	err = emailService.SendPasswordResetEmail(testEmail, resetLink)
	if err != nil {
		t.Errorf("Échec de l'envoi de l'email de réinitialisation : %v", err)
	} else {
		t.Log("Email de réinitialisation de mot de passe envoyé avec succès")
	}
}
