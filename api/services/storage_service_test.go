package services

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSanitizeFilename_Normal(t *testing.T) {
	assert.Equal(t, "photo.jpg", sanitizeFilename("photo.jpg"))
}

func TestSanitizeFilename_PathTraversal(t *testing.T) {
	result := sanitizeFilename("../../etc/passwd")
	assert.Equal(t, "passwd", result)
}

func TestSanitizeFilename_NestedPath(t *testing.T) {
	result := sanitizeFilename("uploads/subdir/image.png")
	assert.Equal(t, "image.png", result)
}

func TestNewLocalStorageStrategy_Success(t *testing.T) {
	dir := t.TempDir()
	svc, err := NewLocalStorageStrategy(dir, "http://localhost:8080")
	require.NoError(t, err)
	assert.NotNil(t, svc)

	uploadDir := filepath.Join(dir, "uploads")
	_, statErr := os.Stat(uploadDir)
	assert.NoError(t, statErr)
}

func TestLocalStorageStrategy_GetBaseURL(t *testing.T) {
	dir := t.TempDir()
	svc, err := NewLocalStorageStrategy(dir, "http://cdn.example.com")
	require.NoError(t, err)
	assert.Equal(t, "http://cdn.example.com", svc.GetBaseURL())
}

func TestStorageService_GetBaseURL(t *testing.T) {
	mock := &mockStorageSvc{}
	mock.On("GetBaseURL").Return("http://storage.example.com")

	svc := NewStorageService(mock)
	assert.Equal(t, "http://storage.example.com", svc.GetBaseURL())
	mock.AssertExpectations(t)
}

func TestStorageService_DeleteFile(t *testing.T) {
	mock := &mockStorageSvc{}
	ctx := context.Background()
	mock.On("DeleteFile", ctx, "http://storage.example.com/img.jpg").Return(nil)

	svc := NewStorageService(mock)
	err := svc.DeleteFile(ctx, "http://storage.example.com/img.jpg")
	require.NoError(t, err)
	mock.AssertExpectations(t)
}

func TestLocalStorageStrategy_DeleteFile_Success(t *testing.T) {
	dir := t.TempDir()
	svc, err := NewLocalStorageStrategy(dir, "http://localhost:8080")
	require.NoError(t, err)

	// Create a real file inside the upload dir
	uploadDir := filepath.Join(dir, "uploads")
	filePath := filepath.Join(uploadDir, "test.jpg")
	require.NoError(t, os.WriteFile(filePath, []byte("data"), 0644))

	fileURL := "http://localhost:8080/uploads/test.jpg"
	err = svc.DeleteFile(context.Background(), fileURL)
	require.NoError(t, err)

	_, statErr := os.Stat(filePath)
	assert.True(t, os.IsNotExist(statErr))
}

func TestLocalStorageStrategy_DeleteFile_NotExist(t *testing.T) {
	dir := t.TempDir()
	svc, err := NewLocalStorageStrategy(dir, "http://localhost:8080")
	require.NoError(t, err)

	err = svc.DeleteFile(context.Background(), "http://localhost:8080/uploads/ghost.jpg")
	assert.Error(t, err)
}
