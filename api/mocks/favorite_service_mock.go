package mocks

import (
	"context"
	"myanimeapi/api/models"
	"reflect"

	"github.com/golang/mock/gomock"
)

// MockFavoriteServiceInterface is a mock of FavoriteServiceInterface interface
type MockFavoriteServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockFavoriteServiceInterfaceMockRecorder
}

// MockFavoriteServiceInterfaceMockRecorder is the mock recorder for MockFavoriteServiceInterface
type MockFavoriteServiceInterfaceMockRecorder struct {
	mock *MockFavoriteServiceInterface
}

// NewMockFavoriteServiceInterface creates a new mock instance
func NewMockFavoriteServiceInterface(ctrl *gomock.Controller) *MockFavoriteServiceInterface {
	mock := &MockFavoriteServiceInterface{ctrl: ctrl}
	mock.recorder = &MockFavoriteServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFavoriteServiceInterface) EXPECT() *MockFavoriteServiceInterfaceMockRecorder {
	return m.recorder
}

// AddFavorite mocks base method
func (m *MockFavoriteServiceInterface) AddFavorite(ctx context.Context, userID uint, animeID uint) (*models.Favorite, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddFavorite", ctx, userID, animeID)
	ret0, _ := ret[0].(*models.Favorite)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddFavorite indicates an expected call of AddFavorite
func (mr *MockFavoriteServiceInterfaceMockRecorder) AddFavorite(ctx, userID, animeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddFavorite", reflect.TypeOf((*MockFavoriteServiceInterface)(nil).AddFavorite), ctx, userID, animeID)
}

// RemoveFavorite mocks base method
func (m *MockFavoriteServiceInterface) RemoveFavorite(ctx context.Context, userID uint, animeID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveFavorite", ctx, userID, animeID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveFavorite indicates an expected call of RemoveFavorite
func (mr *MockFavoriteServiceInterfaceMockRecorder) RemoveFavorite(ctx, userID, animeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveFavorite", reflect.TypeOf((*MockFavoriteServiceInterface)(nil).RemoveFavorite), ctx, userID, animeID)
}

// GetFavorites mocks base method
func (m *MockFavoriteServiceInterface) GetFavorites(ctx context.Context, userID uint) ([]models.Favorite, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFavorites", ctx, userID)
	ret0, _ := ret[0].([]models.Favorite)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFavorites indicates an expected call of GetFavorites
func (mr *MockFavoriteServiceInterfaceMockRecorder) GetFavorites(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFavorites", reflect.TypeOf((*MockFavoriteServiceInterface)(nil).GetFavorites), ctx, userID)
}
