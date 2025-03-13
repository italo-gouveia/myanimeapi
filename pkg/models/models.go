// pkg/models/models.go
package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// UserCredentials for authentication
type UserCredentials struct {
	Username string `json:"username" validate:"required,min=3,max=50" example:"john_doe"`     // Username
	Password string `json:"password" validate:"required,min=5,max=100" example:"password123"` // Password
}

// Claims for JWT
type Claims struct {
	UserID    uint  `json:"user_id"`
	ExpiresAt int64 `json:"exp"`
	jwt.RegisteredClaims
}

// User represents a user in the system
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`                                                                 // Review ID
	CreatedAt time.Time `json:"created_at" example:"2025-02-20T19:27:00Z"`                                                        // Creation timestamp
	UpdatedAt time.Time `json:"updated_at" example:"2025-02-20T19:27:00Z"`                                                        // Update timestamp
	Username  string    `json:"username" gorm:"unique;not null" validate:"required,min=3,max=50" example:"john_doe"`              // Username
	Email     string    `json:"email,omitempty" gorm:"unique" validate:"required,email,min=5,max=100" example:"john@example.com"` // Email (omitted unless necessary)
	Password  string    `json:"-" gorm:"not null" validate:"required,min=5,max=100" example:"password123"`                        // Password (never serialized)
	IsAdmin   bool      `json:"is_admin" gorm:"default:false" example:"false"`                                                    // IsAdmin

	Reviews []Review `json:"reviews,omitempty"` // Relationship with Review (omitted unless necessary)
}

// Anime represents an anime entry
type Anime struct {
	ID          uint      `json:"id" gorm:"primaryKey" example:"1"`                                         // Review ID
	CreatedAt   time.Time `json:"created_at" example:"2025-02-20T19:27:00Z"`                                // Creation timestamp
	UpdatedAt   time.Time `json:"updated_at" example:"2025-02-20T19:27:00Z"`                                // Update timestamp
	Title       string    `json:"title" gorm:"not null" validate:"required,min=3,max=100" example:"Naruto"` // Title
	Description string    `json:"description" validate:"max=500" example:"A story about ninjas."`           // Description
	Rating      float32   `json:"rating" validate:"gte=0,lte=10" example:"8.5"`                             // Rating                          // Rating

	Reviews []Review `json:"reviews,omitempty"` // Relationship with Review (omitted unless necessary)
}

// Review represents a review for an anime
type Review struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`                                                     // Review ID
	CreatedAt time.Time `json:"created_at" example:"2025-02-20T19:27:00Z"`                                            // Creation timestamp
	UpdatedAt time.Time `json:"updated_at" example:"2025-02-20T19:27:00Z"`                                            // Update timestamp
	UserID    uint      `json:"userId" gorm:"not null" example:"1"`                                                   // User ID
	AnimeID   uint      `json:"animeId" gorm:"not null;constraint:OnDelete:CASCADE;" example:"1"`                     // Add cascade delete constraint
	Content   string    `json:"content" gorm:"not null" validate:"required,max=500" example:"This anime is amazing!"` // Review content
	Rating    int       `json:"rating" gorm:"not null" validate:"required,gte=0,lte=10" example:"9"`                  // Rating (0-10)

	User  User  `gorm:"foreignKey:UserID" json:"-" validate:"-"`                               // Exclude User from JSON and validation
	Anime Anime `gorm:"foreignKey:AnimeID;constraint:OnDelete:CASCADE;" json:"-" validate:"-"` // Add cascade delete constraint
}

type ReviewResponse struct {
	ID        uint          `json:"id"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	UserID    uint          `json:"userId"`
	AnimeID   uint          `json:"animeId"`
	Content   string        `json:"content"`
	Rating    int           `json:"rating"`
	User      UserResponse  `json:"user"`
	Anime     AnimeResponse `json:"anime"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
}

type AnimeResponse struct {
	ID          uint    `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Rating      float32 `json:"rating"`
}
