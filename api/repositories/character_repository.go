package repositories

import (
	"context"
	"fmt"
	"net/http"

	"myanimeapi/api/models"
	"myanimeapi/internal/db"
	"myanimeapi/internal/errors"
	"myanimeapi/internal/logger"
)

// CharacterRepositoryImpl implements CharacterRepository.
type CharacterRepositoryImpl struct {
	db     db.DBInterface
	logger *logger.Logger
}

// NewCharacterRepository creates a new CharacterRepositoryImpl.
func NewCharacterRepository(dbConn db.DBInterface) CharacterRepository {
	return &CharacterRepositoryImpl{db: dbConn, logger: logger.New()}
}

// GetByID retrieves a character by ID, preloading its anime list.
func (r *CharacterRepositoryImpl) GetByID(ctx context.Context, id uint) (interface{}, error) {
	r.logger.WithField("id", id).Info("Retrieving character by ID")

	var character models.Character
	if err := r.db.WithContext(ctx).Preload("Animes").First(&character, id).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{"id": id, "error": err.Error()}).Error("Character not found")
		return nil, errors.NewError(errors.ErrResourceNotFound, "Character not found",
			fmt.Sprintf("Character with ID %d not found", id),
			http.StatusNotFound, map[string]interface{}{"id": id}, err)
	}

	r.logger.WithField("id", id).Info("Successfully retrieved character")
	return &character, nil
}

// GetAll retrieves all characters with pagination.
func (r *CharacterRepositoryImpl) GetAll(ctx context.Context, page, limit int) ([]interface{}, int64, error) {
	r.logger.WithFields(map[string]interface{}{"page": page, "limit": limit}).Info("Retrieving all characters")

	var characters []models.Character
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.Character{}).Count(&total).Error; err != nil {
		r.logger.WithField("error", err.Error()).Error("Failed to count characters")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to count characters",
			err.Error(), http.StatusInternalServerError, nil, err)
	}

	offset := (page - 1) * limit
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&characters).Error; err != nil {
		r.logger.WithField("error", err.Error()).Error("Failed to retrieve characters")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve characters",
			err.Error(), http.StatusInternalServerError, nil, err)
	}

	result := make([]interface{}, len(characters))
	for i := range characters {
		c := characters[i]
		result[i] = &c
	}

	r.logger.WithFields(map[string]interface{}{"count": len(characters), "total": total}).Info("Successfully retrieved characters")
	return result, total, nil
}

// Create inserts a new character.
func (r *CharacterRepositoryImpl) Create(ctx context.Context, entity interface{}) error {
	r.logger.Info("Creating character")

	character, ok := entity.(*models.Character)
	if !ok {
		return errors.NewError(errors.ErrInvalidInput, "Invalid entity type",
			"Expected *models.Character", http.StatusBadRequest, nil, nil)
	}

	if err := r.db.WithContext(ctx).Create(character).Error; err != nil {
		r.logger.WithField("error", err.Error()).Error("Failed to create character")
		return errors.NewError(errors.ErrInternalServer, "Failed to create character",
			err.Error(), http.StatusInternalServerError, nil, err)
	}

	r.logger.WithField("id", character.ID).Info("Successfully created character")
	return nil
}

// Update saves an existing character.
func (r *CharacterRepositoryImpl) Update(ctx context.Context, entity interface{}) error {
	r.logger.Info("Updating character")

	character, ok := entity.(*models.Character)
	if !ok {
		return errors.NewError(errors.ErrInvalidInput, "Invalid entity type",
			"Expected *models.Character", http.StatusBadRequest, nil, nil)
	}

	if err := r.db.WithContext(ctx).Save(character).Error; err != nil {
		r.logger.WithField("error", err.Error()).Error("Failed to update character")
		return errors.NewError(errors.ErrInternalServer, "Failed to update character",
			err.Error(), http.StatusInternalServerError, nil, err)
	}

	r.logger.WithField("id", character.ID).Info("Successfully updated character")
	return nil
}

// Delete removes a character by ID.
func (r *CharacterRepositoryImpl) Delete(ctx context.Context, id uint) error {
	r.logger.WithField("id", id).Info("Deleting character")

	if err := r.db.WithContext(ctx).Delete(&models.Character{}, id).Error; err != nil {
		r.logger.WithField("error", err.Error()).Error("Failed to delete character")
		return errors.NewError(errors.ErrInternalServer, "Failed to delete character",
			err.Error(), http.StatusInternalServerError, map[string]interface{}{"id": id}, err)
	}

	r.logger.WithField("id", id).Info("Successfully deleted character")
	return nil
}

// GetByAnimeID retrieves paginated characters for a given anime.
func (r *CharacterRepositoryImpl) GetByAnimeID(ctx context.Context, animeID uint, page, limit int) ([]models.Character, int64, error) {
	r.logger.WithFields(map[string]interface{}{"anime_id": animeID, "page": page, "limit": limit}).Info("Retrieving characters for anime")

	var characters []models.Character
	var total int64

	base := r.db.WithContext(ctx).
		Joins("JOIN anime_characters ON characters.id = anime_characters.character_id").
		Where("anime_characters.anime_id = ?", animeID)

	if err := base.Model(&models.Character{}).Count(&total).Error; err != nil {
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to count characters",
			err.Error(), http.StatusInternalServerError, nil, err)
	}

	offset := (page - 1) * limit
	if err := base.Offset(offset).Limit(limit).Find(&characters).Error; err != nil {
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve characters",
			err.Error(), http.StatusInternalServerError, map[string]interface{}{"anime_id": animeID}, err)
	}

	r.logger.WithFields(map[string]interface{}{"count": len(characters), "total": total}).Info("Successfully retrieved characters for anime")
	return characters, total, nil
}

// GetByName retrieves characters whose name contains the given string (case-insensitive).
func (r *CharacterRepositoryImpl) GetByName(ctx context.Context, name string, page, limit int) ([]models.Character, int64, error) {
	r.logger.WithFields(map[string]interface{}{"name": name, "page": page, "limit": limit}).Info("Searching characters by name")

	var characters []models.Character
	var total int64
	offset := (page - 1) * limit

	query := r.db.WithContext(ctx).Where("name ILIKE ?", "%"+name+"%")

	if err := query.Model(&models.Character{}).Count(&total).Error; err != nil {
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to count characters",
			err.Error(), http.StatusInternalServerError, nil, err)
	}

	if err := query.Offset(offset).Limit(limit).Find(&characters).Error; err != nil {
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to search characters",
			err.Error(), http.StatusInternalServerError, nil, err)
	}

	return characters, total, nil
}

// AddToAnime creates rows in the anime_characters junction for each character ID.
func (r *CharacterRepositoryImpl) AddToAnime(ctx context.Context, animeID uint, characterIDs []uint) error {
	r.logger.WithFields(map[string]interface{}{"anime_id": animeID, "count": len(characterIDs)}).Info("Adding characters to anime")

	anime := models.Anime{}
	anime.ID = animeID

	characters := make([]models.Character, len(characterIDs))
	for i, cid := range characterIDs {
		characters[i] = models.Character{}
		characters[i].ID = cid
	}

	if err := r.db.WithContext(ctx).Model(&anime).Association("Characters").Append(characters); err != nil {
		r.logger.WithField("error", err.Error()).Error("Failed to add characters to anime")
		return errors.NewError(errors.ErrInternalServer, "Failed to add characters to anime",
			err.Error(), http.StatusInternalServerError, nil, err)
	}

	r.logger.WithField("anime_id", animeID).Info("Successfully added characters to anime")
	return nil
}

// RemoveFromAnime removes a character from an anime's association list.
func (r *CharacterRepositoryImpl) RemoveFromAnime(ctx context.Context, animeID, characterID uint) error {
	r.logger.WithFields(map[string]interface{}{"anime_id": animeID, "character_id": characterID}).Info("Removing character from anime")

	anime := models.Anime{}
	anime.ID = animeID
	character := models.Character{}
	character.ID = characterID

	if err := r.db.WithContext(ctx).Model(&anime).Association("Characters").Delete(character); err != nil {
		r.logger.WithField("error", err.Error()).Error("Failed to remove character from anime")
		return errors.NewError(errors.ErrInternalServer, "Failed to remove character from anime",
			err.Error(), http.StatusInternalServerError, nil, err)
	}

	r.logger.Info("Successfully removed character from anime")
	return nil
}
