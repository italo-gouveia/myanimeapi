package utils

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"myanimeapi/api/models"
	"myanimeapi/internal/config"
	testconfig "myanimeapi/tests/config"
)

// SetupTestDatabase creates a test database connection
func SetupTestDatabase(t *testing.T) *gorm.DB {
	t.Helper()

	// Use test environment variables
	dbHost := getEnv("TEST_DB_HOST", "localhost")
	dbPort := getEnv("TEST_DB_PORT", "5433")
	dbUser := getEnv("TEST_DB_USER", "test")
	dbPass := getEnv("TEST_DB_PASSWORD", "test")
	dbName := getEnv("TEST_DB_NAME", "test_db")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	// Test connection
	sqlDB, err := db.DB()
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = sqlDB.PingContext(ctx)
	require.NoError(t, err)

	// Run migrations
	runMigrations(t, db)

	return db
}

// LoadTestConfig loads test configuration
func LoadTestConfig(t *testing.T) *config.Config {
	t.Helper()

	// Load test-specific configuration
	testCfg := testconfig.LoadTestConfig()
	return testCfg.ToAppConfig()
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func runMigrations(t *testing.T, db *gorm.DB) {
	// Run your migrations here
	err := db.AutoMigrate(
		&models.User{},
		&models.Anime{},
		&models.Review{},
		&models.Genre{},
		&models.Tag{},
		&models.Favorite{},
	)
	require.NoError(t, err)
}

// CleanupTestDatabase removes all test data
func CleanupTestDatabase(t *testing.T, db *gorm.DB) {
	t.Helper()

	// Delete all data from tables in reverse dependency order
	err := db.Exec("DELETE FROM favorites").Error
	require.NoError(t, err)

	err = db.Exec("DELETE FROM reviews").Error
	require.NoError(t, err)

	err = db.Exec("DELETE FROM anime_genres").Error
	require.NoError(t, err)

	err = db.Exec("DELETE FROM anime_tags").Error
	require.NoError(t, err)

	err = db.Exec("DELETE FROM animes").Error
	require.NoError(t, err)

	err = db.Exec("DELETE FROM genres").Error
	require.NoError(t, err)

	err = db.Exec("DELETE FROM tags").Error
	require.NoError(t, err)

	err = db.Exec("DELETE FROM users").Error
	require.NoError(t, err)
}
