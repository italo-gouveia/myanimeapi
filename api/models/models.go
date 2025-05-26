// Package models defines the data structures used in the MyAnimeAPI application.
// It provides models for:
//   - Users and authentication (User, UserCredentials, Claims)
//   - Anime and related entities (Anime, Review, Genre, Tag)
//   - Request/Response structures for API endpoints
//   - Common utilities (JSON, BaseModel)
//
// These models are used for:
//   - Database interactions through GORM
//   - Request/response payloads in API endpoints
//   - JWT authentication and authorization
//   - Data validation using struct tags
//
// Each model includes:
//   - JSON serialization tags
//   - GORM database tags
//   - Validation rules
//   - Example values for documentation
//   - Comprehensive godoc comments
package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JSON represents a JSON object that can be stored in the database
type JSON map[string]interface{}

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

// AuthResponse represents the response payload for successful authentication.
// It includes the JWT token for subsequent authenticated requests.
//
// Example:
//
//	{
//	  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
//	}
type AuthResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // JWT token for authentication
}

// Claims represents the JWT claims used for authentication.
// It includes the user ID and expiration time.
type Claims struct {
	UserID    uint  `json:"user_id" example:"1"`      // User ID in the JWT claims
	ExpiresAt int64 `json:"exp" example:"1698765432"` // Expiration time of the JWT token
	jwt.RegisteredClaims
}

// MediaAttachment represents a media file attached to a review.
// It includes fields for the file type, URL, and relationship with the review.
//
// Example:
//
//	{
//	  "id": 1,
//	  "review_id": 1,
//	  "type": "image",
//	  "url": "https://example.com/media/image1.jpg",
//	  "created_at": "2025-02-20T19:27:00Z"
//	}
type MediaAttachment struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`                                                    // Media attachment ID
	ReviewID  uint      `json:"review_id" gorm:"not null;constraint:OnDelete:CASCADE;" example:"1"`                  // Review ID with cascade delete constraint
	Type      string    `json:"type" gorm:"not null" validate:"required,oneof=image video" example:"image"`          // Type of media (image or video)
	URL       string    `json:"url" gorm:"not null" validate:"required,url" example:"https://example.com/image.jpg"` // URL to the media file
	CreatedAt time.Time `json:"created_at" example:"2025-02-20T19:27:00Z"`                                           // Creation timestamp
	UpdatedAt time.Time `json:"updated_at" example:"2025-02-20T19:27:00Z"`                                           // Update timestamp

	Review Review `gorm:"foreignKey:ReviewID" json:"-" validate:"-"` // Exclude Review from JSON and validation
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
//	  "rating": 9,
//	  "media_attachments": [
//	    {
//	      "id": 1,
//	      "type": "image",
//	      "url": "https://example.com/image.jpg"
//	    }
//	  ]
//	}
type Review struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`                                                     // Review ID
	CreatedAt time.Time `json:"created_at" example:"2025-02-20T19:27:00Z"`                                            // Creation timestamp
	UpdatedAt time.Time `json:"updated_at" example:"2025-02-20T19:27:00Z"`                                            // Update timestamp
	UserID    uint      `json:"userId" gorm:"not null" example:"1"`                                                   // User ID
	AnimeID   uint      `json:"animeId" gorm:"not null;constraint:OnDelete:CASCADE;" example:"1"`                     // Anime ID with cascade delete constraint
	Content   string    `json:"content" gorm:"not null" validate:"required,max=500" example:"This anime is amazing!"` // Review content
	Rating    int       `json:"rating" gorm:"not null" validate:"required,gte=0,lte=10" example:"9"`                  // Rating (0-10)

	User             User              `gorm:"foreignKey:UserID" json:"-" validate:"-"`                               // Exclude User from JSON and validation
	Anime            Anime             `gorm:"foreignKey:AnimeID;constraint:OnDelete:CASCADE;" json:"-" validate:"-"` // Exclude Anime from JSON and validation
	MediaAttachments []MediaAttachment `json:"media_attachments,omitempty" gorm:"foreignKey:ReviewID"`                // Associated media attachments
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
//	  },
//	  "media_attachments": [
//	    {
//	      "id": 1,
//	      "type": "image",
//	      "url": "https://example.com/image.jpg"
//	    }
//	  ]
//	}
type ReviewResponse struct {
	ID               uint              `json:"id" example:"1"`                            // Unique identifier for the review
	CreatedAt        time.Time         `json:"created_at" example:"2025-02-20T19:27:00Z"` // Timestamp when the review was created
	UpdatedAt        time.Time         `json:"updated_at" example:"2025-02-20T19:27:00Z"` // Timestamp when the review was last updated
	UserID           uint              `json:"userId" example:"1"`                        // ID of the user who created the review
	AnimeID          uint              `json:"animeId" example:"1"`                       // ID of the anime being reviewed
	Content          string            `json:"content" example:"This anime is amazing!"`  // Content of the review
	Rating           int               `json:"rating" example:"9"`                        // Rating given in the review (0-10)
	User             UserResponse      `json:"user"`                                      // User details
	Anime            AnimeResponse     `json:"anime"`                                     // Anime details
	MediaAttachments []MediaAttachment `json:"media_attachments,omitempty"`               // Associated media attachments
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

// ReviewCreateRequest represents the request payload for creating a review.
// It includes fields for the user ID, anime ID, review content, rating, and media attachments.
//
// Example:
//
//	{
//	  "userId": 1,
//	  "animeId": 1,
//	  "content": "This anime is amazing!",
//	  "rating": 9,
//	  "media_attachments": [
//	    {
//	      "type": "image",
//	      "url": "https://example.com/image.jpg"
//	    }
//	  ]
//	}
type ReviewCreateRequest struct {
	UserID           uint              `json:"userId" validate:"required" example:"1"`                               // ID of the user creating the review (required)
	AnimeID          uint              `json:"animeId" validate:"required" example:"1"`                              // ID of the anime being reviewed (required)
	Content          string            `json:"content" validate:"required,max=500" example:"This anime is amazing!"` // Content of the review (required, max 500 characters)
	Rating           int               `json:"rating" validate:"required,gte=0,lte=10" example:"9"`                  // Rating given in the review (0-10, required)
	MediaAttachments []MediaAttachment `json:"media_attachments,omitempty" validate:"omitempty,dive"`                // Media attachments for the review (optional)
}

// ReviewUpdateRequest represents the request payload for updating a review.
// It includes optional fields that can be updated, with validation rules for each field.
//
// Example:
//
//	{
//	  "content": "Updated review content",
//	  "rating": 8,
//	  "media_attachments": [
//	    {
//	      "type": "image",
//	      "url": "https://example.com/new-image.jpg"
//	    }
//	  ]
//	}
type ReviewUpdateRequest struct {
	Content          string            `json:"content" validate:"omitempty,max=500" example:"Updated review content"` // Updated content of the review (optional, max 500 characters)
	Rating           int               `json:"rating" validate:"omitempty,gte=0,lte=10" example:"8"`                  // Updated rating (optional, 0-10)
	MediaAttachments []MediaAttachment `json:"media_attachments,omitempty" validate:"omitempty,dive"`                 // Updated media attachments (optional)
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

// PasswordChangeRequest represents the request payload for changing a user's password.
// It includes fields for the current password and the new password.
//
// Example:
//
//	{
//	  "current_password": "old_password",
//	  "new_password": "new_password"
//	}
type PasswordChangeRequest struct {
	CurrentPassword string `json:"current_password" validate:"required,min=5,max=100" example:"old_password"` // Current password (required, 5-100 characters)
	NewPassword     string `json:"new_password" validate:"required,min=5,max=100" example:"new_password"`     // New password (required, 5-100 characters)
}

// Response represents a generic API response structure.
// It provides a standardized format for all API responses with:
//   - Status indicator (success/error)
//   - Descriptive message
//   - Optional data payload
//
// The Response type is used across all API endpoints to ensure
// consistent response formatting and error handling.
//
// Example:
//
//	{
//	  "status": "success",
//	  "message": "Operation completed successfully",
//	  "data": {
//	    "id": 1,
//	    "username": "john_doe"
//	  }
//	}
//
// TODO: In the future, consider implementing a more comprehensive generic response type
// that includes pagination metadata, error details, and other common response fields.
// This would standardize API responses across all endpoints and make it easier to
// add new features like pagination, filtering, and sorting
// Error Response Example:
//
//	{
//	  "status": "error",
//	  "message": "Invalid input parameters",
//	  "data": {
//	    "field": "username",
//	    "error": "must be between 3 and 50 characters"
//	  }
//	}
type Response struct {
	Status  string      `json:"status" example:"success"`               // Status of the response (success, error)
	Message string      `json:"message" example:"Operation successful"` // Message describing the result
	Data    interface{} `json:"data"`                                   // Data payload (can be any type)
}

// SocialLinks represents a user's social media links.
// It includes fields for various social media platforms.
//
// Example:
//
//	{
//	  "twitter": "https://twitter.com/johndoe",
//	  "instagram": "https://instagram.com/johndoe"
//	}
type SocialLinks struct {
	Twitter   string `json:"twitter,omitempty" validate:"omitempty,url" example:"https://twitter.com/johndoe"`     // Twitter profile URL
	Instagram string `json:"instagram,omitempty" validate:"omitempty,url" example:"https://instagram.com/johndoe"` // Instagram profile URL
}
