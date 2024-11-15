package auth

import (
	"fmt"
	"log"
	"AuthService/pkg/user"
	"AuthService/pkg/db"
	"AuthService/pkg/email"
	"AuthService/pkg/google"
	"AuthService/pkg/jwt"
	"AuthService/pkg/security"
	"errors"
	"gorm.io/gorm"
)

// AuthService structure
type AuthService struct {
	DBManager     *db.DBManagerService
	EmailService  *email.EmailService
	GoogleAuthService *google.GoogleAuthService
	JWTService    *jwt.JWTService
	Security      *security.SecurityService
}

// Initialize démarre le service d'authentification
func Initialize() (*AuthService, error) {
	fmt.Println("Initializing AuthService...")

	// Initialisation des services nécessaires
	dbManager, err := db.NewDBManagerService()
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation du service DBManager : %v", err)
	}
	
	emailService, err := email.NewEmailService()
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation du service EmailService : %v", err)
	}

	googleService, err := google.NewGoogleAuthService()
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation du service GoogleAuthService : %v", err)
	}

	jwtService, err := jwt.NewJWTService()
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation du service JWTService : %v", err)
	}

	securityService, err := security.NewSecurityService()
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation du service SecurityService : %v", err)
	}

	authService := &AuthService{
		DBManager:     dbManager,
		EmailService:  emailService,
		GoogleAuthService: googleService,
		JWTService:    jwtService,
		Security:      securityService,
	}
	return authService, nil
}

func (s *AuthService) LoginWithEmail(u user.User, password string) (string, error) {
	// 1. Vérifier si l'utilisateur existe dans la base de données
	var existingUser user.User
	err := s.DBManager.DB.Where("email = ?", u.Email).First(&existingUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("utilisateur introuvable avec cet email : %s", u.Email)
		}
		return "", fmt.Errorf("erreur lors de la recherche de l'utilisateur : %v", err)
	}

	// 2. Comparer le mot de passe fourni avec le mot de passe haché dans la base de données
	if !s.Security.ComparePasswords(existingUser.Password, password) {
		return "", errors.New("mot de passe incorrect")
	}

	// 3. Générer un token JWT pour l'utilisateur
	token, err := s.JWTService.GenerateToken(existingUser.Username)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la génération du token JWT : %v", err)
	}

	// 4. Retourner le token JWT généré
	return token, nil
}


func (s *AuthService) RegisterWithEmail(u user.User) (bool, error) {
	// 1. Vérifier si l'utilisateur existe déjà
	var existingUser user.User
	err := s.DBManager.DB.Where("email = ?", u.Email).First(&existingUser).Error
	if err == nil {
		return false, errors.New("l'utilisateur avec cet email existe déjà")
	}

	// 2. Hacher le mot de passe
	u.Password = s.Security.HashPassword(u.Password)

    // 3. Enregistrer l'utilisateur dans la base de données
	err = u.CreateUser(s.DBManager.DB)
	if err != nil {
		return false, fmt.Errorf("erreur lors de la création de l'utilisateur : %v", err)
	}

	// 4. Envoyer un email de vérification
	// err = s.EmailService.SendEmailVerification(u.Email)
	err = s.EmailService.SendEmailVerification("alizeamasse@gmail.com")
	if err != nil {
		return false, fmt.Errorf("erreur lors de l'envoi de l'email de vérification : %v", err)
	}

	return true, nil
}

func (s *AuthService) ForgotPassword(email string) {
	// Envoie un email de réinitialisation de mot de passe au client.
	// 1. Vérifie si l'adresse email existe dans la base de données.
	// 2. Génère un token de réinitialisation sécurisé avec une durée d'expiration.
	// 3. Stocke le token associé à l'utilisateur pour validation ultérieure.
	// 4. Utilise le service EmailService pour envoyer un lien de réinitialisation.
}

func (s *AuthService) ResetPassword(token string, newPassword string) {
	// Réinitialise le mot de passe de l'utilisateur à l'aide d'un token.
	// 1. Vérifie la validité et l'expiration du token fourni.
	// 2. Recherche l'utilisateur associé au token.
	// 3. Hache le nouveau mot de passe avec le service SecurityService.
	// 4. Met à jour le mot de passe dans la base de données.
	// 5. Invalide le token après utilisation pour des raisons de sécurité.
}

func (s *AuthService) LoginWithGoogle() {
	// Logique de connexion avec Google
}

func (s *AuthService) ValidateToken(token string) {
	// Valide un token JWT pour vérifier l'authenticité de l'utilisateur.
	// 1. Vérifie la signature du token avec la clé secrète.
	// 2. Vérifie la validité et l'expiration du token.
	// 3. Extrait les informations utilisateur du token.
	// 4. Renvoie l'identifiant de l'utilisateur pour une utilisation ultérieure.
}

func (s *AuthService) GeneratePasswordResetToken() {
	// Génère un token de réinitialisation de mot de passe sécurisé pour l'utilisateur.
	// 1. Crée un token unique avec une durée d'expiration.
	// 2. Stocke le token associé à l'utilisateur pour validation ultérieure.
	// 3. Renvoie le token généré pour être envoyé par e-mail.
}

func (s *AuthService) Logout() {
	// Logique de déconnexion
}






