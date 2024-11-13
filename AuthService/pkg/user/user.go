package user

import (
	"time" 
)

type User struct {
    ID             int    `gorm:"primaryKey;autoIncrement"`
    Username       string `gorm:"unique;not null"`
    Password       string `gorm:"not null"`
    FirstName      string
    LastName       string
    Email          string `gorm:"unique;not null"`
    IsEmailVerified bool
    GoogleID       string
    PhoneNumber    string
	BirthDate      time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Méthode pour hacher le mot de passe
func (u *User) HashPassword() {
    // Logique de hachage de mot de passe
}

// Méthode pour valider le mot de passe
func (u *User) ValidatePassword() {
    // Logique de validation de mot de passe
}
