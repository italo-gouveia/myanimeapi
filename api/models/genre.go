package models

import (
	"gorm.io/gorm"
)

// Genre represents an anime genre
type Genre struct {
	gorm.Model
	Name  string  `json:"name" gorm:"unique;not null"`
	Anime []Anime `json:"anime,omitempty" gorm:"many2many:anime_genres;"`
}

// GenreResponse represents the response format for a genre
type GenreResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// ToResponse converts a Genre to a GenreResponse
func (g *Genre) ToResponse() GenreResponse {
	return GenreResponse{
		ID:   g.ID,
		Name: g.Name,
	}
}
