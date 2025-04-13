package repositories

import (
	"context"
	"myanimeapi/api/models"
	"reflect"

	"github.com/golang/mock/gomock"
)

// MockAnimeRepository is a mock implementation of AnimeRepository
type MockAnimeRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAnimeRepositoryMockRecorder
}

// MockAnimeRepositoryMockRecorder is the mock recorder for MockAnimeRepository
type MockAnimeRepositoryMockRecorder struct {
	mock *MockAnimeRepository
}

// GetReviewsForAnime provides a mock function with given fields: ctx, animeID, page, limit
func (m *MockAnimeRepository) GetReviewsForAnime(ctx context.Context, animeID uint, page, limit int) ([]models.Review, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReviewsForAnime", ctx, animeID, page, limit)
	ret0, _ := ret[0].([]models.Review)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetReviewsForAnime indicates an expected call of GetReviewsForAnime
func (mr *MockAnimeRepositoryMockRecorder) GetReviewsForAnime(ctx, animeID, page, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReviewsForAnime", reflect.TypeOf((*MockAnimeRepository)(nil).GetReviewsForAnime), ctx, animeID, page, limit)
}
