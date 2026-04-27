package httphandler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"myanimeapi/api/mocks"
	"myanimeapi/api/models"
	apperrors "myanimeapi/internal/errors"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// ---------------------------------------------------------------------------
// Suite definition
// ---------------------------------------------------------------------------

type GenreHandlerSuite struct {
	suite.Suite
	svc     *mocks.MockGenreServiceInterface
	handler *GenreHandler
}

func (s *GenreHandlerSuite) SetupTest() {
	s.svc = new(mocks.MockGenreServiceInterface)
	s.handler = NewGenreHandler(s.svc)
}

func (s *GenreHandlerSuite) TearDownTest() {
	s.svc.AssertExpectations(s.T())
}

func TestGenreHandlerSuite(t *testing.T) {
	suite.Run(t, new(GenreHandlerSuite))
}

// ---------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------

func (s *GenreHandlerSuite) TestCreate_Success() {
	payload := &models.Genre{Name: "Action"}

	s.svc.EXPECT().CreateGenre(mock.Anything, mock.AnythingOfType("*models.Genre")).
		Run(func(ctx context.Context, g *models.Genre) { g.ID = 1 }).
		Return(nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/genres", jsonBuf(payload))
	ctx := withPayload(testCtx(), payload)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	s.handler.CreateGenreHandler(rr, req)

	s.Equal(http.StatusCreated, rr.Code)
	body := bodyJSON(s.T(), rr)
	s.Equal(float64(1), body["id"])
	s.Equal("Action", body["name"])
}

func (s *GenreHandlerSuite) TestCreate_MissingPayload() {
	req := httptest.NewRequest(http.MethodPost, "/genres", jsonBuf(nil))
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.CreateGenreHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

func (s *GenreHandlerSuite) TestCreate_Conflict() {
	payload := &models.Genre{Name: "Action"}

	s.svc.EXPECT().CreateGenre(mock.Anything, mock.Anything).
		Return(apperrors.NewError(apperrors.ErrConflict, "Genre already exists", "", http.StatusConflict, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodPost, "/genres", jsonBuf(payload))
	ctx := withPayload(testCtx(), payload)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	s.handler.CreateGenreHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusConflict, apperrors.ErrConflict)
}

func (s *GenreHandlerSuite) TestCreate_InternalError() {
	payload := &models.Genre{Name: "Action"}

	s.svc.EXPECT().CreateGenre(mock.Anything, mock.Anything).
		Return(apperrors.NewError(apperrors.ErrInternalServer, "DB error", "", http.StatusInternalServerError, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodPost, "/genres", jsonBuf(payload))
	ctx := withPayload(testCtx(), payload)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	s.handler.CreateGenreHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

// ---------------------------------------------------------------------------
// Get
// ---------------------------------------------------------------------------

func (s *GenreHandlerSuite) TestGet_Success() {
	s.svc.EXPECT().GetGenreByID(mock.Anything, uint(1)).
		Return(&models.Genre{ID: 1, Name: "Action"}, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/genres/1", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.GetGenreHandler(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	body := bodyJSON(s.T(), rr)
	s.Equal(float64(1), body["id"])
}

func (s *GenreHandlerSuite) TestGet_InvalidID() {
	req := httptest.NewRequest(http.MethodGet, "/genres/abc", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rr := httptest.NewRecorder()
	s.handler.GetGenreHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *GenreHandlerSuite) TestGet_NotFound() {
	s.svc.EXPECT().GetGenreByID(mock.Anything, uint(999)).
		Return(nil, apperrors.NewError(apperrors.ErrResourceNotFound, "Genre not found", "", http.StatusNotFound, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodGet, "/genres/999", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	rr := httptest.NewRecorder()
	s.handler.GetGenreHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusNotFound, apperrors.ErrResourceNotFound)
}

func (s *GenreHandlerSuite) TestGet_InternalError() {
	s.svc.EXPECT().GetGenreByID(mock.Anything, uint(1)).
		Return(nil, apperrors.NewError(apperrors.ErrInternalServer, "DB error", "", http.StatusInternalServerError, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodGet, "/genres/1", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.GetGenreHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

// ---------------------------------------------------------------------------
// GetAll  (response is a raw JSON array, NOT wrapped)
// ---------------------------------------------------------------------------

func (s *GenreHandlerSuite) TestGetAll_Success() {
	s.svc.EXPECT().GetAllGenres(mock.Anything, 1, 10).
		Return([]models.Genre{{ID: 1, Name: "Action"}, {ID: 2, Name: "Comedy"}}, int64(2), nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/genres", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.GetAllGenresHandler(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	var arr []interface{}
	require.NoError(s.T(), json.Unmarshal(rr.Body.Bytes(), &arr))
	s.Len(arr, 2)
}

func (s *GenreHandlerSuite) TestGetAll_InternalError() {
	s.svc.EXPECT().GetAllGenres(mock.Anything, 1, 10).
		Return(nil, int64(0), apperrors.NewError(apperrors.ErrInternalServer, "DB error", "", http.StatusInternalServerError, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodGet, "/genres", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.GetAllGenresHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

// ---------------------------------------------------------------------------
// Update
// ---------------------------------------------------------------------------

func (s *GenreHandlerSuite) TestUpdate_Success() {
	payload := &models.Genre{Name: "Updated Action"}

	s.svc.EXPECT().UpdateGenre(mock.Anything, mock.AnythingOfType("*models.Genre")).Return(nil).Once()

	req := httptest.NewRequest(http.MethodPut, "/genres/1", jsonBuf(payload))
	ctx := withPayload(testCtx(), payload)
	req = req.WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.UpdateGenreHandler(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	body := bodyJSON(s.T(), rr)
	s.Equal("Updated Action", body["name"])
}

func (s *GenreHandlerSuite) TestUpdate_InvalidID() {
	payload := &models.Genre{Name: "X"}

	req := httptest.NewRequest(http.MethodPut, "/genres/abc", jsonBuf(payload))
	ctx := withPayload(testCtx(), payload)
	req = req.WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rr := httptest.NewRecorder()
	s.handler.UpdateGenreHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *GenreHandlerSuite) TestUpdate_MissingPayload() {
	req := httptest.NewRequest(http.MethodPut, "/genres/1", jsonBuf(nil))
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.UpdateGenreHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

func (s *GenreHandlerSuite) TestUpdate_NotFound() {
	payload := &models.Genre{Name: "X"}

	s.svc.EXPECT().UpdateGenre(mock.Anything, mock.Anything).
		Return(apperrors.NewError(apperrors.ErrResourceNotFound, "Genre not found", "", http.StatusNotFound, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodPut, "/genres/999", jsonBuf(payload))
	ctx := withPayload(testCtx(), payload)
	req = req.WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	rr := httptest.NewRecorder()
	s.handler.UpdateGenreHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusNotFound, apperrors.ErrResourceNotFound)
}

// ---------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------

func (s *GenreHandlerSuite) TestDelete_Success() {
	s.svc.EXPECT().DeleteGenre(mock.Anything, uint(1)).Return(nil).Once()

	req := httptest.NewRequest(http.MethodDelete, "/genres/1", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.DeleteGenreHandler(rr, req)

	s.Equal(http.StatusNoContent, rr.Code)
	s.Empty(rr.Body.String())
}

func (s *GenreHandlerSuite) TestDelete_InvalidID() {
	req := httptest.NewRequest(http.MethodDelete, "/genres/abc", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rr := httptest.NewRecorder()
	s.handler.DeleteGenreHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *GenreHandlerSuite) TestDelete_NotFound() {
	s.svc.EXPECT().DeleteGenre(mock.Anything, uint(999)).
		Return(apperrors.NewError(apperrors.ErrResourceNotFound, "Genre not found", "", http.StatusNotFound, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodDelete, "/genres/999", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	rr := httptest.NewRecorder()
	s.handler.DeleteGenreHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusNotFound, apperrors.ErrResourceNotFound)
}

func (s *GenreHandlerSuite) TestDelete_InternalError() {
	s.svc.EXPECT().DeleteGenre(mock.Anything, uint(1)).
		Return(apperrors.NewError(apperrors.ErrInternalServer, "DB error", "", http.StatusInternalServerError, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodDelete, "/genres/1", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.DeleteGenreHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

// ---------------------------------------------------------------------------
// Search  (response has a "genres" key with an array)
// ---------------------------------------------------------------------------

func (s *GenreHandlerSuite) TestSearch_Success() {
	s.svc.EXPECT().SearchGenres(mock.Anything, "Action", mock.Anything, mock.Anything).
		Return([]models.Genre{{ID: 1, Name: "Action"}}, int64(1), nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/genres/search?query=Action&page=1&limit=10", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.SearchGenresHandler(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	body := bodyJSON(s.T(), rr)
	s.NotNil(body["genres"])
}

func (s *GenreHandlerSuite) TestSearch_MissingQuery() {
	req := httptest.NewRequest(http.MethodGet, "/genres/search", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.SearchGenresHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *GenreHandlerSuite) TestSearch_InternalError() {
	s.svc.EXPECT().SearchGenres(mock.Anything, "Action", mock.Anything, mock.Anything).
		Return(nil, int64(0), apperrors.NewError(apperrors.ErrInternalServer, "DB error", "", http.StatusInternalServerError, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodGet, "/genres/search?query=Action", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.SearchGenresHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

// ---------------------------------------------------------------------------
// BulkCreate  (payload type BulkCreateGenresRequest defined in genre_handler.go)
// ---------------------------------------------------------------------------

func (s *GenreHandlerSuite) TestBulkCreate_Success() {
	payload := &BulkCreateGenresRequest{Genres: []models.Genre{{Name: "Action"}, {Name: "Comedy"}}}

	s.svc.EXPECT().BulkCreateGenres(mock.Anything, mock.Anything).Return(nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/genres/bulk", jsonBuf(payload))
	ctx := withPayload(testCtx(), payload)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	s.handler.BulkCreateGenresHandler(rr, req)

	s.Equal(http.StatusCreated, rr.Code)
	s.Empty(rr.Body.String())
}

func (s *GenreHandlerSuite) TestBulkCreate_MissingPayload() {
	req := httptest.NewRequest(http.MethodPost, "/genres/bulk", jsonBuf(nil))
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.BulkCreateGenresHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

func (s *GenreHandlerSuite) TestBulkCreate_InternalError() {
	payload := &BulkCreateGenresRequest{Genres: []models.Genre{{Name: "Action"}}}

	s.svc.EXPECT().BulkCreateGenres(mock.Anything, mock.Anything).
		Return(apperrors.NewError(apperrors.ErrInternalServer, "DB error", "", http.StatusInternalServerError, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodPost, "/genres/bulk", jsonBuf(payload))
	ctx := withPayload(testCtx(), payload)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	s.handler.BulkCreateGenresHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

// ---------------------------------------------------------------------------
// BulkDelete  (payload type BulkDeleteGenresRequest defined in genre_handler.go)
// ---------------------------------------------------------------------------

func (s *GenreHandlerSuite) TestBulkDelete_Success() {
	payload := &BulkDeleteGenresRequest{IDs: []uint{1, 2, 3}}

	s.svc.EXPECT().BulkDeleteGenres(mock.Anything, mock.Anything).Return(nil).Once()

	req := httptest.NewRequest(http.MethodDelete, "/genres/bulk", jsonBuf(payload))
	ctx := withPayload(testCtx(), payload)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	s.handler.BulkDeleteGenresHandler(rr, req)

	s.Equal(http.StatusNoContent, rr.Code)
	s.Empty(rr.Body.String())
}

func (s *GenreHandlerSuite) TestBulkDelete_MissingPayload() {
	req := httptest.NewRequest(http.MethodDelete, "/genres/bulk", jsonBuf(nil))
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.BulkDeleteGenresHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

func (s *GenreHandlerSuite) TestBulkDelete_InternalError() {
	payload := &BulkDeleteGenresRequest{IDs: []uint{1}}

	s.svc.EXPECT().BulkDeleteGenres(mock.Anything, mock.Anything).
		Return(apperrors.NewError(apperrors.ErrInternalServer, "DB error", "", http.StatusInternalServerError, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodDelete, "/genres/bulk", jsonBuf(payload))
	ctx := withPayload(testCtx(), payload)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	s.handler.BulkDeleteGenresHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}
