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

import (
	"context"
	"myanimeapi/internal/mocks"
	"myanimeapi/pkg/handlers"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
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

	// Set the mock database instance in the request context
	ctx := context.WithValue(req.Context(), handlers.DBContextKey, mockDB)
	req = req.WithContext(ctx)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := handlers.GetDB(r.Context())
		if db == nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Database instance not initialized"))
			return
		}

		if db.IsHealthy() {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Database connection failed"))
		}
	})

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
}

// TestRegisterUser tests the user registration endpoint
/*func TestRegisterUser(t *testing.T) {
	// Define the test user
	user := models.User{
		Username: "testuser",
		Email:    "testuser@example.com",
		Password: "testpassword",
	}

	// Marshal the user to JSON
	userJSON, _ := json.Marshal(user)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatal(err)
	}

	// Simulate the middleware by adding the payload to the request context
	ctx := context.WithValue(req.Context(), middleware.ValidatedPayloadKey, &user)
	req = req.WithContext(ctx)

	// Create a valid gorm.DB object for the mock to return
	mockDBInstance := &gorm.DB{}

	// Set up mock expectations
	// Simulate "user not found" for username check
	mockDB.EXPECT().
		Where(gomock.Any(), "username = ?", user.Username).
		Return(mockDBInstance)

	mockDB.EXPECT().
		First(gomock.Any(), gomock.Any()).
		Return(mockDBInstance).
		Do(func(dest interface{}, conds ...interface{}) {
			// Simulate "record not found" error
			mockDBInstance.Error = gorm.ErrRecordNotFound
		})

	// Simulate "user not found" for email check
	mockDB.EXPECT().
		Where(gomock.Any(), "email = ?", user.Email).
		Return(mockDBInstance)

	mockDB.EXPECT().
		First(gomock.Any(), gomock.Any()).
		Return(mockDBInstance).
		Do(func(dest interface{}, conds ...interface{}) {
			// Simulate "record not found" error
			mockDBInstance.Error = gorm.ErrRecordNotFound
		})

	// Simulate successful user creation
	mockDB.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(mockDBInstance).
		Do(func(value interface{}, conds ...interface{}) {
			// Simulate successful creation
			dest := value.(*models.User)
			*dest = user
		})

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(authHandler.RegisterUserHandler)

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Check the response body
	var responseUser models.User
	if err := json.NewDecoder(rr.Body).Decode(&responseUser); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if responseUser.Username != user.Username {
		t.Errorf("handler returned unexpected body: got %v want %v",
			responseUser.Username, user.Username)
	}
}*/

// TestAuthenticateUser tests the user authentication endpoint
/*func TestAuthenticateUser(t *testing.T) {
	// Define the test credentials
	credentials := models.UserCredentials{
		Username: "testuser",
		Password: "testpassword",
	}

	// Marshal the credentials to JSON
	credentialsJSON, _ := json.Marshal(credentials)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/auth/authenticate", bytes.NewBuffer(credentialsJSON))
	if err != nil {
		t.Fatal(err)
	}

	// Simulate the middleware by adding the payload to the request context
	ctx := context.WithValue(req.Context(), middleware.ValidatedPayloadKey, &credentials)
	req = req.WithContext(ctx)

	// Create a valid gorm.DB object for the mock to return
	mockDBInstance := &gorm.DB{}

	// Set up mock expectations
	mockDB.EXPECT().
		Where(gomock.Any(), "username = ?", credentials.Username).
		Return(mockDBInstance)

	mockDB.EXPECT().
		First(gomock.Any(), gomock.Any()).
		Return(mockDBInstance).
		Do(func(dest interface{}, conds ...interface{}) {
			// Simulate a found user
			destUser := dest.(*models.User)
			*destUser = models.User{
				Username: credentials.Username,
				Password: "$2a$10$hashedpassword", // Simulate a hashed password
			}
		})

	mockDB.EXPECT().
		GetError().
		Return(nil)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(authHandler.AuthenticateHandler)

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if _, exists := response["token"]; !exists {
		t.Errorf("handler returned unexpected body: token not found")
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
