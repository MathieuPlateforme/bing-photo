package handlers

import (
	"net/http"
	"AuthService/services/auth" 
	"AuthService/pkg/google"
	"golang.org/x/oauth2"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"AuthService/models"
	"AuthService/pkg/email"
	"strings"
	"AuthService/pkg/jwt"
)

// AuthHandlers encapsule les services nécessaires
// @title AuthService API
// @version 1.0
// @description API pour gérer l'authentification
// @contact.name Support API
// @contact.email support@authservice.com
// @host localhost:8080
// @BasePath /
type AuthHandlers struct {
	AuthService       *auth.AuthService
	EmailService      *email.EmailService
	GoogleAuthService *google.GoogleAuthService
	JWTService        *jwt.JWTService
}

func NewAuthHandlers(service *auth.AuthService) (*AuthHandlers, error) {
	googleAuthService, err := google.NewGoogleAuthService()
	if err != nil {
		return nil, err
	}

	emailService, err := email.NewEmailService()
	if err != nil {
		return nil, err
	}

	JWTService, err := jwt.NewJWTService()
	if err != nil {
		return nil, err
	}

	return &AuthHandlers{
		AuthService:       service,
		GoogleAuthService: googleAuthService,
		EmailService:      emailService,
		JWTService:        JWTService,
	}, nil
}

// LoginWithEmailHandler gère la connexion par email et mot de passe
// @Summary Connexion utilisateur
// @Description Permet à un utilisateur de se connecter en utilisant son email et mot de passe
// @Tags Authentification
// @Accept json
// @Produce json
// @Param user body models.LoginRequest true "Informations de connexion"
// @Success 200 {string} string "Jeton JWT généré"
// @Failure 400 {string} string "Paramètres invalides"
// @Failure 401 {string} string "Non autorisé"
// @Router /login [post]
func (h *AuthHandlers) LoginWithEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" || req.Password == "" {
		http.Error(w, "Paramètres manquants ou invalides", http.StatusBadRequest)
		return
	}

	token, err := h.AuthService.LoginWithEmail(models.User{Email: req.Email}, req.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de la connexion : %v", err), http.StatusUnauthorized)
		return
	}

	response := map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// RegisterWithEmailHandler gère l'inscription d'un nouvel utilisateur
// @Summary Inscription utilisateur
// @Description Crée un compte utilisateur avec email et mot de passe
// @Tags Authentification
// @Accept json
// @Produce json
// @Param user body models.RegisterRequest true "Informations de l'utilisateur"
// @Success 201 {string} string "Inscription réussie"
// @Failure 400 {string} string "Requête invalide"
// @Failure 409 {string} string "Conflit - Utilisateur déjà existant"
// @Router /register [post]
func (h *AuthHandlers) RegisterWithEmailHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erreur de lecture de la requête", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var newUser models.User
	err = json.Unmarshal(body, &newUser)
	if err != nil {
		http.Error(w, "Erreur de parsing JSON", http.StatusBadRequest)
		return
	}

	success, err := h.AuthService.RegisterWithEmail(newUser)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de l'inscription : %v", err), http.StatusInternalServerError)
		return
	}

	if success {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Inscription réussie"))
	} else {
		http.Error(w, "Erreur lors de l'inscription", http.StatusConflict)
	}
}

// ForgotPasswordHandler envoie un email de réinitialisation de mot de passe
// @Summary Réinitialisation du mot de passe
// @Description Envoie un email pour réinitialiser le mot de passe d'un utilisateur
// @Tags Authentification
// @Accept json
// @Produce json
// @Param user body models.ForgotPasswordRequest true "Informations de l'utilisateur"
// @Success 200 {string} string "Email envoyé avec succès"
// @Failure 400 {string} string "Email invalide"
// @Failure 500 {string} string "Erreur interne"
// @Router /forgot-password [post]
func (h *AuthHandlers) ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email string `json:"email"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" {
		http.Error(w, "Email manquant ou invalide", http.StatusBadRequest)
		return
	}

	err = h.AuthService.ForgotPassword(req.Email)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de l'envoi de l'email : %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Email de réinitialisation envoyé avec succès"}`))
}


// ResetPasswordHandler gère la réinitialisation de mot de passe
// @Summary Réinitialisation de mot de passe
// @Description Permet de réinitialiser un mot de passe via un jeton
// @Tags Authentification
// @Accept json
// @Produce json
// @Param user body models.ResetPasswordRequest true "Informations de l'utilisateur"
// @Success 200 {string} string "Mot de passe réinitialisé avec succès"
// @Failure 400 {string} string "Paramètres invalides"
// @Failure 500 {string} string "Erreur interne"
// @Router /reset-password [post]
func (h *AuthHandlers) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifier que la méthode est POST
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Décoder le corps de la requête JSON
	var req struct {
		Email       string `json:"email"`
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" || req.Token == "" || req.NewPassword == "" {
		http.Error(w, "Paramètres manquants ou invalides", http.StatusBadRequest)
		return
	}

	// Appeler la fonction de service
	err = h.AuthService.ResetPassword(req.Email, req.Token, req.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Répondre avec un message de succès
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message": "Mot de passe réinitialisé avec succès"}`))
}

// LoginWithGoogleHandler gère la connexion avec Google
// @Summary Connexion via Google
// @Description Authentifie un utilisateur via Google OAuth2
// @Tags Authentification
// @Success 302 {string} string "Redirection vers Google Auth"
// @Router /login-google [get]
func (h *AuthHandlers) LoginWithGoogleHandler(w http.ResponseWriter, r *http.Request) {
	authURL := h.GoogleAuthService.AuthenticateWithGoogle()
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// GoogleAuthCallbackHandler gère le retour de l'authentification Google
// @Summary Callback Google Auth
// @Description Traite le retour de Google après l'authentification OAuth2
// @Tags Authentification
// @Success 200 {string} string "Informations utilisateur"
// @Failure 400 {string} string "Code ou état invalide"
// @Failure 500 {string} string "Erreur interne"
// @Router /oauth2/callback [get]
func (h *AuthHandlers) GoogleAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code d'authentification manquant", http.StatusBadRequest)
		return
	}

	state := r.URL.Query().Get("state")
	if state != "state-token" {
		http.Error(w, "État d'authentification invalide", http.StatusBadRequest)
		return
	}

	token, err := h.GoogleAuthService.Config.Exchange(oauth2.NoContext, code)
	if err != nil {
		http.Error(w, "Échec de l'échange de jeton : "+err.Error(), http.StatusInternalServerError)
		return
	}

	userInfo, err := h.GoogleAuthService.GetGoogleUserProfile(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Informations utilisateur : %+v", userInfo)
}


// ValidateTokenHandler vérifie la validité d'un token JWT
// @Summary Validation du token JWT
// @Description Vérifie si le token JWT est valide
// @Tags Authentification
// @Success 200 {object} map[string]interface{}
// @Failure 401 {string} string "Token invalide"
// @Router /validate-token [get]
func (h *AuthHandlers) ValidateTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifie que la méthode est GET
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Extraire le token JWT de l'en-tête Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Token manquant ou invalide", http.StatusUnauthorized)
		return
	}

	// Supprimer le préfixe "Bearer "
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Valider le token
	claims, err := h.JWTService.VerifyToken(tokenString)
	if err != nil {
		http.Error(w, fmt.Sprintf("Token invalide : %v", err), http.StatusUnauthorized)
		return
	}

	// Répondre avec les informations du token
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message": "Token valide",
		"claims":   claims,
	}
	json.NewEncoder(w).Encode(response)
}

// LogoutHandler gère la déconnexion
// @Summary Déconnexion utilisateur
// @Description Déconnecte l'utilisateur actuel
// @Tags Authentification
// @Success 200 {string} string "Déconnexion réussie"
// @Failure 401 {string} string "Non autorisé"
// @Router /logout [post]
func (h *AuthHandlers) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifie que la méthode est POST
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Extraire le token JWT de l'en-tête Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Token manquant ou invalide", http.StatusUnauthorized)
		return
	}

	// Supprimer le préfixe "Bearer "
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Appeler le service d'authentification pour la déconnexion
	err := h.AuthService.Logout(token)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de la déconnexion : %v", err), http.StatusInternalServerError)
		return
	}

	// Répondre avec un message de succès
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message": "Déconnexion réussie"}`))
}
