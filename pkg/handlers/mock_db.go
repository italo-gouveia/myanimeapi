// pkg/handlers/mock_db.go
package handlers

import (
	"context"

	"myanimeapi/pkg/models"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockDB is a mock implementation of DBInterface
type MockDB struct {
	mock.Mock
}

// WithContext mocks the WithContext method
func (m *MockDB) WithContext(ctx context.Context) *gorm.DB {
	args := m.Called(ctx)
	return args.Get(0).(*gorm.DB)
}

// First mocks the First method
func (m *MockDB) First(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(ctx, dest, conds) // args is used
	if destPtr, ok := dest.(**models.Anime); ok {
		*destPtr = &models.Anime{ID: 1, Title: "Naruto"}
	}
	return args.Get(0).(*gorm.DB) // Use args to return the mocked *gorm.DB
}

// Find mocks the Find method
func (m *MockDB) Find(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(ctx, dest, conds) // args is used
	if destPtr, ok := dest.(*[]models.Review); ok {
		*destPtr = []models.Review{
			{ID: 1, AnimeID: 1, Content: "Great anime!"},
			{ID: 2, AnimeID: 1, Content: "Awesome!"},
		}
	}
	return args.Get(0).(*gorm.DB) // Use args to return the mocked *gorm.DB
}

// Where mocks the Where method
func (m *MockDB) Where(ctx context.Context, query interface{}, args ...interface{}) *gorm.DB {
	mockArgs := m.Called(ctx, query, args) // mockArgs is used
	return mockArgs.Get(0).(*gorm.DB)
}

// Create mocks the Create method
func (m *MockDB) Create(ctx context.Context, value interface{}) *gorm.DB {
	args := m.Called(ctx, value) // args is used
	return args.Get(0).(*gorm.DB)
}

// Save mocks the Save method
func (m *MockDB) Save(ctx context.Context, value interface{}) *gorm.DB {
	args := m.Called(ctx, value) // args is used
	return args.Get(0).(*gorm.DB)
}

// Delete mocks the Delete method
func (m *MockDB) Delete(ctx context.Context, value interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(ctx, value, conds) // args is used
	return args.Get(0).(*gorm.DB)
}

// Preload mocks the Preload method
func (m *MockDB) Preload(column string, ctx context.Context, conditions ...interface{}) *gorm.DB {
	args := m.Called(column, ctx, conditions) // args is used
	return args.Get(0).(*gorm.DB)
}

// Offset mocks the Offset method
func (m *MockDB) Offset(offset int) *gorm.DB {
	args := m.Called(offset) // args is used
	return args.Get(0).(*gorm.DB)
}

// Limit mocks the Limit method
func (m *MockDB) Limit(limit int) *gorm.DB {
	args := m.Called(limit) // args is used
	return args.Get(0).(*gorm.DB)
}

// Unscoped mocks the Unscoped method
func (m *MockDB) Unscoped(ctx context.Context) *gorm.DB {
	args := m.Called(ctx) // args is used
	return args.Get(0).(*gorm.DB)
}

// Begin mocks the Begin method
func (m *MockDB) Begin(ctx context.Context) *gorm.DB {
	args := m.Called(ctx) // args is used
	return args.Get(0).(*gorm.DB)
}

// Commit mocks the Commit method
func (m *MockDB) Commit(ctx context.Context) *gorm.DB {
	args := m.Called(ctx) // args is used
	return args.Get(0).(*gorm.DB)
}

// Rollback mocks the Rollback method
func (m *MockDB) Rollback(ctx context.Context) *gorm.DB {
	args := m.Called(ctx) // args is used
	return args.Get(0).(*gorm.DB)
}

// GetError mocks the GetError method
func (m *MockDB) GetError() error {
	args := m.Called() // args is used
	return args.Error(0)
}
