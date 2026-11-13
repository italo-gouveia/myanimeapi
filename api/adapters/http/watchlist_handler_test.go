package httphandler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"myanimeapi/api/mocks"
	"myanimeapi/api/models"
	apperrors "myanimeapi/internal/errors"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type WatchlistHandlerSuite struct {
	suite.Suite
	svc     *mocks.MockWatchlistServiceInterface
	handler *WatchlistHandler
}

func (s *WatchlistHandlerSuite) SetupTest() {
	s.svc = new(mocks.MockWatchlistServiceInterface)
	s.handler = NewWatchlistHandler(s.svc)
}

func (s *WatchlistHandlerSuite) TearDownTest() {
	s.svc.AssertExpectations(s.T())
}

func TestWatchlistHandlerSuite(t *testing.T) {
	suite.Run(t, new(WatchlistHandlerSuite))
}

// ── GetWatchlistHandler ───────────────────────────────────────────────────────

func (s *WatchlistHandlerSuite) TestGetWatchlist_Success_Empty() {
	s.svc.EXPECT().GetWatchlist(mock.Anything, uint(1)).Return([]models.WatchlistEntry{}, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/watchlist", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.GetWatchlistHandler(rr, req)

	s.Equal(http.StatusOK, rr.Code)
}

func (s *WatchlistHandlerSuite) TestGetWatchlist_Success_WithEntries() {
	entries := []models.WatchlistEntry{
		{UserID: 1, AnimeID: 10, Status: models.WatchlistStatusWatching},
		{UserID: 1, AnimeID: 20, Status: models.WatchlistStatusCompleted},
	}
	s.svc.EXPECT().GetWatchlist(mock.Anything, uint(1)).Return(entries, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/watchlist", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.GetWatchlistHandler(rr, req)

	s.Equal(http.StatusOK, rr.Code)
}

func (s *WatchlistHandlerSuite) TestGetWatchlist_Unauthenticated() {
	req := httptest.NewRequest(http.MethodGet, "/watchlist", nil)
	req = req.WithContext(noAuthCtx())
	rr := httptest.NewRecorder()
	s.handler.GetWatchlistHandler(rr, req)

	s.Equal(http.StatusUnauthorized, rr.Code)
}

func (s *WatchlistHandlerSuite) TestGetWatchlist_ServiceError() {
	s.svc.EXPECT().GetWatchlist(mock.Anything, uint(1)).
		Return(nil, apperrors.NewError(apperrors.ErrInternalServer, "db err", "", http.StatusInternalServerError, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodGet, "/watchlist", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.GetWatchlistHandler(rr, req)

	s.Equal(http.StatusInternalServerError, rr.Code)
}

func (s *WatchlistHandlerSuite) TestGetWatchlist_ServiceError_Generic() {
	s.svc.EXPECT().GetWatchlist(mock.Anything, uint(1)).
		Return(nil, apperrors.NewError(apperrors.ErrInternalServer, "db generic", "", http.StatusInternalServerError, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodGet, "/watchlist", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.GetWatchlistHandler(rr, req)

	s.Equal(http.StatusInternalServerError, rr.Code)
}

// ── UpsertWatchlistHandler ────────────────────────────────────────────────────

func (s *WatchlistHandlerSuite) TestUpsertWatchlist_Success() {
	entry := &models.WatchlistEntry{UserID: 1, AnimeID: 5, Status: models.WatchlistStatusWatching}
	s.svc.EXPECT().UpsertWatchlistEntry(mock.Anything, uint(1), uint(5), models.WatchlistStatusWatching).
		Return(entry, nil).Once()

	body := `{"status":"watching"}`
	req := httptest.NewRequest(http.MethodPut, "/watchlist/5", strings.NewReader(body))
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"animeId": "5"})
	rr := httptest.NewRecorder()
	s.handler.UpsertWatchlistHandler(rr, req)

	s.Equal(http.StatusOK, rr.Code)
}

func (s *WatchlistHandlerSuite) TestUpsertWatchlist_Unauthenticated() {
	req := httptest.NewRequest(http.MethodPut, "/watchlist/5", strings.NewReader(`{"status":"watching"}`))
	req = req.WithContext(noAuthCtx())
	req = mux.SetURLVars(req, map[string]string{"animeId": "5"})
	rr := httptest.NewRecorder()
	s.handler.UpsertWatchlistHandler(rr, req)

	s.Equal(http.StatusUnauthorized, rr.Code)
}

func (s *WatchlistHandlerSuite) TestUpsertWatchlist_InvalidAnimeID() {
	req := httptest.NewRequest(http.MethodPut, "/watchlist/bad", strings.NewReader(`{"status":"watching"}`))
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"animeId": "bad"})
	rr := httptest.NewRecorder()
	s.handler.UpsertWatchlistHandler(rr, req)

	s.Equal(http.StatusBadRequest, rr.Code)
}

func (s *WatchlistHandlerSuite) TestUpsertWatchlist_BadJSON() {
	req := httptest.NewRequest(http.MethodPut, "/watchlist/5", strings.NewReader("{bad"))
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"animeId": "5"})
	rr := httptest.NewRecorder()
	s.handler.UpsertWatchlistHandler(rr, req)

	s.Equal(http.StatusBadRequest, rr.Code)
}

func (s *WatchlistHandlerSuite) TestUpsertWatchlist_InvalidStatus() {
	req := httptest.NewRequest(http.MethodPut, "/watchlist/5", strings.NewReader(`{"status":"invalid_status"}`))
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"animeId": "5"})
	rr := httptest.NewRecorder()
	s.handler.UpsertWatchlistHandler(rr, req)

	s.Equal(http.StatusBadRequest, rr.Code)
}

func (s *WatchlistHandlerSuite) TestUpsertWatchlist_ServiceError() {
	s.svc.EXPECT().UpsertWatchlistEntry(mock.Anything, uint(1), uint(5), models.WatchlistStatusWatching).
		Return(nil, apperrors.NewError(apperrors.ErrInternalServer, "db err", "", http.StatusInternalServerError, nil, nil)).Once()

	body := `{"status":"watching"}`
	req := httptest.NewRequest(http.MethodPut, "/watchlist/5", strings.NewReader(body))
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"animeId": "5"})
	rr := httptest.NewRecorder()
	s.handler.UpsertWatchlistHandler(rr, req)

	s.Equal(http.StatusInternalServerError, rr.Code)
}

// ── RemoveWatchlistHandler ────────────────────────────────────────────────────

func (s *WatchlistHandlerSuite) TestRemoveWatchlist_Success() {
	s.svc.EXPECT().RemoveWatchlistEntry(mock.Anything, uint(1), uint(5)).Return(nil).Once()

	req := httptest.NewRequest(http.MethodDelete, "/watchlist/5", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"animeId": "5"})
	rr := httptest.NewRecorder()
	s.handler.RemoveWatchlistHandler(rr, req)

	s.Equal(http.StatusNoContent, rr.Code)
}

func (s *WatchlistHandlerSuite) TestRemoveWatchlist_Unauthenticated() {
	req := httptest.NewRequest(http.MethodDelete, "/watchlist/5", nil)
	req = req.WithContext(noAuthCtx())
	req = mux.SetURLVars(req, map[string]string{"animeId": "5"})
	rr := httptest.NewRecorder()
	s.handler.RemoveWatchlistHandler(rr, req)

	s.Equal(http.StatusUnauthorized, rr.Code)
}

func (s *WatchlistHandlerSuite) TestRemoveWatchlist_InvalidAnimeID() {
	req := httptest.NewRequest(http.MethodDelete, "/watchlist/bad", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"animeId": "bad"})
	rr := httptest.NewRecorder()
	s.handler.RemoveWatchlistHandler(rr, req)

	s.Equal(http.StatusBadRequest, rr.Code)
}

func (s *WatchlistHandlerSuite) TestRemoveWatchlist_ServiceError() {
	s.svc.EXPECT().RemoveWatchlistEntry(mock.Anything, uint(1), uint(5)).
		Return(apperrors.NewError(apperrors.ErrResourceNotFound, "not found", "", http.StatusNotFound, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodDelete, "/watchlist/5", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"animeId": "5"})
	rr := httptest.NewRecorder()
	s.handler.RemoveWatchlistHandler(rr, req)

	s.Equal(http.StatusNotFound, rr.Code)
}
