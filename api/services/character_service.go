package services

import (
	"context"
	"net/http"

	"myanimeapi/api/models"
	"myanimeapi/api/repositories"
	"myanimeapi/internal/errors"
	"myanimeapi/internal/logger"
)

// CharacterServiceInterface defines the character business-logic contract.
type CharacterServiceInterface interface {
	GetCharacterByID(ctx context.Context, id uint) (*models.Character, error)
	GetAllCharacters(ctx context.Context, page, limit int) ([]*models.Character, int64, error)
	GetCharactersByAnime(ctx context.Context, animeID uint, page, limit int) ([]*models.Character, int64, error)
	GetCharactersByName(ctx context.Context, name string, page, limit int) ([]*models.Character, int64, error)
	CreateCharacter(ctx context.Context, character *models.Character) error
	UpdateCharacter(ctx context.Context, character *models.Character) error
	DeleteCharacter(ctx context.Context, id uint) error
	AddCharactersToAnime(ctx context.Context, animeID uint, characterIDs []uint) error
	RemoveCharacterFromAnime(ctx context.Context, animeID, characterID uint) error
}

// CharacterService implements CharacterServiceInterface.
type CharacterService struct {
	characterRepo repositories.CharacterRepository
	logger        *logger.Logger
}

// NewCharacterService creates a new CharacterService.
func NewCharacterService(characterRepo repositories.CharacterRepository) *CharacterService {
	return &CharacterService{characterRepo: characterRepo, logger: logger.New()}
}

// GetCharacterByID retrieves a single character by its ID.
func (s *CharacterService) GetCharacterByID(ctx context.Context, id uint) (*models.Character, error) {
	s.logger.WithField("id", id).Info("Getting character by ID")

	result, err := s.characterRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return result.(*models.Character), nil
}

// GetAllCharacters returns a paginated list of all characters.
func (s *CharacterService) GetAllCharacters(ctx context.Context, page, limit int) ([]*models.Character, int64, error) {
	s.logger.WithFields(map[string]interface{}{"page": page, "limit": limit}).Info("Getting all characters")

	results, total, err := s.characterRepo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, 0, err
	}

	characters := make([]*models.Character, len(results))
	for i, r := range results {
		characters[i] = r.(*models.Character)
	}
	return characters, total, nil
}

// GetCharactersByAnime returns paginated characters for a given anime.
func (s *CharacterService) GetCharactersByAnime(ctx context.Context, animeID uint, page, limit int) ([]*models.Character, int64, error) {
	s.logger.WithField("anime_id", animeID).Info("Getting characters for anime")

	chars, total, err := s.characterRepo.GetByAnimeID(ctx, animeID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	out := make([]*models.Character, len(chars))
	for i := range chars {
		out[i] = &chars[i]
	}
	return out, total, nil
}

// GetCharactersByName returns characters matching a name substring.
func (s *CharacterService) GetCharactersByName(ctx context.Context, name string, page, limit int) ([]*models.Character, int64, error) {
	s.logger.WithField("name", name).Info("Searching characters by name")

	chars, total, err := s.characterRepo.GetByName(ctx, name, page, limit)
	if err != nil {
		return nil, 0, err
	}

	out := make([]*models.Character, len(chars))
	for i := range chars {
		out[i] = &chars[i]
	}
	return out, total, nil
}

// CreateCharacter persists a new character.
func (s *CharacterService) CreateCharacter(ctx context.Context, character *models.Character) error {
	s.logger.WithField("name", character.Name).Info("Creating character")
	return s.characterRepo.Create(ctx, character)
}

// UpdateCharacter updates an existing character.
func (s *CharacterService) UpdateCharacter(ctx context.Context, character *models.Character) error {
	s.logger.WithField("id", character.ID).Info("Updating character")
	return s.characterRepo.Update(ctx, character)
}

// DeleteCharacter deletes a character and its anime associations.
func (s *CharacterService) DeleteCharacter(ctx context.Context, id uint) error {
	s.logger.WithField("id", id).Info("Deleting character")

	if _, err := s.characterRepo.GetByID(ctx, id); err != nil {
		return errors.NewError(errors.ErrResourceNotFound, "Character not found",
			"No character with that ID exists", http.StatusNotFound,
			map[string]interface{}{"id": id}, err)
	}
	return s.characterRepo.Delete(ctx, id)
}

// AddCharactersToAnime associates characters with an anime.
func (s *CharacterService) AddCharactersToAnime(ctx context.Context, animeID uint, characterIDs []uint) error {
	s.logger.WithFields(map[string]interface{}{"anime_id": animeID, "count": len(characterIDs)}).Info("Adding characters to anime")
	return s.characterRepo.AddToAnime(ctx, animeID, characterIDs)
}

// RemoveCharacterFromAnime removes a character from an anime's cast.
func (s *CharacterService) RemoveCharacterFromAnime(ctx context.Context, animeID, characterID uint) error {
	s.logger.WithFields(map[string]interface{}{"anime_id": animeID, "character_id": characterID}).Info("Removing character from anime")
	return s.characterRepo.RemoveFromAnime(ctx, animeID, characterID)
}
