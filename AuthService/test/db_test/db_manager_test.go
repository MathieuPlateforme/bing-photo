package db_test

import (
	"AuthService/pkg/db"
	"testing"
	"github.com/joho/godotenv"
	"log"
)

func init() {
    // Charger les variables d'environnement
    if err := godotenv.Load("../../.env"); err != nil {
        log.Fatalf("Erreur lors du chargement des variables d'environnement : %v", err)
    }
}

// TestNewDBManagerService vérifie si le service DBManager est initialisé correctement
func TestNewDBManagerService(t *testing.T) {
	// Initialiser le service de gestion de la base de données
	dbManager, err := db.NewDBManagerService()
	if err != nil {
		t.Fatalf("Erreur lors de l'initialisation de DBManagerService : %v", err)
	}

	if dbManager.DB == nil {
		t.Fatalf("La connexion à la base de données est nulle")
	}
}

// TestAutoMigrate vérifie si la migration de la base de données s'exécute sans erreur
func TestAutoMigrate(t *testing.T) {
	// Initialiser le service de gestion de la base de données
	dbManager, err := db.NewDBManagerService()
	if err != nil {
		t.Fatalf("Erreur lors de l'initialisation de DBManagerService : %v", err)
	}

	// Tester la migration automatique
	err = dbManager.AutoMigrate()
	if err != nil {
		t.Fatalf("Erreur lors de la migration automatique : %v", err)
	}
}

// TestBeginTransaction vérifie si une transaction peut être commencée sans erreur
func TestBeginTransaction(t *testing.T) {
	// Initialiser le service de gestion de la base de données
	dbManager, err := db.NewDBManagerService()
	if err != nil {
		t.Fatalf("Erreur lors de l'initialisation de DBManagerService : %v", err)
	}

	// Commencer une transaction
	tx, err := dbManager.BeginTransaction()
	if err != nil {
		t.Fatalf("Erreur lors de la création de la transaction : %v", err)
	}

	// Vérifier si l'objet transaction est valide
	if tx == nil {
		t.Fatalf("L'objet transaction est nul")
	}
}