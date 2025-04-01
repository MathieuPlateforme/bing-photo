package main

import (
	"AuthService/middleware"
	"AuthService/pkg/jwt"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"AuthService/models"
	proto "AuthService/proto"
	"AuthService/services/auth"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authServer struct {
	proto.UnimplementedAuthServiceServer
	authService *auth.AuthService
	JWTService  *jwt.JWTService
}

func (s *authServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	// Delegate to the LoginWithEmail function in AuthService
	token, err := s.authService.LoginWithEmail(models.User{
		Email: req.Email,
	}, req.Password)

	if err != nil {
		return &proto.LoginResponse{Message: "Login failed"}, err
	}

	return &proto.LoginResponse{Token: token, Message: "Login successful"}, nil
}

// RegisterWithEmail handles user registration
func (s *authServer) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	success, err := s.authService.RegisterWithEmail(models.User{
		Email:    req.Email,
		Password: req.Password,
		Username: req.Username,
	})

	if err != nil || !success {
		return &proto.RegisterResponse{Message: "Registration failed"}, err
	}

	// Appeler GalleryService pour synchroniser l'utilisateur
	err = s.syncWithGalleryService(ctx, req.Email, req.Username)
	if err != nil {
		return &proto.RegisterResponse{Message: "Failed to sync with GalleryService"}, err
	}

	return &proto.RegisterResponse{Message: "Registration successful"}, nil
}

func (s *authServer) syncWithGalleryService(ctx context.Context, email string, username string) error {
	// Construire la requête pour GalleryService
	url := "http://gallery-service:50052/users"
	payload := map[string]string{
		"email":    email,
		"username": username,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("erreur lors de la construction du payload : %v", err)
	}

	// Envoyer la requête HTTP POST
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("erreur lors de la création de la requête : %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erreur lors de l'envoi de la requête à GalleryService : %v", err)
	}
	defer resp.Body.Close()

	// Vérifier la réponse
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("échec de la synchronisation avec GalleryService, statut : %d", resp.StatusCode)
	}

	log.Println("Utilisateur synchronisé avec succès dans GalleryService")
	return nil
}

func (s *authServer) ForgotPassword(ctx context.Context, req *proto.ForgotPasswordRequest) (*proto.ForgotPasswordResponse, error) {
	err := s.authService.ForgotPassword(req.Email)

	if err != nil {
		return &proto.ForgotPasswordResponse{Message: "Email not found"}, err
	}

	return &proto.ForgotPasswordResponse{Message: "Email succesfully sent"}, nil
}

func (s *authServer) ResetPassword(ctx context.Context, req *proto.ResetPasswordRequest) (*proto.ResetPasswordResponse, error) {
	err := s.authService.ResetPassword(req.Email, req.Token, req.NewPassword)

	if err != nil {
		return &proto.ResetPasswordResponse{}, err
	}

	return &proto.ResetPasswordResponse{}, nil
}

func (s *authServer) Logout(ctx context.Context, req *proto.LogoutRequest) (*proto.LogoutResponse, error) {
	err := s.authService.Logout(req.Token)

	if err != nil {
		return &proto.LogoutResponse{}, err
	}

	return &proto.LogoutResponse{}, nil
}
func (s *authServer) LoginWithGoogle(ctx context.Context, req *proto.GoogleAuthRequest) (*proto.GoogleAuthResponse, error) {
	authUrl, err := s.authService.GoogleAuthService.AuthenticateWithGoogle()

	if err != nil {
		return &proto.GoogleAuthResponse{AuthUrl: "Login failed"}, err
	}

	return &proto.GoogleAuthResponse{AuthUrl: authUrl}, nil
}
func (s *authServer) GoogleAuthCallback(ctx context.Context, req *proto.GoogleAuthCallbackRequest) (*proto.GoogleAuthCallbackResponse, error) {
	token, err := s.authService.GoogleAuthService.Config.Exchange(oauth2.NoContext, req.Code)
	userInfo, err := s.authService.GoogleAuthService.GetGoogleUserProfile(token)

	if err != nil {
		return &proto.GoogleAuthCallbackResponse{Message: "Login failed"}, err
	}

	userInfoJson, err := json.Marshal(userInfo)
	if err != nil {
		return &proto.GoogleAuthCallbackResponse{Message: "Failed to parse user info"}, err
	}
	return &proto.GoogleAuthCallbackResponse{UserInfo: string(userInfoJson), Message: "Login successful"}, nil
}
func main() {

	JWTService, err := jwt.NewJWTService()
	if err != nil {
		log.Fatalf("Failed to initialize JWTService: %v", err)
	}
	// Initialize AuthService (and other services as needed)
	authService, err := auth.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize AuthService: %v", err)
	}

	// Définir les méthodes nécessitant une vérification d'authentification
	methodsToIntercept := map[string]bool{
	}

	// Create gRPC server
	server := grpc.NewServer(grpc.UnaryInterceptor(middleware.AuthInterceptor(JWTService, methodsToIntercept)))
	proto.RegisterAuthServiceServer(server, &authServer{
		authService: authService,
		JWTService:  JWTService,
	})
	if err := authService.DBManager.AutoMigrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Listen on a specific port
	listener, err := net.Listen("tcp", ":50051") // gRPC port for AuthService
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("AuthService gRPC server is running on port 50051")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *authServer) ValidateToken(ctx context.Context, req *proto.ValidateTokenRequest) (*proto.ValidateTokenResponse, error) {
	log.Printf("Token reçu pour validation : %s", req.Token)

	// Vérifier le token avec le service JWT
	claims, err := s.JWTService.VerifyToken(req.Token)
	if err != nil {
		log.Printf("Erreur lors de la validation du token : %v", err)
		return nil, status.Errorf(codes.Unauthenticated, "Token invalide : %v", err)
	}

	log.Printf("Token valide pour l'utilisateur : %v", claims["username"])

	// Réponse avec succès
	return &proto.ValidateTokenResponse{
		Message: "Token valide",
	}, nil
}
