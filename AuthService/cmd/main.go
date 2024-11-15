package main

import (
	"log"
	"fmt"
	"net/http"
	"AuthService/services/auth"
	"AuthService/pkg/db"
	"AuthService/services/handlers" 
)

func main() {
	// Initialiser le service d'authentification
	authService, err := auth.Initialize()
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation d'AuthService : %v", err)
	}
	fmt.Println("AuthService initialisé avec succès :", authService)

	// Initialiser le service DBManager
	dbManager, err := db.NewDBManagerService()
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation du service DBManager : %v", err)
	}

	// Exécuter la migration de la base de données
	if err := dbManager.AutoMigrate(); err != nil {
		log.Fatalf("Erreur lors de la migration de la base de données : %v", err)
	}

	// Initialiser les handlers avec AuthService
	authHandlers := handlers.NewAuthHandlers(authService)

	// Configurer les routes avec les handlers
	http.HandleFunc("/login", authHandlers.LoginWithEmailHandler)
	http.HandleFunc("/register", authHandlers.RegisterWithEmailHandler)
	http.HandleFunc("/forgot-password", authHandlers.ForgotPasswordHandler)
	http.HandleFunc("/reset-password", authHandlers.ResetPasswordHandler)
	http.HandleFunc("/login-google", authHandlers.LoginWithGoogleHandler)
	http.HandleFunc("/validate-token", authHandlers.ValidateTokenHandler)
	http.HandleFunc("/logout", authHandlers.LogoutHandler)

	// Démarrer le serveur HTTP
	log.Fatal(http.ListenAndServe(":8080", nil))
}
