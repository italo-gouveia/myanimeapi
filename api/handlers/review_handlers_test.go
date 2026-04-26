package handlers

import (
	"bytes"
	"mime/multipart"
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

// multipartReq builds an empty multipart/form-data request so ParseMultipartForm succeeds.
func multipartReq(t *testing.T, method, target string) *http.Request {
	t.Helper()
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	if err := w.Close(); err != nil {
		t.Fatalf("failed to close multipart writer: %v", err)
	}
	req := httptest.NewRequest(method, target, body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req
}

// sampleReview returns a minimal Review with associated User and Anime for handler tests.
func sampleReview(id, userID, animeID uint) *models.Review {
	return &models.Review{
		ID:      id,
		UserID:  userID,
		AnimeID: animeID,
		Content: "Great anime!",
		Rating:  9,
		User:    models.User{BaseModel: models.BaseModel{ID: userID}, Username: "testuser"},
		Anime:   models.Anime{ID: animeID, Title: "Naruto"},
	}
}

type ReviewHandlerSuite struct {
	suite.Suite
	reviewSvc  *mocks.MockReviewServiceInterface
	storageSvc *mocks.MockStorageServiceInterface
	handler    *ReviewHandler
}

func (s *ReviewHandlerSuite) SetupTest() {
	s.reviewSvc = new(mocks.MockReviewServiceInterface)
	s.storageSvc = new(mocks.MockStorageServiceInterface)
	s.handler = NewReviewHandler(s.reviewSvc, s.storageSvc)
}

func (s *ReviewHandlerSuite) TearDownTest() {
	s.reviewSvc.AssertExpectations(s.T())
	s.storageSvc.AssertExpectations(s.T())
}

func TestReviewHandlerSuite(t *testing.T) {
	suite.Run(t, new(ReviewHandlerSuite))
}

// --- Get ---

func (s *ReviewHandlerSuite) TestGet_Success() {
	s.reviewSvc.EXPECT().GetReviewByID(mock.Anything, uint(1)).
		Return(sampleReview(1, 1, 1), nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/reviews/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.GetReviewHandler(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	body := bodyJSON(s.T(), rr)
	s.Equal(float64(1), body["id"])
}

func (s *ReviewHandlerSuite) TestGet_InvalidID() {
	req := httptest.NewRequest(http.MethodGet, "/reviews/abc", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rr := httptest.NewRecorder()
	s.handler.GetReviewHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *ReviewHandlerSuite) TestGet_NotFound() {
	s.reviewSvc.EXPECT().GetReviewByID(mock.Anything, uint(999)).
		Return(nil, apperrors.NewError(apperrors.ErrResourceNotFound, "Review not found", "", http.StatusNotFound, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodGet, "/reviews/999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	rr := httptest.NewRecorder()
	s.handler.GetReviewHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusNotFound, apperrors.ErrResourceNotFound)
}

// --- Create ---

func (s *ReviewHandlerSuite) TestCreate_Success() {
	s.reviewSvc.EXPECT().CreateReview(mock.Anything, mock.AnythingOfType("*models.Review")).
		Return(nil).Once()
	s.reviewSvc.EXPECT().GetReviewByID(mock.Anything, mock.Anything).
		Return(sampleReview(1, 1, 1), nil).Once()

	req := multipartReq(s.T(), http.MethodPost, "/reviews")
	ctx := withPayload(testCtx(), &models.ReviewCreateRequest{AnimeID: 1, Content: "Great anime!", Rating: 9})
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	s.handler.CreateReviewHandler(rr, req)

	s.Equal(http.StatusCreated, rr.Code)
	body := bodyJSON(s.T(), rr)
	s.Equal(float64(1), body["id"])
}

func (s *ReviewHandlerSuite) TestCreate_MissingPayload() {
	req := multipartReq(s.T(), http.MethodPost, "/reviews")
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.CreateReviewHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

func (s *ReviewHandlerSuite) TestCreate_ServiceError() {
	s.reviewSvc.EXPECT().CreateReview(mock.Anything, mock.AnythingOfType("*models.Review")).
		Return(apperrors.NewError(apperrors.ErrInternalServer, "DB error", "", http.StatusInternalServerError, nil, nil)).Once()

	req := multipartReq(s.T(), http.MethodPost, "/reviews")
	ctx := withPayload(testCtx(), &models.ReviewCreateRequest{AnimeID: 1, Content: "Great anime!", Rating: 9})
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	s.handler.CreateReviewHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

// --- Update ---

func (s *ReviewHandlerSuite) TestUpdate_Success() {
	s.reviewSvc.EXPECT().GetReviewByID(mock.Anything, uint(1)).
		Return(sampleReview(1, 1, 1), nil).Once()
	s.reviewSvc.EXPECT().UpdateReview(mock.Anything, mock.AnythingOfType("*models.Review")).
		Return(nil).Once()
	s.reviewSvc.EXPECT().GetReviewByID(mock.Anything, uint(1)).
		Return(sampleReview(1, 1, 1), nil).Once()

	req := multipartReq(s.T(), http.MethodPut, "/reviews/1")
	ctx := withPayload(testCtx(), &models.ReviewUpdateRequest{Content: "Updated!", Rating: 8})
	req = req.WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.UpdateReviewHandler(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	body := bodyJSON(s.T(), rr)
	s.Equal(float64(1), body["id"])
}

func (s *ReviewHandlerSuite) TestUpdate_InvalidID() {
	req := multipartReq(s.T(), http.MethodPut, "/reviews/abc")
	ctx := withPayload(testCtx(), &models.ReviewUpdateRequest{Content: "Updated!"})
	req = req.WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rr := httptest.NewRecorder()
	s.handler.UpdateReviewHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *ReviewHandlerSuite) TestUpdate_MissingPayload() {
	req := multipartReq(s.T(), http.MethodPut, "/reviews/1")
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.UpdateReviewHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

func (s *ReviewHandlerSuite) TestUpdate_NotFound() {
	s.reviewSvc.EXPECT().GetReviewByID(mock.Anything, uint(999)).
		Return(nil, apperrors.NewError(apperrors.ErrResourceNotFound, "Review not found", "", http.StatusNotFound, nil, nil)).Once()

	req := multipartReq(s.T(), http.MethodPut, "/reviews/999")
	ctx := withPayload(testCtx(), &models.ReviewUpdateRequest{Content: "Updated!"})
	req = req.WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	rr := httptest.NewRecorder()
	s.handler.UpdateReviewHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusNotFound, apperrors.ErrResourceNotFound)
}

func (s *ReviewHandlerSuite) TestUpdate_Forbidden() {
	// Review belongs to user 2, but request is from user 1
	s.reviewSvc.EXPECT().GetReviewByID(mock.Anything, uint(1)).
		Return(sampleReview(1, 2, 1), nil).Once()

	req := multipartReq(s.T(), http.MethodPut, "/reviews/1")
	ctx := withPayload(testCtx(), &models.ReviewUpdateRequest{Content: "Updated!"})
	req = req.WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.UpdateReviewHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusForbidden, apperrors.ErrForbidden)
}

// --- Delete ---

func (s *ReviewHandlerSuite) TestDelete_Success() {
	s.reviewSvc.EXPECT().GetReviewByID(mock.Anything, uint(1)).
		Return(sampleReview(1, 1, 1), nil).Once()
	s.reviewSvc.EXPECT().DeleteReview(mock.Anything, uint(1)).Return(nil).Once()

	req := httptest.NewRequest(http.MethodDelete, "/reviews/1", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.DeleteReviewHandler(rr, req)

	s.Equal(http.StatusNoContent, rr.Code)
	s.Empty(rr.Body.String())
}

func (s *ReviewHandlerSuite) TestDelete_InvalidID() {
	req := httptest.NewRequest(http.MethodDelete, "/reviews/abc", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rr := httptest.NewRecorder()
	s.handler.DeleteReviewHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *ReviewHandlerSuite) TestDelete_NotFound() {
	s.reviewSvc.EXPECT().GetReviewByID(mock.Anything, uint(999)).
		Return(nil, apperrors.NewError(apperrors.ErrResourceNotFound, "Review not found", "", http.StatusNotFound, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodDelete, "/reviews/999", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	rr := httptest.NewRecorder()
	s.handler.DeleteReviewHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusNotFound, apperrors.ErrResourceNotFound)
}

func (s *ReviewHandlerSuite) TestDelete_Forbidden() {
	s.reviewSvc.EXPECT().GetReviewByID(mock.Anything, uint(1)).
		Return(sampleReview(1, 2, 1), nil).Once()

	req := httptest.NewRequest(http.MethodDelete, "/reviews/1", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.DeleteReviewHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusForbidden, apperrors.ErrForbidden)
}

func (s *ReviewHandlerSuite) TestDelete_InternalError() {
	s.reviewSvc.EXPECT().GetReviewByID(mock.Anything, uint(1)).
		Return(sampleReview(1, 1, 1), nil).Once()
	s.reviewSvc.EXPECT().DeleteReview(mock.Anything, uint(1)).
		Return(apperrors.NewError(apperrors.ErrInternalServer, "DB error", "", http.StatusInternalServerError, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodDelete, "/reviews/1", nil)
	req = req.WithContext(testCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.DeleteReviewHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

func (s *ReviewHandlerSuite) TestDelete_Unauthenticated() {
	req := httptest.NewRequest(http.MethodDelete, "/reviews/1", nil)
	req = req.WithContext(noAuthCtx())
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	s.handler.DeleteReviewHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusUnauthorized, apperrors.ErrUnauthorized)
}
