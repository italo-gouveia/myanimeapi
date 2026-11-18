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

func newReviewSvc(
	reviewRepo *mockReviewRepo,
	userRepo *mockUserRepo,
	animeRepo *mockAnimeRepo,
	storageSvc *mockStorageSvc,
) *ReviewService {
	return &ReviewService{
		reviewRepo: reviewRepo,
		userRepo:   userRepo,
		animeRepo:  animeRepo,
		storageSvc: storageSvc,
		cache:      cache.NewNoopCache(),
		logger:     newTestLogger(),
	}
}

func TestReviewService_GetReviewByID_Success(t *testing.T) {
	ctx := context.Background()
	reviewRepo := &mockReviewRepo{}
	svc := newReviewSvc(reviewRepo, &mockUserRepo{}, &mockAnimeRepo{}, &mockStorageSvc{})

	review := &models.Review{UserID: 1, AnimeID: 10, Rating: 8}
	review.ID = 1
	reviewRepo.On("GetByID", ctx, uint(1)).Return(review, nil)

	result, err := svc.GetReviewByID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, uint(1), result.UserID)
}

func TestReviewService_GetReviewByID_NotFound(t *testing.T) {
	ctx := context.Background()
	reviewRepo := &mockReviewRepo{}
	svc := newReviewSvc(reviewRepo, &mockUserRepo{}, &mockAnimeRepo{}, &mockStorageSvc{})

	reviewRepo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("not found"))

	_, err := svc.GetReviewByID(ctx, 1)
	require.Error(t, err)
}

func TestReviewService_GetAllReviews_Success(t *testing.T) {
	ctx := context.Background()
	reviewRepo := &mockReviewRepo{}
	svc := newReviewSvc(reviewRepo, &mockUserRepo{}, &mockAnimeRepo{}, &mockStorageSvc{})

	r1 := models.Review{UserID: 1, Rating: 8}
	r1.ID = 1
	r2 := models.Review{UserID: 2, Rating: 9}
	r2.ID = 2
	reviewRepo.On("GetAll", ctx, 1, 10).Return([]interface{}{r1, r2}, int64(2), nil)

	result, total, err := svc.GetAllReviews(ctx, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, result, 2)
}

func TestReviewService_GetAllReviews_Error(t *testing.T) {
	ctx := context.Background()
	reviewRepo := &mockReviewRepo{}
	svc := newReviewSvc(reviewRepo, &mockUserRepo{}, &mockAnimeRepo{}, &mockStorageSvc{})

	reviewRepo.On("GetAll", ctx, 1, 10).Return(nil, int64(0), errors.New("db error"))

	_, _, err := svc.GetAllReviews(ctx, 1, 10)
	require.Error(t, err)
}

func TestReviewService_CreateReview_Success(t *testing.T) {
	ctx := context.Background()
	reviewRepo := &mockReviewRepo{}
	userRepo := &mockUserRepo{}
	animeRepo := &mockAnimeRepo{}
	svc := newReviewSvc(reviewRepo, userRepo, animeRepo, &mockStorageSvc{})

	user := &models.User{Username: "testuser"}
	user.ID = 1
	anime := &models.Anime{Title: "Naruto"}
	anime.ID = 10

	userRepo.On("GetByID", ctx, uint(1)).Return(user, nil)
	animeRepo.On("GetByID", ctx, uint(10)).Return(anime, nil)
	reviewRepo.On("GetByUserAndAnime", ctx, uint(1), uint(10)).Return(nil, errors.New("not found"))
	reviewRepo.On("Create", ctx, mock.Anything).Return(nil)

	review := &models.Review{UserID: 1, AnimeID: 10, Rating: 8}
	err := svc.CreateReview(ctx, review)
	require.NoError(t, err)
}

func TestReviewService_CreateReview_UserNotFound(t *testing.T) {
	ctx := context.Background()
	userRepo := &mockUserRepo{}
	svc := newReviewSvc(&mockReviewRepo{}, userRepo, &mockAnimeRepo{}, &mockStorageSvc{})

	userRepo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("not found"))

	review := &models.Review{UserID: 1, AnimeID: 10, Rating: 8}
	err := svc.CreateReview(ctx, review)
	require.Error(t, err)
}

func TestReviewService_CreateReview_AnimeNotFound(t *testing.T) {
	ctx := context.Background()
	userRepo := &mockUserRepo{}
	animeRepo := &mockAnimeRepo{}
	svc := newReviewSvc(&mockReviewRepo{}, userRepo, animeRepo, &mockStorageSvc{})

	user := &models.User{Username: "testuser"}
	user.ID = 1
	userRepo.On("GetByID", ctx, uint(1)).Return(user, nil)
	animeRepo.On("GetByID", ctx, uint(10)).Return(nil, errors.New("not found"))

	review := &models.Review{UserID: 1, AnimeID: 10, Rating: 8}
	err := svc.CreateReview(ctx, review)
	require.Error(t, err)
}

func TestReviewService_CreateReview_AlreadyReviewed(t *testing.T) {
	ctx := context.Background()
	reviewRepo := &mockReviewRepo{}
	userRepo := &mockUserRepo{}
	animeRepo := &mockAnimeRepo{}
	svc := newReviewSvc(reviewRepo, userRepo, animeRepo, &mockStorageSvc{})

	user := &models.User{Username: "testuser"}
	user.ID = 1
	anime := &models.Anime{Title: "Naruto"}
	anime.ID = 10
	existing := &models.Review{UserID: 1, AnimeID: 10}

	userRepo.On("GetByID", ctx, uint(1)).Return(user, nil)
	animeRepo.On("GetByID", ctx, uint(10)).Return(anime, nil)
	reviewRepo.On("GetByUserAndAnime", ctx, uint(1), uint(10)).Return(existing, nil)

	review := &models.Review{UserID: 1, AnimeID: 10, Rating: 8}
	err := svc.CreateReview(ctx, review)
	require.Error(t, err)
}

func TestReviewService_CreateReview_InvalidRating(t *testing.T) {
	ctx := context.Background()
	reviewRepo := &mockReviewRepo{}
	userRepo := &mockUserRepo{}
	animeRepo := &mockAnimeRepo{}
	svc := newReviewSvc(reviewRepo, userRepo, animeRepo, &mockStorageSvc{})

	user := &models.User{Username: "testuser"}
	user.ID = 1
	anime := &models.Anime{Title: "Naruto"}
	anime.ID = 10

	userRepo.On("GetByID", ctx, uint(1)).Return(user, nil)
	animeRepo.On("GetByID", ctx, uint(10)).Return(anime, nil)
	reviewRepo.On("GetByUserAndAnime", ctx, uint(1), uint(10)).Return(nil, errors.New("not found"))

	review := &models.Review{UserID: 1, AnimeID: 10, Rating: 11}
	err := svc.CreateReview(ctx, review)
	require.Error(t, err)
}

func TestReviewService_UpdateReview_Success(t *testing.T) {
	ctx := context.Background()
	reviewRepo := &mockReviewRepo{}
	svc := newReviewSvc(reviewRepo, &mockUserRepo{}, &mockAnimeRepo{}, &mockStorageSvc{})

	existing := &models.Review{UserID: 1, AnimeID: 10, Rating: 8}
	existing.ID = 1
	reviewRepo.On("GetByID", ctx, uint(1)).Return(existing, nil)
	reviewRepo.On("Update", ctx, mock.Anything).Return(nil)

	review := &models.Review{UserID: 1, AnimeID: 10, Rating: 9}
	review.ID = 1
	err := svc.UpdateReview(ctx, review)
	require.NoError(t, err)
}

func TestReviewService_UpdateReview_NotFound(t *testing.T) {
	ctx := context.Background()
	reviewRepo := &mockReviewRepo{}
	svc := newReviewSvc(reviewRepo, &mockUserRepo{}, &mockAnimeRepo{}, &mockStorageSvc{})

	reviewRepo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("not found"))

	review := &models.Review{UserID: 1, AnimeID: 10, Rating: 9}
	review.ID = 1
	err := svc.UpdateReview(ctx, review)
	require.Error(t, err)
}

func TestReviewService_UpdateReview_InvalidRating(t *testing.T) {
	ctx := context.Background()
	reviewRepo := &mockReviewRepo{}
	svc := newReviewSvc(reviewRepo, &mockUserRepo{}, &mockAnimeRepo{}, &mockStorageSvc{})

	existing := &models.Review{UserID: 1, AnimeID: 10, Rating: 8}
	existing.ID = 1
	reviewRepo.On("GetByID", ctx, uint(1)).Return(existing, nil)

	review := &models.Review{UserID: 1, AnimeID: 10, Rating: 0}
	review.ID = 1
	err := svc.UpdateReview(ctx, review)
	require.Error(t, err)
}

func TestReviewService_DeleteReview_Success(t *testing.T) {
	ctx := context.Background()
	reviewRepo := &mockReviewRepo{}
	svc := newReviewSvc(reviewRepo, &mockUserRepo{}, &mockAnimeRepo{}, &mockStorageSvc{})

	review := &models.Review{UserID: 1, AnimeID: 10, Rating: 8}
	review.ID = 1
	reviewRepo.On("GetByID", ctx, uint(1)).Return(review, nil)
	reviewRepo.On("Delete", ctx, uint(1)).Return(nil)

	err := svc.DeleteReview(ctx, 1)
	require.NoError(t, err)
}

func TestReviewService_DeleteReview_NotFound(t *testing.T) {
	ctx := context.Background()
	reviewRepo := &mockReviewRepo{}
	svc := newReviewSvc(reviewRepo, &mockUserRepo{}, &mockAnimeRepo{}, &mockStorageSvc{})

	reviewRepo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("not found"))

	err := svc.DeleteReview(ctx, 1)
	require.Error(t, err)
}

func TestReviewService_GetReviewsByUserID_Success(t *testing.T) {
	ctx := context.Background()
	reviewRepo := &mockReviewRepo{}
	userRepo := &mockUserRepo{}
	svc := newReviewSvc(reviewRepo, userRepo, &mockAnimeRepo{}, &mockStorageSvc{})

	user := &models.User{Username: "testuser"}
	user.ID = 1
	reviews := []models.Review{{UserID: 1, Rating: 8}}

	userRepo.On("GetByID", ctx, uint(1)).Return(user, nil)
	reviewRepo.On("GetByUserID", ctx, uint(1), 1, 10).Return(reviews, int64(1), nil)

	result, total, err := svc.GetReviewsByUserID(ctx, 1, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, result, 1)
}

func TestReviewService_GetReviewsByUserID_UserNotFound(t *testing.T) {
	ctx := context.Background()
	userRepo := &mockUserRepo{}
	svc := newReviewSvc(&mockReviewRepo{}, userRepo, &mockAnimeRepo{}, &mockStorageSvc{})

	userRepo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("not found"))

	_, _, err := svc.GetReviewsByUserID(ctx, 1, 1, 10)
	require.Error(t, err)
}

func TestReviewService_GetReviewsByAnimeID_Success(t *testing.T) {
	ctx := context.Background()
	reviewRepo := &mockReviewRepo{}
	animeRepo := &mockAnimeRepo{}
	svc := newReviewSvc(reviewRepo, &mockUserRepo{}, animeRepo, &mockStorageSvc{})

	anime := &models.Anime{Title: "Naruto"}
	anime.ID = 10
	reviews := []models.Review{{UserID: 1, AnimeID: 10, Rating: 8}}

	animeRepo.On("GetByID", ctx, uint(10)).Return(anime, nil)
	reviewRepo.On("GetByAnimeID", ctx, uint(10), 1, 10).Return(reviews, int64(1), nil)

	result, total, err := svc.GetReviewsByAnimeID(ctx, 10, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, result, 1)
}

func TestReviewService_GetReviewsByAnimeID_AnimeNotFound(t *testing.T) {
	ctx := context.Background()
	animeRepo := &mockAnimeRepo{}
	svc := newReviewSvc(&mockReviewRepo{}, &mockUserRepo{}, animeRepo, &mockStorageSvc{})

	animeRepo.On("GetByID", ctx, uint(10)).Return(nil, errors.New("not found"))

	_, _, err := svc.GetReviewsByAnimeID(ctx, 10, 1, 10)
	require.Error(t, err)
}

func TestNewReviewService(t *testing.T) {
	svc := NewReviewService(&mockReviewRepo{}, &mockUserRepo{}, &mockAnimeRepo{}, &mockStorageSvc{}, cache.NewNoopCache())
	require.NotNil(t, svc)
}
