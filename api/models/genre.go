package models

// Genre represents a category or classification for anime.
// It includes the genre name and its relationship with anime.
//
// Example:
//
//	{
//	  "id": 1,
//	  "name": "Action",
//	  "anime": [
//	    {
//	      "id": 1,
//	      "title": "Naruto",
//	      "description": "A story about ninjas.",
//	      "rating": 8.5
//	    }
//	  ]
//	}
type Genre struct {
	BaseModel
	ID    uint    `json:"id" gorm:"primaryKey" example:"1"`               // Unique identifier for the genre
	Name  string  `json:"name" gorm:"not null;unique" example:"Action"`   // Name of the genre
	Anime []Anime `json:"anime,omitempty" gorm:"many2many:anime_genres;"` // Associated anime entries
}

// GenreResponse represents the response format for a genre.
// It includes the genre details and associated anime information.
//
// Example:
//
//	{
//	  "id": 1,
//	  "name": "Action",
//	  "anime": [
//	    {
//	      "id": 1,
//	      "title": "Naruto",
//	      "description": "A story about ninjas.",
//	      "rating": 8.5
//	    }
//	  ]
//	}
type GenreResponse struct {
	ID    uint    `json:"id" example:"1"`        // Unique identifier for the genre
	Name  string  `json:"name" example:"Action"` // Name of the genre
	Anime []Anime `json:"anime,omitempty"`       // Associated anime entries
}

// ToResponse converts a Genre model to a GenreResponse.
// This method is used to serialize genre data for API responses.
func (g Genre) ToResponse() GenreResponse {
	return GenreResponse{
		ID:    g.ID,
		Name:  g.Name,
		Anime: g.Anime,
	}
}
