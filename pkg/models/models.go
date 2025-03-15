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
	Username string `json:"username" validate:"required,min=3,max=50" example:"john_doe"`     // Username for authentication
	Password string `json:"password" validate:"required,min=5,max=100" example:"password123"` // Password for authentication
}

// Claims represents the JWT claims used for authentication.
// It includes the user ID and expiration time.
type Claims struct {
	UserID    uint  `json:"user_id" example:"1"`      // User ID in the JWT claims
	ExpiresAt int64 `json:"exp" example:"1698765432"` // Expiration time of the JWT token
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
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`                                                                 // Unique identifier for the user
	CreatedAt time.Time `json:"created_at" example:"2025-02-20T19:27:00Z"`                                                        // Timestamp when the user was created
	UpdatedAt time.Time `json:"updated_at" example:"2025-02-20T19:27:00Z"`                                                        // Timestamp when the user was last updated
	Username  string    `json:"username" gorm:"unique;not null" validate:"required,min=3,max=50" example:"john_doe"`              // Unique username for the user
	Email     string    `json:"email,omitempty" gorm:"unique" validate:"required,email,min=5,max=100" example:"john@example.com"` // Unique Email address of the user (omitted unless necessary)
	Password  string    `json:"-" gorm:"not null" validate:"required,min=5,max=100" example:"password123"`                        // Password of the user (never serialized)
	IsAdmin   bool      `json:"is_admin" gorm:"default:false" example:"false"`                                                    // Indicates if the user has admin privileges

	Reviews []Review `json:"reviews,omitempty"` // List of reviews created by the user (omitted unless necessary)
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
	ID          uint      `json:"id" gorm:"primaryKey" example:"1"`                                         // Unique identifier for the anime
	CreatedAt   time.Time `json:"created_at" example:"2025-02-20T19:27:00Z"`                                // Timestamp when the anime was created
	UpdatedAt   time.Time `json:"updated_at" example:"2025-02-20T19:27:00Z"`                                // Timestamp when the anime was last updated
	Title       string    `json:"title" gorm:"not null" validate:"required,min=3,max=100" example:"Naruto"` // Title of the anime
	Description string    `json:"description" validate:"max=500" example:"A story about ninjas."`           // Description of the anime
	Rating      float32   `json:"rating" validate:"gte=0,lte=10" example:"8.5"`                             // Average rating of the anime

	Reviews []Review `json:"reviews,omitempty"` // List of reviews for the anime (omitted unless necessary)
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
	ID        uint      `json:"id" example:"1"`                            // Unique identifier for the user
	Username  string    `json:"username" example:"john_doe"`               // Username of the user
	IsAdmin   bool      `json:"is_admin" example:"false"`                  // Indicates if the user has admin privileges
	CreatedAt time.Time `json:"created_at" example:"2025-02-20T19:27:00Z"` // Timestamp when the user was created
	UpdatedAt time.Time `json:"updated_at" example:"2025-02-20T19:27:00Z"` // Timestamp when the user was last updated
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
  ID          uint      `json:"id" example:"1"`                              // Unique identifier for the anime
	CreatedAt   time.Time `json:"created_at" example:"2025-02-20T19:27:00Z"`   // Timestamp when the anime was created
	UpdatedAt   time.Time `json:"updated_at" example:"2025-02-20T19:27:00Z"`   // Timestamp when the anime was last updated
	Title       string    `json:"title" example:"Naruto"`                      // Title of the anime
	Description string    `json:"description" example:"A story about ninjas."` // Description of the anime
	Rating      float32   `json:"rating" example:"8.5"`                        // Average rating of the anime

	User  User  `gorm:"foreignKey:UserID" json:"-" validate:"-"`                               // User who created the review (excluded from JSON)
	Anime Anime `gorm:"foreignKey:AnimeID;constraint:OnDelete:CASCADE;" json:"-" validate:"-"` // Anime being reviewed (excluded from JSON)
}

// AnimeCreateRequest represents the request payload for creating an anime.
// It includes fields for the title, description, and rating of the anime.
//
// Example:
//
//	{
//	  "title": "Naruto",
//	  "description": "A story about ninjas.",
//	  "rating": 8.5
//	}
type AnimeCreateRequest struct {
	Title       string  `json:"title" validate:"required,min=3,max=100" example:"Naruto"`       // Title of the anime (required, 3-100 characters)
	Description string  `json:"description" validate:"max=500" example:"A story about ninjas."` // Description of the anime (optional, max 500 characters)
	Rating      float32 `json:"rating" validate:"gte=0,lte=10" example:"8.5"`                   // Rating of the anime (0-10)
}

// ReviewCreateRequest represents the request payload for creating a review.
// It includes fields for the user ID, anime ID, review content, and rating.
//
// Example:
//
//	{
//	  "userId": 1,
//	  "animeId": 1,
//	  "content": "This anime is amazing!",
//	  "rating": 9
//	}
type ReviewCreateRequest struct {
	UserID  uint   `json:"userId" validate:"required" example:"1"`                               // ID of the user creating the review (required)
	AnimeID uint   `json:"animeId" validate:"required" example:"1"`                              // ID of the anime being reviewed (required)
	Content string `json:"content" validate:"required,max=500" example:"This anime is amazing!"` // Content of the review (required, max 500 characters)
	Rating  int    `json:"rating" validate:"required,gte=0,lte=10" example:"9"`                  // Rating given in the review (0-10, required)
}

// UserCreateRequest represents the request payload for creating a user.
// It includes fields for the username, email, and password.
//
// Example:
//
//	{
//	  "username": "john_doe",
//	  "email": "john@example.com",
//	  "password": "password123"
//	}
type UserCreateRequest struct {
  Username string `json:"username" validate:"required,min=3,max=50" example:"john_doe"`             // Username for the new user (required, 3-50 characters)
	Email    string `json:"email" validate:"required,email,min=5,max=100" example:"john@example.com"` // Email address for the new user (required, valid email format, 5-100 characters)
	Password string `json:"password" validate:"required,min=5,max=100" example:"password123"`         // Password for the new user (required, 5-100 characters)
}
