package main

import (
	"context"
	"log"
	"net"

	"AuthService/models"
	proto "AuthService/proto"
	"AuthService/services/auth"

	"google.golang.org/grpc"
)

type authServer struct {
	proto.UnimplementedAuthServiceServer
	authService *auth.AuthService
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
		Email:     req.Email,
		Password:  req.Password,
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})

	if err != nil || !success {
		return &proto.RegisterResponse{Message: "Registration failed"}, err
	}

	return &proto.RegisterResponse{Message: "Registration successful"}, nil
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

func main() {
	// Initialize AuthService (and other services as needed)
	authService, err := auth.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize AuthService: %v", err)
	}

	// Create gRPC server
	server := grpc.NewServer()
	proto.RegisterAuthServiceServer(server, &authServer{authService: authService})

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
