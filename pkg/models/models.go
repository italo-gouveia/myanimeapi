package models

import (
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

// UserCredentials for authentication
type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Claims for JWT
type Claims struct {
	UserID    uint  `json:"user_id"`
	ExpiresAt int64 `json:"exp"`
	jwt.StandardClaims
}

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string
	IsAdmin  bool // Add this field to indicate if a user is an admin
}

type Anime struct {
	gorm.Model
	Title       string
	Description string
	Rating      float32
}

type Review struct {
	gorm.Model
	UserID  uint
	AnimeID uint
	Content string
	Rating  int
}
