package main

import (
	"AuthService/middleware"
	"AuthService/pkg/jwt"
	"context"
	"fmt"
	"log"
	"net"

	"AuthService/models"
	proto "AuthService/proto"
	"AuthService/services/auth"
	"AuthService/utils"

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
	// Connect to GalleryService via gRPC
	conn, err := grpc.Dial("gallery-service:50052", grpc.WithInsecure()) // Replace with TLS in production
	if err != nil {
		return fmt.Errorf("failed to connect to GalleryService: %v", err)
	}
	defer conn.Close()

	// Create a gRPC client
	client := proto.NewUserServiceClient(conn)

	// Call the CreateUser gRPC method
	_, err = client.CreateUser(ctx, &proto.CreateUserRequest{
		Email:    email,
		Username: username,
	})
	if err != nil {
		return fmt.Errorf("failed to sync with GalleryService: %v", err)
	}

	log.Println("User successfully synced with GalleryService via gRPC")
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
	// ✅ Récupérer les claims à partir du contexte (via le middleware JWT)
	claims, err := s.JWTService.VerifyTokenFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Token invalide : %v", err)
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "username introuvable dans les claims")
	}

	token, err := s.JWTService.ExtractTokenFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Token introuvable dans le contexte : %v", err)
	}

	err = s.authService.RevokeToken(token, username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Erreur de déconnexion : %v", err)
	}

	return &proto.LogoutResponse{Message: "Déconnexion réussie"}, nil
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
	if err != nil {
		return &proto.GoogleAuthCallbackResponse{Message: "Échec de l'échange du code"}, err
	}

	userInfo, err := s.authService.GoogleAuthService.GetGoogleUserProfile(token)
	if err != nil {
		return &proto.GoogleAuthCallbackResponse{Message: "Échec de la récupération du profil Google"}, err
	}

	// Appelle la méthode que tu viens d’ajouter
	jwtToken, err := s.authService.LoginOrCreateGoogleUser(userInfo)
	if err != nil {
		return &proto.GoogleAuthCallbackResponse{Message: "Login ou création échouée"}, err
	}

	return &proto.GoogleAuthCallbackResponse{
		UserInfo: userInfo.Email,
		Message:  "Login Google réussi",
		Token:    jwtToken,
	}, nil
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
		"/proto.AuthService/Logout":       true,
		"/proto.AuthService/ValidateToken": true,
		"/proto.AuthService/UpdateUser":   true,
		"/proto.AuthService/GetMe":   true,
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

func (s *authServer) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	log.Printf("Mise à jour du profil utilisateur ID: %d", userID)

	// Récupérer l'utilisateur existant
	var user models.User
	if err := s.authService.DBManager.DB.First(&user, userID).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "Utilisateur introuvable : %v", err)
	}

	// Mise à jour conditionnelle des champs
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Picture != "" {
		user.Picture = req.Picture
	}

	// Sauvegarder les modifications
	if err := s.authService.DBManager.DB.Save(&user).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "Échec de la mise à jour de l'utilisateur : %v", err)
	}

	log.Printf("Utilisateur %d mis à jour avec succès", userID)

	return &proto.UpdateUserResponse{
		Message: "Utilisateur mis à jour avec succès",
	}, nil
}

// func (s *authServer) GetMe(ctx context.Context, req *proto.GetMeRequest) (*proto.GetMeResponse, error) {

// 	// Réponse avec succès
// 	return &proto.GetMeResponse{
// 		Email:     "",
// 		Username:  "",
// 		FirstName: "",
// 		LastName:  "",
// 		Picture:   "",
// 	}, nil
// }

func (s *authServer) GetMe(ctx context.Context, req *proto.GetMeRequest) (*proto.GetMeResponse, error) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		log.Printf("GetMe: impossible d'extraire userID du contexte : %v", err)
		return nil, status.Errorf(codes.Unauthenticated, "Token invalide")
	}

	var user models.User
	if err := user.GetUserByID(s.authService.DBManager.DB, userID); err != nil {
		log.Printf("GetMe: utilisateur introuvable : %v", err)
		return nil, status.Errorf(codes.NotFound, "Utilisateur non trouvé")
	}

	return &proto.GetMeResponse{
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Picture:   user.Picture,
	}, nil
}