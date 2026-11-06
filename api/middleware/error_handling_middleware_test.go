package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	apperrors "myanimeapi/internal/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
	}{
		{
			name:       "unauthorized error",
			err:        apperrors.NewError(apperrors.ErrUnauthorized, "unauthorized", "", http.StatusUnauthorized, nil, nil),
			wantStatus: http.StatusUnauthorized,
			wantCode:   apperrors.ErrUnauthorized,
		},
		{
			name:       "forbidden error",
			err:        apperrors.NewError(apperrors.ErrForbidden, "forbidden", "", http.StatusForbidden, nil, nil),
			wantStatus: http.StatusForbidden,
			wantCode:   apperrors.ErrForbidden,
		},
		{
			name:       "not found error",
			err:        apperrors.NewNotFoundError("Anime", "1", apperrors.ErrResourceNotFound),
			wantStatus: http.StatusNotFound,
			wantCode:   apperrors.ErrResourceNotFound,
		},
		{
			name:       "invalid input error",
			err:        apperrors.NewError(apperrors.ErrInvalidInput, "invalid", "", http.StatusBadRequest, nil, nil),
			wantStatus: http.StatusBadRequest,
			wantCode:   apperrors.ErrInvalidInput,
		},
		{
			name:       "generic error",
			err:        fmt.Errorf("some unexpected error"),
			wantStatus: http.StatusInternalServerError,
			wantCode:   apperrors.ErrInternalServer,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/test", nil)

			assert.NotPanics(t, func() { HandleError(w, r, tc.err) })
			assert.Equal(t, tc.wantStatus, w.Code)
		})
	}
}

func TestGetErrorMetrics(t *testing.T) {
	m := GetErrorMetrics()
	require.NotNil(t, m)
}
