package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"GalleryService/internal/api"
	"GalleryService/internal/db"
	"GalleryService/internal/services"

	"github.com/joho/godotenv"
)

func main() {
	// Charger les variables d'environnement
	if err := godotenv.Load(); err != nil {
		log.Println("Avertissement : Impossible de charger le fichier .env, utilisation des variables système.")
	}

	// Initialiser le gestionnaire de base de données
	dbManager, err := db.NewDBManagerService()
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation de la base de données : %v", err)
	}
	defer func() {
		log.Println("Fermeture de la connexion à la base de données...")
		dbManager.CloseConnection()
	}()

	// Effectuer la migration des modèles
	if err := dbManager.AutoMigrate(); err != nil {
		log.Fatalf("Erreur lors de la migration des modèles : %v", err)
	}

	// Initialiser le S3Service
	s3Service := services.NewS3Service("http://my-s3-clone:9090")

	// Initialiser le routeur HTTP
	router := api.NewRouter(dbManager, s3Service)

	// Configurer le serveur HTTP
	server := &http.Server{
		Addr:    ":50052",
		Handler: router,
	}

	// Canal pour gérer les signaux système (interruption ou arrêt)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Démarrer le serveur HTTP dans une goroutine
	go func() {
		log.Println("GalleryService démarré sur le port 50052...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erreur lors de l'exécution du serveur HTTP : %v", err)
		}
	}()

	// Attendre un signal d'arrêt
	<-stop
	log.Println("Signal reçu, arrêt du service GalleryService...")

	// Arrêter gracieusement le serveur HTTP
	if err := server.Close(); err != nil {
		log.Fatalf("Erreur lors de l'arrêt du serveur HTTP : %v", err)
	}

	log.Println("Service GalleryService arrêté avec succès.")
}
