package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"ApiGateway/handlers"
	proto "ApiGateway/proto"
)

type apiGateway struct {
	authClient proto.AuthServiceClient
}

func connectToService(address string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func main() {

	authServiceAddress := os.Getenv("AUTH_SERVICE")

	authConn, err := connectToService(authServiceAddress)
	if err != nil {
		log.Fatalf("Failed to connect to AuthService: %v", err)
	}
	defer authConn.Close()

	authClient := proto.NewAuthServiceClient(authConn)
	authHandler := handlers.NewApiGateway(authClient)

	r := mux.NewRouter()

	r.HandleFunc("/login", authHandler.LoginHandler).Methods("POST")
	r.HandleFunc("/register", authHandler.RegisterHandler).Methods("POST")
	r.HandleFunc("/google", authHandler.GoogleHandler).Methods("POST")
	r.HandleFunc("/oauth2/callback", authHandler.GoogleHandler).Methods("POST")
	r.HandleFunc("/forgot-password", authHandler.ForgotPasswordHandler).Methods("POST")
	r.HandleFunc("/reset-password", authHandler.ResetPasswordHandler).Methods("POST")
	r.HandleFunc("/logout", authHandler.LogoutHandler).Methods("POST")

	server := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("API Gateway is running on port 8080")
	log.Fatal(server.ListenAndServe())
}
