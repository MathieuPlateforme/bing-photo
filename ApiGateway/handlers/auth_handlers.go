package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

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

	response := map[string]string{"Token": res.Token}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
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

	response := map[string]string{"Message": res.Message}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
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
	response := map[string]string{"Message": res.Message}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
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

	response := map[string]string{"Message": res.Message}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
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
	response := map[string]string{"Message": res.Message}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
