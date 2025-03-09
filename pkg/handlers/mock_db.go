// pkg/handlers/mock_db.go
package handlers

import (
	"context"
	"myanimeapi/pkg/models"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// MockDB is a mock implementation of DBInterface
type MockDB struct {
	mock.Mock
}

// First is a mock method for querying the database
/*func (m *MockDB) First(out interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(out, where)
	if len(args) > 0 {
		return args.Get(0).(*gorm.DB)
	}
	return nil
}*/

/*func (m *MockDB) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(dest, conds)
	return args.Get(0).(*gorm.DB)
}*/

/*
	func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
		args = append([]interface{}{query}, args...)
		result := m.Called(args...)
		return result.Get(0).(*gorm.DB)
	}
*/

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

/*
	func (m *MockDB) Offset(offset int) *gorm.DB {
		args := m.Called(offset)
		return args.Get(0).(*gorm.DB)
	}
*/
func (m *MockDB) Where(ctx context.Context, query interface{}, args ...interface{}) *gorm.DB {
	// Call the mock's Called method with the arguments
	m.Called(ctx, query, args)
	// Return a properly initialized *gorm.DB object
	return &gorm.DB{
		Statement: &gorm.Statement{
			DB:      &gorm.DB{},
			Clauses: make(map[string]clause.Clause), // Initialize the Clauses map
		},
		Config: &gorm.Config{}, // Initialize the Config field
	}
}

func (m *MockDB) Offset(offset int) *gorm.DB {
	// Call the mock's Called method with the argument
	m.Called(offset)
	// Return a properly initialized *gorm.DB object
	return &gorm.DB{
		Statement: &gorm.Statement{
			DB:      &gorm.DB{},
			Clauses: make(map[string]clause.Clause), // Initialize the Clauses map
		},
		Config: &gorm.Config{}, // Initialize the Config field
	}
}

func (m *MockDB) Limit(limit int) *gorm.DB {
	// Call the mock's Called method with the argument
	m.Called(limit)
	// Return a properly initialized *gorm.DB object
	return &gorm.DB{
		Statement: &gorm.Statement{
			DB:      &gorm.DB{},
			Clauses: make(map[string]clause.Clause), // Initialize the Clauses map
		},
		Config: &gorm.Config{}, // Initialize the Config field
	}
}

/*func (m *MockDB) Find(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB {
	// Call the mock's Called method with the arguments
	m.Called(ctx, dest, conds)
	// Return a properly initialized *gorm.DB object
	return &gorm.DB{
		Statement: &gorm.Statement{
			DB:      &gorm.DB{},
			Clauses: make(map[string]clause.Clause), // Initialize the Clauses map
		},
		Config: &gorm.Config{}, // Initialize the Config field
	}
}*/

func (m *MockDB) Find(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB {
	// Call the mock's Called method with the arguments
	m.Called(ctx, dest, conds)
	// Set the destination value (if applicable)
	if destPtr, ok := dest.(*[]models.Review); ok {
		*destPtr = []models.Review{
			{ID: 1, AnimeID: 1, Content: "Great anime!"},
			{ID: 2, AnimeID: 1, Content: "Awesome!"},
		}
	}
	// Return a properly initialized *gorm.DB object
	return &gorm.DB{
		Statement: &gorm.Statement{
			DB:      &gorm.DB{},
			Clauses: make(map[string]clause.Clause), // Initialize the Clauses map
		},
		Config: &gorm.Config{}, // Initialize the Config field
	}
}

func (m *MockDB) First(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB {
	// Call the mock's Called method with the arguments
	m.Called(ctx, dest, conds)
	// Set the destination value
	if destPtr, ok := dest.(**models.Anime); ok {
		*destPtr = &models.Anime{ID: 1, Title: "Naruto"} // Set the destination to a test anime
	}
	// Return a properly initialized *gorm.DB object
	return &gorm.DB{
		Statement: &gorm.Statement{
			DB:      &gorm.DB{},
			Clauses: make(map[string]clause.Clause), // Initialize the Clauses map
		},
		Config: &gorm.Config{}, // Initialize the Config field
	}
}
