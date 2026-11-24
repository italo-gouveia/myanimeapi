package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"myanimeapi/api/adapters/cache"
	"myanimeapi/api/models"
	"myanimeapi/api/repositories"
)

func newGenreSvc(repo *mockGenreRepo) *GenreService {
	return &GenreService{
		genreRepo: repo,
		logger:    newTestLogger(),
		cache:     cache.NewNoopCache(),
	}
}

func TestGenreService_GetGenreByID_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockGenreRepo{}
	svc := newGenreSvc(repo)

	genre := &models.Genre{Name: "Action"}
	genre.ID = 1
	repo.On("GetByID", ctx, uint(1)).Return(genre, nil)

	result, err := svc.GetGenreByID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, "Action", result.Name)
}

func TestGenreService_GetGenreByID_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockGenreRepo{}
	svc := newGenreSvc(repo)

	repo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("not found"))

	_, err := svc.GetGenreByID(ctx, 1)
	require.Error(t, err)
}

func TestGenreService_GetAllGenres_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockGenreRepo{}
	svc := newGenreSvc(repo)

	g1 := &models.Genre{Name: "Action"}
	g1.ID = 1
	g2 := &models.Genre{Name: "Comedy"}
	g2.ID = 2
	ifaces := []interface{}{g1, g2}
	repo.On("GetAll", ctx, 1, 10).Return(ifaces, int64(2), nil)

	result, total, err := svc.GetAllGenres(ctx, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, result, 2)
}

func TestGenreService_GetAllGenres_Error(t *testing.T) {
	ctx := context.Background()
	repo := &mockGenreRepo{}
	svc := newGenreSvc(repo)

	repo.On("GetAll", ctx, 1, 10).Return(nil, int64(0), errors.New("db error"))

	_, _, err := svc.GetAllGenres(ctx, 1, 10)
	require.Error(t, err)
}

func TestGenreService_CreateGenre_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockGenreRepo{}
	svc := newGenreSvc(repo)

	repo.On("GetAll", ctx, 1, 1000).Return([]interface{}{}, int64(0), nil)
	repo.On("Create", ctx, mock.Anything).Return(nil)

	genre := &models.Genre{Name: "NewGenre"}
	err := svc.CreateGenre(ctx, genre)
	require.NoError(t, err)
}

func TestGenreService_CreateGenre_NameConflict(t *testing.T) {
	ctx := context.Background()
	repo := &mockGenreRepo{}
	svc := newGenreSvc(repo)

	existing := &models.Genre{Name: "Action"}
	existing.ID = 1
	repo.On("GetAll", ctx, 1, 1000).Return([]interface{}{existing}, int64(1), nil)

	genre := &models.Genre{Name: "Action"}
	err := svc.CreateGenre(ctx, genre)
	require.Error(t, err)
}

func TestGenreService_CreateGenre_CreateError(t *testing.T) {
	ctx := context.Background()
	repo := &mockGenreRepo{}
	svc := newGenreSvc(repo)

	repo.On("GetAll", ctx, 1, 1000).Return([]interface{}{}, int64(0), nil)
	repo.On("Create", ctx, mock.Anything).Return(errors.New("db error"))

	genre := &models.Genre{Name: "NewGenre"}
	err := svc.CreateGenre(ctx, genre)
	require.Error(t, err)
}

func TestGenreService_UpdateGenre_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockGenreRepo{}
	svc := newGenreSvc(repo)

	existing := &models.Genre{Name: "Action"}
	existing.ID = 1
	repo.On("GetByID", ctx, uint(1)).Return(existing, nil)
	repo.On("Update", ctx, mock.Anything).Return(nil)

	genre := &models.Genre{Name: "Action"}
	genre.ID = 1
	err := svc.UpdateGenre(ctx, genre)
	require.NoError(t, err)
}

func TestGenreService_UpdateGenre_NameConflict(t *testing.T) {
	ctx := context.Background()
	repo := &mockGenreRepo{}
	svc := newGenreSvc(repo)

	existing := &models.Genre{Name: "Action"}
	existing.ID = 1
	other := &models.Genre{Name: "Comedy"}
	other.ID = 2
	repo.On("GetByID", ctx, uint(1)).Return(existing, nil)
	repo.On("GetAll", ctx, 1, 1000).Return([]interface{}{existing, other}, int64(2), nil)

	genre := &models.Genre{Name: "Comedy"}
	genre.ID = 1
	err := svc.UpdateGenre(ctx, genre)
	require.Error(t, err)
}

func TestGenreService_DeleteGenre_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockGenreRepo{}
	svc := newGenreSvc(repo)

	existing := &models.Genre{Name: "Action"}
	existing.ID = 1
	repo.On("GetByID", ctx, uint(1)).Return(existing, nil)
	repo.On("Delete", ctx, uint(1)).Return(nil)

	err := svc.DeleteGenre(ctx, 1)
	require.NoError(t, err)
}

func TestGenreService_DeleteGenre_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockGenreRepo{}
	svc := newGenreSvc(repo)

	repo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("not found"))

	err := svc.DeleteGenre(ctx, 1)
	require.Error(t, err)
}

func TestGenreService_GetGenresByIDs_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockGenreRepo{}
	svc := newGenreSvc(repo)

	genres := []models.Genre{{Name: "Action"}, {Name: "Comedy"}}
	repo.On("GetByIDs", ctx, []uint{1, 2}).Return(genres, nil)

	result, err := svc.GetGenresByIDs(ctx, []uint{1, 2})
	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestGenreService_GetGenresByIDs_Error(t *testing.T) {
	ctx := context.Background()
	repo := &mockGenreRepo{}
	svc := newGenreSvc(repo)

	repo.On("GetByIDs", ctx, []uint{1}).Return(nil, errors.New("db error"))

	_, err := svc.GetGenresByIDs(ctx, []uint{1})
	require.Error(t, err)
}

func TestGenreService_SearchGenres_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockGenreRepo{}
	svc := newGenreSvc(repo)

	searchResult := repositories.SearchResult{
		Genres: []models.Genre{{Name: "Action"}},
		Total:  1,
	}
	repo.On("Search", ctx, "act", 1, 10).Return(searchResult)

	result, total, err := svc.SearchGenres(ctx, "act", 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, result, 1)
}

func TestGenreService_SearchGenres_Error(t *testing.T) {
	ctx := context.Background()
	repo := &mockGenreRepo{}
	svc := newGenreSvc(repo)

	searchResult := repositories.SearchResult{Error: errors.New("search failed")}
	repo.On("Search", ctx, "act", 1, 10).Return(searchResult)

	_, _, err := svc.SearchGenres(ctx, "act", 1, 10)
	require.Error(t, err)
}

func TestGenreService_BulkCreateGenres_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockGenreRepo{}
	svc := newGenreSvc(repo)

	genres := []models.Genre{{Name: "A"}, {Name: "B"}}
	repo.On("BulkCreate", ctx, genres).Return(nil)

	err := svc.BulkCreateGenres(ctx, genres)
	require.NoError(t, err)
}

func TestGenreService_BulkCreateGenres_Error(t *testing.T) {
	ctx := context.Background()
	repo := &mockGenreRepo{}
	svc := newGenreSvc(repo)

	genres := []models.Genre{{Name: "A"}}
	repo.On("BulkCreate", ctx, genres).Return(errors.New("db error"))

	err := svc.BulkCreateGenres(ctx, genres)
	require.Error(t, err)
}

func TestGenreService_BulkDeleteGenres_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockGenreRepo{}
	svc := newGenreSvc(repo)

	ids := []uint{1, 2}
	repo.On("BulkDelete", ctx, ids).Return(nil)

	err := svc.BulkDeleteGenres(ctx, ids)
	require.NoError(t, err)
}

func TestGenreService_BulkDeleteGenres_Error(t *testing.T) {
	ctx := context.Background()
	repo := &mockGenreRepo{}
	svc := newGenreSvc(repo)

	ids := []uint{1}
	repo.On("BulkDelete", ctx, ids).Return(errors.New("db error"))

	err := svc.BulkDeleteGenres(ctx, ids)
	require.Error(t, err)
}

func TestGenreService_UpdateGenre_UpdateError(t *testing.T) {
	ctx := context.Background()
	repo := &mockGenreRepo{}
	svc := newGenreSvc(repo)

	existing := &models.Genre{Name: "Action"}
	existing.ID = 1
	repo.On("GetByID", ctx, uint(1)).Return(existing, nil)
	repo.On("Update", ctx, mock.Anything).Return(errors.New("db error"))

	genre := &models.Genre{Name: "Action"}
	genre.ID = 1
	err := svc.UpdateGenre(ctx, genre)
	require.Error(t, err)
}

func TestGenreService_DeleteGenre_DeleteError(t *testing.T) {
	ctx := context.Background()
	repo := &mockGenreRepo{}
	svc := newGenreSvc(repo)

	existing := &models.Genre{Name: "Action"}
	existing.ID = 1
	repo.On("GetByID", ctx, uint(1)).Return(existing, nil)
	repo.On("Delete", ctx, uint(1)).Return(errors.New("db error"))

	err := svc.DeleteGenre(ctx, 1)
	require.Error(t, err)
}

func TestNewGenreService(t *testing.T) {
	svc := NewGenreService(&mockGenreRepo{}, cache.NewNoopCache())
	require.NotNil(t, svc)
}

func TestGenreService_GetGenreByID_CacheHit(t *testing.T) {
	ctx := context.Background()
	repo := &mockGenreRepo{}
	mc := &mockCache{}

	data := []byte(`{"ID":1,"Name":"Action","CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","Animes":null}`)
	mc.On("Get", ctx, cache.KeyGenre(uint(1))).Return(data, nil)

	svc := &GenreService{genreRepo: repo, logger: newTestLogger(), cache: mc}
	result, err := svc.GetGenreByID(ctx, 1)
	require.NoError(t, err)
	assert.NotNil(t, result)
	mc.AssertExpectations(t)
	repo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
}
