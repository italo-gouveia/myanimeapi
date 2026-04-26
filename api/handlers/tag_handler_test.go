package handlers

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

type TagHandlerSuite struct {
	suite.Suite
	svc     *mocks.MockTagServiceInterface
	handler *TagHandler
}

func (s *TagHandlerSuite) SetupTest() {
	s.svc = new(mocks.MockTagServiceInterface)
	s.handler = NewTagHandler(s.svc)
}

func (s *TagHandlerSuite) TearDownTest() {
	s.svc.AssertExpectations(s.T())
}

func TestTagHandlerSuite(t *testing.T) {
	suite.Run(t, new(TagHandlerSuite))
}

// Create tests

func (s *TagHandlerSuite) TestCreate_Success() {
	s.svc.EXPECT().CreateTag(mock.Anything, mock.AnythingOfType("*models.Tag")).
		Run(func(ctx context.Context, tag *models.Tag) { tag.ID = 1 }).
		Return(nil).Once()

	payload := &models.TagCreateRequest{Name: "Action"}
	req := httptest.NewRequest(http.MethodPost, "/tags", jsonBuf(payload))
	req = req.WithContext(withPayload(testCtx(), payload))
	rr := httptest.NewRecorder()
	s.handler.CreateTagHandler(rr, req)

	s.Equal(http.StatusCreated, rr.Code)
	body := bodyJSON(s.T(), rr)
	s.Equal(float64(1), body["id"])
	s.Equal("Action", body["name"])
}

func (s *TagHandlerSuite) TestCreate_MissingPayload() {
	req := httptest.NewRequest(http.MethodPost, "/tags", jsonBuf(nil))
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.CreateTagHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

func (s *TagHandlerSuite) TestCreate_Conflict() {
	s.svc.EXPECT().CreateTag(mock.Anything, mock.Anything).
		Return(apperrors.NewError(apperrors.ErrConflict, "Tag already exists", "", http.StatusConflict, nil, nil)).Once()

	payload := &models.TagCreateRequest{Name: "Action"}
	req := httptest.NewRequest(http.MethodPost, "/tags", jsonBuf(payload))
	req = req.WithContext(withPayload(testCtx(), payload))
	rr := httptest.NewRecorder()
	s.handler.CreateTagHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusConflict, apperrors.ErrConflict)
}

func (s *TagHandlerSuite) TestCreate_InternalError() {
	s.svc.EXPECT().CreateTag(mock.Anything, mock.Anything).
		Return(apperrors.NewError(apperrors.ErrInternalServer, "DB error", "", http.StatusInternalServerError, nil, nil)).Once()

	payload := &models.TagCreateRequest{Name: "Action"}
	req := httptest.NewRequest(http.MethodPost, "/tags", jsonBuf(payload))
	req = req.WithContext(withPayload(testCtx(), payload))
	rr := httptest.NewRecorder()
	s.handler.CreateTagHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

// Get tests

func (s *TagHandlerSuite) TestGet_Success() {
	s.svc.EXPECT().GetTagByID(mock.Anything, uint(1)).
		Return(&models.Tag{ID: 1, Name: "Action"}, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/tags/1", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.GetTagHandler(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	body := bodyJSON(s.T(), rr)
	s.Equal(float64(1), body["id"])
}

func (s *TagHandlerSuite) TestGet_InvalidID() {
	req := httptest.NewRequest(http.MethodGet, "/tags/abc", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rr := httptest.NewRecorder()
	s.handler.GetTagHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *TagHandlerSuite) TestGet_NotFound() {
	s.svc.EXPECT().GetTagByID(mock.Anything, uint(999)).
		Return(nil, apperrors.NewError(apperrors.ErrResourceNotFound, "Tag not found", "", http.StatusNotFound, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodGet, "/tags/999", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	rr := httptest.NewRecorder()
	s.handler.GetTagHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusNotFound, apperrors.ErrResourceNotFound)
}

func (s *TagHandlerSuite) TestGet_InternalError() {
	s.svc.EXPECT().GetTagByID(mock.Anything, uint(1)).
		Return(nil, apperrors.NewError(apperrors.ErrInternalServer, "DB error", "", http.StatusInternalServerError, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodGet, "/tags/1", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.GetTagHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

// GetAll tests

func (s *TagHandlerSuite) TestGetAll_Success() {
	s.svc.EXPECT().GetAllTags(mock.Anything, 1, 10).
		Return([]models.Tag{{ID: 1, Name: "Action"}, {ID: 2, Name: "Comedy"}}, int64(2), nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/tags?page=1&limit=10", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.GetAllTagsHandler(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	body := bodyJSON(s.T(), rr)
	s.Equal(float64(2), body["total"])
}

func (s *TagHandlerSuite) TestGetAll_InvalidPagination() {
	req := httptest.NewRequest(http.MethodGet, "/tags?page=0&limit=10", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.GetAllTagsHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *TagHandlerSuite) TestGetAll_EmptyResult() {
	s.svc.EXPECT().GetAllTags(mock.Anything, 1, 10).
		Return([]models.Tag{}, int64(0), nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/tags?page=1&limit=10", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.GetAllTagsHandler(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	body := bodyJSON(s.T(), rr)
	s.Equal(float64(0), body["total"])
}

func (s *TagHandlerSuite) TestGetAll_InternalError() {
	s.svc.EXPECT().GetAllTags(mock.Anything, 1, 10).
		Return(nil, int64(0), apperrors.NewError(apperrors.ErrInternalServer, "DB error", "", http.StatusInternalServerError, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodGet, "/tags?page=1&limit=10", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.GetAllTagsHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

// Update tests

func (s *TagHandlerSuite) TestUpdate_Success() {
	s.svc.EXPECT().UpdateTag(mock.Anything, mock.AnythingOfType("*models.Tag")).Return(nil).Once()

	payload := &models.TagUpdateRequest{Name: "Updated Action"}
	req := httptest.NewRequest(http.MethodPut, "/tags/1", jsonBuf(payload))
	req = req.WithContext(withPayload(testCtx(), payload))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.UpdateTagHandler(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	body := bodyJSON(s.T(), rr)
	s.Equal("Updated Action", body["name"])
}

func (s *TagHandlerSuite) TestUpdate_InvalidID() {
	payload := &models.TagUpdateRequest{Name: "X"}
	req := httptest.NewRequest(http.MethodPut, "/tags/abc", jsonBuf(payload))
	req = req.WithContext(withPayload(testCtx(), payload))
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rr := httptest.NewRecorder()
	s.handler.UpdateTagHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *TagHandlerSuite) TestUpdate_MissingPayload() {
	req := httptest.NewRequest(http.MethodPut, "/tags/1", jsonBuf(nil))
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.UpdateTagHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

func (s *TagHandlerSuite) TestUpdate_NotFound() {
	s.svc.EXPECT().UpdateTag(mock.Anything, mock.Anything).
		Return(apperrors.NewError(apperrors.ErrResourceNotFound, "Tag not found", "", http.StatusNotFound, nil, nil)).Once()

	payload := &models.TagUpdateRequest{Name: "X"}
	req := httptest.NewRequest(http.MethodPut, "/tags/999", jsonBuf(payload))
	req = req.WithContext(withPayload(testCtx(), payload))
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	rr := httptest.NewRecorder()
	s.handler.UpdateTagHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusNotFound, apperrors.ErrResourceNotFound)
}

// Delete tests

func (s *TagHandlerSuite) TestDelete_Success() {
	s.svc.EXPECT().DeleteTag(mock.Anything, uint(1)).Return(nil).Once()

	req := httptest.NewRequest(http.MethodDelete, "/tags/1", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.DeleteTagHandler(rr, req)

	s.Equal(http.StatusNoContent, rr.Code)
	s.Empty(rr.Body.String())
}

func (s *TagHandlerSuite) TestDelete_InvalidID() {
	req := httptest.NewRequest(http.MethodDelete, "/tags/abc", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rr := httptest.NewRecorder()
	s.handler.DeleteTagHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *TagHandlerSuite) TestDelete_NotFound() {
	s.svc.EXPECT().DeleteTag(mock.Anything, uint(999)).
		Return(apperrors.NewError(apperrors.ErrResourceNotFound, "Tag not found", "", http.StatusNotFound, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodDelete, "/tags/999", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	rr := httptest.NewRecorder()
	s.handler.DeleteTagHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusNotFound, apperrors.ErrResourceNotFound)
}

func (s *TagHandlerSuite) TestDelete_InternalError() {
	s.svc.EXPECT().DeleteTag(mock.Anything, uint(1)).
		Return(apperrors.NewError(apperrors.ErrInternalServer, "DB error", "", http.StatusInternalServerError, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodDelete, "/tags/1", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.DeleteTagHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}
