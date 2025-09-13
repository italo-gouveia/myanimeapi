# 🧪 **MyAnimeAPI Tests**

This directory contains the complete test structure for MyAnimeAPI, including unit, integration, and E2E tests.

## 🏗️ **Directory Structure**

```
tests/
├── unit/                    # Unit tests
│   ├── models/             # Model tests
│   ├── services/           # Service tests
│   ├── repositories/       # Repository tests
│   ├── handlers/           # Handler tests
│   └── middleware/         # Middleware tests
├── integration/            # Integration tests
│   ├── database/           # Database tests
│   ├── api/                # API tests
│   └── auth/               # Authentication tests
├── e2e/                    # End-to-end tests
│   ├── api/                # E2E API tests
│   └── frontend/           # E2E frontend tests
├── fixtures/               # Reusable test data
│   ├── models/             # Model fixtures
│   ├── database/           # Database fixtures
│   └── http/               # HTTP request fixtures
├── mocks/                  # Generated and custom mocks
├── utils/                  # Test utilities
├── suites/                 # Organized test suites
└── config/                 # Test-specific configuration
```

## 🚀 **How to Run Tests**

### **Run All Tests**
```bash
make test
```

### **Run Specific Tests**
```bash
# Unit tests only
make test-unit

# Integration tests only
make test-integration

# E2E tests only
make test-e2e

# Tests with coverage
make test-coverage

# Parallel tests
make test-parallel
```

### **Run Specific Test**
```bash
make test-specific TEST_NAME=TestAnimeService_CreateAnime
```

## 🔧 **Environment Configuration**

### **Environment Variables for Tests**
```bash
# Database
export TEST_DB_HOST=localhost
export TEST_DB_PORT=5432
export TEST_DB_USER=test
export TEST_DB_PASSWORD=test
export TEST_DB_NAME=test_db

# Server
export TEST_SERVER_PORT=8081
export TEST_LOG_LEVEL=debug
```

### **Docker Compose for Tests**
```bash
# Start test database
docker-compose -f docker-compose.test.yml up -d

# Stop test database
docker-compose -f docker-compose.test.yml down
```

## 📝 **Naming Conventions**

### **Test Files**
- `*_test.go` for all test files
- `Test[StructName]_[MethodName]_[Scenario]` for function names

### **Examples**
```go
func TestAnimeService_CreateAnime_Success(t *testing.T)
func TestAnimeService_CreateAnime_InvalidInput(t *testing.T)
func TestAnimeService_CreateAnime_DatabaseError(t *testing.T)
```

## 🎭 **Fixtures and Mocks**

### **Fixtures**
Fixtures are located in `tests/fixtures/models/` and provide reusable test data:

```go
// Create default user
user := modelsfixtures.ActiveUser()

// Create user with specific options
user := modelsfixtures.NewUserFixture(
    modelsfixtures.WithUsername("admin"),
    modelsfixtures.WithEmail("admin@example.com"),
)
```

### **Mocks**
Mocks are managed by `MockManager` in `tests/mocks/mock_manager.go`:

```go
suite := suites.NewBaseSuite(t)
suite.Mocks.UserRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
```

## 🧪 **Test Types**

### **1. Unit Tests**
- Test components in isolation
- Use mocks for external dependencies
- Run quickly
- Located in `tests/unit/`

### **2. Integration Tests**
- Test interaction between components
- Use real database (test)
- Run slower
- Located in `tests/integration/`

### **3. E2E Tests**
- Test complete application flows
- Use real HTTP server
- Run slower
- Located in `tests/e2e/`

## 🏗️ **BaseSuite**

`BaseSuite` provides common configuration for all tests:

```go
suite := suites.NewBaseSuite(t)
defer suite.CleanupDatabase()

// Use test database
err := suite.DB.Create(fixture).Error

// Use router for E2E tests
suite.Router.ServeHTTP(rec, req)
```

## 📊 **Test Coverage**

### **Generate Coverage Report**
```bash
make test-coverage
```

### **View Coverage**
```bash
# Open in browser
open coverage.html

# Or use go tool cover
go tool cover -html=coverage.out
```

## 🔄 **CI/CD**

Tests are automatically executed in GitHub Actions:

- **Push/Pull Request**: Runs all tests
- **Coverage**: Generates coverage report
- **Upload**: Sends coverage to Codecov

## 🐛 **Test Debugging**

### **Run Test with Verbose**
```bash
go test -v ./tests/unit/services/
```

### **Run Specific Test with Verbose**
```bash
go test -v -run TestAnimeService_CreateAnime ./tests/unit/services/
```

### **Run with Race Detector**
```bash
go test -race ./tests/...
```

## 📚 **Best Practices**

### **1. AAA Structure**
```go
func TestExample(t *testing.T) {
    // Arrange - Prepare data and mocks
    mockRepo := NewMockRepository(t)
    service := NewService(mockRepo)
    
    // Act - Execute the action
    result, err := service.Execute(input)
    
    // Assert - Verify results
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

### **2. Automatic Cleanup**
```go
suite := suites.NewBaseSuite(t)
defer suite.CleanupDatabase() // Cleans automatically
```

### **3. Parallel Tests**
```go
func TestExample(t *testing.T) {
    t.Parallel() // Runs in parallel when possible
    // ... test
}
```

### **4. Helpers**
```go
func createTestUser(t *testing.T) *models.User {
    t.Helper() // Marks as helper
    return modelsfixtures.ActiveUser()
}
```

## 🚨 **Troubleshooting**

### **Database Connection Error**
```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Check environment variables
echo $TEST_DB_HOST
echo $TEST_DB_PORT
```

### **Mock Error**
```bash
# Regenerate mocks
make test-mocks

# Check if mockgen is installed
which mockgen
```

### **Failing Tests**
```bash
# Clear test cache
make test-clean

# Run with verbose for more details
go test -v ./tests/...
```

## 📞 **Support**

For test-related questions:
1. Check this README
2. Consult the architecture documentation
3. Check existing examples
4. Open an issue in the repository

---

**Remember**: Tests are essential for code quality! 🚀
