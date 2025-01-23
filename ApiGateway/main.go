// package main

// import (
// 	"context"
// 	"log"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/gorilla/mux"
// 	"github.com/rs/cors"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"

// 	"ApiGateway/handlers"
// 	proto "ApiGateway/proto"
// )

// type apiGateway struct {
// 	authClient proto.AuthServiceClient
// }

// func connectToService(address string) (*grpc.ClientConn, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		return nil, err
// 	}

// 	return conn, nil
// }

// func main() {

// 	authServiceAddress := os.Getenv("AUTH_SERVICE")

// 	authConn, err := connectToService(authServiceAddress)
// 	if err != nil {
// 		log.Fatalf("Failed to connect to AuthService: %v", err)
// 	}
// 	defer authConn.Close()

// 	authClient := proto.NewAuthServiceClient(authConn)
// 	authHandler := handlers.NewApiGateway(authClient)

// 	r := mux.NewRouter()

// 	c := cors.New(cors.Options{
// 		AllowedOrigins: []string{"*"},
// 		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
// 	})

// 	r.Use(func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			w.Header().Set("Content-Type", "application/json")
// 			next.ServeHTTP(w, r)
// 		})
// 	})

// 	r.Use(c.Handler)

// 	r.HandleFunc("/login", authHandler.LoginHandler).Methods("POST")
// 	r.HandleFunc("/register", authHandler.RegisterHandler).Methods("POST")
// 	r.HandleFunc("/forgot-password", authHandler.ForgotPasswordHandler).Methods("POST")
// 	r.HandleFunc("/reset-password", authHandler.ResetPasswordHandler).Methods("POST")
// 	r.HandleFunc("/logout", authHandler.LogoutHandler).Methods("POST")

// 	server := &http.Server{
// 		Handler:      r,
// 		Addr:         ":8080",
// 		WriteTimeout: 15 * time.Second,
// 		ReadTimeout:  15 * time.Second,
// 	}

// 	log.Println("API Gateway is running on port 8080")
// 	log.Fatal(server.ListenAndServe())
// }


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

// connectToService établit une connexion gRPC avec un service donné
func connectToService(address string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// jsonMiddleware ajoute l'en-tête Content-Type: application/json à toutes les réponses
func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Lecture de l'adresse du service Auth depuis les variables d'environnement
	authServiceAddress := os.Getenv("AUTH_SERVICE")
	if authServiceAddress == "" {
		log.Fatal("AUTH_SERVICE environment variable is not set")
	}

	// Connexion au service Auth
	authConn, err := connectToService(authServiceAddress)
	if err != nil {
		log.Fatalf("Failed to connect to AuthService: %v", err)
	}
	defer authConn.Close()

	// Initialisation du client gRPC
	authClient := proto.NewAuthServiceClient(authConn)
	authHandler := handlers.NewApiGateway(authClient)

	// Configuration du routeur
	r := mux.NewRouter()

	// Configuration CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Autorise toutes les origines (attention en production)
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true, // Permet les cookies et les identifiants
	})

	// Application des middlewares
	r.Use(jsonMiddleware) // Pour les réponses JSON
	r.Use(c.Handler)      // Middleware CORS

	// Définition des routes
	r.HandleFunc("/login", authHandler.LoginHandler).Methods("POST")
	r.HandleFunc("/register", authHandler.RegisterHandler).Methods("POST")
	r.HandleFunc("/forgot-password", authHandler.ForgotPasswordHandler).Methods("POST")
	r.HandleFunc("/reset-password", authHandler.ResetPasswordHandler).Methods("POST")
	r.HandleFunc("/logout", authHandler.LogoutHandler).Methods("POST")

	// Configuration du serveur
	server := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("API Gateway is running on port 8080")
	log.Fatal(server.ListenAndServe())
}
