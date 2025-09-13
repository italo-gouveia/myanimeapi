package config

import (
	"os"
	"strconv"

	"myanimeapi/internal/config"
)

// TestConfig provides test-specific configuration
type TestConfig struct {
	Database DatabaseTestConfig
	Server   ServerTestConfig
	Logging  LoggingTestConfig
}

// DatabaseTestConfig provides database configuration for tests
type DatabaseTestConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

// ServerTestConfig provides server configuration for tests
type ServerTestConfig struct {
	Port int
}

// LoggingTestConfig provides logging configuration for tests
type LoggingTestConfig struct {
	Level string
}

// LoadTestConfig loads test configuration from environment variables
func LoadTestConfig() *TestConfig {
	return &TestConfig{
		Database: DatabaseTestConfig{
			Host:     getEnv("TEST_DB_HOST", "localhost"),
			Port:     getEnvAsInt("TEST_DB_PORT", 5432),
			User:     getEnv("TEST_DB_USER", "test"),
			Password: getEnv("TEST_DB_PASSWORD", "test"),
			Name:     getEnv("TEST_DB_NAME", "test_db"),
			SSLMode:  getEnv("TEST_DB_SSLMODE", "disable"),
		},
		Server: ServerTestConfig{
			Port: getEnvAsInt("TEST_SERVER_PORT", 8081),
		},
		Logging: LoggingTestConfig{
			Level: getEnv("TEST_LOG_LEVEL", "debug"),
		},
	}
}

// ToAppConfig converts TestConfig to the application's config.Config
func (tc *TestConfig) ToAppConfig() *config.Config {
	return &config.Config{
		Database: config.DatabaseConfig{
			Host:     tc.Database.Host,
			Port:     tc.Database.Port,
			User:     tc.Database.User,
			Password: tc.Database.Password,
			Name:     tc.Database.Name,
		},
		Server: config.ServerConfig{
			Port: tc.Server.Port,
		},
		Logging: config.LoggingConfig{
			Level: tc.Logging.Level,
		},
	}
}

// GetDatabaseDSN returns the database connection string for tests
func (tc *TestConfig) GetDatabaseDSN() string {
	return "host=" + tc.Database.Host +
		" port=" + strconv.Itoa(tc.Database.Port) +
		" user=" + tc.Database.User +
		" password=" + tc.Database.Password +
		" dbname=" + tc.Database.Name +
		" sslmode=" + tc.Database.SSLMode
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// IsCIEnvironment checks if running in CI environment
func IsCIEnvironment() bool {
	return os.Getenv("CI") == "true" || os.Getenv("GITHUB_ACTIONS") == "true"
}

// GetTestDatabaseName returns a unique database name for tests
func GetTestDatabaseName() string {
	baseName := getEnv("TEST_DB_NAME", "test_db")
	if IsCIEnvironment() {
		// In CI, use a unique name to avoid conflicts
		return baseName + "_" + os.Getenv("GITHUB_RUN_ID")
	}
	return baseName
}
