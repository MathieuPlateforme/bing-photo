package main

import (
	"log"
	"fmt"
	"net/http"
	"AuthService/services/auth"
	"AuthService/pkg/db"
	"AuthService/services/handlers"
	"AuthService/pkg/google"
	"time"
	_ "AuthService/docs"   
	"github.com/swaggo/http-swagger"         
)

// @title AuthService API
// @version 1.0
// @description API de gestion de l'authentification pour Bing Photo
// @contact.name Support API
// @contact.email support@authservice.com
// @host localhost:8080
// @BasePath /

func waitForDatabaseConnection(dbManager *db.DBManagerService) {
	maxRetries := 20
	for i := 0; i < maxRetries; i++ {
		err := dbManager.Ping(dbManager.DB)
		if err == nil {
			return
		}
		log.Println("Attente de la disponibilité de la base de données...")
		time.Sleep(5 * time.Second)
	}
	log.Fatal("La base de données n'est pas disponible après plusieurs tentatives")
}

func main() {
	
	// Initialiser le service DBManager
	dbManager, err := db.NewDBManagerService()
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation du service DBManager : %v", err)
	}
	waitForDatabaseConnection(dbManager)
	
	// Initialiser le service d'authentification
	authService, err := auth.Initialize()
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation d'AuthService : %v", err)
	}
	fmt.Println("AuthService initialisé avec succès :", authService)


	// Initialiser le service GoogleAuthService
	googleAuthService, err := google.NewGoogleAuthService()
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation de GoogleAuthService : %v", err)
	}
	fmt.Println("GoogleAuthService initialisé avec succès :", googleAuthService)

	// Initialiser AuthHandlers avec AuthService 
	authHandlers, err := handlers.NewAuthHandlers(authService)
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation des gestionnaires : %v", err)
	}
	fmt.Println("AuthHandlers initialisé avec succès :", authHandlers)

	
	// Exécuter la migration de la base de données
	if err := dbManager.AutoMigrate(); err != nil {
		log.Fatalf("Erreur lors de la migration de la base de données : %v", err)
	}

	// Configurer les routes avec les handlers
	http.HandleFunc("/login", authHandlers.LoginWithEmailHandler)
	http.HandleFunc("/register", authHandlers.RegisterWithEmailHandler)
	http.HandleFunc("/forgot-password", authHandlers.ForgotPasswordHandler)
	http.HandleFunc("/reset-password", authHandlers.ResetPasswordHandler)
	http.HandleFunc("/login-google", authHandlers.LoginWithGoogleHandler)
	http.HandleFunc("/oauth2/callback",authHandlers.GoogleAuthCallbackHandler)
	http.HandleFunc("/validate-token", authHandlers.ValidateTokenHandler)
	http.HandleFunc("/logout", authHandlers.LogoutHandler)
	http.Handle("/swagger/", httpSwagger.WrapHandler)


	// Démarrer le serveur HTTP
	log.Fatal(http.ListenAndServe(":8080", nil))
}
