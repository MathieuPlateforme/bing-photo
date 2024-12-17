package auth

import (
	"fmt"
	"log"
	"AuthService/models"
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

func (s *AuthService) LoginWithEmail(u models.User, password string) (string, error) {
	// 1. Vérifier si l'utilisateur existe dans la base de données
	var existingUser models.User
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


func (s *AuthService) RegisterWithEmail(u models.User) (bool, error) {
	// 1. Vérifier si l'utilisateur existe déjà
	var existingUser models.User
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

func (s *AuthService) ForgotPassword(email string) error {
	// 1. Vérifier si l'adresse email existe dans la base de données
	var existingUser models.User
	err := existingUser.GetUserByEmail(s.DBManager.DB, email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("utilisateur non trouvé")
		}
		return fmt.Errorf("erreur lors de la récupération de l'utilisateur : %v", err)
	}

	// 2. Générer un token de réinitialisation sécurisé
	token := s.Security.GenerateSecureToken()

	// 3. Mettre à jour le token dans la base de données
	err = existingUser.UpdateResetToken(s.DBManager.DB, token)
	if err != nil {
		return fmt.Errorf("erreur lors de la mise à jour du token : %v", err)
	}

	// 4. Construire le lien de réinitialisation
	resetLink := fmt.Sprintf("http://localhost:8080/reset-password?token=%s", token)

	// 5. Envoyer l'email
	err = s.EmailService.SendPasswordResetEmail(email, resetLink)
	if err != nil {
		return fmt.Errorf("erreur lors de l'envoi de l'email de réinitialisation : %v", err)
	}

	return nil
}


func (s *AuthService) ResetPassword(email, token, newPassword string) error {
	// Vérifier si l'utilisateur existe
	var existingUser models.User
	err := existingUser.GetUserByEmail(s.DBManager.DB, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("utilisateur non trouvé")
		}
		return fmt.Errorf("erreur lors de la récupération de l'utilisateur : %v", err)
	}

	// Valider le token de réinitialisation
	if existingUser.ResetToken != token {
		return fmt.Errorf("token invalide ou expiré")
	}

	// 3. Hache le nouveau mot de passe avec le service SecurityService.
	hashedPassword := s.Security.HashPassword(newPassword)
	// 4. Met à jour le mot de passe dans la base de données.
	err = existingUser.UpdatePassword(s.DBManager.DB, hashedPassword)
	if err != nil {
		return fmt.Errorf("erreur lors de la mise à jour du mot de passe : %v", err)
	}
	// 5. Invalide le token après utilisation pour des raisons de sécurité.
	err = existingUser.UpdateResetToken(s.DBManager.DB, "")
	if err != nil {
		return fmt.Errorf("erreur lors de l'invalidation du token : %v", err)
	}

	return nil
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

func (s *AuthService) Logout(token string) error {
	// Vérifie si le token JWT est valide
	claims, err := s.JWTService.VerifyToken(token)
	if err != nil {
		return fmt.Errorf("token invalide ou expiré")
	}

	// Extraire le nom d'utilisateur des claims
	username, ok := claims["username"].(string)
	if !ok {
		return fmt.Errorf("erreur lors de l'extraction du nom d'utilisateur")
	}

	// Invalider le token en l'ajoutant à une liste de révocation
	var revokedToken models.RevokedToken
	err = revokedToken.RevokeToken(s.DBManager.DB, token, username)
	if err != nil {
		log.Printf("Erreur lors de l'invalidation du token : %v", err)
		return fmt.Errorf("erreur lors de l'invalidation du token")
	}

	return nil
}








