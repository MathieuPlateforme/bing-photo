package email_test

import (
	"AuthService/pkg/email"
	"testing"
)

// TestSendEmailVerification teste la fonction SendEmailVerification
func TestSendEmailVerification(t *testing.T) {

	// Simuler une adresse email de test
	testEmail := "alizeamasse@gmail.com"

	// Exécuter la fonction
	email.SendEmailVerification(testEmail)

	t.Log("Email de vérification envoyé avec succès")
}

// TestSendPasswordResetEmail teste la fonction SendPasswordResetEmail
func TestSendPasswordResetEmail(t *testing.T) {

	// Simuler une adresse email de test
	testEmail := "alizeamasse@gmail.com"

	// Exécuter la fonction
	email.SendPasswordResetEmail(testEmail)

	t.Log("Email de réinitialisation de mot de passe envoyé avec succès")
}
