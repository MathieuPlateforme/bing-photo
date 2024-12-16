package models

import (
	"fmt"
	"log"
	"time"
	"gorm.io/gorm"
)

// Modèle RevokedToken
type RevokedToken struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Token     string    `gorm:"unique;not null"`
	Username  string    `gorm:"not null"`
	RevokedAt time.Time `gorm:"autoCreateTime"`
}

// Méthode pour révoquer un token
func (r *RevokedToken) RevokeToken(db *gorm.DB, token string, username string) error {
	// Crée une instance de RevokedToken
	r.Token = token
	r.Username = username
	r.RevokedAt = time.Now()

	// Enregistre le token dans la base de données
	if err := db.Create(r).Error; err != nil {
		log.Printf("Erreur lors de l'invalidation du token : %v", err)
		return fmt.Errorf("erreur lors de l'invalidation du token")
	}
	return nil
}


