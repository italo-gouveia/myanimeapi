// pkg/models/models.go
// Package models defines the data structures used in the MyAnimeAPI application.
// It includes models for users, anime, reviews, and related responses.
// These models are used for database interactions, request/response payloads, and JWT claims.
package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// UserCredentials represents the credentials used for user authentication.
// It includes a username and password.
//
// Example:
//
//	{
//	  "username": "john_doe",
//	  "password": "password123"
//	}
type UserCredentials struct {
	Username string `json:"username" validate:"required,min=3,max=50" example:"john_doe"`     // Username
	Password string `json:"password" validate:"required,min=5,max=100" example:"password123"` // Password
}

// Claims represents the JWT claims used for authentication.
// It includes the user ID and expiration time.
type Claims struct {
	UserID    uint  `json:"user_id"` // User ID
	ExpiresAt int64 `json:"exp"`     // Expiration time
	jwt.RegisteredClaims
}

// User represents a user in the system.
// It includes fields for user details, authentication, and relationships with reviews.
//
// Example:
//
//	{
//	  "id": 1,
//	  "created_at": "2025-02-20T19:27:00Z",
//	  "updated_at": "2025-02-20T19:27:00Z",
//	  "username": "john_doe",
//	  "email": "john@example.com",
//	  "is_admin": false
//	}
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`                                                                 // User ID
	CreatedAt time.Time `json:"created_at" example:"2025-02-20T19:27:00Z"`                                                        // Creation timestamp
	UpdatedAt time.Time `json:"updated_at" example:"2025-02-20T19:27:00Z"`                                                        // Update timestamp
	Username  string    `json:"username" gorm:"unique;not null" validate:"required,min=3,max=50" example:"john_doe"`              // Username
	Email     string    `json:"email,omitempty" gorm:"unique" validate:"required,email,min=5,max=100" example:"john@example.com"` // Email (omitted unless necessary)
	Password  string    `json:"-" gorm:"not null" validate:"required,min=5,max=100" example:"password123"`                        // Password (never serialized)
	IsAdmin   bool      `json:"is_admin" gorm:"default:false" example:"false"`                                                    // IsAdmin

	Reviews []Review `json:"reviews,omitempty"` // Relationship with Review (omitted unless necessary)
}

// Anime represents an anime entry in the system.
// It includes fields for anime details and relationships with reviews.
//
// Example:
//
//	{
//	  "id": 1,
//	  "created_at": "2025-02-20T19:27:00Z",
//	  "updated_at": "2025-02-20T19:27:00Z",
//	  "title": "Naruto",
//	  "description": "A story about ninjas.",
//	  "rating": 8.5
//	}
type Anime struct {
	ID          uint      `json:"id" gorm:"primaryKey" example:"1"`                                         // Anime ID
	CreatedAt   time.Time `json:"created_at" example:"2025-02-20T19:27:00Z"`                                // Creation timestamp
	UpdatedAt   time.Time `json:"updated_at" example:"2025-02-20T19:27:00Z"`                                // Update timestamp
	Title       string    `json:"title" gorm:"not null" validate:"required,min=3,max=100" example:"Naruto"` // Title
	Description string    `json:"description" validate:"max=500" example:"A story about ninjas."`           // Description
	Rating      float32   `json:"rating" validate:"gte=0,lte=10" example:"8.5"`                             // Rating

	Reviews []Review `json:"reviews,omitempty"` // Relationship with Review (omitted unless necessary)
}

// Review represents a review for an anime.
// It includes fields for review content, ratings, and relationships with users and anime.
//
// Example:
//
//	{
//	  "id": 1,
//	  "created_at": "2025-02-20T19:27:00Z",
//	  "updated_at": "2025-02-20T19:27:00Z",
//	  "userId": 1,
//	  "animeId": 1,
//	  "content": "This anime is amazing!",
//	  "rating": 9
//	}
type Review struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`                                                     // Review ID
	CreatedAt time.Time `json:"created_at" example:"2025-02-20T19:27:00Z"`                                            // Creation timestamp
	UpdatedAt time.Time `json:"updated_at" example:"2025-02-20T19:27:00Z"`                                            // Update timestamp
	UserID    uint      `json:"userId" gorm:"not null" example:"1"`                                                   // User ID
	AnimeID   uint      `json:"animeId" gorm:"not null;constraint:OnDelete:CASCADE;" example:"1"`                     // Anime ID with cascade delete constraint
	Content   string    `json:"content" gorm:"not null" validate:"required,max=500" example:"This anime is amazing!"` // Review content
	Rating    int       `json:"rating" gorm:"not null" validate:"required,gte=0,lte=10" example:"9"`                  // Rating (0-10)

	User  User  `gorm:"foreignKey:UserID" json:"-" validate:"-"`                               // Exclude User from JSON and validation
	Anime Anime `gorm:"foreignKey:AnimeID;constraint:OnDelete:CASCADE;" json:"-" validate:"-"` // Exclude Anime from JSON and validation
}

// ReviewResponse represents a review response that includes user and anime details.
// It is used for serializing review data in API responses.
//
// Example:
//
//	{
//	  "id": 1,
//	  "created_at": "2025-02-20T19:27:00Z",
//	  "updated_at": "2025-02-20T19:27:00Z",
//	  "userId": 1,
//	  "animeId": 1,
//	  "content": "This anime is amazing!",
//	  "rating": 9,
//	  "user": {
//	    "id": 1,
//	    "username": "john_doe",
//	    "is_admin": false
//	  },
//	  "anime": {
//	    "id": 1,
//	    "title": "Naruto",
//	    "description": "A story about ninjas.",
//	    "rating": 8.5
//	  }
//	}
type ReviewResponse struct {
	ID        uint          `json:"id"`         // Review ID
	CreatedAt time.Time     `json:"created_at"` // Creation timestamp
	UpdatedAt time.Time     `json:"updated_at"` // Update timestamp
	UserID    uint          `json:"userId"`     // User ID
	AnimeID   uint          `json:"animeId"`    // Anime ID
	Content   string        `json:"content"`    // Review content
	Rating    int           `json:"rating"`     // Rating
	User      UserResponse  `json:"user"`       // User details
	Anime     AnimeResponse `json:"anime"`      // Anime details
}

// UserResponse represents a simplified user response.
// It is used for serializing user data in API responses.
//
// Example:
//
//	{
//	  "id": 1,
//	  "username": "john_doe",
//	  "is_admin": false
//	}
type UserResponse struct {
	ID       uint   `json:"id"`       // User ID
	Username string `json:"username"` // Username
	IsAdmin  bool   `json:"is_admin"` // IsAdmin
}

// AnimeResponse represents a simplified anime response.
// It is used for serializing anime data in API responses.
//
// Example:
//
//	{
//	  "id": 1,
//	  "title": "Naruto",
//	  "description": "A story about ninjas.",
//	  "rating": 8.5
//	}
type AnimeResponse struct {
	ID          uint    `json:"id"`          // Anime ID
	Title       string  `json:"title"`       // Title
	Description string  `json:"description"` // Description
	Rating      float32 `json:"rating"`      // Rating
}
