package models

import "time"

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

	Reviews []Review `json:"reviews,omitempty"`                               // List of reviews for the anime (omitted unless necessary)
	Genres  []Genre  `json:"genres,omitempty" gorm:"many2many:anime_genres;"` // List of genres for the anime
	Tags    []Tag    `json:"tags,omitempty" gorm:"many2many:anime_tags;"`     // List of tags for the anime
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

	User   User            `gorm:"foreignKey:UserID" json:"-" validate:"-"`                               // User who created the review (excluded from JSON)
	Anime  Anime           `gorm:"foreignKey:AnimeID;constraint:OnDelete:CASCADE;" json:"-" validate:"-"` // Anime being reviewed (excluded from JSON)
	Genres []GenreResponse `json:"genres,omitempty"`
	Tags   []TagResponse   `json:"tags,omitempty"`
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
	GenreIDs    []uint  `json:"genre_ids" example:"[1,2,3]"`                                    // List of genre IDs to associate with the anime
	TagIDs      []uint  `json:"tag_ids" example:"[1,2,3]"`                                      // List of tag IDs to associate with the anime
}

// ToResponse converts an Anime to an AnimeResponse
func (a *Anime) ToResponse() AnimeResponse {
	genreResponses := make([]GenreResponse, len(a.Genres))
	for i, genre := range a.Genres {
		genreResponses[i] = genre.ToResponse()
	}

	tagResponses := make([]TagResponse, len(a.Tags))
	for i, tag := range a.Tags {
		tagResponses[i] = tag.ToResponse()
	}

	return AnimeResponse{
		ID:          a.ID,
		Title:       a.Title,
		Description: a.Description,
		Rating:      a.Rating,
		Genres:      genreResponses,
		Tags:        tagResponses,
	}
}
