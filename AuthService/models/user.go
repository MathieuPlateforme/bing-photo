package models

import (
	"AuthService/pkg/security"
	"errors"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          int    `gorm:"primaryKey;autoIncrement"`
	Username    string `gorm:"not null"`
	Password    string `gorm:"not null"`
	FirstName   string
	LastName    string
	Email       string `gorm:"unique;not null"`
	GoogleID    string
	PhoneNumber string
	BirthDate   time.Time
	ResetToken  string
	Picture		string `gorm:"column:picture"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (u *User) CreateUser(db *gorm.DB) error {
	// Vérifier si l'utilisateur existe déjà
	var existingUser User
	if err := db.Where("email = ?", u.Email).First(&existingUser).Error; err == nil {
		return errors.New("l'utilisateur avec cet email existe déjà")
	}

	// Créer le nouvel utilisateur
	return db.Create(&u).Error
}

func (u *User) GetUserByEmail(db *gorm.DB, email string) error {
	if err := db.Where("email = ?", email).First(u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("utilisateur non trouvé")
		}
		return err
	}
	return nil
}

// GetUserByID récupère un utilisateur par ID
func (u *User) GetUserByID(db *gorm.DB, id uint) error {
	if err := db.First(u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("utilisateur non trouvé")
		}
		return err
	}
	return nil
}

func (u *User) UpdatePassword(db *gorm.DB, newPassword string) error {
	u.Password = newPassword
	u.UpdatedAt = time.Now()
	return db.Save(&u).Error
}

// Méthode pour valider le mot de passe
func (u *User) ValidatePassword(db *gorm.DB, password string, securityService *security.SecurityService) bool {
	return securityService.ComparePasswords(u.Password, password)
}

func (u *User) UpdateResetToken(db *gorm.DB, token string) error {
	u.ResetToken = token
	u.UpdatedAt = time.Now()
	return db.Save(&u).Error
}

func (u *User) UpdateUser(db *gorm.DB, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	return db.Model(u).Updates(updates).Error
}
