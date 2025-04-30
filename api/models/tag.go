package models

import (
	"gorm.io/gorm"
)

// Tag represents an anime tag
type Tag struct {
	gorm.Model
	Name  string  `json:"name" gorm:"unique;not null"`
	Anime []Anime `json:"anime,omitempty" gorm:"many2many:anime_tags;"`
}

// TagResponse represents the response format for a tag
type TagResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// ToResponse converts a Tag to a TagResponse
func (t *Tag) ToResponse() TagResponse {
	return TagResponse{
		ID:   t.ID,
		Name: t.Name,
	}
}
