// Package smoke contains end-to-end smoke tests that spin up a real HTTP
// server wired with mocked services and verify the critical happy paths with
// actual HTTP round-trips — no database required.
package smoke_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"myanimeapi/api/handlers"
	"myanimeapi/api/mocks"
	"myanimeapi/api/models"
	apperrors "myanimeapi/internal/errors"
	"myanimeapi/internal/logger"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// ---------------------------------------------------------------------------
// Suite
// ---------------------------------------------------------------------------

// SmokeSuite verifies the most critical happy-path and error scenarios via
// real HTTP round-trips. Each test gets fresh mocks and its own httptest.Server
// so expectations are always fully isolated.
type SmokeSuite struct {
	suite.Suite

	// rebuilt per test
	server   *httptest.Server
	client   *http.Client
	authSvc  *mocks.MockAuthServiceInterface
	animeSvc *mocks.MockAnimeServiceInterface
	genreSvc *mocks.MockGenreServiceInterface
	tagSvc   *mocks.MockTagServiceInterface
}

func TestSmokeSuite(t *testing.T) {
	suite.Run(t, new(SmokeSuite))
}

// SetupTest spins up a fresh server with fresh mocks before each test.
func (s *SmokeSuite) SetupTest() {
	s.authSvc = new(mocks.MockAuthServiceInterface)
	s.animeSvc = new(mocks.MockAnimeServiceInterface)
	s.genreSvc = new(mocks.MockGenreServiceInterface)
	s.tagSvc = new(mocks.MockTagServiceInterface)

	log := logger.New()

	router := mux.NewRouter()
	v1 := router.PathPrefix("/v1").Subrouter()

	handlers.NewAuthHandler(s.authSvc, log).RegisterAuthRoutes(v1)
	handlers.NewAnimeHandler(s.animeSvc).RegisterAnimeRoutes(v1)
	handlers.NewGenreHandler(s.genreSvc).RegisterGenreRoutes(v1)
	handlers.NewTagHandler(s.tagSvc).RegisterTagRoutes(v1)

	s.server = httptest.NewServer(router)
	s.client = s.server.Client()
}

// TearDownTest asserts all mock expectations and closes the server.
func (s *SmokeSuite) TearDownTest() {
	s.server.Close()
	s.authSvc.AssertExpectations(s.T())
	s.animeSvc.AssertExpectations(s.T())
	s.genreSvc.AssertExpectations(s.T())
	s.tagSvc.AssertExpectations(s.T())
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func (s *SmokeSuite) url(path string) string {
	return s.server.URL + "/v1" + path
}

func (s *SmokeSuite) postJSON(path string, body interface{}) *http.Response {
	s.T().Helper()
	b, err := json.Marshal(body)
	s.Require().NoError(err)
	resp, err := s.client.Post(s.url(path), "application/json", bytes.NewReader(b))
	s.Require().NoError(err)
	return resp
}

func (s *SmokeSuite) getURL(path string) *http.Response {
	s.T().Helper()
	resp, err := s.client.Get(s.url(path))
	s.Require().NoError(err)
	return resp
}

func (s *SmokeSuite) decodeBody(resp *http.Response, dst interface{}) {
	s.T().Helper()
	s.Require().NoError(json.NewDecoder(resp.Body).Decode(dst))
}

// ---------------------------------------------------------------------------
// Auth
// ---------------------------------------------------------------------------

func (s *SmokeSuite) TestAuth_Register_HappyPath() {
	s.authSvc.EXPECT().
		RegisterUser(mock.Anything, mock.MatchedBy(func(u *models.User) bool {
			return u.Username == "smokeuser" && u.Email == "smoke@example.com"
		})).
		Return(nil).Once()

	resp := s.postJSON("/auth/register", map[string]string{
		"username": "smokeuser",
		"email":    "smoke@example.com",
		"password": "Password123!",
	})
	defer func() { _ = resp.Body.Close() }()

	s.Equal(http.StatusCreated, resp.StatusCode)
	s.Contains(resp.Header.Get("Content-Type"), "application/json")
}

func (s *SmokeSuite) TestAuth_Register_BadJSON() {
	resp, err := s.client.Post(s.url("/auth/register"), "application/json", bytes.NewReader([]byte("bad json{")))
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()

	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *SmokeSuite) TestAuth_Authenticate_HappyPath() {
	s.authSvc.EXPECT().
		AuthenticateUser(mock.Anything, mock.MatchedBy(func(c *models.UserCredentials) bool {
			return c.Username == "smokeuser" && c.Password == "Password123!"
		})).
		Return("jwt-smoke-token", nil).Once()

	resp := s.postJSON("/auth/authenticate", map[string]string{
		"username": "smokeuser",
		"password": "Password123!",
	})
	defer func() { _ = resp.Body.Close() }()

	s.Equal(http.StatusOK, resp.StatusCode)

	var body map[string]interface{}
	s.decodeBody(resp, &body)
	data, _ := body["data"].(map[string]interface{})
	s.Equal("jwt-smoke-token", data["token"])
}

func (s *SmokeSuite) TestAuth_Authenticate_BadJSON() {
	resp, err := s.client.Post(s.url("/auth/authenticate"), "application/json", bytes.NewReader([]byte("{oops")))
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()

	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

// ---------------------------------------------------------------------------
// Anime
// ---------------------------------------------------------------------------

func (s *SmokeSuite) TestAnime_GetAll_HappyPath() {
	animes := []*models.Anime{
		{ID: 1, Title: "Naruto", Episodes: 220, Status: "Completed", Rating: 9.0},
		{ID: 2, Title: "One Piece", Episodes: 1000, Status: "Ongoing", Rating: 9.5},
	}

	s.animeSvc.EXPECT().
		GetAllAnimes(mock.Anything, 1, 10).
		Return(animes, int64(2), nil).Once()

	resp := s.getURL("/animes?page=1&limit=10")
	defer func() { _ = resp.Body.Close() }()

	s.Equal(http.StatusOK, resp.StatusCode)

	var body map[string]interface{}
	s.decodeBody(resp, &body)
	data, _ := body["data"].([]interface{})
	s.Len(data, 2)
	s.Equal(float64(2), body["total"])
}

func (s *SmokeSuite) TestAnime_GetAll_DefaultPagination() {
	s.animeSvc.EXPECT().
		GetAllAnimes(mock.Anything, 1, 100).
		Return([]*models.Anime{}, int64(0), nil).Once()

	resp := s.getURL("/animes")
	defer func() { _ = resp.Body.Close() }()

	s.Equal(http.StatusOK, resp.StatusCode)
}

func (s *SmokeSuite) TestAnime_GetByID_HappyPath() {
	anime := &models.Anime{ID: 1, Title: "Naruto", Episodes: 220, Status: "Completed", Rating: 9.0}

	s.animeSvc.EXPECT().
		GetAnimeByID(mock.Anything, uint(1)).
		Return(anime, nil).Once()

	resp := s.getURL("/animes/1")
	defer func() { _ = resp.Body.Close() }()

	s.Equal(http.StatusOK, resp.StatusCode)

	var body map[string]interface{}
	s.decodeBody(resp, &body)
	s.Equal("Naruto", body["title"])
	s.Equal(float64(1), body["id"])
}

func (s *SmokeSuite) TestAnime_GetByID_NotFound() {
	s.animeSvc.EXPECT().
		GetAnimeByID(mock.Anything, uint(999)).
		Return(nil, apperrors.NewError(
			apperrors.ErrResourceNotFound,
			"Anime not found",
			"No anime with that ID",
			http.StatusNotFound,
			nil, nil,
		)).Once()

	resp := s.getURL("/animes/999")
	defer func() { _ = resp.Body.Close() }()

	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *SmokeSuite) TestAnime_GetByID_InvalidID() {
	resp := s.getURL("/animes/notanumber")
	defer func() { _ = resp.Body.Close() }()

	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

// ---------------------------------------------------------------------------
// Genre
// ---------------------------------------------------------------------------

func (s *SmokeSuite) TestGenre_GetAll_HappyPath() {
	genres := []models.Genre{
		{ID: 1, Name: "Action"},
		{ID: 2, Name: "Adventure"},
	}

	// The handler always calls GetAllGenres(ctx, 1, 10)
	s.genreSvc.EXPECT().
		GetAllGenres(mock.Anything, 1, 10).
		Return(genres, int64(2), nil).Once()

	resp := s.getURL("/genres")
	defer func() { _ = resp.Body.Close() }()

	s.Equal(http.StatusOK, resp.StatusCode)

	var list []interface{}
	s.decodeBody(resp, &list)
	s.Len(list, 2)
}

// ---------------------------------------------------------------------------
// Tag
// ---------------------------------------------------------------------------

func (s *SmokeSuite) TestTag_GetAll_HappyPath() {
	tags := []models.Tag{
		{ID: 1, Name: "Shounen"},
		{ID: 2, Name: "Isekai"},
	}

	// The handler calls GetAllTags with default pagination (page=1, limit=100)
	s.tagSvc.EXPECT().
		GetAllTags(mock.Anything, 1, 100).
		Return(tags, int64(2), nil).Once()

	resp := s.getURL("/tags")
	defer func() { _ = resp.Body.Close() }()

	s.Equal(http.StatusOK, resp.StatusCode)

	var body map[string]interface{}
	s.decodeBody(resp, &body)
	data, _ := body["data"].([]interface{})
	s.Len(data, 2)
}
