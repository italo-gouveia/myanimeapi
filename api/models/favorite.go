package models

import "time"

// Favorite represents a user's favorite anime entry.
// It includes fields for the user ID, anime ID, and timestamps.
//
// Example:
//
//	{
//	  "id": 1,
//	  "user_id": 1,
//	  "anime_id": 1,
//	  "created_at": "2025-02-20T19:27:00Z",
//	  "updated_at": "2025-02-20T19:27:00Z"
//	}
type Favorite struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`          // Unique identifier for the favorite entry
	UserID    uint      `json:"user_id" gorm:"not null" example:"1"`       // ID of the user who favorited the anime
	AnimeID   uint      `json:"anime_id" gorm:"not null" example:"1"`      // ID of the favorited anime
	CreatedAt time.Time `json:"created_at" example:"2025-02-20T19:27:00Z"` // Timestamp when the favorite was created
	UpdatedAt time.Time `json:"updated_at" example:"2025-02-20T19:27:00Z"` // Timestamp when the favorite was last updated
	Anime     Anime     `json:"anime,omitempty" gorm:"foreignKey:AnimeID"` // Associated anime details
	User      User      `json:"user,omitempty" gorm:"foreignKey:UserID"`   // Associated user details
}

// FavoriteCreateRequest represents the request payload for creating a favorite.
// It includes the anime ID to be favorited.
//
// Example:
//
//	{
//	  "anime_id": 1
//	}
type FavoriteCreateRequest struct {
	AnimeID uint `json:"anime_id" validate:"required" example:"1"` // ID of the anime to favorite
}

// FavoriteResponse represents the response format for a favorite entry.
// It includes the favorite details and associated anime information.
//
// Example:
//
//	{
//	  "id": 1,
//	  "user_id": 1,
//	  "anime_id": 1,
//	  "created_at": "2025-02-20T19:27:00Z",
//	  "updated_at": "2025-02-20T19:27:00Z",
//	  "anime": {
//	    "id": 1,
//	    "title": "Naruto",
//	    "description": "A story about ninjas.",
//	    "rating": 8.5
//	  }
//	}
type FavoriteResponse struct {
	ID        uint      `json:"id" example:"1"`                            // Unique identifier for the favorite entry
	UserID    uint      `json:"user_id" example:"1"`                       // ID of the user who favorited the anime
	AnimeID   uint      `json:"anime_id" example:"1"`                      // ID of the favorited anime
	CreatedAt time.Time `json:"created_at" example:"2025-02-20T19:27:00Z"` // Timestamp when the favorite was created
	UpdatedAt time.Time `json:"updated_at" example:"2025-02-20T19:27:00Z"` // Timestamp when the favorite was last updated
	Anime     Anime     `json:"anime,omitempty"`                           // Associated anime details
}

// TableName specifies the table name for the Favorite model.
// This method is used by GORM to determine the database table name.
func (Favorite) TableName() string {
	return "favorites"
}

// ToResponse converts a Favorite model to a FavoriteResponse.
func (f Favorite) ToResponse() FavoriteResponse {
	return FavoriteResponse{
		ID:        f.ID,
		UserID:    f.UserID,
		AnimeID:   f.AnimeID,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
		Anime:     f.Anime,
	}
}
