package httphandler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteJSON_Success(t *testing.T) {
	w := httptest.NewRecorder()
	payload := map[string]string{"key": "value"}

	ok := writeJSON(w, http.StatusOK, payload, nil)

	assert.True(t, ok)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), `"key"`)
}

func TestWriteJSON_Created(t *testing.T) {
	w := httptest.NewRecorder()

	ok := writeJSON(w, http.StatusCreated, map[string]int{"id": 1}, nil)

	assert.True(t, ok)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestWriteJSON_UnencodableValue(t *testing.T) {
	w := httptest.NewRecorder()
	// json.Encoder.Encode writes the status before encoding, so by the time
	// encoding fails the status is already committed and we can only send a
	// second response (which is what writeJSON does on error).
	// A channel cannot be marshalled to JSON.
	ch := make(chan int)

	ok := writeJSON(w, http.StatusOK, ch, nil)

	assert.False(t, ok)
}

func TestRespondInvalidPayload_SetsResponse(t *testing.T) {
	w := httptest.NewRecorder()

	respondInvalidPayload(w, nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), msgInvalidPayload)
}
