package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	proto "ApiGateway/proto"
)

type apiGateway struct {
	authClient proto.AuthServiceClient
}

func (g *apiGateway) loginHandler(w http.ResponseWriter, r *http.Request) {
	req := &proto.LoginRequest{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	res, err := g.authClient.Login(context.Background(), req)
	if err != nil {
		http.Error(w, "Login failed", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Token: " + res.Token))
}

func (g *apiGateway) registerHandler(w http.ResponseWriter, r *http.Request) {
	req := &proto.RegisterRequest{
		Email:     r.FormValue("email"),
		Password:  r.FormValue("password"),
		Username:  r.FormValue("username"),
		FirstName: r.FormValue("first_name"),
		LastName:  r.FormValue("last_name"),
	}
	res, err := g.authClient.Register(context.Background(), req)
	if err != nil {
		http.Error(w, "Register failed", http.StatusInternalServerError)
		return
	}
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

	http.HandleFunc("/login", gateway.loginHandler)
	http.HandleFunc("/register", gateway.registerHandler)

	log.Println("API Gateway is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
