package suites

import (
	"testing"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"myanimeapi/api/routes"
	"myanimeapi/internal/config"
	appdb "myanimeapi/internal/db"
	"myanimeapi/tests/mocks"
	"myanimeapi/tests/utils"
)

// BaseSuite provides common setup for all test suites
type BaseSuite struct {
	T            *testing.T
	Router       *mux.Router
	DB           *gorm.DB
	Mocks        *mocks.MockManager
	Config       *config.Config
	CleanupFuncs []func()
}

// NewBaseSuite creates a new test suite with all dependencies
func NewBaseSuite(t *testing.T) *BaseSuite {
	t.Helper()

	suite := &BaseSuite{
		T:            t,
		CleanupFuncs: make([]func(), 0),
	}

	suite.setupTestEnvironment()
	suite.setupMocks()
	suite.setupApplication()

	t.Cleanup(suite.cleanup)

	return suite
}

// setupTestEnvironment initializes the test environment
func (s *BaseSuite) setupTestEnvironment() {
	// Load test configuration
	s.Config = utils.LoadTestConfig(s.T)

	// Setup test database
	s.DB = utils.SetupTestDatabase(s.T)
	s.addCleanup(func() {
		utils.CleanupTestDatabase(s.T, s.DB)
	})
}

// setupMocks initializes all mocks
func (s *BaseSuite) setupMocks() {
	s.Mocks = mocks.NewMockManager(s.T)
	s.Mocks.SetupDefaultExpectations()
}

// setupApplication initializes the application components
func (s *BaseSuite) setupApplication() {
	// Setup router
	s.Router = mux.NewRouter()

	// Register application routes so E2E tests hit real handlers
	dbWrapper := appdb.NewGormDB(s.DB)
	routes.RegisterRoutes(s.Router, "", dbWrapper, nil, "test")
}

// addCleanup adds a cleanup function to be executed after the test
func (s *BaseSuite) addCleanup(cleanup func()) {
	s.CleanupFuncs = append(s.CleanupFuncs, cleanup)
}

// cleanup executes all cleanup functions
func (s *BaseSuite) cleanup() {
	for _, cleanup := range s.CleanupFuncs {
		cleanup()
	}
}

// GetTestDB returns the test database instance
func (s *BaseSuite) GetTestDB() *gorm.DB {
	return s.DB
}

// GetTestRouter returns the test router instance
func (s *BaseSuite) GetTestRouter() *mux.Router {
	return s.Router
}

// GetTestConfig returns the test configuration
func (s *BaseSuite) GetTestConfig() *config.Config {
	return s.Config
}
