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
	ID     int    `gorm:"primaryKey"`
	Name   string
	UserID int
}

type Media struct {
	ID        int    `gorm:"primaryKey"`
	MediaPath string
	Type      string
	AlbumID   int
	PHash     string 
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
