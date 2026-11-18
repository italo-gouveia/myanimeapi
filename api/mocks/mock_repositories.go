package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"myanimeapi/api/models"
	"myanimeapi/api/repositories"
)

// MockUserRepository is a mock for repositories.UserRepository.
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (interface{}, error) {
	ret := m.Called(ctx, id)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0), ret.Error(1)
}

func (m *MockUserRepository) GetAll(ctx context.Context, page, limit int) ([]interface{}, int64, error) {
	ret := m.Called(ctx, page, limit)
	if ret.Get(0) == nil {
		return nil, ret.Get(1).(int64), ret.Error(2)
	}
	return ret.Get(0).([]interface{}), ret.Get(1).(int64), ret.Error(2)
}

func (m *MockUserRepository) Create(ctx context.Context, entity interface{}) error {
	ret := m.Called(ctx, entity)
	return ret.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, entity interface{}) error {
	ret := m.Called(ctx, entity)
	return ret.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	ret := m.Called(ctx, id)
	return ret.Error(0)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	ret := m.Called(ctx, username)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*models.User), ret.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	ret := m.Called(ctx, email)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*models.User), ret.Error(1)
}

// MockAnimeRepository is a mock for repositories.AnimeRepository.
type MockAnimeRepository struct {
	mock.Mock
}

func (m *MockAnimeRepository) GetByID(ctx context.Context, id uint) (interface{}, error) {
	ret := m.Called(ctx, id)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0), ret.Error(1)
}

func (m *MockAnimeRepository) GetAll(ctx context.Context, page, limit int) ([]interface{}, int64, error) {
	ret := m.Called(ctx, page, limit)
	if ret.Get(0) == nil {
		return nil, ret.Get(1).(int64), ret.Error(2)
	}
	return ret.Get(0).([]interface{}), ret.Get(1).(int64), ret.Error(2)
}

func (m *MockAnimeRepository) Create(ctx context.Context, entity interface{}) error {
	ret := m.Called(ctx, entity)
	return ret.Error(0)
}

func (m *MockAnimeRepository) Update(ctx context.Context, entity interface{}) error {
	ret := m.Called(ctx, entity)
	return ret.Error(0)
}

func (m *MockAnimeRepository) Delete(ctx context.Context, id uint) error {
	ret := m.Called(ctx, id)
	return ret.Error(0)
}

func (m *MockAnimeRepository) GetByTitle(ctx context.Context, title string, page, limit int) ([]models.Anime, int64, error) {
	ret := m.Called(ctx, title, page, limit)
	if ret.Get(0) == nil {
		return nil, ret.Get(1).(int64), ret.Error(2)
	}
	return ret.Get(0).([]models.Anime), ret.Get(1).(int64), ret.Error(2)
}

func (m *MockAnimeRepository) GetByGenre(ctx context.Context, genre string, page, limit int) ([]models.Anime, int64, error) {
	ret := m.Called(ctx, genre, page, limit)
	if ret.Get(0) == nil {
		return nil, ret.Get(1).(int64), ret.Error(2)
	}
	return ret.Get(0).([]models.Anime), ret.Get(1).(int64), ret.Error(2)
}

func (m *MockAnimeRepository) GetReviewsForAnime(ctx context.Context, animeID uint, page, limit int) ([]models.Review, int64, error) {
	ret := m.Called(ctx, animeID, page, limit)
	if ret.Get(0) == nil {
		return nil, ret.Get(1).(int64), ret.Error(2)
	}
	return ret.Get(0).([]models.Review), ret.Get(1).(int64), ret.Error(2)
}

func (m *MockAnimeRepository) GetWithFilters(ctx context.Context, filter models.AnimeFilter, page, limit int) ([]interface{}, int64, error) {
	ret := m.Called(ctx, filter, page, limit)
	if ret.Get(0) == nil {
		return nil, ret.Get(1).(int64), ret.Error(2)
	}
	return ret.Get(0).([]interface{}), ret.Get(1).(int64), ret.Error(2)
}

// MockReviewRepository is a mock for repositories.ReviewRepository.
type MockReviewRepository struct {
	mock.Mock
}

func (m *MockReviewRepository) GetByID(ctx context.Context, id uint) (interface{}, error) {
	ret := m.Called(ctx, id)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0), ret.Error(1)
}

func (m *MockReviewRepository) GetAll(ctx context.Context, page, limit int) ([]interface{}, int64, error) {
	ret := m.Called(ctx, page, limit)
	if ret.Get(0) == nil {
		return nil, ret.Get(1).(int64), ret.Error(2)
	}
	return ret.Get(0).([]interface{}), ret.Get(1).(int64), ret.Error(2)
}

func (m *MockReviewRepository) Create(ctx context.Context, entity interface{}) error {
	ret := m.Called(ctx, entity)
	return ret.Error(0)
}

func (m *MockReviewRepository) Update(ctx context.Context, entity interface{}) error {
	ret := m.Called(ctx, entity)
	return ret.Error(0)
}

func (m *MockReviewRepository) Delete(ctx context.Context, id uint) error {
	ret := m.Called(ctx, id)
	return ret.Error(0)
}

func (m *MockReviewRepository) GetByUserID(ctx context.Context, userID uint, page, limit int) ([]models.Review, int64, error) {
	ret := m.Called(ctx, userID, page, limit)
	if ret.Get(0) == nil {
		return nil, ret.Get(1).(int64), ret.Error(2)
	}
	return ret.Get(0).([]models.Review), ret.Get(1).(int64), ret.Error(2)
}

func (m *MockReviewRepository) GetByAnimeID(ctx context.Context, animeID uint, page, limit int) ([]models.Review, int64, error) {
	ret := m.Called(ctx, animeID, page, limit)
	if ret.Get(0) == nil {
		return nil, ret.Get(1).(int64), ret.Error(2)
	}
	return ret.Get(0).([]models.Review), ret.Get(1).(int64), ret.Error(2)
}

func (m *MockReviewRepository) GetByUserAndAnime(ctx context.Context, userID, animeID uint) (*models.Review, error) {
	ret := m.Called(ctx, userID, animeID)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*models.Review), ret.Error(1)
}

// MockFavoriteRepository is a mock for repositories.FavoriteRepository.
type MockFavoriteRepository struct {
	mock.Mock
}

func (m *MockFavoriteRepository) GetByID(ctx context.Context, id uint) (interface{}, error) {
	ret := m.Called(ctx, id)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0), ret.Error(1)
}

func (m *MockFavoriteRepository) GetAll(ctx context.Context, page, limit int) ([]interface{}, int64, error) {
	ret := m.Called(ctx, page, limit)
	if ret.Get(0) == nil {
		return nil, ret.Get(1).(int64), ret.Error(2)
	}
	return ret.Get(0).([]interface{}), ret.Get(1).(int64), ret.Error(2)
}

func (m *MockFavoriteRepository) Create(ctx context.Context, entity interface{}) error {
	ret := m.Called(ctx, entity)
	return ret.Error(0)
}

func (m *MockFavoriteRepository) Update(ctx context.Context, entity interface{}) error {
	ret := m.Called(ctx, entity)
	return ret.Error(0)
}

func (m *MockFavoriteRepository) Delete(ctx context.Context, id uint) error {
	ret := m.Called(ctx, id)
	return ret.Error(0)
}

func (m *MockFavoriteRepository) GetByUserID(ctx context.Context, userID uint) ([]models.Favorite, error) {
	ret := m.Called(ctx, userID)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).([]models.Favorite), ret.Error(1)
}

func (m *MockFavoriteRepository) GetByUserAndAnime(ctx context.Context, userID, animeID uint) (*models.Favorite, error) {
	ret := m.Called(ctx, userID, animeID)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*models.Favorite), ret.Error(1)
}

func (m *MockFavoriteRepository) DeleteByUserAndAnime(ctx context.Context, userID, animeID uint) error {
	ret := m.Called(ctx, userID, animeID)
	return ret.Error(0)
}

// MockGenreRepository is a mock for repositories.GenreRepository.
type MockGenreRepository struct {
	mock.Mock
}

func (m *MockGenreRepository) GetByID(ctx context.Context, id uint) (interface{}, error) {
	ret := m.Called(ctx, id)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0), ret.Error(1)
}

func (m *MockGenreRepository) GetAll(ctx context.Context, page, limit int) ([]interface{}, int64, error) {
	ret := m.Called(ctx, page, limit)
	if ret.Get(0) == nil {
		return nil, ret.Get(1).(int64), ret.Error(2)
	}
	return ret.Get(0).([]interface{}), ret.Get(1).(int64), ret.Error(2)
}

func (m *MockGenreRepository) Create(ctx context.Context, entity interface{}) error {
	ret := m.Called(ctx, entity)
	return ret.Error(0)
}

func (m *MockGenreRepository) Update(ctx context.Context, entity interface{}) error {
	ret := m.Called(ctx, entity)
	return ret.Error(0)
}

func (m *MockGenreRepository) Delete(ctx context.Context, id uint) error {
	ret := m.Called(ctx, id)
	return ret.Error(0)
}

func (m *MockGenreRepository) GetByIDs(ctx context.Context, ids []uint) ([]models.Genre, error) {
	ret := m.Called(ctx, ids)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).([]models.Genre), ret.Error(1)
}

func (m *MockGenreRepository) Search(ctx context.Context, query string, page, limit int) repositories.SearchResult {
	ret := m.Called(ctx, query, page, limit)
	return ret.Get(0).(repositories.SearchResult)
}

func (m *MockGenreRepository) BulkCreate(ctx context.Context, genres []models.Genre) error {
	ret := m.Called(ctx, genres)
	return ret.Error(0)
}

func (m *MockGenreRepository) BulkDelete(ctx context.Context, ids []uint) error {
	ret := m.Called(ctx, ids)
	return ret.Error(0)
}

// MockTagRepository is a mock for repositories.TagRepository.
type MockTagRepository struct {
	mock.Mock
}

func (m *MockTagRepository) GetByID(ctx context.Context, id uint) (interface{}, error) {
	ret := m.Called(ctx, id)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0), ret.Error(1)
}

func (m *MockTagRepository) GetAll(ctx context.Context, page, limit int) ([]interface{}, int64, error) {
	ret := m.Called(ctx, page, limit)
	if ret.Get(0) == nil {
		return nil, ret.Get(1).(int64), ret.Error(2)
	}
	return ret.Get(0).([]interface{}), ret.Get(1).(int64), ret.Error(2)
}

func (m *MockTagRepository) Create(ctx context.Context, entity interface{}) error {
	ret := m.Called(ctx, entity)
	return ret.Error(0)
}

func (m *MockTagRepository) Update(ctx context.Context, entity interface{}) error {
	ret := m.Called(ctx, entity)
	return ret.Error(0)
}

func (m *MockTagRepository) Delete(ctx context.Context, id uint) error {
	ret := m.Called(ctx, id)
	return ret.Error(0)
}

func (m *MockTagRepository) GetByIDs(ctx context.Context, ids []uint) ([]models.Tag, error) {
	ret := m.Called(ctx, ids)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).([]models.Tag), ret.Error(1)
}
