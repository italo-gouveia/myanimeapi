package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"myanimeapi/api/models"
)

func newFavoriteSvc(favRepo *mockFavoriteRepo, userRepo *mockUserRepo, animeRepo *mockAnimeRepo) *FavoriteService {
	return &FavoriteService{
		favoriteRepo: favRepo,
		userRepo:     userRepo,
		animeRepo:    animeRepo,
		logger:       newTestLogger(),
	}
}

func TestFavoriteService_AddFavorite_Success(t *testing.T) {
	ctx := context.Background()
	favRepo := &mockFavoriteRepo{}
	userRepo := &mockUserRepo{}
	animeRepo := &mockAnimeRepo{}
	svc := newFavoriteSvc(favRepo, userRepo, animeRepo)

	user := &models.User{Username: "testuser"}
	user.ID = 1
	anime := &models.Anime{Title: "Naruto"}
	anime.ID = 10
	favorite := &models.Favorite{UserID: 1, AnimeID: 10}
	favorite.ID = 100

	userRepo.On("GetByID", ctx, uint(1)).Return(user, nil)
	animeRepo.On("GetByID", ctx, uint(10)).Return(anime, nil)
	favRepo.On("GetByUserAndAnime", ctx, uint(1), uint(10)).Return(nil, errors.New("not found"))
	favRepo.On("Create", ctx, mock.Anything).Return(nil)
	favRepo.On("GetByID", ctx, uint(0)).Return(favorite, nil)

	result, err := svc.AddFavorite(ctx, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, uint(1), result.UserID)
	assert.Equal(t, uint(10), result.AnimeID)
}

func TestFavoriteService_AddFavorite_UserNotFound(t *testing.T) {
	ctx := context.Background()
	svc := newFavoriteSvc(&mockFavoriteRepo{}, &mockUserRepo{}, &mockAnimeRepo{})

	svc.userRepo.(*mockUserRepo).On("GetByID", ctx, uint(1)).Return(nil, errors.New("user not found"))

	_, err := svc.AddFavorite(ctx, 1, 10)
	require.Error(t, err)
}

func TestFavoriteService_AddFavorite_AnimeNotFound(t *testing.T) {
	ctx := context.Background()
	favRepo := &mockFavoriteRepo{}
	userRepo := &mockUserRepo{}
	animeRepo := &mockAnimeRepo{}
	svc := newFavoriteSvc(favRepo, userRepo, animeRepo)

	user := &models.User{Username: "testuser"}
	user.ID = 1
	userRepo.On("GetByID", ctx, uint(1)).Return(user, nil)
	animeRepo.On("GetByID", ctx, uint(10)).Return(nil, errors.New("anime not found"))

	_, err := svc.AddFavorite(ctx, 1, 10)
	require.Error(t, err)
}

func TestFavoriteService_AddFavorite_AlreadyExists(t *testing.T) {
	ctx := context.Background()
	favRepo := &mockFavoriteRepo{}
	userRepo := &mockUserRepo{}
	animeRepo := &mockAnimeRepo{}
	svc := newFavoriteSvc(favRepo, userRepo, animeRepo)

	user := &models.User{Username: "testuser"}
	user.ID = 1
	anime := &models.Anime{Title: "Naruto"}
	anime.ID = 10
	existing := &models.Favorite{UserID: 1, AnimeID: 10}

	userRepo.On("GetByID", ctx, uint(1)).Return(user, nil)
	animeRepo.On("GetByID", ctx, uint(10)).Return(anime, nil)
	favRepo.On("GetByUserAndAnime", ctx, uint(1), uint(10)).Return(existing, nil)

	_, err := svc.AddFavorite(ctx, 1, 10)
	require.Error(t, err)
}

func TestFavoriteService_AddFavorite_CreateError(t *testing.T) {
	ctx := context.Background()
	favRepo := &mockFavoriteRepo{}
	userRepo := &mockUserRepo{}
	animeRepo := &mockAnimeRepo{}
	svc := newFavoriteSvc(favRepo, userRepo, animeRepo)

	user := &models.User{Username: "testuser"}
	user.ID = 1
	anime := &models.Anime{Title: "Naruto"}
	anime.ID = 10

	userRepo.On("GetByID", ctx, uint(1)).Return(user, nil)
	animeRepo.On("GetByID", ctx, uint(10)).Return(anime, nil)
	favRepo.On("GetByUserAndAnime", ctx, uint(1), uint(10)).Return(nil, errors.New("not found"))
	favRepo.On("Create", ctx, mock.Anything).Return(errors.New("db error"))

	_, err := svc.AddFavorite(ctx, 1, 10)
	require.Error(t, err)
}

func TestFavoriteService_RemoveFavorite_Success(t *testing.T) {
	ctx := context.Background()
	favRepo := &mockFavoriteRepo{}
	svc := newFavoriteSvc(favRepo, &mockUserRepo{}, &mockAnimeRepo{})

	existing := &models.Favorite{UserID: 1, AnimeID: 10}
	favRepo.On("GetByUserAndAnime", ctx, uint(1), uint(10)).Return(existing, nil)
	favRepo.On("DeleteByUserAndAnime", ctx, uint(1), uint(10)).Return(nil)

	err := svc.RemoveFavorite(ctx, 1, 10)
	require.NoError(t, err)
}

func TestFavoriteService_RemoveFavorite_NotFound(t *testing.T) {
	ctx := context.Background()
	favRepo := &mockFavoriteRepo{}
	svc := newFavoriteSvc(favRepo, &mockUserRepo{}, &mockAnimeRepo{})

	favRepo.On("GetByUserAndAnime", ctx, uint(1), uint(10)).Return(nil, errors.New("not found"))

	err := svc.RemoveFavorite(ctx, 1, 10)
	require.Error(t, err)
}

func TestFavoriteService_RemoveFavorite_DeleteError(t *testing.T) {
	ctx := context.Background()
	favRepo := &mockFavoriteRepo{}
	svc := newFavoriteSvc(favRepo, &mockUserRepo{}, &mockAnimeRepo{})

	existing := &models.Favorite{UserID: 1, AnimeID: 10}
	favRepo.On("GetByUserAndAnime", ctx, uint(1), uint(10)).Return(existing, nil)
	favRepo.On("DeleteByUserAndAnime", ctx, uint(1), uint(10)).Return(errors.New("db error"))

	err := svc.RemoveFavorite(ctx, 1, 10)
	require.Error(t, err)
}

func TestFavoriteService_GetFavorites_Success(t *testing.T) {
	ctx := context.Background()
	favRepo := &mockFavoriteRepo{}
	userRepo := &mockUserRepo{}
	svc := newFavoriteSvc(favRepo, userRepo, &mockAnimeRepo{})

	user := &models.User{Username: "testuser"}
	user.ID = 1
	favorites := []models.Favorite{{UserID: 1, AnimeID: 10}}

	userRepo.On("GetByID", ctx, uint(1)).Return(user, nil)
	favRepo.On("GetByUserID", ctx, uint(1)).Return(favorites, nil)

	result, err := svc.GetFavorites(ctx, 1)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestFavoriteService_GetFavorites_UserNotFound(t *testing.T) {
	ctx := context.Background()
	userRepo := &mockUserRepo{}
	svc := newFavoriteSvc(&mockFavoriteRepo{}, userRepo, &mockAnimeRepo{})

	userRepo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("not found"))

	_, err := svc.GetFavorites(ctx, 1)
	require.Error(t, err)
}

func TestFavoriteService_GetFavorites_RepoError(t *testing.T) {
	ctx := context.Background()
	favRepo := &mockFavoriteRepo{}
	userRepo := &mockUserRepo{}
	svc := newFavoriteSvc(favRepo, userRepo, &mockAnimeRepo{})

	user := &models.User{Username: "testuser"}
	user.ID = 1
	userRepo.On("GetByID", ctx, uint(1)).Return(user, nil)
	favRepo.On("GetByUserID", ctx, uint(1)).Return(nil, errors.New("db error"))

	_, err := svc.GetFavorites(ctx, 1)
	require.Error(t, err)
}
