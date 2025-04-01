package handlers

import (
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
// @Param request body proto.LogoutRequest true "Token to invalidate"
// @Success 200 {object} map[string]string "Logout successful"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Logout failed"
// @Router /auth/logout [post]
func (g *ApiGateway) LogoutHandler(w http.ResponseWriter, r *http.Request) {

	req := &proto.LogoutRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Failed to parse request: %v\n", err)
		return
	}

	res, err := g.AuthClient.Logout(context.Background(), req)
	if err != nil {
		http.Error(w, "Logout failed", http.StatusInternalServerError)
		log.Printf("Logout error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]string{"Message": res.Message}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
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

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message: " + res.AuthUrl))
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
// @Router /auth/google/callback [post]
func (g *ApiGateway) GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	req := &proto.GoogleAuthCallbackRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Failed to parse request: %v\n", err)
		return
	}

	res, err := g.AuthClient.GoogleAuthCallback(context.Background(), req)
	if err != nil {
		http.Error(w, "Google callback failed"+err.Error(), http.StatusInternalServerError)
		log.Printf("Google callback error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Token: " + res.Message + "user: " + res.UserInfo))
}

// ValidateTokenHandler godoc
// @Summary Validate Token
// @Description Validates a JWT token
// @Tags Auth
// @Produce plain
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {string} string "Token valid"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Token validation failed"
// @Router /auth/validateToken [post]
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
