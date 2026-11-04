package httphandler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"myanimeapi/api/mocks"
	"myanimeapi/api/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newETLHandler(t *testing.T) (*ETLHandler, *mocks.MockETLServiceInterface) {
	t.Helper()
	svc := new(mocks.MockETLServiceInterface)
	return NewETLHandler(svc), svc
}

func TestSyncJikanHandler_Success(t *testing.T) {
	h, svc := newETLHandler(t)
	result := services.SyncResult{Imported: 10, Updated: 5, Pages: 5}
	svc.EXPECT().SyncFromJikan(mock.Anything, 5).Return(result, nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/etl/sync", nil)
	rr := httptest.NewRecorder()
	h.SyncJikanHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	svc.AssertExpectations(t)
}

func TestSyncJikanHandler_CustomPages(t *testing.T) {
	h, svc := newETLHandler(t)
	result := services.SyncResult{Imported: 3, Pages: 2}
	svc.EXPECT().SyncFromJikan(mock.Anything, 2).Return(result, nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/etl/sync?pages=2", nil)
	rr := httptest.NewRecorder()
	h.SyncJikanHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	svc.AssertExpectations(t)
}

func TestSyncJikanHandler_PagesCapped(t *testing.T) {
	h, svc := newETLHandler(t)
	result := services.SyncResult{Imported: 50, Pages: 20}
	// pages=99 should be capped to 20
	svc.EXPECT().SyncFromJikan(mock.Anything, 20).Return(result, nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/etl/sync?pages=99", nil)
	rr := httptest.NewRecorder()
	h.SyncJikanHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	svc.AssertExpectations(t)
}

func TestSyncJikanHandler_ServiceError(t *testing.T) {
	h, svc := newETLHandler(t)
	svc.EXPECT().SyncFromJikan(mock.Anything, 5).Return(services.SyncResult{}, fmt.Errorf("jikan unavailable")).Once()

	req := httptest.NewRequest(http.MethodPost, "/etl/sync", nil)
	rr := httptest.NewRecorder()
	h.SyncJikanHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	svc.AssertExpectations(t)
}

func TestSyncJikanHandler_InvalidPagesParam(t *testing.T) {
	h, svc := newETLHandler(t)
	result := services.SyncResult{Pages: 5}
	// invalid "pages" value falls back to default 5
	svc.EXPECT().SyncFromJikan(mock.Anything, 5).Return(result, nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/etl/sync?pages=abc", nil)
	rr := httptest.NewRecorder()
	h.SyncJikanHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	svc.AssertExpectations(t)
}
