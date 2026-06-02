package models

import "time"

// Anime represents an anime entry in the system.
// It includes details about the anime and its relationships with reviews, genres, and tags.
//
// Example:
//
//	{
//	  "id": 1,
//	  "title": "Naruto",
//	  "description": "A story about ninjas.",
//	  "rating": 8.5,
//	  "episodes": 220,
//	  "status": "Completed",
//	  "start_date": "2002-10-03T00:00:00Z",
//	  "end_date": "2007-02-08T00:00:00Z",
//	  "created_at": "2025-02-20T19:27:00Z",
//	  "updated_at": "2025-02-20T19:27:00Z",
//	  "reviews": [
//	    {
//	      "id": 1,
//	      "content": "Great anime!",
//	      "rating": 9.0,
//	      "user_id": 1
//	    }
//	  ],
//	  "genres": [
//	    {
//	      "id": 1,
//	      "name": "Action"
//	    }
//	  ],
//	  "tags": [
//	    {
//	      "id": 1,
//	      "name": "Ninja"
//	    }
//	  ]
//	}
type Anime struct {
	ID          uint      `json:"id" gorm:"primaryKey" example:"1"`                                                    // Unique identifier for the anime
	Title       string    `json:"title" gorm:"not null;index:idx_animes_title" example:"Naruto"`                       // Title of the anime
	Description string    `json:"description" example:"A story about ninjas."`                                         // Description of the anime
	Rating      float64   `json:"rating" example:"8.5"`                                                                // Average rating of the anime
	Episodes    int       `json:"episodes" example:"220"`                                                              // Number of episodes
	Status      string    `json:"status" gorm:"index:idx_animes_status" example:"Completed"`                          // Current status of the anime
	StartDate   time.Time `json:"start_date" example:"2002-10-03T00:00:00Z"`                                           // Date when the anime started airing
	EndDate     time.Time `json:"end_date" example:"2007-02-08T00:00:00Z"`                                             // Date when the anime finished airing
	CoverURL    *string   `json:"cover_url,omitempty" gorm:"column:cover_url" example:"https://cdn.myanimelist.net/…"` // Cover image URL (populated by ETL)
	MALId       *int      `json:"mal_id,omitempty" gorm:"column:mal_id;uniqueIndex:idx_animes_mal_id" example:"20"`    // MyAnimeList ID (used for ETL deduplication)
	CreatedAt   time.Time `json:"created_at" example:"2025-02-20T19:27:00Z"`                                           // Timestamp when the anime was added
	UpdatedAt   time.Time `json:"updated_at" example:"2025-02-20T19:27:00Z"`                                           // Timestamp when the anime was last updated
	Reviews     []Review  `json:"reviews,omitempty" gorm:"foreignKey:AnimeID"`                                         // Associated reviews
	Genres      []Genre   `json:"genres,omitempty" gorm:"many2many:anime_genres"`                                      // Associated genres
	Tags        []Tag     `json:"tags,omitempty" gorm:"many2many:anime_tags"`                                          // Associated tags
}

// AnimeResponse represents the response format for an anime entry.
// It includes the anime details and associated information.
//
// Example:
//
//	{
//	  "id": 1,
//	  "title": "Naruto",
//	  "description": "A story about ninjas.",
//	  "rating": 8.5,
//	  "episodes": 220,
//	  "status": "Completed",
//	  "start_date": "2002-10-03T00:00:00Z",
//	  "end_date": "2007-02-08T00:00:00Z",
//	  "created_at": "2025-02-20T19:27:00Z",
//	  "updated_at": "2025-02-20T19:27:00Z",
//	  "reviews": [
//	    {
//	      "id": 1,
//	      "content": "Great anime!",
//	      "rating": 9.0,
//	      "user_id": 1
//	    }
//	  ],
//	  "genres": [
//	    {
//	      "id": 1,
//	      "name": "Action"
//	    }
//	  ],
//	  "tags": [
//	    {
//	      "id": 1,
//	      "name": "Ninja"
//	    }
//	  ]
//	}
type AnimeResponse struct {
	ID          uint      `json:"id" example:"1"`                                                      // Unique identifier for the anime
	Title       string    `json:"title" example:"Naruto"`                                              // Title of the anime
	Description string    `json:"description" example:"A story about ninjas."`                         // Description of the anime
	Rating      float64   `json:"rating" example:"8.5"`                                                // Average rating of the anime
	Episodes    int       `json:"episodes" example:"220"`                                              // Number of episodes
	Status      string    `json:"status" example:"Completed"`                                          // Current status of the anime
	StartDate   time.Time `json:"start_date" example:"2002-10-03T00:00:00Z"`                           // Date when the anime started airing
	EndDate     time.Time `json:"end_date" example:"2007-02-08T00:00:00Z"`                             // Date when the anime finished airing
	CoverURL    *string   `json:"cover_url,omitempty" example:"https://cdn.myanimelist.net/images/…"`  // Cover image URL
	MALId       *int      `json:"mal_id,omitempty" example:"20"`                                       // MyAnimeList ID
	CreatedAt   time.Time `json:"created_at" example:"2025-02-20T19:27:00Z"`                           // Timestamp when the anime was added
	UpdatedAt   time.Time `json:"updated_at" example:"2025-02-20T19:27:00Z"`                           // Timestamp when the anime was last updated
	Reviews     []AnimeReviewEntry `json:"reviews,omitempty"`                                          // Associated reviews (includes username)
	Genres      []Genre   `json:"genres,omitempty"`                                                    // Associated genres
	Tags        []Tag     `json:"tags,omitempty"`                                                      // Associated tags
}

// AnimeCreateRequest represents the request payload for creating a new anime.
// It includes the required fields for creating an anime entry.
//
// Example:
//
//	{
//	  "title": "Naruto",
//	  "description": "A story about ninjas.",
//	  "episodes": 220,
//	  "status": "Completed",
//	  "start_date": "2002-10-03T00:00:00Z",
//	  "end_date": "2007-02-08T00:00:00Z"
//	}
type AnimeCreateRequest struct {
	Title       string    `json:"title" validate:"required" example:"Naruto"`                    // Title of the anime
	Description string    `json:"description" example:"A story about ninjas."`                   // Description of the anime
	Episodes    int       `json:"episodes" validate:"required" example:"220"`                    // Number of episodes
	Status      string    `json:"status" validate:"required" example:"Completed"`                // Current status of the anime
	StartDate   time.Time `json:"start_date" validate:"required" example:"2002-10-03T00:00:00Z"` // Date when the anime started airing
	EndDate     time.Time `json:"end_date" example:"2007-02-08T00:00:00Z"`                       // Date when the anime finished airing
	Rating      float64   `json:"rating" validate:"required,gte=0,lte=10" example:"8.5"`         // Initial rating for the anime
	GenreIDs    []uint    `json:"genre_ids" validate:"omitempty,dive,gt=0" example:"1,2,3"`      // IDs of genres to associate with the anime
	TagIDs      []uint    `json:"tag_ids" validate:"omitempty,dive,gt=0" example:"1,2,3"`        // IDs of tags to associate with the anime
}

// AnimeUpdateRequest represents the request payload for updating an existing anime.
// It includes the fields that can be updated for an anime entry.
//
// Example:
//
//	{
//	  "title": "Naruto Shippuden",
//	  "description": "The continuation of Naruto's story.",
//	  "episodes": 500,
//	  "status": "Completed",
//	  "start_date": "2007-02-15T00:00:00Z",
//	  "end_date": "2017-03-23T00:00:00Z",
//	  "rating": 8.7
//	}
type AnimeUpdateRequest struct {
	Title       string    `json:"title" validate:"omitempty" example:"Naruto Shippuden"`          // Title of the anime
	Description string    `json:"description" example:"The continuation of Naruto's story."`      // Description of the anime
	Episodes    int       `json:"episodes" validate:"omitempty" example:"500"`                    // Number of episodes
	Status      string    `json:"status" validate:"omitempty" example:"Completed"`                // Current status of the anime
	StartDate   time.Time `json:"start_date" validate:"omitempty" example:"2007-02-15T00:00:00Z"` // Date when the anime started airing
	EndDate     time.Time `json:"end_date" example:"2017-03-23T00:00:00Z"`                        // Date when the anime finished airing
	Rating      float64   `json:"rating" validate:"omitempty,gte=0,lte=10" example:"8.7"`         // Updated rating for the anime
}

// AnimeGenresRequest represents the payload for adding/removing genres from an anime.
type AnimeGenresRequest struct {
	GenreIDs []uint `json:"genre_ids" validate:"required,min=1"`
}

// AnimeTagsRequest represents the payload for adding/removing tags from an anime.
type AnimeTagsRequest struct {
	TagIDs []uint `json:"tag_ids" validate:"required,min=1"`
}

// AnimeReviewEntry is the review shape embedded in AnimeResponse.
// It includes the reviewer's username (populated via User preload).
type AnimeReviewEntry struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"userId"`
	Username  string    `json:"username,omitempty"`
	Content   string    `json:"content"`
	Rating    int       `json:"rating"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converts an Anime model to an AnimeResponse.
func (a Anime) ToResponse() AnimeResponse {
	reviews := make([]AnimeReviewEntry, 0, len(a.Reviews))
	for _, r := range a.Reviews {
		reviews = append(reviews, AnimeReviewEntry{
			ID:        r.ID,
			UserID:    r.UserID,
			Username:  r.User.Username, // populated when Reviews.User is preloaded
			Content:   r.Content,
			Rating:    r.Rating,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		})
	}
	return AnimeResponse{
		ID:          a.ID,
		Title:       a.Title,
		Description: a.Description,
		Rating:      a.Rating,
		Episodes:    a.Episodes,
		Status:      a.Status,
		StartDate:   a.StartDate,
		EndDate:     a.EndDate,
		CoverURL:    a.CoverURL,
		MALId:       a.MALId,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
		Reviews:     reviews,
		Genres:      a.Genres,
		Tags:        a.Tags,
	}
}

// AnimeListResponse represents a paginated list of anime responses.
// It includes the list of animes and pagination metadata.
//
// Example:
//
//	{
//	  "animes": [
//	    {
//	      "id": 1,
//	      "title": "Naruto",
//	      "description": "A story about ninjas.",
//	      "rating": 8.5,
//	      "episodes": 220,
//	      "status": "Completed",
//	      "start_date": "2002-10-03T00:00:00Z",
//	      "end_date": "2007-02-08T00:00:00Z",
//	      "created_at": "2025-02-20T19:27:00Z",
//	      "updated_at": "2025-02-20T19:27:00Z"
//	    }
//	  ],
//	  "total": 1,
//	  "page": 1,
//	  "limit": 10
//	}
type AnimeListResponse struct {
	Animes []AnimeResponse `json:"animes"`             // List of anime responses
	Total  int64           `json:"total" example:"1"`  // Total number of animes
	Page   int             `json:"page" example:"1"`   // Current page number
	Limit  int             `json:"limit" example:"10"` // Number of items per page
}
