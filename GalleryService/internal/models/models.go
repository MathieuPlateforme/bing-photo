package models

import (
	"time"
)

type User struct {
	ID               uint      `gorm:"primaryKey"`
	Email            string    `gorm:"unique;not null"`
	Username         string    `gorm:"not null"`
	PrivateAlbumPin  string    `gorm:"default:null"` 
	PrivateAlbumID   uint      `gorm:"default:null"`
	PrivateAlbum     *Album     `gorm:"foreignKey:PrivateAlbumID"`
	MainAlbumID   	 uint      `gorm:"default:null"`
	MainAlbum        *Album     `gorm:"foreignKey:PrivateAlbumID"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type Album struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"unique;not null"`
	UserID      uint      `gorm:"not null"`
	BucketName  string    `gorm:"not null"`
	IsPrivate   bool      `gorm:"default:false"` 
	IsMain   	bool      `gorm:"default:false"` 
	Description string    
	Media       []Media     `gorm:"foreignKey:AlbumID"` 
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Champ pour indiquer si le bucket existe dans S3 
	ExistsInS3 bool `gorm:"-"`
}

type Media struct {
	ID         uint   `gorm:"primaryKey"`
	AlbumID    uint   `gorm:"not null"`
	Album      *Album  `gorm:"foreignKey:AlbumID"`
	Path       string `gorm:"not null"`
	Name       string `gorm:"not null"`
	Type       string
	IsFavorite bool   `gorm:"default:false"`
	Hash 	   *string `gorm:"column:hash;not null"`
	FileSize   uint   `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
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
	ID             uint      `gorm:"primaryKey"`
	MediaID        uint      `gorm:"not null"`
	Code           string    `gorm:"unique;not null"` 
	IsPrivate      bool      `gorm:"default:false"`   
	Pin            string    `gorm:"-"`              
	PinHash        string    `gorm:"default:null"`    
	ExpirationDate time.Time `gorm:"default:null"`    
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type UserAccess struct {
	ID     int    `gorm:"primaryKey"`
	Name   string
	UserID int
}
