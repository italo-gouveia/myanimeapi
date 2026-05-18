package httphandler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"myanimeapi/api/mocks"
	"myanimeapi/api/models"
	apperrors "myanimeapi/internal/errors"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// ---------------------------------------------------------------------------
// Suite definition
// ---------------------------------------------------------------------------

type AnimeHandlerSuite struct {
	suite.Suite
	svc     *mocks.MockAnimeServiceInterface
	handler *AnimeHandler
}

func (s *AnimeHandlerSuite) SetupTest() {
	s.svc = new(mocks.MockAnimeServiceInterface)
	s.handler = NewAnimeHandler(s.svc)
}

func (s *AnimeHandlerSuite) TearDownTest() {
	s.svc.AssertExpectations(s.T())
}

func TestAnimeHandlerSuite(t *testing.T) {
	suite.Run(t, new(AnimeHandlerSuite))
}

// ---------------------------------------------------------------------------
// Get
// ---------------------------------------------------------------------------

func (s *AnimeHandlerSuite) TestGet_Success() {
	s.svc.EXPECT().GetAnimeByID(mock.Anything, uint(1)).
		Return(&models.Anime{ID: 1, Title: "Naruto", Status: "Completed"}, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/animes/1", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.GetAnimeHandler(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	body := bodyJSON(s.T(), rr)
	s.Equal(float64(1), body["id"])
	s.Equal("Naruto", body["title"])
}

func (s *AnimeHandlerSuite) TestGet_InvalidID() {
	req := httptest.NewRequest(http.MethodGet, "/animes/abc", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rr := httptest.NewRecorder()
	s.handler.GetAnimeHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *AnimeHandlerSuite) TestGet_NotFound() {
	s.svc.EXPECT().GetAnimeByID(mock.Anything, uint(999)).
		Return(nil, apperrors.NewError(apperrors.ErrResourceNotFound, "Anime not found", "", http.StatusNotFound, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodGet, "/animes/999", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	rr := httptest.NewRecorder()
	s.handler.GetAnimeHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusNotFound, apperrors.ErrResourceNotFound)
}

func (s *AnimeHandlerSuite) TestGet_InternalError() {
	s.svc.EXPECT().GetAnimeByID(mock.Anything, uint(1)).
		Return(nil, apperrors.NewError(apperrors.ErrInternalServer, "DB error", "", http.StatusInternalServerError, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodGet, "/animes/1", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.GetAnimeHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

// ---------------------------------------------------------------------------
// GetAll
// ---------------------------------------------------------------------------

func (s *AnimeHandlerSuite) TestGetAll_Success() {
	s.svc.EXPECT().GetAllAnimes(mock.Anything, 1, 10, models.AnimeFilter{}).
		Return([]*models.Anime{
			{ID: 1, Title: "Naruto"},
			{ID: 2, Title: "One Piece"},
		}, int64(2), nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/animes?page=1&limit=10", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.GetAllAnimesHandler(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	body := bodyJSON(s.T(), rr)
	s.Equal(float64(2), body["total"])
	s.NotNil(body["data"])
}

func (s *AnimeHandlerSuite) TestGetAll_WithFilters() {
	expected := models.AnimeFilter{
		Status:    "Completed",
		Genre:     "Action",
		RatingMin: 7.5,
		SortBy:    "rating",
		SortOrder: "desc",
	}
	s.svc.EXPECT().GetAllAnimes(mock.Anything, 1, 10, expected).
		Return([]*models.Anime{{ID: 1, Title: "Naruto", Status: "Completed", Rating: 8.5}}, int64(1), nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/animes?page=1&limit=10&status=Completed&genre=Action&rating_min=7.5&sort_by=rating&order=desc", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.GetAllAnimesHandler(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	body := bodyJSON(s.T(), rr)
	s.Equal(float64(1), body["total"])
}

func (s *AnimeHandlerSuite) TestGetAll_InvalidFilter() {
	req := httptest.NewRequest(http.MethodGet, "/animes?status=Invalid", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.GetAllAnimesHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *AnimeHandlerSuite) TestGetAll_InvalidPagination() {
	req := httptest.NewRequest(http.MethodGet, "/animes?page=0&limit=10", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.GetAllAnimesHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *AnimeHandlerSuite) TestGetAll_InternalError() {
	s.svc.EXPECT().GetAllAnimes(mock.Anything, 1, 10, models.AnimeFilter{}).
		Return(nil, int64(0), apperrors.NewError(apperrors.ErrInternalServer, "DB error", "", http.StatusInternalServerError, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodGet, "/animes?page=1&limit=10", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.GetAllAnimesHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

// ---------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------

func (s *AnimeHandlerSuite) TestCreate_Success() {
	payload := &models.AnimeCreateRequest{
		Title:    "New Anime",
		Episodes: 12,
		Status:   "Completed",
		Rating:   8.5,
		GenreIDs: []uint{1},
		TagIDs:   []uint{1},
	}

	s.svc.EXPECT().CreateAnime(mock.Anything, mock.AnythingOfType("*models.Anime")).
		Run(func(ctx context.Context, a *models.Anime) { a.ID = 1 }).
		Return(nil).Once()
	s.svc.EXPECT().AddGenresToAnime(mock.Anything, uint(1), []uint{1}).Return(nil).Once()
	s.svc.EXPECT().AddTagsToAnime(mock.Anything, uint(1), []uint{1}).Return(nil).Once()
	s.svc.EXPECT().GetAnimeByID(mock.Anything, uint(1)).
		Return(&models.Anime{ID: 1, Title: "New Anime"}, nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/animes", jsonBuf(payload))
	ctx := withPayload(adminCtx(), payload)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	s.handler.CreateAnimeHandler(rr, req)

	s.Equal(http.StatusCreated, rr.Code)
	body := bodyJSON(s.T(), rr)
	s.Equal(float64(1), body["id"])
}

func (s *AnimeHandlerSuite) TestCreate_MissingPayload() {
	req := httptest.NewRequest(http.MethodPost, "/animes", jsonBuf(nil))
	req = req.WithContext(adminCtx())
	rr := httptest.NewRecorder()
	s.handler.CreateAnimeHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

func (s *AnimeHandlerSuite) TestCreate_ServiceError() {
	payload := &models.AnimeCreateRequest{
		Title:    "New Anime",
		Episodes: 12,
		Status:   "Completed",
		Rating:   8.5,
		GenreIDs: []uint{1},
		TagIDs:   []uint{1},
	}

	s.svc.EXPECT().CreateAnime(mock.Anything, mock.Anything).
		Return(apperrors.NewError(apperrors.ErrInternalServer, "DB error", "", http.StatusInternalServerError, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodPost, "/animes", jsonBuf(payload))
	ctx := withPayload(adminCtx(), payload)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	s.handler.CreateAnimeHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

// ---------------------------------------------------------------------------
// Update
// ---------------------------------------------------------------------------

func (s *AnimeHandlerSuite) TestUpdate_Success() {
	payload := &models.AnimeUpdateRequest{Title: "Updated Title", Status: "Ongoing"}

	s.svc.EXPECT().GetAnimeByID(mock.Anything, uint(1)).
		Return(&models.Anime{ID: 1, Title: "Old Title", Genres: []models.Genre{}, Tags: []models.Tag{}}, nil).Once()
	s.svc.EXPECT().UpdateAnime(mock.Anything, mock.Anything).Return(nil).Once()

	req := httptest.NewRequest(http.MethodPut, "/animes/1", jsonBuf(payload))
	ctx := withPayload(adminCtx(), payload)
	req = req.WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.UpdateAnimeHandler(rr, req)

	s.Equal(http.StatusOK, rr.Code)
}

func (s *AnimeHandlerSuite) TestUpdate_InvalidID() {
	payload := &models.AnimeUpdateRequest{Title: "Updated Title", Status: "Ongoing"}

	req := httptest.NewRequest(http.MethodPut, "/animes/abc", jsonBuf(payload))
	ctx := withPayload(adminCtx(), payload)
	req = req.WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rr := httptest.NewRecorder()
	s.handler.UpdateAnimeHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *AnimeHandlerSuite) TestUpdate_MissingPayload() {
	req := httptest.NewRequest(http.MethodPut, "/animes/1", jsonBuf(nil))
	req = req.WithContext(adminCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.UpdateAnimeHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

func (s *AnimeHandlerSuite) TestUpdate_NotFound() {
	payload := &models.AnimeUpdateRequest{Title: "Updated Title", Status: "Ongoing"}

	s.svc.EXPECT().GetAnimeByID(mock.Anything, uint(999)).
		Return(nil, apperrors.NewError(apperrors.ErrResourceNotFound, "Anime not found", "", http.StatusNotFound, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodPut, "/animes/999", jsonBuf(payload))
	ctx := withPayload(adminCtx(), payload)
	req = req.WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	rr := httptest.NewRecorder()
	s.handler.UpdateAnimeHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusNotFound, apperrors.ErrResourceNotFound)
}

// ---------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------

func (s *AnimeHandlerSuite) TestDelete_Success() {
	s.svc.EXPECT().DeleteAnime(mock.Anything, uint(1)).Return(nil).Once()

	req := httptest.NewRequest(http.MethodDelete, "/animes/1", nil)
	req = req.WithContext(adminCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.DeleteAnimeHandler(rr, req)

	s.Equal(http.StatusNoContent, rr.Code)
	s.Empty(rr.Body.String())
}

func (s *AnimeHandlerSuite) TestDelete_InvalidID() {
	req := httptest.NewRequest(http.MethodDelete, "/animes/abc", nil)
	req = req.WithContext(adminCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rr := httptest.NewRecorder()
	s.handler.DeleteAnimeHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *AnimeHandlerSuite) TestDelete_NotFound() {
	s.svc.EXPECT().DeleteAnime(mock.Anything, uint(999)).
		Return(apperrors.NewError(apperrors.ErrResourceNotFound, "Anime not found", "", http.StatusNotFound, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodDelete, "/animes/999", nil)
	req = req.WithContext(adminCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	rr := httptest.NewRecorder()
	s.handler.DeleteAnimeHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusNotFound, apperrors.ErrResourceNotFound)
}

// ---------------------------------------------------------------------------
// SearchByTitle
// ---------------------------------------------------------------------------

func (s *AnimeHandlerSuite) TestSearchByTitle_Success() {
	s.svc.EXPECT().GetAllAnimes(mock.Anything, 1, 10, models.AnimeFilter{Title: "Naruto"}).
		Return([]*models.Anime{{ID: 1, Title: "Naruto"}}, int64(1), nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/animes/search?title=Naruto&page=1&limit=10", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.GetAnimesByTitleHandler(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	body := bodyJSON(s.T(), rr)
	s.NotNil(body["data"])
}

func (s *AnimeHandlerSuite) TestSearchByTitle_NoResults() {
	s.svc.EXPECT().GetAllAnimes(mock.Anything, 1, 10, models.AnimeFilter{Title: "Unknown"}).
		Return([]*models.Anime{}, int64(0), nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/animes/search?title=Unknown&page=1&limit=10", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.GetAnimesByTitleHandler(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	body := bodyJSON(s.T(), rr)
	s.NotNil(body["data"])
}

func (s *AnimeHandlerSuite) TestSearchByTitle_InternalError() {
	s.svc.EXPECT().GetAllAnimes(mock.Anything, 1, 10, models.AnimeFilter{Title: "Naruto"}).
		Return(nil, int64(0), apperrors.NewError(apperrors.ErrInternalServer, "DB error", "", http.StatusInternalServerError, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodGet, "/animes/search?title=Naruto&page=1&limit=10", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.GetAnimesByTitleHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

func (s *AnimeHandlerSuite) TestSearchByTitle_WithSort() {
	s.svc.EXPECT().GetAllAnimes(mock.Anything, 1, 10, models.AnimeFilter{Title: "Naruto", SortBy: "rating", SortOrder: "desc"}).
		Return([]*models.Anime{{ID: 1, Title: "Naruto", Rating: 9.0}}, int64(1), nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/animes/search?title=Naruto&page=1&limit=10&sort_by=rating&order=desc", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.GetAnimesByTitleHandler(rr, req)

	s.Equal(http.StatusOK, rr.Code)
}

func (s *AnimeHandlerSuite) TestSearchByTitle_InvalidSort() {
	req := httptest.NewRequest(http.MethodGet, "/animes/search?title=Naruto&sort_by=invalid_col", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.GetAnimesByTitleHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

// ---------------------------------------------------------------------------
// AddGenres
// ---------------------------------------------------------------------------

func (s *AnimeHandlerSuite) TestAddGenres_Success() {
	payload := &models.AnimeGenresRequest{GenreIDs: []uint{1, 2}}

	s.svc.EXPECT().AddGenresToAnime(mock.Anything, uint(1), []uint{1, 2}).Return(nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/animes/1/genres", jsonBuf(payload))
	ctx := withPayload(adminCtx(), payload)
	req = req.WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.AddGenresToAnimeHandler(rr, req)

	s.Equal(http.StatusNoContent, rr.Code)
	s.Empty(rr.Body.String())
}

func (s *AnimeHandlerSuite) TestAddGenres_InvalidID() {
	payload := &models.AnimeGenresRequest{GenreIDs: []uint{1}}

	req := httptest.NewRequest(http.MethodPost, "/animes/abc/genres", jsonBuf(payload))
	ctx := withPayload(adminCtx(), payload)
	req = req.WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rr := httptest.NewRecorder()
	s.handler.AddGenresToAnimeHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *AnimeHandlerSuite) TestAddGenres_MissingPayload() {
	req := httptest.NewRequest(http.MethodPost, "/animes/1/genres", jsonBuf(nil))
	req = req.WithContext(adminCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.AddGenresToAnimeHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

func (s *AnimeHandlerSuite) TestAddGenres_ServiceError() {
	payload := &models.AnimeGenresRequest{GenreIDs: []uint{1}}
	s.svc.EXPECT().AddGenresToAnime(mock.Anything, uint(1), []uint{1}).
		Return(apperrors.NewError(apperrors.ErrInternalServer, "DB error", "", http.StatusInternalServerError, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodPost, "/animes/1/genres", jsonBuf(payload))
	ctx := withPayload(adminCtx(), payload)
	req = req.WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.AddGenresToAnimeHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

// ---------------------------------------------------------------------------
// RemoveGenres
// ---------------------------------------------------------------------------

func (s *AnimeHandlerSuite) TestRemoveGenres_Success() {
	payload := &models.AnimeGenresRequest{GenreIDs: []uint{1, 2}}
	s.svc.EXPECT().RemoveGenresFromAnime(mock.Anything, uint(1), []uint{1, 2}).Return(nil).Once()

	req := httptest.NewRequest(http.MethodDelete, "/animes/1/genres", jsonBuf(payload))
	ctx := withPayload(adminCtx(), payload)
	req = req.WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.RemoveGenresFromAnimeHandler(rr, req)

	s.Equal(http.StatusNoContent, rr.Code)
}

func (s *AnimeHandlerSuite) TestRemoveGenres_InvalidID() {
	payload := &models.AnimeGenresRequest{GenreIDs: []uint{1}}

	req := httptest.NewRequest(http.MethodDelete, "/animes/abc/genres", jsonBuf(payload))
	ctx := withPayload(adminCtx(), payload)
	req = req.WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rr := httptest.NewRecorder()
	s.handler.RemoveGenresFromAnimeHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *AnimeHandlerSuite) TestRemoveGenres_MissingPayload() {
	req := httptest.NewRequest(http.MethodDelete, "/animes/1/genres", jsonBuf(nil))
	req = req.WithContext(adminCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.RemoveGenresFromAnimeHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

func (s *AnimeHandlerSuite) TestRemoveGenres_ServiceError() {
	payload := &models.AnimeGenresRequest{GenreIDs: []uint{1}}
	s.svc.EXPECT().RemoveGenresFromAnime(mock.Anything, uint(1), []uint{1}).
		Return(apperrors.NewError(apperrors.ErrResourceNotFound, "Anime not found", "", http.StatusNotFound, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodDelete, "/animes/1/genres", jsonBuf(payload))
	ctx := withPayload(adminCtx(), payload)
	req = req.WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.RemoveGenresFromAnimeHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusNotFound, apperrors.ErrResourceNotFound)
}

// ---------------------------------------------------------------------------
// AddTags
// ---------------------------------------------------------------------------

func (s *AnimeHandlerSuite) TestAddTags_Success() {
	payload := &models.AnimeTagsRequest{TagIDs: []uint{1, 2}}
	s.svc.EXPECT().AddTagsToAnime(mock.Anything, uint(1), []uint{1, 2}).Return(nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/animes/1/tags", jsonBuf(payload))
	ctx := withPayload(adminCtx(), payload)
	req = req.WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.AddTagsToAnimeHandler(rr, req)

	s.Equal(http.StatusNoContent, rr.Code)
}

func (s *AnimeHandlerSuite) TestAddTags_InvalidID() {
	payload := &models.AnimeTagsRequest{TagIDs: []uint{1}}

	req := httptest.NewRequest(http.MethodPost, "/animes/abc/tags", jsonBuf(payload))
	ctx := withPayload(adminCtx(), payload)
	req = req.WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rr := httptest.NewRecorder()
	s.handler.AddTagsToAnimeHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *AnimeHandlerSuite) TestAddTags_MissingPayload() {
	req := httptest.NewRequest(http.MethodPost, "/animes/1/tags", jsonBuf(nil))
	req = req.WithContext(adminCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.AddTagsToAnimeHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

func (s *AnimeHandlerSuite) TestAddTags_ServiceError() {
	payload := &models.AnimeTagsRequest{TagIDs: []uint{1}}
	s.svc.EXPECT().AddTagsToAnime(mock.Anything, uint(1), []uint{1}).
		Return(apperrors.NewError(apperrors.ErrInternalServer, "DB error", "", http.StatusInternalServerError, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodPost, "/animes/1/tags", jsonBuf(payload))
	ctx := withPayload(adminCtx(), payload)
	req = req.WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.AddTagsToAnimeHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

// ---------------------------------------------------------------------------
// RemoveTags
// ---------------------------------------------------------------------------

func (s *AnimeHandlerSuite) TestRemoveTags_Success() {
	payload := &models.AnimeTagsRequest{TagIDs: []uint{1, 2}}
	s.svc.EXPECT().RemoveTagsFromAnime(mock.Anything, uint(1), []uint{1, 2}).Return(nil).Once()

	req := httptest.NewRequest(http.MethodDelete, "/animes/1/tags", jsonBuf(payload))
	ctx := withPayload(adminCtx(), payload)
	req = req.WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.RemoveTagsFromAnimeHandler(rr, req)

	s.Equal(http.StatusNoContent, rr.Code)
}

func (s *AnimeHandlerSuite) TestRemoveTags_InvalidID() {
	payload := &models.AnimeTagsRequest{TagIDs: []uint{1}}

	req := httptest.NewRequest(http.MethodDelete, "/animes/abc/tags", jsonBuf(payload))
	ctx := withPayload(adminCtx(), payload)
	req = req.WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rr := httptest.NewRecorder()
	s.handler.RemoveTagsFromAnimeHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *AnimeHandlerSuite) TestRemoveTags_MissingPayload() {
	req := httptest.NewRequest(http.MethodDelete, "/animes/1/tags", jsonBuf(nil))
	req = req.WithContext(adminCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.RemoveTagsFromAnimeHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

func (s *AnimeHandlerSuite) TestRemoveTags_ServiceError() {
	payload := &models.AnimeTagsRequest{TagIDs: []uint{1}}
	s.svc.EXPECT().RemoveTagsFromAnime(mock.Anything, uint(1), []uint{1}).
		Return(apperrors.NewError(apperrors.ErrResourceNotFound, "Anime not found", "", http.StatusNotFound, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodDelete, "/animes/1/tags", jsonBuf(payload))
	ctx := withPayload(adminCtx(), payload)
	req = req.WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.RemoveTagsFromAnimeHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusNotFound, apperrors.ErrResourceNotFound)
}

// ---------------------------------------------------------------------------
// GetByGenre
// ---------------------------------------------------------------------------

func (s *AnimeHandlerSuite) TestGetByGenre_Success() {
	s.svc.EXPECT().GetAnimesByGenre(mock.Anything, "Action", 1, 10).
		Return([]*models.Anime{{ID: 1, Title: "Naruto"}, {ID: 2, Title: "Bleach"}}, int64(2), nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/animes/genre/Action?page=1&limit=10", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"genre": "Action"})
	rr := httptest.NewRecorder()
	s.handler.GetAnimesByGenreHandler(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	body := bodyJSON(s.T(), rr)
	s.Equal(float64(2), body["total"])
}

func (s *AnimeHandlerSuite) TestGetByGenre_InvalidPagination() {
	req := httptest.NewRequest(http.MethodGet, "/animes/genre/Action?page=abc", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"genre": "Action"})
	rr := httptest.NewRecorder()
	s.handler.GetAnimesByGenreHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *AnimeHandlerSuite) TestGetByGenre_ServiceError() {
	s.svc.EXPECT().GetAnimesByGenre(mock.Anything, "Action", 1, 10).
		Return(nil, int64(0), apperrors.NewError(apperrors.ErrInternalServer, "DB error", "", http.StatusInternalServerError, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodGet, "/animes/genre/Action?page=1&limit=10", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"genre": "Action"})
	rr := httptest.NewRecorder()
	s.handler.GetAnimesByGenreHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}
