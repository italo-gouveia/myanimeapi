package httphandler

import (
	"encoding/json"
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

type FavoriteHandlerSuite struct {
	suite.Suite
	svc     *mocks.MockFavoriteServiceInterface
	handler *FavoriteHandler
}

func (s *FavoriteHandlerSuite) SetupTest() {
	s.svc = new(mocks.MockFavoriteServiceInterface)
	s.handler = NewFavoriteHandler(s.svc)
}

func (s *FavoriteHandlerSuite) TearDownTest() {
	s.svc.AssertExpectations(s.T())
}

func TestFavoriteHandlerSuite(t *testing.T) {
	suite.Run(t, new(FavoriteHandlerSuite))
}

// --- Add ---

func (s *FavoriteHandlerSuite) TestAdd_Success() {
	s.svc.EXPECT().AddFavorite(mock.Anything, uint(1), uint(1)).
		Return(&models.Favorite{ID: 1, UserID: 1, AnimeID: 1}, nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/favorites/1", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"anime_id": "1"})
	rr := httptest.NewRecorder()
	s.handler.AddFavoriteHandler(rr, req)

	s.Equal(http.StatusCreated, rr.Code)
	body := bodyJSON(s.T(), rr)
	s.Equal(float64(1), body["id"])
}

func (s *FavoriteHandlerSuite) TestAdd_InvalidAnimeID() {
	req := httptest.NewRequest(http.MethodPost, "/favorites/abc", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"anime_id": "abc"})
	rr := httptest.NewRecorder()
	s.handler.AddFavoriteHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *FavoriteHandlerSuite) TestAdd_AlreadyFavorited() {
	s.svc.EXPECT().AddFavorite(mock.Anything, uint(1), uint(1)).
		Return(nil, apperrors.NewError(apperrors.ErrConflict, "Already favorited", "", http.StatusConflict, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodPost, "/favorites/1", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"anime_id": "1"})
	rr := httptest.NewRecorder()
	s.handler.AddFavoriteHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusConflict, apperrors.ErrConflict)
}

func (s *FavoriteHandlerSuite) TestAdd_InternalError() {
	s.svc.EXPECT().AddFavorite(mock.Anything, uint(1), uint(1)).
		Return(nil, apperrors.NewError(apperrors.ErrInternalServer, "DB error", "", http.StatusInternalServerError, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodPost, "/favorites/1", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"anime_id": "1"})
	rr := httptest.NewRecorder()
	s.handler.AddFavoriteHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

func (s *FavoriteHandlerSuite) TestAdd_Unauthenticated() {
	req := httptest.NewRequest(http.MethodPost, "/favorites/1", nil)
	req = req.WithContext(noAuthCtx())
	req = mux.SetURLVars(req, map[string]string{"anime_id": "1"})
	rr := httptest.NewRecorder()
	s.handler.AddFavoriteHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusUnauthorized, apperrors.ErrUnauthorized)
}

// --- Remove ---

func (s *FavoriteHandlerSuite) TestRemove_Success() {
	s.svc.EXPECT().RemoveFavorite(mock.Anything, uint(1), uint(1)).Return(nil).Once()

	req := httptest.NewRequest(http.MethodDelete, "/favorites/1", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"anime_id": "1"})
	rr := httptest.NewRecorder()
	s.handler.RemoveFavoriteHandler(rr, req)

	s.Equal(http.StatusNoContent, rr.Code)
	s.Empty(rr.Body.String())
}

func (s *FavoriteHandlerSuite) TestRemove_InvalidAnimeID() {
	req := httptest.NewRequest(http.MethodDelete, "/favorites/abc", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"anime_id": "abc"})
	rr := httptest.NewRecorder()
	s.handler.RemoveFavoriteHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *FavoriteHandlerSuite) TestRemove_NotFound() {
	s.svc.EXPECT().RemoveFavorite(mock.Anything, uint(1), uint(999)).
		Return(apperrors.NewError(apperrors.ErrResourceNotFound, "Favorite not found", "", http.StatusNotFound, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodDelete, "/favorites/999", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"anime_id": "999"})
	rr := httptest.NewRecorder()
	s.handler.RemoveFavoriteHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusNotFound, apperrors.ErrResourceNotFound)
}

func (s *FavoriteHandlerSuite) TestRemove_InternalError() {
	s.svc.EXPECT().RemoveFavorite(mock.Anything, uint(1), uint(1)).
		Return(apperrors.NewError(apperrors.ErrInternalServer, "DB error", "", http.StatusInternalServerError, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodDelete, "/favorites/1", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"anime_id": "1"})
	rr := httptest.NewRecorder()
	s.handler.RemoveFavoriteHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

// --- GetFavorites ---

func (s *FavoriteHandlerSuite) TestGetFavorites_Success() {
	s.svc.EXPECT().GetFavorites(mock.Anything, uint(1)).
		Return([]models.Favorite{
			{ID: 1, UserID: 1, AnimeID: 1},
			{ID: 2, UserID: 1, AnimeID: 2},
		}, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/favorites", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.GetFavoritesHandler(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	var data []interface{}
	s.NoError(json.Unmarshal(rr.Body.Bytes(), &data))
	s.Len(data, 2)
}

func (s *FavoriteHandlerSuite) TestGetFavorites_Empty() {
	s.svc.EXPECT().GetFavorites(mock.Anything, uint(1)).
		Return([]models.Favorite{}, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/favorites", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.GetFavoritesHandler(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	var data []interface{}
	s.NoError(json.Unmarshal(rr.Body.Bytes(), &data))
	s.Len(data, 0)
}

func (s *FavoriteHandlerSuite) TestGetFavorites_InternalError() {
	s.svc.EXPECT().GetFavorites(mock.Anything, uint(1)).
		Return(nil, apperrors.NewError(apperrors.ErrInternalServer, "DB error", "", http.StatusInternalServerError, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodGet, "/favorites", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.GetFavoritesHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

func (s *FavoriteHandlerSuite) TestGetFavorites_Unauthenticated() {
	req := httptest.NewRequest(http.MethodGet, "/favorites", nil)
	req = req.WithContext(noAuthCtx())
	rr := httptest.NewRecorder()
	s.handler.GetFavoritesHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusUnauthorized, apperrors.ErrUnauthorized)
}
