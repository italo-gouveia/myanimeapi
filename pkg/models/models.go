// pkg/models/models.go
package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// UserCredentials for authentication
type UserCredentials struct {
	Username string `json:"username" validate:"required,min=3,max=50" example:"john_doe"`     // Username for authentication
	Password string `json:"password" validate:"required,min=5,max=100" example:"password123"` // Password for authentication
}

// Claims for JWT
type Claims struct {
	UserID    uint  `json:"user_id" example:"1"`      // User ID in the JWT claims
	ExpiresAt int64 `json:"exp" example:"1698765432"` // Expiration time of the JWT token
	jwt.RegisteredClaims
}

// User represents a user in the system
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`                                                                 // Unique identifier for the user
	CreatedAt time.Time `json:"created_at" example:"2025-02-20T19:27:00Z"`                                                        // Timestamp when the user was created
	UpdatedAt time.Time `json:"updated_at" example:"2025-02-20T19:27:00Z"`                                                        // Timestamp when the user was last updated
	Username  string    `json:"username" gorm:"unique;not null" validate:"required,min=3,max=50" example:"john_doe"`              // Unique username for the user
	Email     string    `json:"email,omitempty" gorm:"unique" validate:"required,email,min=5,max=100" example:"john@example.com"` // Email address of the user (omitted unless necessary)
	Password  string    `json:"-" gorm:"not null" validate:"required,min=5,max=100" example:"password123"`                        // Password of the user (never serialized)
	IsAdmin   bool      `json:"is_admin" gorm:"default:false" example:"false"`                                                    // Indicates if the user has admin privileges

	Reviews []Review `json:"reviews,omitempty"` // List of reviews created by the user (omitted unless necessary)
}

// Anime represents an anime entry
type Anime struct {
	ID          uint      `json:"id" gorm:"primaryKey" example:"1"`                                         // Unique identifier for the anime
	CreatedAt   time.Time `json:"created_at" example:"2025-02-20T19:27:00Z"`                                // Timestamp when the anime was created
	UpdatedAt   time.Time `json:"updated_at" example:"2025-02-20T19:27:00Z"`                                // Timestamp when the anime was last updated
	Title       string    `json:"title" gorm:"not null" validate:"required,min=3,max=100" example:"Naruto"` // Title of the anime
	Description string    `json:"description" validate:"max=500" example:"A story about ninjas."`           // Description of the anime
	Rating      float32   `json:"rating" validate:"gte=0,lte=10" example:"8.5"`                             // Average rating of the anime

	Reviews []Review `json:"reviews,omitempty"` // List of reviews for the anime (omitted unless necessary)
}

// Review represents a review for an anime
type Review struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`                                                     // Unique identifier for the review
	CreatedAt time.Time `json:"created_at" example:"2025-02-20T19:27:00Z"`                                            // Timestamp when the review was created
	UpdatedAt time.Time `json:"updated_at" example:"2025-02-20T19:27:00Z"`                                            // Timestamp when the review was last updated
	UserID    uint      `json:"userId" gorm:"not null" example:"1"`                                                   // ID of the user who created the review
	AnimeID   uint      `json:"animeId" gorm:"not null;constraint:OnDelete:CASCADE;" example:"1"`                     // ID of the anime being reviewed
	Content   string    `json:"content" gorm:"not null" validate:"required,max=500" example:"This anime is amazing!"` // Content of the review
	Rating    int       `json:"rating" gorm:"not null" validate:"required,gte=0,lte=10" example:"9"`                  // Rating given in the review (0-10)

	User  User  `gorm:"foreignKey:UserID" json:"-" validate:"-"`                               // User who created the review (excluded from JSON)
	Anime Anime `gorm:"foreignKey:AnimeID;constraint:OnDelete:CASCADE;" json:"-" validate:"-"` // Anime being reviewed (excluded from JSON)
}

// ReviewResponse represents the response payload for a review
type ReviewResponse struct {
	ID        uint          `json:"id" example:"1"`                            // Unique identifier for the review
	CreatedAt time.Time     `json:"created_at" example:"2025-02-20T19:27:00Z"` // Timestamp when the review was created
	UpdatedAt time.Time     `json:"updated_at" example:"2025-02-20T19:27:00Z"` // Timestamp when the review was last updated
	UserID    uint          `json:"userId" example:"1"`                        // ID of the user who created the review
	AnimeID   uint          `json:"animeId" example:"1"`                       // ID of the anime being reviewed
	Content   string        `json:"content" example:"This anime is amazing!"`  // Content of the review
	Rating    int           `json:"rating" example:"9"`                        // Rating given in the review (0-10)
	User      UserResponse  `json:"user"`                                      // User details
	Anime     AnimeResponse `json:"anime"`                                     // Anime details
}

// UserResponse represents a response for a user, excluding sensitive information
type UserResponse struct {
	ID       uint   `json:"id" example:"1"`              // Unique identifier for the user
	Username string `json:"username" example:"john_doe"` // Username of the user
	IsAdmin  bool   `json:"is_admin" example:"false"`    // Indicates if the user has admin privileges
}

// AnimeResponse represents the response payload for an anime
type AnimeResponse struct {
	ID          uint    `json:"id" example:"1"`                              // Unique identifier for the anime
	Title       string  `json:"title" example:"Naruto"`                      // Title of the anime
	Description string  `json:"description" example:"A story about ninjas."` // Description of the anime
	Rating      float32 `json:"rating" example:"8.5"`                        // Average rating of the anime
}

// AnimeCreateRequest represents the request payload for creating an anime
type AnimeCreateRequest struct {
	Title       string  `json:"title" validate:"required,min=3,max=100" example:"Naruto"`       // Title of the anime
	Description string  `json:"description" validate:"max=500" example:"A story about ninjas."` // Description of the anime
	Rating      float32 `json:"rating" validate:"gte=0,lte=10" example:"8.5"`                   // Rating of the anime
}

// ReviewCreateRequest represents the request payload for creating a review
type ReviewCreateRequest struct {
	UserID  uint   `json:"userId" validate:"required" example:"1"`                               // ID of the user creating the review
	AnimeID uint   `json:"animeId" validate:"required" example:"1"`                              // ID of the anime being reviewed
	Content string `json:"content" validate:"required,max=500" example:"This anime is amazing!"` // Content of the review
	Rating  int    `json:"rating" validate:"required,gte=0,lte=10" example:"9"`                  // Rating given in the review (0-10)
}
