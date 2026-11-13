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
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// testCtx returns a context with a valid user ID — simulates an authenticated (non-admin) request.
func testCtx() context.Context {
	return context.WithValue(context.Background(), middleware.UserContextKey, uint(1))
}

// adminCtx returns a context with a valid user ID and admin flag — simulates an admin request.
func adminCtx() context.Context {
	ctx := context.WithValue(context.Background(), middleware.UserContextKey, uint(1))
	return context.WithValue(ctx, middleware.IsAdminContextKey, true)
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

// newSQLiteTestDB opens an isolated in-memory SQLite database, runs AutoMigrate on the
// provided models, and returns the ready-to-use *gorm.DB.
// MaxOpenConns(1) keeps all GORM pool connections on the same :memory: instance.
func newSQLiteTestDB(t *testing.T, migrate ...interface{}) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.SetMaxOpenConns(1)

	if len(migrate) > 0 {
		require.NoError(t, db.AutoMigrate(migrate...))
	}
	return db
}
