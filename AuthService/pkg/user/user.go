package user

import (
	"time" 
    "gorm.io/gorm"
)

type User struct {
    ID             int    `gorm:"primaryKey;autoIncrement"`
    Username       string `gorm:"unique;not null"`
    Password       string `gorm:"not null"`
    FirstName      string
    LastName       string
    Email          string `gorm:"unique;not null"`
    GoogleID       string
    PhoneNumber    string
	BirthDate      time.Time
    ResetToken     string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}


// Méthode pour valider le mot de passe
func (u *User) ValidatePassword() {
    // Logique de validation de mot de passe
}

func (u *User) UpdateResetToken(db *gorm.DB, token string) error {
    u.ResetToken = token
    return db.Save(&u).Error
}
