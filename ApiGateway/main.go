package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
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

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
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

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		Debug:          true,
	})

	r.Use(c.Handler)

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			enableCors(&w)
			next.ServeHTTP(w, r)
		})
	})

	r.Use(c.Handler)

	r.HandleFunc("/login", authHandler.LoginHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/register", authHandler.RegisterHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/google", authHandler.GoogleHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/oauth2/callback", authHandler.GoogleCallbackHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/forgot-password", authHandler.ForgotPasswordHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/reset-password", authHandler.ResetPasswordHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/logout", authHandler.LogoutHandler).Methods("POST", "OPTIONS")

	server := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("API Gateway is running on port 8080")
	log.Fatal(server.ListenAndServe())
}
