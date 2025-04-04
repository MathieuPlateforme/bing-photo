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

	_ "ApiGateway/docs"
	"ApiGateway/handlers"
	proto "ApiGateway/proto"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Bing Photo API Gateway
// @version         1.0
// @description     This is the API Gateway for Bing Photo project.
// @contact.name    Your Name
// @contact.email   your@email.com
// @host            localhost:8081
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Format : Bearer <votre_token>

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
	galleryServiceAddress := os.Getenv("GALLERY_SERVICE")

	// Connect to Auth Service
	authConn, err := connectToService(authServiceAddress)
	if err != nil {
		log.Fatalf("Failed to connect to AuthService: %v", err)
	}
	defer authConn.Close()

	// Connect to Gallery Service
	galleryConn, err := connectToService(galleryServiceAddress)
	if err != nil {
		log.Fatalf("Failed to connect to GalleryService: %v", err)
	}
	defer galleryConn.Close()

	// Initialize handlers
	authClient := proto.NewAuthServiceClient(authConn)
	authHandler := handlers.NewApiGateway(authClient)

	// Initialize gallery service clients
	albumClient := proto.NewAlbumServiceClient(galleryConn)
	mediaClient := proto.NewMediaServiceClient(galleryConn)
	userClient := proto.NewUserServiceClient(galleryConn)
	galleryHandler := handlers.NewGalleryGateway(albumClient, mediaClient, userClient)

	r := mux.NewRouter()

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	// Swagger documentation route
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Auth routes
	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", authHandler.LoginHandler).Methods("POST", "OPTIONS")
	auth.HandleFunc("/register", authHandler.RegisterHandler).Methods("POST", "OPTIONS")
	auth.HandleFunc("/google", authHandler.GoogleHandler).Methods("GET", "OPTIONS")
	auth.HandleFunc("/oauth2/callback", authHandler.GoogleCallbackHandler).Methods("GET", "POST", "OPTIONS")
	auth.HandleFunc("/forgot-password", authHandler.ForgotPasswordHandler).Methods("POST", "OPTIONS")
	auth.HandleFunc("/reset-password", authHandler.ResetPasswordHandler).Methods("POST", "OPTIONS")
	auth.HandleFunc("/logout", authHandler.LogoutHandler).Methods("POST", "OPTIONS")
	auth.HandleFunc("/validateToken", authHandler.ValidateTokenHandler).Methods("POST")
	auth.HandleFunc("/update-user", authHandler.UpdateUserHandler).Methods("PUT", "OPTIONS")
	auth.HandleFunc("/get-me", authHandler.GetMeHandler).Methods("GET", "OPTIONS")

	// Album routes
	r.HandleFunc("/albums", galleryHandler.CreateAlbumHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/albums/user", galleryHandler.GetAlbumsByUserHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/albums/{id}", galleryHandler.UpdateAlbumHandler).Methods("PUT", "OPTIONS")
	r.HandleFunc("/albums/{id}", galleryHandler.DeleteAlbumHandler).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/albums/type", galleryHandler.GetPrivateAlbumHandler).Methods("GET", "OPTIONS")

	// Media routes
	r.HandleFunc("/media", galleryHandler.AddMediaHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/media/user", galleryHandler.GetMediaByUserHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/media/{id}/private", galleryHandler.MarkAsPrivateHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/media/private", galleryHandler.GetPrivateMediaHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/media/{id}/download", galleryHandler.DownloadMediaHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/media/{id}", galleryHandler.DeleteMediaHandler).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/media/similar", galleryHandler.DetectSimilarMediaHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/media/album/{id}", galleryHandler.GetMediaByAlbumHandler).Methods("GET", "OPTIONS")


	// User routes
	r.HandleFunc("/users", galleryHandler.CreateUserHandler).Methods("POST", "OPTIONS")

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			enableCors(&w)
			next.ServeHTTP(w, r)
		})
	})

	r.Use(c.Handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API Gateway starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
