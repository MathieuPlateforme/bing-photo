package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	proto "ApiGateway/proto"
)

type apiGateway struct {
	authClient proto.AuthServiceClient
}

func (g *apiGateway) loginHandler(w http.ResponseWriter, r *http.Request) {
	req := &proto.LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Failed to parse request: %v\n", err)
		return
	}

	res, err := g.authClient.Login(context.Background(), req)
	if err != nil {
		http.Error(w, "Login failed", http.StatusInternalServerError)
		log.Printf("Login error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Token: " + res.Token))
}

func (g *apiGateway) registerHandler(w http.ResponseWriter, r *http.Request) {

	req := &proto.RegisterRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Failed to parse request: %v\n", err)
		return
	}

	res, err := g.authClient.Register(context.Background(), req)
	if err != nil {
		http.Error(w, "Register failed", http.StatusInternalServerError)
		log.Printf("Register error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message: " + res.Message))
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
	gateway := &apiGateway{authClient: authClient}

	r := mux.NewRouter()

	r.HandleFunc("/login", gateway.loginHandler).Methods("POST")
	r.HandleFunc("/register", gateway.registerHandler).Methods("POST")

	server := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("API Gateway is running on port 8080")
	log.Fatal(server.ListenAndServe())
}
