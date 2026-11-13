package httphandler

import (
	"context"
	"fmt"
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

type CharacterHandlerSuite struct {
	suite.Suite
	svc     *mocks.MockCharacterServiceInterface
	handler *CharacterHandler
}

func (s *CharacterHandlerSuite) SetupTest() {
	s.svc = new(mocks.MockCharacterServiceInterface)
	s.handler = NewCharacterHandler(s.svc)
}

func (s *CharacterHandlerSuite) TearDownTest() {
	s.svc.AssertExpectations(s.T())
}

func TestCharacterHandlerSuite(t *testing.T) {
	suite.Run(t, new(CharacterHandlerSuite))
}

// ── paginationParams ──────────────────────────────────────────────────────────

func (s *CharacterHandlerSuite) TestPaginationParams_Defaults() {
	req := httptest.NewRequest(http.MethodGet, "/characters", nil)
	page, limit := paginationParams(req)
	s.Equal(1, page)
	s.Equal(10, limit)
}

func (s *CharacterHandlerSuite) TestPaginationParams_Custom() {
	req := httptest.NewRequest(http.MethodGet, "/characters?page=3&limit=25", nil)
	page, limit := paginationParams(req)
	s.Equal(3, page)
	s.Equal(25, limit)
}

func (s *CharacterHandlerSuite) TestPaginationParams_Invalid() {
	req := httptest.NewRequest(http.MethodGet, "/characters?page=abc&limit=-1", nil)
	page, limit := paginationParams(req)
	s.Equal(1, page)
	s.Equal(10, limit)
}

// ── pathID ────────────────────────────────────────────────────────────────────

func (s *CharacterHandlerSuite) TestPathID_Valid() {
	req := httptest.NewRequest(http.MethodGet, "/characters/5", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "5"})
	w := httptest.NewRecorder()
	id, ok := pathID(w, req, "id")
	s.True(ok)
	s.Equal(uint(5), id)
}

func (s *CharacterHandlerSuite) TestPathID_Invalid() {
	req := httptest.NewRequest(http.MethodGet, "/characters/abc", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	w := httptest.NewRecorder()
	_, ok := pathID(w, req, "id")
	s.False(ok)
	s.Equal(http.StatusBadRequest, w.Code)
}

// ── writeServiceError ─────────────────────────────────────────────────────────

func (s *CharacterHandlerSuite) TestWriteServiceError_AppError() {
	w := httptest.NewRecorder()
	appErr := apperrors.NewError(apperrors.ErrResourceNotFound, "not found", "", http.StatusNotFound, nil, nil)
	s.handler.writeServiceError(w, appErr, "fallback")
	s.Equal(http.StatusNotFound, w.Code)
}

func (s *CharacterHandlerSuite) TestWriteServiceError_GenericError() {
	w := httptest.NewRecorder()
	s.handler.writeServiceError(w, fmt.Errorf("db failure"), "fallback msg")
	s.Equal(http.StatusInternalServerError, w.Code)
}

// ── GetCharactersHandler ──────────────────────────────────────────────────────

func (s *CharacterHandlerSuite) TestGetCharacters_Success() {
	chars := []*models.Character{{Name: "Naruto"}, {Name: "Sasuke"}}
	s.svc.EXPECT().GetAllCharacters(mock.Anything, 1, 10).Return(chars, int64(2), nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/characters", nil)
	rr := httptest.NewRecorder()
	s.handler.GetCharactersHandler(rr, req)

	s.Equal(http.StatusOK, rr.Code)
}

func (s *CharacterHandlerSuite) TestGetCharacters_ServiceError() {
	s.svc.EXPECT().GetAllCharacters(mock.Anything, 1, 10).
		Return(nil, int64(0), apperrors.NewError(apperrors.ErrInternalServer, "db err", "", http.StatusInternalServerError, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodGet, "/characters", nil)
	rr := httptest.NewRecorder()
	s.handler.GetCharactersHandler(rr, req)

	s.Equal(http.StatusInternalServerError, rr.Code)
}

// ── GetCharacterHandler ───────────────────────────────────────────────────────

func (s *CharacterHandlerSuite) TestGetCharacter_Success() {
	char := &models.Character{Name: "Naruto"}
	char.ID = 1
	s.svc.EXPECT().GetCharacterByID(mock.Anything, uint(1)).Return(char, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/characters/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.GetCharacterHandler(rr, req)

	s.Equal(http.StatusOK, rr.Code)
}

func (s *CharacterHandlerSuite) TestGetCharacter_InvalidID() {
	req := httptest.NewRequest(http.MethodGet, "/characters/abc", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rr := httptest.NewRecorder()
	s.handler.GetCharacterHandler(rr, req)

	s.Equal(http.StatusBadRequest, rr.Code)
}

func (s *CharacterHandlerSuite) TestGetCharacter_NotFound() {
	s.svc.EXPECT().GetCharacterByID(mock.Anything, uint(99)).
		Return(nil, apperrors.NewError(apperrors.ErrResourceNotFound, "not found", "", http.StatusNotFound, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodGet, "/characters/99", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "99"})
	rr := httptest.NewRecorder()
	s.handler.GetCharacterHandler(rr, req)

	s.Equal(http.StatusNotFound, rr.Code)
}

// ── SearchCharactersHandler ───────────────────────────────────────────────────

func (s *CharacterHandlerSuite) TestSearchCharacters_Success() {
	chars := []*models.Character{{Name: "Naruto"}}
	s.svc.EXPECT().GetCharactersByName(mock.Anything, "Naruto", 1, 10).Return(chars, int64(1), nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/characters/search?name=Naruto", nil)
	rr := httptest.NewRecorder()
	s.handler.SearchCharactersHandler(rr, req)

	s.Equal(http.StatusOK, rr.Code)
}

func (s *CharacterHandlerSuite) TestSearchCharacters_MissingName() {
	req := httptest.NewRequest(http.MethodGet, "/characters/search", nil)
	rr := httptest.NewRecorder()
	s.handler.SearchCharactersHandler(rr, req)

	s.Equal(http.StatusBadRequest, rr.Code)
}

func (s *CharacterHandlerSuite) TestSearchCharacters_ServiceError() {
	s.svc.EXPECT().GetCharactersByName(mock.Anything, "X", 1, 10).
		Return(nil, int64(0), fmt.Errorf("db error")).Once()

	req := httptest.NewRequest(http.MethodGet, "/characters/search?name=X", nil)
	rr := httptest.NewRecorder()
	s.handler.SearchCharactersHandler(rr, req)

	s.Equal(http.StatusInternalServerError, rr.Code)
}

// ── GetAnimeCharactersHandler ─────────────────────────────────────────────────

func (s *CharacterHandlerSuite) TestGetAnimeCharacters_Success() {
	chars := []*models.Character{{Name: "Naruto"}}
	s.svc.EXPECT().GetCharactersByAnime(mock.Anything, uint(1), 1, 10).Return(chars, int64(1), nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/animes/1/characters", nil)
	req = mux.SetURLVars(req, map[string]string{"animeId": "1"})
	rr := httptest.NewRecorder()
	s.handler.GetAnimeCharactersHandler(rr, req)

	s.Equal(http.StatusOK, rr.Code)
}

func (s *CharacterHandlerSuite) TestGetAnimeCharacters_InvalidID() {
	req := httptest.NewRequest(http.MethodGet, "/animes/bad/characters", nil)
	req = mux.SetURLVars(req, map[string]string{"animeId": "bad"})
	rr := httptest.NewRecorder()
	s.handler.GetAnimeCharactersHandler(rr, req)

	s.Equal(http.StatusBadRequest, rr.Code)
}

func (s *CharacterHandlerSuite) TestGetAnimeCharacters_ServiceError() {
	s.svc.EXPECT().GetCharactersByAnime(mock.Anything, uint(1), 1, 10).
		Return(nil, int64(0), fmt.Errorf("db error")).Once()

	req := httptest.NewRequest(http.MethodGet, "/animes/1/characters", nil)
	req = mux.SetURLVars(req, map[string]string{"animeId": "1"})
	rr := httptest.NewRecorder()
	s.handler.GetAnimeCharactersHandler(rr, req)

	s.Equal(http.StatusInternalServerError, rr.Code)
}

// ── CreateCharacterHandler ────────────────────────────────────────────────────

func (s *CharacterHandlerSuite) TestCreateCharacter_Success() {
	s.svc.EXPECT().CreateCharacter(mock.Anything, mock.AnythingOfType("*models.Character")).
		Run(func(ctx context.Context, c *models.Character) { c.ID = 1 }).
		Return(nil).Once()

	body := `{"name":"Naruto","description":"Hero","voice_actor":"Junko"}`
	req := httptest.NewRequest(http.MethodPost, "/characters", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	s.handler.CreateCharacterHandler(rr, req)

	s.Equal(http.StatusCreated, rr.Code)
}

func (s *CharacterHandlerSuite) TestCreateCharacter_BadJSON() {
	req := httptest.NewRequest(http.MethodPost, "/characters", strings.NewReader("{invalid"))
	rr := httptest.NewRecorder()
	s.handler.CreateCharacterHandler(rr, req)

	s.Equal(http.StatusBadRequest, rr.Code)
}

func (s *CharacterHandlerSuite) TestCreateCharacter_ServiceError() {
	s.svc.EXPECT().CreateCharacter(mock.Anything, mock.Anything).
		Return(apperrors.NewError(apperrors.ErrConflict, "conflict", "", http.StatusConflict, nil, nil)).Once()

	body := `{"name":"Naruto"}`
	req := httptest.NewRequest(http.MethodPost, "/characters", strings.NewReader(body))
	rr := httptest.NewRecorder()
	s.handler.CreateCharacterHandler(rr, req)

	s.Equal(http.StatusConflict, rr.Code)
}

// ── UpdateCharacterHandler ────────────────────────────────────────────────────

func (s *CharacterHandlerSuite) TestUpdateCharacter_Success() {
	char := &models.Character{Name: "Old Name"}
	char.ID = 1
	s.svc.EXPECT().GetCharacterByID(mock.Anything, uint(1)).Return(char, nil).Once()
	s.svc.EXPECT().UpdateCharacter(mock.Anything, mock.Anything).Return(nil).Once()

	body := `{"name":"New Name","description":"Updated"}`
	req := httptest.NewRequest(http.MethodPut, "/characters/1", strings.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.UpdateCharacterHandler(rr, req)

	s.Equal(http.StatusOK, rr.Code)
}

func (s *CharacterHandlerSuite) TestUpdateCharacter_InvalidID() {
	req := httptest.NewRequest(http.MethodPut, "/characters/abc", strings.NewReader(`{}`))
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rr := httptest.NewRecorder()
	s.handler.UpdateCharacterHandler(rr, req)

	s.Equal(http.StatusBadRequest, rr.Code)
}

func (s *CharacterHandlerSuite) TestUpdateCharacter_NotFound() {
	s.svc.EXPECT().GetCharacterByID(mock.Anything, uint(1)).
		Return(nil, apperrors.NewError(apperrors.ErrResourceNotFound, "not found", "", http.StatusNotFound, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodPut, "/characters/1", strings.NewReader(`{"name":"X"}`))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.UpdateCharacterHandler(rr, req)

	s.Equal(http.StatusNotFound, rr.Code)
}

func (s *CharacterHandlerSuite) TestUpdateCharacter_BadJSON() {
	char := &models.Character{Name: "Name"}
	char.ID = 1
	s.svc.EXPECT().GetCharacterByID(mock.Anything, uint(1)).Return(char, nil).Once()

	req := httptest.NewRequest(http.MethodPut, "/characters/1", strings.NewReader("{bad"))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.UpdateCharacterHandler(rr, req)

	s.Equal(http.StatusBadRequest, rr.Code)
}

// ── DeleteCharacterHandler ────────────────────────────────────────────────────

func (s *CharacterHandlerSuite) TestDeleteCharacter_Success() {
	s.svc.EXPECT().DeleteCharacter(mock.Anything, uint(1)).Return(nil).Once()

	req := httptest.NewRequest(http.MethodDelete, "/characters/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.DeleteCharacterHandler(rr, req)

	s.Equal(http.StatusNoContent, rr.Code)
}

func (s *CharacterHandlerSuite) TestDeleteCharacter_InvalidID() {
	req := httptest.NewRequest(http.MethodDelete, "/characters/x", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "x"})
	rr := httptest.NewRecorder()
	s.handler.DeleteCharacterHandler(rr, req)

	s.Equal(http.StatusBadRequest, rr.Code)
}

func (s *CharacterHandlerSuite) TestDeleteCharacter_ServiceError() {
	s.svc.EXPECT().DeleteCharacter(mock.Anything, uint(1)).
		Return(apperrors.NewError(apperrors.ErrResourceNotFound, "not found", "", http.StatusNotFound, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodDelete, "/characters/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.DeleteCharacterHandler(rr, req)

	s.Equal(http.StatusNotFound, rr.Code)
}

// ── AddAnimeCharactersHandler ─────────────────────────────────────────────────

func (s *CharacterHandlerSuite) TestAddAnimeCharacters_Success() {
	s.svc.EXPECT().AddCharactersToAnime(mock.Anything, uint(1), []uint{2, 3}).Return(nil).Once()

	body := `{"character_ids":[2,3]}`
	req := httptest.NewRequest(http.MethodPost, "/animes/1/characters", strings.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"animeId": "1"})
	rr := httptest.NewRecorder()
	s.handler.AddAnimeCharactersHandler(rr, req)

	s.Equal(http.StatusNoContent, rr.Code)
}

func (s *CharacterHandlerSuite) TestAddAnimeCharacters_InvalidAnimeID() {
	req := httptest.NewRequest(http.MethodPost, "/animes/bad/characters", strings.NewReader(`{}`))
	req = mux.SetURLVars(req, map[string]string{"animeId": "bad"})
	rr := httptest.NewRecorder()
	s.handler.AddAnimeCharactersHandler(rr, req)

	s.Equal(http.StatusBadRequest, rr.Code)
}

func (s *CharacterHandlerSuite) TestAddAnimeCharacters_BadJSON() {
	req := httptest.NewRequest(http.MethodPost, "/animes/1/characters", strings.NewReader("{bad"))
	req = mux.SetURLVars(req, map[string]string{"animeId": "1"})
	rr := httptest.NewRecorder()
	s.handler.AddAnimeCharactersHandler(rr, req)

	s.Equal(http.StatusBadRequest, rr.Code)
}

func (s *CharacterHandlerSuite) TestAddAnimeCharacters_ServiceError() {
	s.svc.EXPECT().AddCharactersToAnime(mock.Anything, uint(1), mock.Anything).
		Return(fmt.Errorf("db error")).Once()

	body := `{"character_ids":[2]}`
	req := httptest.NewRequest(http.MethodPost, "/animes/1/characters", strings.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"animeId": "1"})
	rr := httptest.NewRecorder()
	s.handler.AddAnimeCharactersHandler(rr, req)

	s.Equal(http.StatusInternalServerError, rr.Code)
}

// ── RemoveAnimeCharacterHandler ───────────────────────────────────────────────

func (s *CharacterHandlerSuite) TestRemoveAnimeCharacter_Success() {
	s.svc.EXPECT().RemoveCharacterFromAnime(mock.Anything, uint(1), uint(2)).Return(nil).Once()

	req := httptest.NewRequest(http.MethodDelete, "/animes/1/characters/2", nil)
	req = mux.SetURLVars(req, map[string]string{"animeId": "1", "characterId": "2"})
	rr := httptest.NewRecorder()
	s.handler.RemoveAnimeCharacterHandler(rr, req)

	s.Equal(http.StatusNoContent, rr.Code)
}

func (s *CharacterHandlerSuite) TestRemoveAnimeCharacter_InvalidAnimeID() {
	req := httptest.NewRequest(http.MethodDelete, "/animes/bad/characters/2", nil)
	req = mux.SetURLVars(req, map[string]string{"animeId": "bad", "characterId": "2"})
	rr := httptest.NewRecorder()
	s.handler.RemoveAnimeCharacterHandler(rr, req)

	s.Equal(http.StatusBadRequest, rr.Code)
}

func (s *CharacterHandlerSuite) TestRemoveAnimeCharacter_InvalidCharacterID() {
	req := httptest.NewRequest(http.MethodDelete, "/animes/1/characters/bad", nil)
	req = mux.SetURLVars(req, map[string]string{"animeId": "1", "characterId": "bad"})
	rr := httptest.NewRecorder()
	s.handler.RemoveAnimeCharacterHandler(rr, req)

	s.Equal(http.StatusBadRequest, rr.Code)
}

func (s *CharacterHandlerSuite) TestRemoveAnimeCharacter_ServiceError() {
	s.svc.EXPECT().RemoveCharacterFromAnime(mock.Anything, uint(1), uint(2)).
		Return(fmt.Errorf("db error")).Once()

	req := httptest.NewRequest(http.MethodDelete, "/animes/1/characters/2", nil)
	req = mux.SetURLVars(req, map[string]string{"animeId": "1", "characterId": "2"})
	rr := httptest.NewRecorder()
	s.handler.RemoveAnimeCharacterHandler(rr, req)

	s.Equal(http.StatusInternalServerError, rr.Code)
}
