// pkg/handlers/mock_db.go
package handlers

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockDB is a mock implementation of DBInterface
type MockDB struct {
	mock.Mock
}

// First is a mock method for querying the database
func (m *MockDB) First(out interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(out, where)
	if len(args) > 0 {
		return args.Get(0).(*gorm.DB)
	}
	return nil
}

func (m *MockDB) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(dest, conds)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	args = append([]interface{}{query}, args...)
	result := m.Called(args...)
	return result.Get(0).(*gorm.DB)
}

func (m *MockDB) Create(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Save(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(value, conds)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Preload(column string, conditions ...interface{}) *gorm.DB {
	args := m.Called(column, conditions)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Offset(offset int) *gorm.DB {
	args := m.Called(offset)
	return args.Get(0).(*gorm.DB)
}
