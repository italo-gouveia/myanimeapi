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
)

func newAnimeSvc(
	animeRepo *mockAnimeRepo,
	genreRepo *mockGenreRepo,
	tagRepo *mockTagRepo,
	reviewRepo *mockReviewRepo,
) *AnimeService {
	return &AnimeService{
		animeRepo:  animeRepo,
		genreRepo:  genreRepo,
		tagRepo:    tagRepo,
		reviewRepo: reviewRepo,
		logger:     newTestLogger(),
		cache:      cache.NewNoopCache(),
	}
}

func TestAnimeService_GetAnimeByID_Success(t *testing.T) {
	ctx := context.Background()
	animeRepo := &mockAnimeRepo{}
	svc := newAnimeSvc(animeRepo, &mockGenreRepo{}, &mockTagRepo{}, &mockReviewRepo{})

	anime := &models.Anime{Title: "Naruto"}
	anime.ID = 1
	animeRepo.On("GetByID", ctx, uint(1)).Return(anime, nil)

	result, err := svc.GetAnimeByID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, "Naruto", result.Title)
}

func TestAnimeService_GetAnimeByID_NotFound(t *testing.T) {
	ctx := context.Background()
	animeRepo := &mockAnimeRepo{}
	svc := newAnimeSvc(animeRepo, &mockGenreRepo{}, &mockTagRepo{}, &mockReviewRepo{})

	animeRepo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("not found"))

	_, err := svc.GetAnimeByID(ctx, 1)
	require.Error(t, err)
}

func TestAnimeService_GetAllAnimes_Success(t *testing.T) {
	ctx := context.Background()
	animeRepo := &mockAnimeRepo{}
	svc := newAnimeSvc(animeRepo, &mockGenreRepo{}, &mockTagRepo{}, &mockReviewRepo{})

	a1 := &models.Anime{Title: "Naruto"}
	a1.ID = 1
	a2 := &models.Anime{Title: "Bleach"}
	a2.ID = 2
	filter := models.AnimeFilter{}
	animeRepo.On("GetWithFilters", ctx, filter, 1, 10).Return([]interface{}{a1, a2}, int64(2), nil)

	result, total, err := svc.GetAllAnimes(ctx, 1, 10, filter)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, result, 2)
}

func TestAnimeService_GetAllAnimes_Error(t *testing.T) {
	ctx := context.Background()
	animeRepo := &mockAnimeRepo{}
	svc := newAnimeSvc(animeRepo, &mockGenreRepo{}, &mockTagRepo{}, &mockReviewRepo{})

	filter := models.AnimeFilter{}
	animeRepo.On("GetWithFilters", ctx, filter, 1, 10).Return(nil, int64(0), errors.New("db error"))

	_, _, err := svc.GetAllAnimes(ctx, 1, 10, filter)
	require.Error(t, err)
}

func TestAnimeService_CreateAnime_Success(t *testing.T) {
	ctx := context.Background()
	animeRepo := &mockAnimeRepo{}
	svc := newAnimeSvc(animeRepo, &mockGenreRepo{}, &mockTagRepo{}, &mockReviewRepo{})

	animeRepo.On("GetByTitle", ctx, "Naruto", 1, 1).Return([]models.Anime{}, int64(0), nil)
	animeRepo.On("Create", ctx, mock.Anything).Return(nil)

	anime := &models.Anime{Title: "Naruto"}
	err := svc.CreateAnime(ctx, anime)
	require.NoError(t, err)
}

func TestAnimeService_CreateAnime_TitleConflict(t *testing.T) {
	ctx := context.Background()
	animeRepo := &mockAnimeRepo{}
	svc := newAnimeSvc(animeRepo, &mockGenreRepo{}, &mockTagRepo{}, &mockReviewRepo{})

	existing := models.Anime{Title: "Naruto"}
	existing.ID = 1
	animeRepo.On("GetByTitle", ctx, "Naruto", 1, 1).Return([]models.Anime{existing}, int64(1), nil)

	anime := &models.Anime{Title: "Naruto"}
	err := svc.CreateAnime(ctx, anime)
	require.Error(t, err)
}

func TestAnimeService_CreateAnime_CreateError(t *testing.T) {
	ctx := context.Background()
	animeRepo := &mockAnimeRepo{}
	svc := newAnimeSvc(animeRepo, &mockGenreRepo{}, &mockTagRepo{}, &mockReviewRepo{})

	animeRepo.On("GetByTitle", ctx, "Naruto", 1, 1).Return([]models.Anime{}, int64(0), nil)
	animeRepo.On("Create", ctx, mock.Anything).Return(errors.New("db error"))

	anime := &models.Anime{Title: "Naruto"}
	err := svc.CreateAnime(ctx, anime)
	require.Error(t, err)
}

func TestAnimeService_UpdateAnime_Success(t *testing.T) {
	ctx := context.Background()
	animeRepo := &mockAnimeRepo{}
	svc := newAnimeSvc(animeRepo, &mockGenreRepo{}, &mockTagRepo{}, &mockReviewRepo{})

	existing := &models.Anime{Title: "Naruto"}
	existing.ID = 1
	animeRepo.On("GetByID", ctx, uint(1)).Return(existing, nil)
	animeRepo.On("Update", ctx, mock.Anything).Return(nil)

	anime := &models.Anime{Title: "Naruto"}
	anime.ID = 1
	err := svc.UpdateAnime(ctx, anime)
	require.NoError(t, err)
}

func TestAnimeService_UpdateAnime_NotFound(t *testing.T) {
	ctx := context.Background()
	animeRepo := &mockAnimeRepo{}
	svc := newAnimeSvc(animeRepo, &mockGenreRepo{}, &mockTagRepo{}, &mockReviewRepo{})

	animeRepo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("not found"))

	anime := &models.Anime{Title: "Naruto"}
	anime.ID = 1
	err := svc.UpdateAnime(ctx, anime)
	require.Error(t, err)
}

func TestAnimeService_DeleteAnime_Success(t *testing.T) {
	ctx := context.Background()
	animeRepo := &mockAnimeRepo{}
	svc := newAnimeSvc(animeRepo, &mockGenreRepo{}, &mockTagRepo{}, &mockReviewRepo{})

	existing := &models.Anime{Title: "Naruto"}
	existing.ID = 1
	animeRepo.On("GetByID", ctx, uint(1)).Return(existing, nil)
	animeRepo.On("Delete", ctx, uint(1)).Return(nil)

	err := svc.DeleteAnime(ctx, 1)
	require.NoError(t, err)
}

func TestAnimeService_DeleteAnime_NotFound(t *testing.T) {
	ctx := context.Background()
	animeRepo := &mockAnimeRepo{}
	svc := newAnimeSvc(animeRepo, &mockGenreRepo{}, &mockTagRepo{}, &mockReviewRepo{})

	animeRepo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("not found"))

	err := svc.DeleteAnime(ctx, 1)
	require.Error(t, err)
}

func TestAnimeService_GetAnimesByTitle_Success(t *testing.T) {
	ctx := context.Background()
	animeRepo := &mockAnimeRepo{}
	svc := newAnimeSvc(animeRepo, &mockGenreRepo{}, &mockTagRepo{}, &mockReviewRepo{})

	animes := []models.Anime{{Title: "Naruto"}}
	animeRepo.On("GetByTitle", ctx, "Naru", 1, 10).Return(animes, int64(1), nil)

	result, total, err := svc.GetAnimesByTitle(ctx, "Naru", 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, result, 1)
}

func TestAnimeService_GetAnimesByTitle_Error(t *testing.T) {
	ctx := context.Background()
	animeRepo := &mockAnimeRepo{}
	svc := newAnimeSvc(animeRepo, &mockGenreRepo{}, &mockTagRepo{}, &mockReviewRepo{})

	animeRepo.On("GetByTitle", ctx, "Naru", 1, 10).Return(nil, int64(0), errors.New("db error"))

	_, _, err := svc.GetAnimesByTitle(ctx, "Naru", 1, 10)
	require.Error(t, err)
}

func TestAnimeService_GetAnimesByGenre_Success(t *testing.T) {
	ctx := context.Background()
	animeRepo := &mockAnimeRepo{}
	svc := newAnimeSvc(animeRepo, &mockGenreRepo{}, &mockTagRepo{}, &mockReviewRepo{})

	animes := []models.Anime{{Title: "Naruto"}}
	animeRepo.On("GetByGenre", ctx, "Action", 1, 10).Return(animes, int64(1), nil)

	result, total, err := svc.GetAnimesByGenre(ctx, "Action", 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, result, 1)
}

func TestAnimeService_GetReviewsForAnime_Success(t *testing.T) {
	ctx := context.Background()
	animeRepo := &mockAnimeRepo{}
	reviewRepo := &mockReviewRepo{}
	svc := newAnimeSvc(animeRepo, &mockGenreRepo{}, &mockTagRepo{}, reviewRepo)

	anime := &models.Anime{Title: "Naruto"}
	anime.ID = 1
	reviews := []models.Review{{UserID: 1, AnimeID: 1, Rating: 8}}

	animeRepo.On("GetByID", ctx, uint(1)).Return(anime, nil)
	reviewRepo.On("GetByAnimeID", ctx, uint(1), 1, 10).Return(reviews, int64(1), nil)

	result, total, err := svc.GetReviewsForAnime(ctx, 1, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, result, 1)
}

func TestAnimeService_GetReviewsForAnime_AnimeNotFound(t *testing.T) {
	ctx := context.Background()
	animeRepo := &mockAnimeRepo{}
	svc := newAnimeSvc(animeRepo, &mockGenreRepo{}, &mockTagRepo{}, &mockReviewRepo{})

	animeRepo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("not found"))

	_, _, err := svc.GetReviewsForAnime(ctx, 1, 1, 10)
	require.Error(t, err)
}

func TestAnimeService_AddGenresToAnime_Success(t *testing.T) {
	ctx := context.Background()
	animeRepo := &mockAnimeRepo{}
	genreRepo := &mockGenreRepo{}
	svc := newAnimeSvc(animeRepo, genreRepo, &mockTagRepo{}, &mockReviewRepo{})

	anime := &models.Anime{Title: "Naruto"}
	anime.ID = 1
	genres := []models.Genre{{Name: "Action"}}

	animeRepo.On("GetByID", ctx, uint(1)).Return(anime, nil)
	genreRepo.On("GetByIDs", ctx, []uint{2}).Return(genres, nil)
	animeRepo.On("Update", ctx, mock.Anything).Return(nil)

	err := svc.AddGenresToAnime(ctx, 1, []uint{2})
	require.NoError(t, err)
}

func TestAnimeService_AddGenresToAnime_AnimeNotFound(t *testing.T) {
	ctx := context.Background()
	animeRepo := &mockAnimeRepo{}
	svc := newAnimeSvc(animeRepo, &mockGenreRepo{}, &mockTagRepo{}, &mockReviewRepo{})

	animeRepo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("not found"))

	err := svc.AddGenresToAnime(ctx, 1, []uint{2})
	require.Error(t, err)
}

func TestAnimeService_RemoveGenresFromAnime_Success(t *testing.T) {
	ctx := context.Background()
	animeRepo := &mockAnimeRepo{}
	svc := newAnimeSvc(animeRepo, &mockGenreRepo{}, &mockTagRepo{}, &mockReviewRepo{})

	actionGenre := models.Genre{Name: "Action"}
	actionGenre.ID = 2
	anime := &models.Anime{Title: "Naruto", Genres: []models.Genre{actionGenre}}
	anime.ID = 1

	animeRepo.On("GetByID", ctx, uint(1)).Return(anime, nil)
	animeRepo.On("Update", ctx, mock.Anything).Return(nil)

	err := svc.RemoveGenresFromAnime(ctx, 1, []uint{2})
	require.NoError(t, err)
}
