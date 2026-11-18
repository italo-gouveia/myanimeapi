package services

import (
	"context"
	"mime/multipart"

	"github.com/stretchr/testify/mock"

	"myanimeapi/api/models"
	"myanimeapi/api/repositories"
)

// mockUserRepo is a testify mock for repositories.UserRepository.
type mockUserRepo struct{ mock.Mock }

func (m *mockUserRepo) GetByID(ctx context.Context, id uint) (interface{}, error) {
	ret := m.Called(ctx, id)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0), ret.Error(1)
}
func (m *mockUserRepo) GetAll(ctx context.Context, page, limit int) ([]interface{}, int64, error) {
	ret := m.Called(ctx, page, limit)
	if ret.Get(0) == nil {
		return nil, ret.Get(1).(int64), ret.Error(2)
	}
	return ret.Get(0).([]interface{}), ret.Get(1).(int64), ret.Error(2)
}
func (m *mockUserRepo) Create(ctx context.Context, entity interface{}) error {
	return m.Called(ctx, entity).Error(0)
}
func (m *mockUserRepo) Update(ctx context.Context, entity interface{}) error {
	return m.Called(ctx, entity).Error(0)
}
func (m *mockUserRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockUserRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	ret := m.Called(ctx, username)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*models.User), ret.Error(1)
}
func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	ret := m.Called(ctx, email)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*models.User), ret.Error(1)
}

// mockAnimeRepo is a testify mock for repositories.AnimeRepository.
type mockAnimeRepo struct{ mock.Mock }

func (m *mockAnimeRepo) GetByID(ctx context.Context, id uint) (interface{}, error) {
	ret := m.Called(ctx, id)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0), ret.Error(1)
}
func (m *mockAnimeRepo) GetAll(ctx context.Context, page, limit int) ([]interface{}, int64, error) {
	ret := m.Called(ctx, page, limit)
	if ret.Get(0) == nil {
		return nil, ret.Get(1).(int64), ret.Error(2)
	}
	return ret.Get(0).([]interface{}), ret.Get(1).(int64), ret.Error(2)
}
func (m *mockAnimeRepo) Create(ctx context.Context, entity interface{}) error {
	return m.Called(ctx, entity).Error(0)
}
func (m *mockAnimeRepo) Update(ctx context.Context, entity interface{}) error {
	return m.Called(ctx, entity).Error(0)
}
func (m *mockAnimeRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockAnimeRepo) GetByTitle(ctx context.Context, title string, page, limit int) ([]models.Anime, int64, error) {
	ret := m.Called(ctx, title, page, limit)
	if ret.Get(0) == nil {
		return nil, ret.Get(1).(int64), ret.Error(2)
	}
	return ret.Get(0).([]models.Anime), ret.Get(1).(int64), ret.Error(2)
}
func (m *mockAnimeRepo) GetByGenre(ctx context.Context, genre string, page, limit int) ([]models.Anime, int64, error) {
	ret := m.Called(ctx, genre, page, limit)
	if ret.Get(0) == nil {
		return nil, ret.Get(1).(int64), ret.Error(2)
	}
	return ret.Get(0).([]models.Anime), ret.Get(1).(int64), ret.Error(2)
}
func (m *mockAnimeRepo) GetReviewsForAnime(ctx context.Context, animeID uint, page, limit int) ([]models.Review, int64, error) {
	ret := m.Called(ctx, animeID, page, limit)
	if ret.Get(0) == nil {
		return nil, ret.Get(1).(int64), ret.Error(2)
	}
	return ret.Get(0).([]models.Review), ret.Get(1).(int64), ret.Error(2)
}
func (m *mockAnimeRepo) GetWithFilters(ctx context.Context, filter models.AnimeFilter, page, limit int) ([]interface{}, int64, error) {
	ret := m.Called(ctx, filter, page, limit)
	if ret.Get(0) == nil {
		return nil, ret.Get(1).(int64), ret.Error(2)
	}
	return ret.Get(0).([]interface{}), ret.Get(1).(int64), ret.Error(2)
}

// mockReviewRepo is a testify mock for repositories.ReviewRepository.
type mockReviewRepo struct{ mock.Mock }

func (m *mockReviewRepo) GetByID(ctx context.Context, id uint) (interface{}, error) {
	ret := m.Called(ctx, id)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0), ret.Error(1)
}
func (m *mockReviewRepo) GetAll(ctx context.Context, page, limit int) ([]interface{}, int64, error) {
	ret := m.Called(ctx, page, limit)
	if ret.Get(0) == nil {
		return nil, ret.Get(1).(int64), ret.Error(2)
	}
	return ret.Get(0).([]interface{}), ret.Get(1).(int64), ret.Error(2)
}
func (m *mockReviewRepo) Create(ctx context.Context, entity interface{}) error {
	return m.Called(ctx, entity).Error(0)
}
func (m *mockReviewRepo) Update(ctx context.Context, entity interface{}) error {
	return m.Called(ctx, entity).Error(0)
}
func (m *mockReviewRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockReviewRepo) GetByUserID(ctx context.Context, userID uint, page, limit int) ([]models.Review, int64, error) {
	ret := m.Called(ctx, userID, page, limit)
	if ret.Get(0) == nil {
		return nil, ret.Get(1).(int64), ret.Error(2)
	}
	return ret.Get(0).([]models.Review), ret.Get(1).(int64), ret.Error(2)
}
func (m *mockReviewRepo) GetByAnimeID(ctx context.Context, animeID uint, page, limit int) ([]models.Review, int64, error) {
	ret := m.Called(ctx, animeID, page, limit)
	if ret.Get(0) == nil {
		return nil, ret.Get(1).(int64), ret.Error(2)
	}
	return ret.Get(0).([]models.Review), ret.Get(1).(int64), ret.Error(2)
}
func (m *mockReviewRepo) GetByUserAndAnime(ctx context.Context, userID, animeID uint) (*models.Review, error) {
	ret := m.Called(ctx, userID, animeID)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*models.Review), ret.Error(1)
}

// mockFavoriteRepo is a testify mock for repositories.FavoriteRepository.
type mockFavoriteRepo struct{ mock.Mock }

func (m *mockFavoriteRepo) GetByID(ctx context.Context, id uint) (interface{}, error) {
	ret := m.Called(ctx, id)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0), ret.Error(1)
}
func (m *mockFavoriteRepo) GetAll(ctx context.Context, page, limit int) ([]interface{}, int64, error) {
	ret := m.Called(ctx, page, limit)
	if ret.Get(0) == nil {
		return nil, ret.Get(1).(int64), ret.Error(2)
	}
	return ret.Get(0).([]interface{}), ret.Get(1).(int64), ret.Error(2)
}
func (m *mockFavoriteRepo) Create(ctx context.Context, entity interface{}) error {
	return m.Called(ctx, entity).Error(0)
}
func (m *mockFavoriteRepo) Update(ctx context.Context, entity interface{}) error {
	return m.Called(ctx, entity).Error(0)
}
func (m *mockFavoriteRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockFavoriteRepo) GetByUserID(ctx context.Context, userID uint) ([]models.Favorite, error) {
	ret := m.Called(ctx, userID)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).([]models.Favorite), ret.Error(1)
}
func (m *mockFavoriteRepo) GetByUserAndAnime(ctx context.Context, userID, animeID uint) (*models.Favorite, error) {
	ret := m.Called(ctx, userID, animeID)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*models.Favorite), ret.Error(1)
}
func (m *mockFavoriteRepo) DeleteByUserAndAnime(ctx context.Context, userID, animeID uint) error {
	return m.Called(ctx, userID, animeID).Error(0)
}

// mockGenreRepo is a testify mock for repositories.GenreRepository.
type mockGenreRepo struct{ mock.Mock }

func (m *mockGenreRepo) GetByID(ctx context.Context, id uint) (interface{}, error) {
	ret := m.Called(ctx, id)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0), ret.Error(1)
}
func (m *mockGenreRepo) GetAll(ctx context.Context, page, limit int) ([]interface{}, int64, error) {
	ret := m.Called(ctx, page, limit)
	if ret.Get(0) == nil {
		return nil, ret.Get(1).(int64), ret.Error(2)
	}
	return ret.Get(0).([]interface{}), ret.Get(1).(int64), ret.Error(2)
}
func (m *mockGenreRepo) Create(ctx context.Context, entity interface{}) error {
	return m.Called(ctx, entity).Error(0)
}
func (m *mockGenreRepo) Update(ctx context.Context, entity interface{}) error {
	return m.Called(ctx, entity).Error(0)
}
func (m *mockGenreRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockGenreRepo) GetByIDs(ctx context.Context, ids []uint) ([]models.Genre, error) {
	ret := m.Called(ctx, ids)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).([]models.Genre), ret.Error(1)
}
func (m *mockGenreRepo) Search(ctx context.Context, query string, page, limit int) repositories.SearchResult {
	ret := m.Called(ctx, query, page, limit)
	return ret.Get(0).(repositories.SearchResult)
}
func (m *mockGenreRepo) BulkCreate(ctx context.Context, genres []models.Genre) error {
	return m.Called(ctx, genres).Error(0)
}
func (m *mockGenreRepo) BulkDelete(ctx context.Context, ids []uint) error {
	return m.Called(ctx, ids).Error(0)
}

// mockTagRepo is a testify mock for repositories.TagRepository.
type mockTagRepo struct{ mock.Mock }

func (m *mockTagRepo) GetByID(ctx context.Context, id uint) (interface{}, error) {
	ret := m.Called(ctx, id)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0), ret.Error(1)
}
func (m *mockTagRepo) GetAll(ctx context.Context, page, limit int) ([]interface{}, int64, error) {
	ret := m.Called(ctx, page, limit)
	if ret.Get(0) == nil {
		return nil, ret.Get(1).(int64), ret.Error(2)
	}
	return ret.Get(0).([]interface{}), ret.Get(1).(int64), ret.Error(2)
}
func (m *mockTagRepo) Create(ctx context.Context, entity interface{}) error {
	return m.Called(ctx, entity).Error(0)
}
func (m *mockTagRepo) Update(ctx context.Context, entity interface{}) error {
	return m.Called(ctx, entity).Error(0)
}
func (m *mockTagRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockTagRepo) GetByIDs(ctx context.Context, ids []uint) ([]models.Tag, error) {
	ret := m.Called(ctx, ids)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).([]models.Tag), ret.Error(1)
}

// mockWatchlistRepo is a testify mock for repositories.WatchlistRepository.
type mockWatchlistRepo struct{ mock.Mock }

func (m *mockWatchlistRepo) GetByID(ctx context.Context, id uint) (interface{}, error) {
	ret := m.Called(ctx, id)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0), ret.Error(1)
}
func (m *mockWatchlistRepo) GetAll(ctx context.Context, page, limit int) ([]interface{}, int64, error) {
	ret := m.Called(ctx, page, limit)
	if ret.Get(0) == nil {
		return nil, ret.Get(1).(int64), ret.Error(2)
	}
	return ret.Get(0).([]interface{}), ret.Get(1).(int64), ret.Error(2)
}
func (m *mockWatchlistRepo) Create(ctx context.Context, entity interface{}) error {
	return m.Called(ctx, entity).Error(0)
}
func (m *mockWatchlistRepo) Update(ctx context.Context, entity interface{}) error {
	return m.Called(ctx, entity).Error(0)
}
func (m *mockWatchlistRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockWatchlistRepo) GetByUserID(ctx context.Context, userID uint) ([]models.WatchlistEntry, error) {
	ret := m.Called(ctx, userID)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).([]models.WatchlistEntry), ret.Error(1)
}
func (m *mockWatchlistRepo) GetByUserAndAnime(ctx context.Context, userID, animeID uint) (*models.WatchlistEntry, error) {
	ret := m.Called(ctx, userID, animeID)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*models.WatchlistEntry), ret.Error(1)
}
func (m *mockWatchlistRepo) UpsertByUserAndAnime(ctx context.Context, userID, animeID uint, status models.WatchlistStatus) (*models.WatchlistEntry, error) {
	ret := m.Called(ctx, userID, animeID, status)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*models.WatchlistEntry), ret.Error(1)
}
func (m *mockWatchlistRepo) DeleteByUserAndAnime(ctx context.Context, userID, animeID uint) error {
	return m.Called(ctx, userID, animeID).Error(0)
}

// mockStorageSvc is a testify mock for StorageServiceInterface.
type mockStorageSvc struct{ mock.Mock }

func (m *mockStorageSvc) SaveFile(file *multipart.FileHeader, directory string) (string, error) {
	ret := m.Called(file, directory)
	return ret.String(0), ret.Error(1)
}
func (m *mockStorageSvc) DeleteFile(ctx context.Context, fileURL string) error {
	return m.Called(ctx, fileURL).Error(0)
}
func (m *mockStorageSvc) GetBaseURL() string {
	return m.Called().String(0)
}
