package main

import (
	"log"
	"net/http"
	"authservice/pkg/auth"
	"authservice/pkg/db"
)

func main() {
	// Initialiser le service d'authentification
	if err := auth.Initialize(); err != nil {
		log.Fatalf("Failed to initialize AuthService: %v", err)
	}
	log.Println("AuthService is running...")

	// Se connecter à la base de données
	db.Connect()

	// Exécuter la migration de la base de données
	db.AutoMigrate()
	
	// Démarrer un serveur HTTP pour maintenir le service actif
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("AuthService is running..."))
	})

	// Écouter sur le port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}
