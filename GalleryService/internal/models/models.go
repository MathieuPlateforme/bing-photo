package models

import (
	"time"
)

type User struct {
	ID       int    `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string
}

type Album struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"unique;not null"`
	UserID      uint      `gorm:"not null"`
	BucketName  string    `gorm:"not null"`
	Description string    
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Champ pour indiquer si le bucket existe dans S3 
	ExistsInS3 bool `gorm:"-"`
}

type Media struct {
	ID         uint   `gorm:"primaryKey"`
	AlbumID    uint   `gorm:"not null"`
	Path       string `gorm:"not null"`
	Name       string `gorm:"not null"`
	Type       string
	IsFavorite bool   `gorm:"default:false"`
	Hash       string `gorm:"not null"`
}

type SimilarGroup struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type SimilarMedia struct {
	ID             uint    `gorm:"primaryKey"`
	SimilarGroupID uint    `gorm:"not null"`
	MediaID        uint    `gorm:"not null"`
	SimilarityScore float64 `gorm:"not null"` 
}


type Access struct {
	ID             int    `gorm:"primaryKey"`
	MediaID        int
	UserAccessID   int
	ExpirationDate time.Time
	Code           string
	Status         string
	Link           string
	Type           string
}

type UserAccess struct {
	ID     int    `gorm:"primaryKey"`
	Name   string
	UserID int
}
