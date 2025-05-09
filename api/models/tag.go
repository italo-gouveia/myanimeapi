package models

// Tag represents a descriptive label or keyword for anime.
// It includes the tag name and its relationship with anime.
//
// Example:
//
//	{
//	  "id": 1,
//	  "name": "Ninja",
//	  "anime": [
//	    {
//	      "id": 1,
//	      "title": "Naruto",
//	      "description": "A story about ninjas.",
//	      "rating": 8.5
//	    }
//	  ]
//	}
type Tag struct {
	BaseModel
	ID    uint    `json:"id" gorm:"primaryKey" example:"1"`             // Unique identifier for the tag
	Name  string  `json:"name" gorm:"not null;unique" example:"Ninja"`  // Name of the tag
	Anime []Anime `json:"anime,omitempty" gorm:"many2many:anime_tags;"` // Associated anime entries
}

// TagResponse represents the response format for a tag.
// It includes the tag details and associated anime information.
//
// Example:
//
//	{
//	  "id": 1,
//	  "name": "Ninja",
//	  "anime": [
//	    {
//	      "id": 1,
//	      "title": "Naruto",
//	      "description": "A story about ninjas.",
//	      "rating": 8.5
//	    }
//	  ]
//	}
type TagResponse struct {
	ID    uint    `json:"id" example:"1"`       // Unique identifier for the tag
	Name  string  `json:"name" example:"Ninja"` // Name of the tag
	Anime []Anime `json:"anime,omitempty"`      // Associated anime entries
}

// ToResponse converts a Tag model to a TagResponse.
// This method is used to serialize tag data for API responses.
func (t *Tag) ToResponse() TagResponse {
	return TagResponse{
		ID:    t.ID,
		Name:  t.Name,
		Anime: t.Anime,
	}
}
