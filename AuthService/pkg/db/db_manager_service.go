package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"AuthService/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBManagerService struct {
	DB *gorm.DB
}

func NewDBManagerService() (*DBManagerService, error) {
	fmt.Println("Initializing DBService...")

	// Charger les variables d'environnement
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")

	// Construire la chaîne de connexion
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Ouvrir la connexion avec GORM
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la connexion à la base de données : %v", err)
	}

	log.Println("Connexion à la base de données réussie")

	return &DBManagerService{DB: db}, nil
}

func (manager *DBManagerService) AutoMigrate() error {
	err := manager.DB.AutoMigrate(&models.User{}, &models.RevokedToken{})
	if err != nil {
		return fmt.Errorf("erreur lors de la migration de la base de données : %v", err)
	}
	log.Println("Migration de la base de données réussie")
	return nil
}

func (manager *DBManagerService) CloseConnection() {
	// Fermer la connexion à la base de données
	db, err := manager.DB.DB()
	if err != nil {
		log.Fatalf("Erreur lors de la fermeture de la connexion à la base de données : %v", err)
	}
	if err := db.Close(); err != nil {
		log.Fatalf("Erreur lors de la fermeture de la connexion à la base de données : %v", err)
	}
	log.Println("Connexion à la base de données fermée")
}

func (manager *DBManagerService) BeginTransaction() (*gorm.DB, error) {
	// Commencer une transaction
	tx := manager.DB.Begin()

	// Vérifier s'il y a eu une erreur lors du démarrage de la transaction
	if tx.Error != nil {
		log.Printf("Erreur lors de la création de la transaction : %v", tx.Error)
		return nil, tx.Error
	}

	// Retourner l'objet transaction
	return tx, nil
}

func (manager *DBManagerService) Ping(tx *gorm.DB) error {
	// Créez un contexte avec un délai d'expiration
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Utilisez le contexte pour exécuter la requête
	err := tx.WithContext(ctx).Exec("SELECT 1").Error
	if err != nil {
		log.Printf("Erreur lors de la vérification de la connexion à la base de données : %v", err)
		return err
	}
	return nil
}
