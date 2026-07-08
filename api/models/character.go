package models

// Character represents a character that appears in one or more anime.
//
// Example:
//
//	{
//	  "id": 1,
//	  "name": "Naruto Uzumaki",
//	  "description": "The main protagonist.",
//	  "voice_actor": "Junko Takeuchi",
//	  "image_url": "https://cdn.example.com/naruto.jpg",
//	  "created_at": "2025-02-20T19:27:00Z",
//	  "updated_at": "2025-02-20T19:27:00Z"
//	}
type Character struct {
	BaseModel
	Name        string  `json:"name" gorm:"not null;index:idx_characters_name" validate:"required,max=150" example:"Naruto Uzumaki"` // Character name
	Description string  `json:"description,omitempty" gorm:"type:text" validate:"omitempty,max=1000" example:"The main protagonist."` // Brief description
	VoiceActor  string  `json:"voice_actor,omitempty" validate:"omitempty,max=150" example:"Junko Takeuchi"`                           // Voice actor name
	ImageURL    *string `json:"image_url,omitempty" validate:"omitempty,url" example:"https://cdn.example.com/naruto.jpg"`             // Image URL

	Animes []Anime `json:"animes,omitempty" gorm:"many2many:anime_characters"` // Anime appearances
}

// CharacterResponse is the API response shape for a Character.
type CharacterResponse struct {
	ID          uint    `json:"id" example:"1"`
	Name        string  `json:"name" example:"Naruto Uzumaki"`
	Description string  `json:"description,omitempty" example:"The main protagonist."`
	VoiceActor  string  `json:"voice_actor,omitempty" example:"Junko Takeuchi"`
	ImageURL    *string `json:"image_url,omitempty" example:"https://cdn.example.com/naruto.jpg"`
}

// CharacterListResponse wraps a paginated character list.
type CharacterListResponse struct {
	Characters []CharacterResponse `json:"characters"`
	Total      int64               `json:"total" example:"100"`
	Page       int                 `json:"page" example:"1"`
	Limit      int                 `json:"limit" example:"10"`
}

// CharacterCreateRequest is the body for POST /characters.
type CharacterCreateRequest struct {
	Name        string  `json:"name" validate:"required,max=150" example:"Naruto Uzumaki"`
	Description string  `json:"description" validate:"omitempty,max=1000" example:"The main protagonist."`
	VoiceActor  string  `json:"voice_actor" validate:"omitempty,max=150" example:"Junko Takeuchi"`
	ImageURL    *string `json:"image_url,omitempty" validate:"omitempty,url" example:"https://cdn.example.com/naruto.jpg"`
}

// CharacterUpdateRequest is the body for PUT /characters/{id}.
type CharacterUpdateRequest struct {
	Name        string  `json:"name" validate:"omitempty,max=150" example:"Naruto Uzumaki"`
	Description string  `json:"description" validate:"omitempty,max=1000" example:"The main protagonist."`
	VoiceActor  string  `json:"voice_actor" validate:"omitempty,max=150" example:"Junko Takeuchi"`
	ImageURL    *string `json:"image_url,omitempty" validate:"omitempty,url" example:"https://cdn.example.com/naruto.jpg"`
}

// AnimeCharacterRequest is the body for POST /animes/{id}/characters.
type AnimeCharacterRequest struct {
	CharacterIDs []uint `json:"character_ids" validate:"required,min=1,dive,gt=0" example:"1,2,3"`
}

// ToResponse converts a Character model to CharacterResponse.
func (c Character) ToResponse() CharacterResponse {
	return CharacterResponse{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		VoiceActor:  c.VoiceActor,
		ImageURL:    c.ImageURL,
	}
}
