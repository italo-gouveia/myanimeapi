package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"myanimeapi/api/models"
)

func newWatchlistSvc(repo *mockWatchlistRepo) *WatchlistService {
	return &WatchlistService{
		watchlistRepo: repo,
		logger:        newTestLogger(),
	}
}

func TestWatchlistService_UpsertWatchlistEntry_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockWatchlistRepo{}
	svc := newWatchlistSvc(repo)

	entry := &models.WatchlistEntry{UserID: 1, AnimeID: 10, Status: models.WatchlistStatusWatching}
	repo.On("UpsertByUserAndAnime", ctx, uint(1), uint(10), models.WatchlistStatusWatching).Return(entry, nil)

	result, err := svc.UpsertWatchlistEntry(ctx, 1, 10, models.WatchlistStatusWatching)
	require.NoError(t, err)
	assert.Equal(t, uint(1), result.UserID)
	assert.Equal(t, models.WatchlistStatusWatching, result.Status)
}

func TestWatchlistService_UpsertWatchlistEntry_Error(t *testing.T) {
	ctx := context.Background()
	repo := &mockWatchlistRepo{}
	svc := newWatchlistSvc(repo)

	repo.On("UpsertByUserAndAnime", ctx, uint(1), uint(10), models.WatchlistStatusWatching).Return(nil, errors.New("db error"))

	_, err := svc.UpsertWatchlistEntry(ctx, 1, 10, models.WatchlistStatusWatching)
	require.Error(t, err)
}

func TestWatchlistService_RemoveWatchlistEntry_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockWatchlistRepo{}
	svc := newWatchlistSvc(repo)

	entry := &models.WatchlistEntry{UserID: 1, AnimeID: 10}
	repo.On("GetByUserAndAnime", ctx, uint(1), uint(10)).Return(entry, nil)
	repo.On("DeleteByUserAndAnime", ctx, uint(1), uint(10)).Return(nil)

	err := svc.RemoveWatchlistEntry(ctx, 1, 10)
	require.NoError(t, err)
}

func TestWatchlistService_RemoveWatchlistEntry_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockWatchlistRepo{}
	svc := newWatchlistSvc(repo)

	repo.On("GetByUserAndAnime", ctx, uint(1), uint(10)).Return(nil, errors.New("not found"))

	err := svc.RemoveWatchlistEntry(ctx, 1, 10)
	require.Error(t, err)
}

func TestWatchlistService_RemoveWatchlistEntry_DeleteError(t *testing.T) {
	ctx := context.Background()
	repo := &mockWatchlistRepo{}
	svc := newWatchlistSvc(repo)

	entry := &models.WatchlistEntry{UserID: 1, AnimeID: 10}
	repo.On("GetByUserAndAnime", ctx, uint(1), uint(10)).Return(entry, nil)
	repo.On("DeleteByUserAndAnime", ctx, uint(1), uint(10)).Return(errors.New("db error"))

	err := svc.RemoveWatchlistEntry(ctx, 1, 10)
	require.Error(t, err)
}

func TestWatchlistService_GetWatchlist_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockWatchlistRepo{}
	svc := newWatchlistSvc(repo)

	entries := []models.WatchlistEntry{
		{UserID: 1, AnimeID: 10, Status: models.WatchlistStatusWatching},
		{UserID: 1, AnimeID: 20, Status: models.WatchlistStatusCompleted},
	}
	repo.On("GetByUserID", ctx, uint(1)).Return(entries, nil)

	result, err := svc.GetWatchlist(ctx, 1)
	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestWatchlistService_GetWatchlist_Error(t *testing.T) {
	ctx := context.Background()
	repo := &mockWatchlistRepo{}
	svc := newWatchlistSvc(repo)

	repo.On("GetByUserID", ctx, uint(1)).Return(nil, errors.New("db error"))

	_, err := svc.GetWatchlist(ctx, 1)
	require.Error(t, err)
}

func TestWatchlistService_GetWatchlistEntry_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockWatchlistRepo{}
	svc := newWatchlistSvc(repo)

	entry := &models.WatchlistEntry{UserID: 1, AnimeID: 10, Status: models.WatchlistStatusPlanToWatch}
	repo.On("GetByUserAndAnime", ctx, uint(1), uint(10)).Return(entry, nil)

	result, err := svc.GetWatchlistEntry(ctx, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, models.WatchlistStatusPlanToWatch, result.Status)
}

func TestWatchlistService_GetWatchlistEntry_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockWatchlistRepo{}
	svc := newWatchlistSvc(repo)

	repo.On("GetByUserAndAnime", ctx, uint(1), uint(10)).Return(nil, errors.New("not found"))

	_, err := svc.GetWatchlistEntry(ctx, 1, 10)
	require.Error(t, err)
}

func TestNewWatchlistService(t *testing.T) {
	repo := &mockWatchlistRepo{}
	svc := NewWatchlistService(repo)
	require.NotNil(t, svc)
}
