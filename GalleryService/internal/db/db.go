package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"GalleryService/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBManagerService struct {
	DB *gorm.DB
}

// NewDBManagerService initialise la connexion à la base de données
func NewDBManagerService() (*DBManagerService, error) {
	log.Println("Initializing DBManagerService...")

	// Charger les variables d'environnement
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")

	if user == "" || password == "" || dbname == "" || host == "" || port == "" {
		return nil, fmt.Errorf("les variables d'environnement pour la base de données ne sont pas correctement définies")
	}

	// Construire la chaîne de connexion
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Ouvrir la connexion avec GORM
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Activer les logs de niveau Info
	})

	if err != nil {
		return nil, fmt.Errorf("erreur lors de la connexion à la base de données : %v", err)
	}

	log.Println("Connexion à la base de données réussie")

	return &DBManagerService{DB: db}, nil
}

// AutoMigrate effectue la migration des modèles
func (manager *DBManagerService) AutoMigrate() error {
	log.Println("Démarrage de la migration des modèles...")
	err := manager.DB.AutoMigrate(
		&models.User{},
		&models.Album{},
		&models.Media{},
		&models.Access{},
		&models.UserAccess{},
		&models.SimilarGroup{},
		&models.SimilarMedia{},
	)
	if err != nil {
		return fmt.Errorf("erreur lors de la migration de la base de données : %v", err)
	}
	log.Println("Migration de la base de données réussie")
	return nil
}

// CloseConnection ferme la connexion à la base de données
func (manager *DBManagerService) CloseConnection() {
	db, err := manager.DB.DB()
	if err != nil {
		log.Fatalf("Erreur lors de la récupération de la connexion à la base de données : %v", err)
	}
	if err := db.Close(); err != nil {
		log.Fatalf("Erreur lors de la fermeture de la connexion à la base de données : %v", err)
	}
	log.Println("Connexion à la base de données fermée")
}

// BeginTransaction démarre une transaction
func (manager *DBManagerService) BeginTransaction() (*gorm.DB, error) {
	tx := manager.DB.Begin()
	if tx.Error != nil {
		log.Printf("Erreur lors de la création de la transaction : %v", tx.Error)
		return nil, tx.Error
	}
	return tx, nil
}

// Ping vérifie la connexion à la base de données
func (manager *DBManagerService) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Vérification de la connexion
	err := manager.DB.WithContext(ctx).Exec("SELECT 1").Error
	if err != nil {
		log.Printf("Erreur lors de la vérification de la connexion à la base de données : %v", err)
		return err
	}
	log.Println("Connexion à la base de données vérifiée avec succès")
	return nil
}
