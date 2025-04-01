package email

import (
	"fmt"
	"net/smtp"
	"os"
)

type EmailService struct {

}

// NewEmailService initialise et retourne une nouvelle instance d'EmailService
func NewEmailService() (*EmailService, error) {
	fmt.Println("Initializing EmailService...")
	return &EmailService{}, nil
}

func (s *EmailService) SendEmailVerification(email string, token string) error {
	// Logique pour envoyer un email de vérification
	from := "alizeamasse@gmail.com"
	password := os.Getenv("APP_MAIL_PASSWORD")

	// Configuration de l'authentification SMTP

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Destinataire

	to := []string{email}

	// Corps du mail

	subject := "Sujet : Vérification de l'adresse email\n"
	body := "Veuillez cliquer sur ce lien pour verifier votre mail : http://localhost:3000/verify?token=" + token
	message := []byte(subject + "\n" + body)

	// Authentification avec le serveur SMTP

	auth := smtp.PlainAuth("", from, password, smtpHost)

	// envoi du mail

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println("Erreur lors de l'envoi de l'email :", err)
		return err
	}

	fmt.Println("Email envoyé avec succès à :", email)
	return nil
}

func (s *EmailService) SendPasswordResetEmail(email, resetLink string) error {
	// Logique pour envoyer un email de réinitialisation de mot de passe
	from := "alizeamasse@gmail.com"
	password := os.Getenv("APP_MAIL_PASSWORD")

	// Configuration de l'authentification SMTP

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	
	// Destinataire

	to := []string{email}

	// Corps du mail

	subject := "Réinitialisation de votre mot de passe"
	body := "Cliquez sur ce lien pour réinitialiser votre mot de passe : " + resetLink
	message := []byte(subject + "\n" + body)

	// Authentification avec le serveur SMTP

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println("Erreur lors de l'envoi de l'email :", err)
		return err
	}

	fmt.Println("Email envoyé avec succès à :", email)
	return nil
}

