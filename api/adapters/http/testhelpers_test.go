package httphandler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"myanimeapi/api/middleware"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testCtx returns a context with a valid user ID — simulates an authenticated request.
func testCtx() context.Context {
	return context.WithValue(context.Background(), middleware.UserContextKey, uint(1))
}

// noAuthCtx returns a bare context with no user ID — simulates an unauthenticated request.
func noAuthCtx() context.Context {
	return context.Background()
}

// withPayload injects a validated payload into the context, bypassing the validation middleware.
func withPayload(ctx context.Context, payload interface{}) context.Context {
	return context.WithValue(ctx, middleware.ValidatedPayloadKey, payload)
}

// assertErrorCode asserts the HTTP status code and the "error.code" field in the JSON body.
func assertErrorCode(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantCode string) {
	t.Helper()
	assert.Equal(t, wantStatus, rr.Code)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body), "body must be valid JSON")
	errObj, ok := body["error"].(map[string]interface{})
	require.True(t, ok, "body must contain 'error' object")
	assert.Equal(t, wantCode, errObj["code"])
}

// bodyJSON decodes the response body as a JSON map.
func bodyJSON(t *testing.T, rr *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	return body
}

// jsonBuf returns a *bytes.Buffer containing the JSON-encoded value (for request bodies).
func jsonBuf(v interface{}) *bytes.Buffer {
	b, _ := json.Marshal(v)
	return bytes.NewBuffer(b)
}
