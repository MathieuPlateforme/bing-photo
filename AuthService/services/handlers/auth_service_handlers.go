package handlers

import (
	"net/http"
	"AuthService/services/auth" 
	"AuthService/pkg/google"
	"golang.org/x/oauth2"
	"fmt"
)

// Déclaration de AuthHandlers pour encapsuler AuthService
type AuthHandlers struct {
	Service *auth.AuthService
	GoogleAuthService *google.GoogleAuthService
}

func NewAuthHandlers(service *auth.AuthService) (*AuthHandlers,error) {

	// Initialiser le service GoogleAuthService
	googleAuthService, err := google.NewGoogleAuthService()
	if err != nil {
		return nil, err
	}

	return &AuthHandlers{Service: service, GoogleAuthService: googleAuthService}, nil
}

func (h *AuthHandlers) LoginWithEmailHandler(w http.ResponseWriter, r *http.Request) {
	// Logique de connexion avec email et mot de passe
	w.Write([]byte("LoginWithEmail endpoint"))
}

func (h *AuthHandlers) RegisterWithEmailHandler(w http.ResponseWriter, r *http.Request) {
	// Logique d'inscription avec email et mot de passe
	w.Write([]byte("RegisterWithEmail endpoint"))
}

func (h *AuthHandlers) ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	// Logique pour envoyer un email de réinitialisation de mot de passe
	w.Write([]byte("ForgotPassword endpoint"))
}

func (h *AuthHandlers) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	// Logique pour réinitialiser le mot de passe de l'utilisateur
	w.Write([]byte("ResetPassword endpoint"))
}

func (h *AuthHandlers) LoginWithGoogleHandler(w http.ResponseWriter, r *http.Request) {
	// Appeler AuthenticateWithGoogle pour obtenir l'URL d'authentification
	authURL := h.GoogleAuthService.AuthenticateWithGoogle()

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
	// Logique pour valider un token JWT
	w.Write([]byte("ValidateToken endpoint"))
}

func (h *AuthHandlers) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Logique pour déconnecter un utilisateur
	w.Write([]byte("Logout endpoint"))
}
