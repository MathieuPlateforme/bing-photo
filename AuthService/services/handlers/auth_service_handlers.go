package handlers

import (
	"AuthService/models"
	"AuthService/pkg/email"
	"AuthService/pkg/google"
	"AuthService/pkg/jwt"
	"AuthService/services/auth"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

// Déclaration de AuthHandlers pour encapsuler AuthService
type AuthHandlers struct {
	AuthService       *auth.AuthService
	EmailService      *email.EmailService
	GoogleAuthService *google.GoogleAuthService
	JWTService        *jwt.JWTService
}

func NewAuthHandlers(service *auth.AuthService) (*AuthHandlers, error) {

	// Initialiser le service GoogleAuthService
	googleAuthService, err := google.NewGoogleAuthService()
	if err != nil {
		return nil, err
	}

	// Initialiser le service EmailService
	emailService, err := email.NewEmailService()
	if err != nil {
		return nil, err
	}

	// Initialiser le service JWTService

	JWTService, err := jwt.NewJWTService()
	if err != nil {
		return nil, err
	}

	return &AuthHandlers{AuthService: service, GoogleAuthService: googleAuthService, EmailService: emailService, JWTService: JWTService}, nil
}

func (h *AuthHandlers) LoginWithEmailHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifie que la méthode est POST
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Décode le corps de la requête
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Lire et décoder la requête
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" || req.Password == "" {
		http.Error(w, "Paramètres manquants ou invalides", http.StatusBadRequest)
		return
	}

	// Appeler le service d'authentification
	token, err := h.AuthService.LoginWithEmail(models.User{Email: req.Email}, req.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de la connexion : %v", err), http.StatusUnauthorized)
		return
	}

	// Répondre avec le token JWT
	response := map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandlers) RegisterWithEmailHandler(w http.ResponseWriter, r *http.Request) {
	// Lire le corps de la requête
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erreur lors de la lecture du corps de la requête", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Décode le JSON reçu en un objet User
	var newUser models.User
	err = json.Unmarshal(body, &newUser)
	if err != nil {
		http.Error(w, "Erreur de parsing JSON", http.StatusBadRequest)
		return
	}

	// Appeler le service d'inscription
	success, err := h.AuthService.RegisterWithEmail(newUser)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de l'inscription : %v", err), http.StatusInternalServerError)
		return
	}

	// Retourner une réponse de succès
	if success {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Inscription réussie"))
	} else {
		http.Error(w, "Erreur lors de l'inscription", http.StatusConflict)
	}
}

func (h *AuthHandlers) ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifier que la méthode est POST
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Décoder le corps de la requête JSON
	var req struct {
		Email string `json:"email"`
	}

	// Décoder la requête entrante
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" {
		http.Error(w, "Email manquant ou invalide", http.StatusBadRequest)
		return
	}

	// Appeler le service ForgotPassword
	err = h.AuthService.ForgotPassword(req.Email)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de l'envoi de l'email de réinitialisation : %v", err), http.StatusInternalServerError)
		return
	}

	// Répondre avec un message de succès
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message": "Email de réinitialisation envoyé avec succès"}`))
}

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

func (h *AuthHandlers) LoginWithGoogleHandler(w http.ResponseWriter, r *http.Request) {
	// Appeler AuthenticateWithGoogle pour obtenir l'URL d'authentification
	authURL, err := h.GoogleAuthService.AuthenticateWithGoogle()
	if err != nil {
		http.Error(w, "Échec de l'échange de jeton : "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Rediriger l'utilisateur vers l'URL d'authentification Google
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (h *AuthHandlers) GoogleAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer le code d'authentification de Google
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code d'authentification manquant", http.StatusBadRequest)
		return
	}

	// Récupérer l'état de l'authentification
	state := r.URL.Query().Get("state")
	if state != "state-token" {
		http.Error(w, "État d'authentification invalide", http.StatusBadRequest)
		return
	}

	// Échanger le code d'authentification contre un jeton d'accès
	token, err := h.GoogleAuthService.Config.Exchange(oauth2.NoContext, code)
	// log.Println("Token FOR TEST : ", token)
	if err != nil {
		http.Error(w, "Échec de l'échange de jeton : "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Obtenir le profil utilisateur en utilisant GoogleAuthService
	userInfo, err := h.GoogleAuthService.GetGoogleUserProfile(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Renvoyer ou traiter les informations utilisateur
	fmt.Fprintf(w, "Informations utilisateur : %+v", userInfo)
}

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
		"claims":  claims,
	}
	json.NewEncoder(w).Encode(response)
}

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
