package main

import (
	"context"
	"log"
	"net/http"
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

	authConn, err := connectToService("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to connect to AuthService: %v", err)
	}
	defer authConn.Close()

	authClient := proto.NewAuthServiceClient(authConn)
	authHandler := handlers.NewApiGateway(authClient)

	r := mux.NewRouter()

	r.HandleFunc("/login", authHandler.LoginHandler).Methods("POST")
	r.HandleFunc("/register", authHandler.RegisterHandler).Methods("POST")

	server := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("API Gateway is running on port 8080")
	log.Fatal(server.ListenAndServe())
}
