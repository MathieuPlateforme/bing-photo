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

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Token: " + res.Token))
}

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

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message: " + res.Message))
}

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

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message: " + res.Message))
}

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

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message: " + res.Message))
}

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
	w.Write([]byte("Message: " + res.Message))
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

