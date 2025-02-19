// pkg/models/models.go
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
	Username string `gorm:"unique;not null"` // Username is required and unique
	Email    string `gorm:"unique"`          // Email is unique but nullable
	Password string `gorm:"not null"`        // Password is required
	IsAdmin  bool   `gorm:"default:false"`   // Default value for IsAdmin

	Reviews []Review // Relationship with Review
}

type Anime struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string
	Rating      float32

	Reviews []Review // Relationship with Review
}

type Review struct {
	gorm.Model
	UserID  uint   `gorm:"not null"`
	AnimeID uint   `gorm:"not null"`
	Content string `gorm:"not null"`
	Rating  int    `gorm:"not null"`

	User  User  `gorm:"foreignKey:UserID"`  // Relationship with User
	Anime Anime `gorm:"foreignKey:AnimeID"` // Relationship with Anime
}
