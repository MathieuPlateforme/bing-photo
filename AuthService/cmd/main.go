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
	// Delegate to the RegisterWithEmail function in AuthService
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
