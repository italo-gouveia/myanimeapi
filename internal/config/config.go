// internal/config/config.go
// Package config provides configuration settings for the application.
// It loads configuration from environment variables and defines the structure for server, database, logging, and API configurations.
// The package uses the `os` package to read environment variables and the `strconv` package to convert environment variables to integers.
// It also provides helper functions to retrieve environment variables with default values.
package config

import (
	"myanimeapi/internal/logger"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

// ServerConfig holds configuration settings for the server.
type ServerConfig struct {
	Host string // Host address for the server
	Port int    // Port number for the server
}

// DatabaseConfig holds configuration settings for the database.
type DatabaseConfig struct {
	Type     string // Type of the database (e.g., postgres, mysql)
	Host     string // Host address for the database
	Port     int    // Port number for the database
	User     string // Username for the database
	Password string // Password for the database
	Name     string // Name of the database
}

// LoggingConfig holds configuration settings for logging.
type LoggingConfig struct {
	Level string // Logging level (e.g., info, debug, error)
	File  string // File path for logging output
}

// APIConfig holds configuration settings for the API.
type APIConfig struct {
	Key string // API key for authentication
}

// Config is the top-level configuration structure that holds all configuration settings.
type Config struct {
	Server   ServerConfig   // Server configuration
	Database DatabaseConfig // Database configuration
	Logging  LoggingConfig  // Logging configuration
	API      APIConfig      // API configuration
}

// LoadConfig loads the configuration from environment variables and returns a Config struct.
// It uses default values for missing environment variables and logs the loaded configuration.
// Example usage:
//
//	cfg := LoadConfig()
//	fmt.Println(cfg.Server.Host)
//
// This will load the configuration and print the server host.
func LoadConfig() *Config {
	var cfg Config

	// Load server configuration
	cfg.Server.Host = getEnv("SERVER_HOST", "localhost")
	cfg.Server.Port = getEnvAsInt("SERVER_PORT", 8080)
	log := logger.Get()
	log.WithFields(logrus.Fields{
		"host": cfg.Server.Host,
		"port": cfg.Server.Port,
	}).Info("Loaded server configuration")

	// Load database configuration
	cfg.Database.Type = getEnv("DB_TYPE", "postgres")
	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnvAsInt("DB_PORT", 5432)
	cfg.Database.User = getEnv("DB_USER", "user")
	cfg.Database.Password = getEnv("DB_PASSWORD", "password")
	cfg.Database.Name = getEnv("DB_NAME", "myanimeapi")
	log.WithFields(logrus.Fields{
		"type": cfg.Database.Type,
		"host": cfg.Database.Host,
		"port": cfg.Database.Port,
		"user": cfg.Database.User,
		"name": cfg.Database.Name,
	}).Info("Loaded database configuration")

	if cfg.Database.Password == "" {
		log.Fatal("Database password must be set in environment variables (DB_PASSWORD)")
	}

	// Load logging configuration
	cfg.Logging.Level = getEnv("LOG_LEVEL", "info")
	cfg.Logging.File = getEnv("LOG_FILE", "/var/log/myapp.log")
	log.WithFields(logrus.Fields{
		"level": cfg.Logging.Level,
		"file":  cfg.Logging.File,
	}).Info("Loaded logging configuration")

	// Load API configuration
	cfg.API.Key = getEnv("API_KEY", "your_api_key")
	log.WithField("key", cfg.API.Key).Info("Loaded API configuration")

	if cfg.API.Key == "" {
		log.Fatal("API key must be set in environment variables (API_KEY)")
	}

	return &cfg
}

// getEnv retrieves the value of an environment variable or returns a default value if the variable is not set.
func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt retrieves the value of an environment variable as an integer or returns a default value if the variable is not set or cannot be converted.
func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsBool (commented out) retrieves the value of an environment variable as a boolean or returns a default value if the variable is not set or cannot be converted.
/*
func getEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		boolValue, err := strconv.ParseBool(value)
		if err == nil {
			return boolValue
		}
	}
	return defaultValue
}
*/

// getEnvAsDuration (commented out) retrieves the value of an environment variable as a duration or returns a default value if the variable is not set or cannot be converted.
/*
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		durationValue, err := time.ParseDuration(value)
		if err == nil {
			return durationValue
		}
	}
	return defaultValue
}
*/
