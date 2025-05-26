package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// StorageServiceInterface defines the interface for storage service operations
type StorageServiceInterface interface {
	// SaveFile saves a file and returns its URL
	SaveFile(file *multipart.FileHeader, directory string) (string, error)
	// DeleteFile deletes a file by its URL
	DeleteFile(ctx context.Context, fileURL string) error
	// GetBaseURL returns the base URL for accessing files
	GetBaseURL() string
}

// StorageStrategy defines the interface for different storage implementations
type StorageStrategy interface {
	// SaveFile saves a file and returns its URL
	SaveFile(file *multipart.FileHeader, directory string) (string, error)
	// DeleteFile deletes a file by its URL
	DeleteFile(ctx context.Context, fileURL string) error
	// GetBaseURL returns the base URL for accessing files
	GetBaseURL() string
}

// LocalStorageStrategy implements StorageStrategy for local file system
type LocalStorageStrategy struct {
	baseDir   string
	baseURL   string
	uploadDir string
}

// NewLocalStorageStrategy creates a new local storage strategy
func NewLocalStorageStrategy(baseDir, baseURL string) (*LocalStorageStrategy, error) {
	uploadDir := filepath.Join(baseDir, "uploads")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	return &LocalStorageStrategy{
		baseDir:   baseDir,
		baseURL:   baseURL,
		uploadDir: uploadDir,
	}, nil
}

// SaveFile implements StorageStrategy for local storage
func (s *LocalStorageStrategy) SaveFile(file *multipart.FileHeader, directory string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Create directory if it doesn't exist
	dirPath := filepath.Join(s.uploadDir, directory)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Create destination file
	dstPath := filepath.Join(dirPath, file.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy file contents
	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	// Return the URL for the saved file
	return fmt.Sprintf("%s/%s/%s", s.baseURL, directory, file.Filename), nil
}

// DeleteFile implements StorageStrategy for local storage
func (s *LocalStorageStrategy) DeleteFile(ctx context.Context, fileURL string) error {
	// Extract the file path from the URL
	path := strings.TrimPrefix(fileURL, s.baseURL)
	filePath := filepath.Join(s.baseDir, path)

	// Delete the file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// GetBaseURL implements StorageStrategy for local storage
func (s *LocalStorageStrategy) GetBaseURL() string {
	return s.baseURL
}

// S3StorageStrategy implements StorageStrategy for AWS S3
type S3StorageStrategy struct {
	s3Client  *s3.S3
	bucket    string
	baseURL   string
	uploadDir string
}

// NewS3StorageStrategy creates a new S3 storage strategy
func NewS3StorageStrategy(region, bucket, baseURL, uploadDir string) (*S3StorageStrategy, error) {
	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewEnvCredentials(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	// Create S3 client
	s3Client := s3.New(sess)

	return &S3StorageStrategy{
		s3Client:  s3Client,
		bucket:    bucket,
		baseURL:   baseURL,
		uploadDir: uploadDir,
	}, nil
}

// SaveFile implements StorageStrategy for S3 storage
func (s *S3StorageStrategy) SaveFile(file *multipart.FileHeader, directory string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Create the S3 key
	key := fmt.Sprintf("%s/%s/%s", s.uploadDir, directory, file.Filename)

	// Upload file to S3
	_, err = s.s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   src,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Return the URL for the saved file
	return fmt.Sprintf("%s/%s", s.baseURL, key), nil
}

// DeleteFile implements StorageStrategy for S3 storage
func (s *S3StorageStrategy) DeleteFile(ctx context.Context, fileURL string) error {
	// Extract the key from the URL
	key := strings.TrimPrefix(fileURL, s.baseURL)

	// Delete the file from S3
	_, err := s.s3Client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}

// GetBaseURL implements StorageStrategy for S3 storage
func (s *S3StorageStrategy) GetBaseURL() string {
	return s.baseURL
}

// StorageService manages file storage operations using the configured strategy
// It implements the StorageServiceInterface.
type StorageService struct {
	strategy StorageStrategy
}

// NewStorageService creates a new storage service with the specified strategy
// It returns a StorageServiceInterface implementation.
func NewStorageService(strategy StorageStrategy) StorageServiceInterface {
	return &StorageService{
		strategy: strategy,
	}
}

// SaveFile saves a file using the configured storage strategy
func (s *StorageService) SaveFile(file *multipart.FileHeader, directory string) (string, error) {
	return s.strategy.SaveFile(file, directory)
}

// DeleteFile deletes a file using the configured storage strategy
func (s *StorageService) DeleteFile(ctx context.Context, fileURL string) error {
	return s.strategy.DeleteFile(ctx, fileURL)
}

// GetBaseURL returns the base URL for accessing files
func (s *StorageService) GetBaseURL() string {
	return s.strategy.GetBaseURL()
}
