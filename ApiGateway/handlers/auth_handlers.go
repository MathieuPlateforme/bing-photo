package handlers

import (
	"ApiGateway/utils"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	proto "ApiGateway/proto"
)

type ApiGateway struct {
	AuthClient proto.AuthServiceClient
}

func NewApiGateway(authClient proto.AuthServiceClient) *ApiGateway {
	return &ApiGateway{AuthClient: authClient}
}

// LoginHandler godoc
// @Summary Login
// @Description Authenticates a user and returns a JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body proto.LoginRequest true "User credentials"
// @Success 200 {object} map[string]string "Token returned"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Login failed"
// @Router /auth/login [post]
func (g *ApiGateway) LoginHandler(w http.ResponseWriter, r *http.Request) {
	req := &proto.LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Failed to parse request: %v\n", err)
		return
	}

	res, err := g.AuthClient.Login(context.Background(), req)
	if err != nil {
		http.Error(w, "Login failed"+err.Error(), http.StatusInternalServerError)
		log.Printf("Login error: %v\n", err)
		return
	}

	response := map[string]string{"Token": res.Token}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// RegisterHandler godoc
// @Summary Register
// @Description Registers a new user and syncs with the gallery service
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body proto.RegisterRequest true "User registration data"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Registration failed"
// @Router /auth/register [post]
func (g *ApiGateway) RegisterHandler(w http.ResponseWriter, r *http.Request) {

	req := &proto.RegisterRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Failed to parse request: %v\n", err)
		return
	}

	res, err := g.AuthClient.Register(context.Background(), req)
	if err != nil {
		http.Error(w, "Register failed"+err.Error(), http.StatusInternalServerError)
		log.Printf("Register error: %v\n", err)
		return
	}

	response := map[string]string{"Message": res.Message}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// ForgotPasswordHandler godoc
// @Summary Forgot Password
// @Description Sends a reset password email
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body proto.ForgotPasswordRequest true "Email for password reset"
// @Success 200 {object} map[string]string "Email successfully sent"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Forgot password failed"
// @Router /auth/forgot-password [post]
func (g *ApiGateway) ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {

	req := &proto.ForgotPasswordRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Failed to parse request: %v\n", err)
		return
	}

	res, err := g.AuthClient.ForgotPassword(context.Background(), req)
	if err != nil {
		http.Error(w, "Forgot password failed"+err.Error(), http.StatusInternalServerError)
		log.Printf("Forgot password error: %v\n", err)
		return
	}
	response := map[string]string{"Message": res.Message}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// ResetPasswordHandler godoc
// @Summary Reset Password
// @Description Resets the user's password using a token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body proto.ResetPasswordRequest true "Reset token and new password"
// @Success 200 {object} map[string]string "Password reset success"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Reset password failed"
// @Router /auth/reset-password [post]
func (g *ApiGateway) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {

	req := &proto.ResetPasswordRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Failed to parse request: %v\n", err)
		return
	}

	res, err := g.AuthClient.ResetPassword(context.Background(), req)
	if err != nil {
		http.Error(w, "Reset password failed: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Reset password error: %v\n", err)
		return
	}

	response := map[string]string{"Message": res.Message}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// LogoutHandler godoc
// @Summary Logout
// @Description Logs the user out by invalidating the token
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body proto.LogoutRequest true "Token to invalidate"
// @Success 200 {object} map[string]string "Logout successful"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Logout failed"
// @Router /auth/logout [post]
func (g *ApiGateway) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx, err := utils.AttachTokenToContext(r)
	if err != nil {
		http.Error(w, "Token manquant ou invalide", http.StatusUnauthorized)
		return
	}

	_, err = g.AuthClient.Logout(ctx, &proto.LogoutRequest{}) // pas besoin de Token dans le body
	if err != nil {
		http.Error(w, "Erreur lors de la déconnexion : "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Déconnexion réussie"})
}

// GoogleHandler godoc
// @Summary Google OAuth
// @Description Generates a Google login URL
// @Tags Auth
// @Produce plain
// @Success 200 {string} string "Google login URL"
// @Failure 500 {string} string "Failed to generate URL"
// @Router /auth/google [get]
func (g *ApiGateway) GoogleHandler(w http.ResponseWriter, r *http.Request) {
	res, err := g.AuthClient.LoginWithGoogle(context.Background(), &proto.GoogleAuthRequest{})
	if err != nil {
		http.Error(w, "Failed to generate URL", http.StatusInternalServerError)
		log.Printf("Failed to generate URL: %v\n", err)
		return
	}

	// ✅ Redirection réelle vers Google
	http.Redirect(w, r, res.AuthUrl, http.StatusSeeOther)
}

// GoogleCallbackHandler godoc
// @Summary Google OAuth Callback
// @Description Handles the OAuth callback after Google login
// @Tags Auth
// @Accept json
// @Produce plain
// @Param request body proto.GoogleAuthCallbackRequest true "Authorization code"
// @Success 200 {string} string "Login success and user info"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Google callback failed"
// @Router /auth/google/callback [get]
func (g *ApiGateway) GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// ✅ Récupérer les valeurs depuis les query params (GET)
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" || state == "" {
		http.Error(w, "Code ou state manquant dans la requête", http.StatusBadRequest)
		log.Printf("Code/state manquant. code: %s, state: %s\n", code, state)
		return
	}

	// ✅ Préparer la requête gRPC avec ces données
	req := &proto.GoogleAuthCallbackRequest{
		Code:  code,
		State: state,
	}

	// ✅ Appeler le service Auth
	res, err := g.AuthClient.GoogleAuthCallback(context.Background(), req)
	if err != nil {
		http.Error(w, "Google callback failed: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Google callback error: %v\n", err)
		return
	}

	// ✅ Rediriger vers le front avec le token JWT
	redirectURL := "http://localhost:3000/overview?token=" + res.Message
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (g *ApiGateway) ValidateTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Extraire le token de l'en-tête Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Token manquant ou invalide", http.StatusUnauthorized)
		log.Println("Authorization header missing or invalid")
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Créer une requête gRPC pour valider le token
	req := &proto.ValidateTokenRequest{
		Token: token,
	}

	// Appeler le client AuthService pour valider le token
	res, err := g.AuthClient.ValidateToken(context.Background(), req)
	if err != nil {
		http.Error(w, "Échec de validation du token", http.StatusInternalServerError)
		log.Printf("Token validation error: %v\n", err)
		return
	}

	// Répondre avec le message du service d'authentification
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message: " + res.Message))
}

// UpdateUserHandler godoc
// @Summary Mettre à jour un utilisateur
// @Description Met à jour les informations d'un utilisateur (nom, prénom, email, photo)
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body proto.UpdateUserRequest true "Champs à mettre à jour"
// @Success 200 {object} proto.UpdateUserResponse
// @Failure 400 {string} string "Requête invalide"
// @Failure 500 {string} string "Erreur interne du serveur"
// @Router /auth/update-user [put]
func (g *ApiGateway) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req proto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Failed to parse request: %v\n", err)
		return
	}

	ctx, err := utils.AttachTokenToContext(r)
	if err != nil {
		http.Error(w, "Token manquant ou invalide", http.StatusUnauthorized)
		log.Printf("Token manquant ou invalide : %v\n", err)
		return
	}

	res, err := g.AuthClient.UpdateUser(ctx, &req)
	if err != nil {
		http.Error(w, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Update user error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// func (g *ApiGateway) GetMeHandler(w http.ResponseWriter, r *http.Request) {
// 	res, err := g.AuthClient.GetMe(context.Background(), &proto.GetMeRequest{})
// 	if err != nil {
// 		http.Error(w, "Failed to get user: "+err.Error(), http.StatusInternalServerError)
// 		log.Printf("Get user error: %v\n", err)
// 		return
// 	}

//		w.WriteHeader(http.StatusOK)
//		json.NewEncoder(w).Encode(res)
//	}
//
// GetMeHandler godoc
// @Summary Obtenir les informations de l'utilisateur connecté
// @Description Retourne les informations du profil de l'utilisateur authentifié
// @Tags Auth
// @Produce json
// @Success 200 {object} proto.GetMeResponse
// @Failure 401 {string} string "Token manquant ou invalide"
// @Failure 500 {string} string "Erreur interne"
// @Security BearerAuth
// @Router /auth/get-me [get]
func (g *ApiGateway) GetMeHandler(w http.ResponseWriter, r *http.Request) {
	ctx, err := utils.AttachTokenToContext(r)
	if err != nil {
		http.Error(w, "Token manquant ou invalide", http.StatusUnauthorized)
		log.Printf("GetMeHandler: token invalide : %v", err)
		return
	}

	res, err := g.AuthClient.GetMe(ctx, &proto.GetMeRequest{})
	if err != nil {
		http.Error(w, "Erreur lors de la récupération de l'utilisateur : "+err.Error(), http.StatusInternalServerError)
		log.Printf("GetMeHandler: erreur gRPC : %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
