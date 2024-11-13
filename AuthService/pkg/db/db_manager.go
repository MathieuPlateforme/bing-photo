package db

import (
	"fmt"
	"log"
	"os"
	_ "github.com/lib/pq" 

	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"authservice/pkg/user"
)

var DB *gorm.DB

func Connect() {
	var err error

    // Logique pour se connecter à la base de données
    user := os.Getenv("POSTGRES_USER")
    password := os.Getenv("POSTGRES_PASSWORD")
    dbname := os.Getenv("POSTGRES_DB")
    host := os.Getenv("POSTGRES_HOST")
    port := os.Getenv("POSTGRES_PORT")

    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

    // Ouvrir une connexion avec GORM
	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("Erreur lors de la connexion à la base de données : %v", err)
	}

	log.Println("Connexion à la base de données réussie")

}

func AutoMigrate() {
    // Exécuter la migration pour créer la table `user`
    err := DB.AutoMigrate(&user.User{})
    if err != nil {
        log.Fatalf("Erreur lors de la migration de la base de données : %v", err)
    }
    log.Println("Migration de la base de données réussie")
}


