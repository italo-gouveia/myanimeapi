// tests/smoke/smoke_test.go
// This package defines the smoke tests for the server.
// It is used to test the server functionality.
// It imports the necessary packages to run the tests.
// It defines the TestMain function to initialize the mock database and handlers before running the tests.
// It defines the TestRegisterUser function to test the user registration endpoint.
// It defines the TestAuthenticateUser function to test the user authentication endpoint.
// It defines the TestGetAnime function to test the endpoint to get an anime by ID.
// It defines the TestCreateReview function to test the endpoint to create a review.
// It uses the myanimeapi/internal/mocks package to create a mock database.
// It uses the myanimeapi/pkg/handlers package to test the handlers.
// It uses the myanimeapi/pkg/middleware package to simulate the middleware.
// It uses the myanimeapi/pkg/models package to define the test models.
// It uses the net/http/httptest package to create a new HTTP request.
// It uses the os package to exit the tests with the correct status code.
// It uses the testing package to define the tests.
// It uses the gomock package to create a new mock controller.
// It uses the gorm.io/gorm package to create a valid gorm.DB object for the mock to return.
// It uses the golang/mock/gomock package to create a new mock controller.
// It uses the encoding/json package to marshal and unmarshal JSON data.
// It uses the bytes package to create a new buffer for the HTTP request.
// It uses the context package to create a new context for the HTTP request.
// It uses the http package to create a new HTTP request.
// It uses the testing package to define the tests.
// It uses the fmt package to print messages.
// It uses the os package to exit the tests with the correct status code.
// It uses the myanimeapi/pkg/handlers package to create the handlers.
package main

/*
import (
	"bytes"
	"context"
	"encoding/json"
	"myanimeapi/api/handlers"
	"myanimeapi/api/mocks"
	"myanimeapi/api/models"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"myanimeapi/internal/errors" // Import the errors package

	"github.com/golang/mock/gomock"
	"gorm.io/gorm"
)

// Define handler variables at the package level
var (
	authHandler   *handlers.AuthHandler
	animeHandler  *handlers.AnimeHandler
	reviewHandler *handlers.ReviewHandler
	mockDB        *mocks.MockDBInterface
	mockCtrl      *gomock.Controller
)

// TestMain initializes the mock database and handlers before running the tests
func TestMain(m *testing.M) {
	// Initialize the mock controller
	mockCtrl = gomock.NewController(nil)
	defer mockCtrl.Finish()

	// Create the mock DB
	mockDB = mocks.NewMockDBInterface(mockCtrl)

	// Initialize handlers with the mock DB
	authHandler = handlers.NewAuthHandler(mockDB)
	animeHandler = handlers.NewAnimeHandler(mockDB)
	reviewHandler = handlers.NewReviewHandler(mockDB)

	// Run the tests
	os.Exit(m.Run())
}

func TestHealthCheck(t *testing.T) {
	// Initialize the mock controller
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Create the mock DB
	mockDB := mocks.NewMockDBInterface(mockCtrl)

	// Set up the expectation for IsHealthy
	mockDB.EXPECT().IsHealthy().Return(true)

	// Create a request to the health check endpoint
	req, err := http.NewRequest("GET", "/v1/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler with the mock database instance
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if mockDB == nil {
			errors.WriteErrorResponse(w, http.StatusServiceUnavailable, errors.ErrServiceUnavailable, "Database instance not initialized", "The database instance is not initialized.")
			return
		}

		if mockDB.IsHealthy() {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte("OK")); err != nil {
				errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to write response", "An internal server error occurred while writing the response.")
				return
			}
		} else {
			errors.WriteErrorResponse(w, http.StatusServiceUnavailable, errors.ErrServiceUnavailable, "Database connection failed", "The database connection is not healthy.")
		}
	})

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := "OK"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// Verify that all expectations were met
	mockCtrl.Finish()
}*/

// TestRegisterUser tests the user registration endpoint
/*func TestRegisterUser(t *testing.T) {
	// Reset mock expectations before each test
	mockCtrl.Finish()
	mockCtrl = gomock.NewController(t)

	tests := []struct {
		name           string
		user           models.User
		setupMocks     func()
		expectedStatus int
		checkResponse  func(*httptest.ResponseRecorder)
	}{
		{
			name: "Successful Registration",
			user: models.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func() {
				// Create a mock DB instance for the chain
				mockDBInstance := &gorm.DB{}

				// Expect username check
				mockDB.EXPECT().
					Where(gomock.Any(), "username = ?", "testuser").
					Return(mockDBInstance)

				mockDB.EXPECT().
					First(gomock.Any(), gomock.Any()).
					Return(&gorm.DB{Error: gorm.ErrRecordNotFound})

				// Expect email check
				mockDB.EXPECT().
					Where(gomock.Any(), "email = ?", "test@example.com").
					Return(mockDBInstance)

				mockDB.EXPECT().
					First(gomock.Any(), gomock.Any()).
					Return(&gorm.DB{Error: gorm.ErrRecordNotFound})

				// Expect user creation
				mockDB.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ interface{}, user interface{}) *gorm.DB {
						u := user.(*models.User)
						u.ID = 1
						u.CreatedAt = time.Now()
						u.UpdatedAt = time.Now()
						return &gorm.DB{}
					})
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(rr *httptest.ResponseRecorder) {
				var response models.User
				err := json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Errorf("Failed to decode response: %v", err)
				}
				if response.ID != 1 {
					t.Errorf("Expected user ID 1, got %d", response.ID)
				}
				if response.Username != "testuser" {
					t.Errorf("Expected username 'testuser', got %s", response.Username)
				}
				if response.Email != "test@example.com" {
					t.Errorf("Expected email 'test@example.com', got %s", response.Email)
				}
			},
		},
		{
			name: "Username Already Exists",
			user: models.User{
				Username: "existinguser",
				Email:    "new@example.com",
				Password: "password123",
			},
			setupMocks: func() {
				// Expect username check to find existing user
				mockDB.EXPECT().
					Where(gomock.Any(), "username = ?", "existinguser").
					Return(&gorm.DB{})

				mockDB.EXPECT().
					First(gomock.Any(), gomock.Any()).
					Return(&gorm.DB{})
			},
			expectedStatus: http.StatusConflict,
			checkResponse: func(rr *httptest.ResponseRecorder) {
				var response errors.ErrorResponse
				err := json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Errorf("Failed to decode error response: %v", err)
				}
				if response.Error.Message != "User with this username already exists" {
					t.Errorf("Expected error message about existing username, got %s", response.Error.Message)
				}
			},
		},
		{
			name: "Email Already Exists",
			user: models.User{
				Username: "newuser",
				Email:    "existing@example.com",
				Password: "password123",
			},
			setupMocks: func() {
				// Expect username check
				mockDB.EXPECT().
					Where(gomock.Any(), "username = ?", "newuser").
					Return(&gorm.DB{})

				mockDB.EXPECT().
					First(gomock.Any(), gomock.Any()).
					Return(&gorm.DB{Error: gorm.ErrRecordNotFound})

				// Expect email check to find existing user
				mockDB.EXPECT().
					Where(gomock.Any(), "email = ?", "existing@example.com").
					Return(&gorm.DB{})

				mockDB.EXPECT().
					First(gomock.Any(), gomock.Any()).
					Return(&gorm.DB{})
			},
			expectedStatus: http.StatusConflict,
			checkResponse: func(rr *httptest.ResponseRecorder) {
				var response errors.ErrorResponse
				err := json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Errorf("Failed to decode error response: %v", err)
				}
				if response.Error.Message != "User with this email already exists" {
					t.Errorf("Expected error message about existing email, got %s", response.Error.Message)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			tt.setupMocks()

			// Create request body
			userJSON, err := json.Marshal(tt.user)
			if err != nil {
				t.Fatalf("Failed to marshal user: %v", err)
			}

			// Create request
			req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(userJSON))
			req.Header.Set("Content-Type", "application/json")

			// Add validated payload to context
			ctx := context.WithValue(req.Context(), middleware.ValidatedPayloadKey, &tt.user)
			req = req.WithContext(ctx)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			handler := http.HandlerFunc(authHandler.RegisterUserHandler)
			handler.ServeHTTP(rr, req)

			// Check status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			// Check response
			tt.checkResponse(rr)
		})
	}
}*/

// TestAuthenticateUser tests the user authentication endpoint
/*func TestAuthenticateUser(t *testing.T) {
	// Initialize the mock controller
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Create the mock DB
	mockDB := mocks.NewMockDBInterface(mockCtrl)

	// Create the auth handler with the mock DB
	authHandler := handlers.NewAuthHandler(mockDB)

	tests := []struct {
		name           string
		credentials    models.UserCredentials
		setupMocks     func()
		expectedStatus int
		checkResponse  func(*httptest.ResponseRecorder)
	}{
		{
			name: "Successful Authentication",
			credentials: models.UserCredentials{
				Username: "testuser",
				Password: "password123",
			},
			setupMocks: func() {
				// Mock user lookup
				mockDB.EXPECT().
					Where(gomock.Any(), "username = ?", "testuser").
					Return(&gorm.DB{})

				mockDB.EXPECT().
					First(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ interface{}, _ ...interface{}) *gorm.DB {
						// Simulate finding the user
						return &gorm.DB{}
					})

				// Mock saving the migrated hash (if needed)
				mockDB.EXPECT().
					Save(gomock.Any(), gomock.Any()).
					Return(&gorm.DB{}).
					AnyTimes()
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(rr *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Errorf("Failed to decode response: %v", err)
				}
				if _, exists := response["token"]; !exists {
					t.Error("Response does not contain a token")
				}
			},
		},
		{
			name: "User Not Found",
			credentials: models.UserCredentials{
				Username: "nonexistentuser",
				Password: "password123",
			},
			setupMocks: func() {
				// Mock user lookup failure
				mockDB.EXPECT().
					Where(gomock.Any(), "username = ?", "nonexistentuser").
					Return(&gorm.DB{})

				mockDB.EXPECT().
					First(gomock.Any(), gomock.Any()).
					Return(&gorm.DB{Error: gorm.ErrRecordNotFound})
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(rr *httptest.ResponseRecorder) {
				var response errors.ErrorResponse
				err := json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Errorf("Failed to decode error response: %v", err)
				}
				if response.Error.Message != "User not found" {
					t.Errorf("Expected error message 'User not found', got %s", response.Error.Message)
				}
			},
		},
		{
			name: "Invalid Password",
			credentials: models.UserCredentials{
				Username: "testuser",
				Password: "wrongpassword",
			},
			setupMocks: func() {
				// Mock user lookup success but password verification failure
				mockDB.EXPECT().
					Where(gomock.Any(), "username = ?", "testuser").
					Return(&gorm.DB{})

				mockDB.EXPECT().
					First(gomock.Any(), gomock.Any()).
					DoAndReturn(func(dest interface{}, _ ...interface{}) *gorm.DB {
						user := dest.(*models.User)
						user.Username = "testuser"
						user.Password = "$2a$10$invalidhashforpassword" // Invalid hash that won't match
						return &gorm.DB{}
					})
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(rr *httptest.ResponseRecorder) {
				var response errors.ErrorResponse
				err := json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Errorf("Failed to decode error response: %v", err)
				}
				if response.Error.Message != "Invalid credentials" {
					t.Errorf("Expected error message 'Invalid credentials', got %s", response.Error.Message)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			tt.setupMocks()

			// Create request body
			credentialsJSON, err := json.Marshal(tt.credentials)
			if err != nil {
				t.Fatalf("Failed to marshal credentials: %v", err)
			}

			// Create request
			req := httptest.NewRequest("POST", "/auth/authenticate", bytes.NewBuffer(credentialsJSON))
			req.Header.Set("Content-Type", "application/json")

			// Add validated payload to context
			ctx := context.WithValue(req.Context(), middleware.ValidatedPayloadKey, &tt.credentials)
			req = req.WithContext(ctx)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			handler := http.HandlerFunc(authHandler.AuthenticateHandler)
			handler.ServeHTTP(rr, req)

			// Check status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			// Check response
			tt.checkResponse(rr)
		})
	}
}*/

// TestGetAnime tests the endpoint to get an anime by ID
/*func TestGetAnime(t *testing.T) {
	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/anime/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a valid gorm.DB object for the mock to return
	mockDBInstance := &gorm.DB{}

	// Set up mock expectations
	mockDB.EXPECT().
		First(gomock.Any(), gomock.Any()).
		Return(mockDBInstance).
		Do(func(dest interface{}, conds ...interface{}) {
			// Simulate a found anime
			destAnime := dest.(*models.Anime)
			*destAnime = models.Anime{
				ID:    1,
				Title: "Naruto",
			}
		})

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(animeHandler.GetAnimeHandler)

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	var anime models.Anime
	if err := json.NewDecoder(rr.Body).Decode(&anime); err != nil {
		t.Fatal(err)
	}
	if anime.ID != 1 {
		t.Errorf("handler returned unexpected body: got %v want %v",
			anime.ID, 1)
	}
}*/

// TestCreateReview tests the endpoint to create a review
/*func TestCreateReview(t *testing.T) {
	// Create a new HTTP request
	review := models.Review{
		UserID:  1,
		AnimeID: 1,
		Content: "Great anime!",
		Rating:  9,
	}
	reviewJSON, _ := json.Marshal(review)
	req, err := http.NewRequest("POST", "/reviews", bytes.NewBuffer(reviewJSON))
	if err != nil {
		t.Fatal(err)
	}

	// Simulate the middleware by adding the payload to the request context
	ctx := context.WithValue(req.Context(), middleware.ValidatedPayloadKey, &review)
	req = req.WithContext(ctx)

	// Create a valid gorm.DB object for the mock to return
	mockDBInstance := &gorm.DB{}

	// Set up mock expectations
	mockDB.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(mockDBInstance).
		Do(func(value interface{}, conds ...interface{}) {
			// Simulate the created review
			destReview := value.(*models.Review)
			*destReview = review
		})

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(reviewHandler.CreateReviewHandler)

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Check the response body
	var responseReview models.Review
	if err := json.NewDecoder(rr.Body).Decode(&responseReview); err != nil {
		t.Fatal(err)
	}
	if responseReview.Content != review.Content {
		t.Errorf("handler returned unexpected body: got %v want %v",
			responseReview.Content, review.Content)
	}
}*/

// TestFavorites tests the favorite-related endpoints
/*func TestFavorites(t *testing.T) {
	// Reset mock expectations before each test
	mockCtrl.Finish()
	mockCtrl = gomock.NewController(t)

	tests := []struct {
		name           string
		setupMocks     func()
		request        func() *http.Request
		expectedStatus int
		checkResponse  func(*httptest.ResponseRecorder)
	}{
		{
			name: "Add Favorite Success",
			setupMocks: func() {
				// Mock anime existence check
				mockDB.EXPECT().
					First(gomock.Any(), gomock.Any()).
					Return(&gorm.DB{}).
					Do(func(dest interface{}, _ ...interface{}) {
						anime := dest.(*models.Anime)
						*anime = models.Anime{ID: 1, Title: "Test Anime"}
					})

				// Mock favorite existence check
				mockDB.EXPECT().
					Where(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&gorm.DB{})

				mockDB.EXPECT().
					First(gomock.Any(), gomock.Any()).
					Return(&gorm.DB{Error: gorm.ErrRecordNotFound})

				// Mock favorite creation
				mockDB.EXPECT().
					Create(gomock.Any()).
					Return(&gorm.DB{}).
					Do(func(value interface{}, _ ...interface{}) {
						favorite := value.(*models.Favorite)
						*favorite = models.Favorite{
							ID:        1,
							UserID:    1,
							AnimeID:   1,
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						}
					})
			},
			request: func() *http.Request {
				req := httptest.NewRequest("POST", "/favorites", bytes.NewBufferString(`{"anime_id": 1}`))
				req.Header.Set("Content-Type", "application/json")
				ctx := context.WithValue(req.Context(), "user", &models.User{ID: 1})
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(rr *httptest.ResponseRecorder) {
				var favorite models.Favorite
				if err := json.NewDecoder(rr.Body).Decode(&favorite); err != nil {
					t.Errorf("Failed to decode response: %v", err)
				}
				if favorite.ID != 1 {
					t.Errorf("Expected favorite ID 1, got %d", favorite.ID)
				}
				if favorite.AnimeID != 1 {
					t.Errorf("Expected anime ID 1, got %d", favorite.AnimeID)
				}
			},
		},
		{
			name: "Get Favorites Success",
			setupMocks: func() {
				mockDB.EXPECT().
					Preload(gomock.Any()).
					Return(mockDB)

				mockDB.EXPECT().
					Where(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(mockDB)

				mockDB.EXPECT().
					Find(gomock.Any()).
					Return(&gorm.DB{}).
					Do(func(dest interface{}, _ ...interface{}) {
						favorites := dest.(*[]models.Favorite)
						*favorites = []models.Favorite{
							{
								ID:        1,
								UserID:    1,
								AnimeID:   1,
								CreatedAt: time.Now(),
								UpdatedAt: time.Now(),
								Anime: models.Anime{
									ID:    1,
									Title: "Test Anime",
								},
							},
						}
					})
			},
			request: func() *http.Request {
				req := httptest.NewRequest("GET", "/favorites", nil)
				ctx := context.WithValue(req.Context(), "user", &models.User{ID: 1})
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(rr *httptest.ResponseRecorder) {
				var favorites []models.Favorite
				if err := json.NewDecoder(rr.Body).Decode(&favorites); err != nil {
					t.Errorf("Failed to decode response: %v", err)
				}
				if len(favorites) != 1 {
					t.Errorf("Expected 1 favorite, got %d", len(favorites))
				}
				if favorites[0].Anime.Title != "Test Anime" {
					t.Errorf("Expected anime title 'Test Anime', got %s", favorites[0].Anime.Title)
				}
			},
		},
		{
			name: "Remove Favorite Success",
			setupMocks: func() {
				mockDB.EXPECT().
					Where(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(mockDB)

				mockDB.EXPECT().
					Delete(gomock.Any()).
					Return(&gorm.DB{})
			},
			request: func() *http.Request {
				req := httptest.NewRequest("DELETE", "/favorites/1", nil)
				ctx := context.WithValue(req.Context(), "user", &models.User{ID: 1})
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusNoContent,
			checkResponse: func(rr *httptest.ResponseRecorder) {
				if rr.Body.Len() != 0 {
					t.Error("Expected empty response body")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/favorites":
					if r.Method == "POST" {
						favoriteHandler.AddFavoriteHandler(w, r)
					} else if r.Method == "GET" {
						favoriteHandler.GetFavoritesHandler(w, r)
					}
				case "/favorites/1":
					if r.Method == "DELETE" {
						favoriteHandler.RemoveFavoriteHandler(w, r)
					}
				}
			})

			handler.ServeHTTP(rr, tt.request())

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			tt.checkResponse(rr)
		})
	}
}*/
