// pkg/models/models.go
package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// UserCredentials for authentication
type UserCredentials struct {
	Username string `json:"username" example:"john_doe"`    // Username
	Password string `json:"password" example:"password123"` // Password
}

// Claims for JWT
type Claims struct {
	UserID    uint  `json:"user_id"`
	ExpiresAt int64 `json:"exp"`
	jwt.StandardClaims
}

// User represents a user in the system
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`                   // Review ID
	CreatedAt time.Time `json:"created_at" example:"2025-02-20T19:27:00Z"`          // Creation timestamp
	UpdatedAt time.Time `json:"updated_at" example:"2025-02-20T19:27:00Z"`          // Update timestamp
	Username  string    `json:"username" gorm:"unique;not null" example:"john_doe"` // Username
	Email     string    `json:"email" gorm:"unique" example:"john@example.com"`     // Email
	Password  string    `json:"password" gorm:"not null" example:"password123"`     // Password
	IsAdmin   bool      `json:"is_admin" gorm:"default:false" example:"false"`      // IsAdmin

	Reviews []Review // Relationship with Review
}

// Anime represents an anime entry
type Anime struct {
	ID          uint      `json:"id" gorm:"primaryKey" example:"1"`            // Review ID
	CreatedAt   time.Time `json:"created_at" example:"2025-02-20T19:27:00Z"`   // Creation timestamp
	UpdatedAt   time.Time `json:"updated_at" example:"2025-02-20T19:27:00Z"`   // Update timestamp
	Title       string    `json:"title" gorm:"not null" example:"Naruto"`      // Title
	Description string    `json:"description" example:"A story about ninjas."` // Description
	Rating      float32   `json:"rating" example:"8.5"`                        // Rating

	Reviews []Review // Relationship with Review
}

// Review represents a review for an anime
type Review struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`                         // Review ID
	CreatedAt time.Time `json:"created_at" example:"2025-02-20T19:27:00Z"`                // Creation timestamp
	UpdatedAt time.Time `json:"updated_at" example:"2025-02-20T19:27:00Z"`                // Update timestamp
	UserID    uint      `json:"user_id" gorm:"not null" example:"1"`                      // User ID
	AnimeID   uint      `json:"anime_id" gorm:"not null" example:"1"`                     // Anime ID
	Content   string    `json:"content" gorm:"not null" example:"This anime is amazing!"` // Review content
	Rating    int       `json:"rating" gorm:"not null" example:"9"`                       // Rating (0-10)

	User  User  `gorm:"foreignKey:UserID"`  // Relationship with User
	Anime Anime `gorm:"foreignKey:AnimeID"` // Relationship with Anime
}
