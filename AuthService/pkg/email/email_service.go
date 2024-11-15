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

func SendEmailVerification(email string) {
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
	body := "Veuillez cliquer sur ce lien pour verifier votre mail : localhost:5050/verify?email=" + email
	message := []byte(subject + "\n" + body)

	// Authentification avec le serveur SMTP

	auth := smtp.PlainAuth("", from, password, smtpHost)

	// envoi du mail

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println("Error sending email: ", err)
	}

	fmt.Println("Email sent successfully to: ", email)
}

func SendPasswordResetEmail(email string) {
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
	body := "Cliquez sur ce lien pour réinitialiser votre mot de passe : [lien de réinitialisation]"
	message := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body)

	// Authentification avec le serveur SMTP

	auth := smtp.PlainAuth("", from, password, smtpHost)

	// envoi du mail

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(message))

	if err != nil {
		fmt.Println("Error sending email: ", err)
	}

	fmt.Println("Email sent successfully to: ", email)
}

