package models

import "time"

// Favorite represents a user's favorite anime entry
// @Description Model representing a user's favorite anime
type Favorite struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	AnimeID   uint      `json:"anime_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Anime     Anime     `json:"anime,omitempty" gorm:"foreignKey:AnimeID"`
	User      User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// FavoriteCreateRequest represents the request body for creating a favorite
// @Description Request model for adding an anime to favorites
type FavoriteCreateRequest struct {
	AnimeID uint `json:"anime_id" validate:"required" example:"1"`
}

// FavoriteResponse represents the response for favorite operations
// @Description Response model for favorite operations
type FavoriteResponse struct {
	ID        uint      `json:"id" example:"1"`
	UserID    uint      `json:"user_id" example:"1"`
	AnimeID   uint      `json:"anime_id" example:"1"`
	CreatedAt time.Time `json:"created_at" example:"2024-04-02T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-04-02T12:00:00Z"`
	Anime     Anime     `json:"anime,omitempty"`
}

// TableName specifies the table name for the Favorite model
// @Description Returns the table name for the Favorite model
// @Return string "Table name for the Favorite model"
func (Favorite) TableName() string {
	return "favorites"
}
