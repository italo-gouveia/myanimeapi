// internal/config/config.go
// Package config provides configuration settings for the application.
// It loads the configuration from environment variables.
// It defines the configuration structure.
// It provides a function to load the configuration.
// It uses the os package to read environment variables.
// It uses the strconv package to convert environment variables to integers.
// It uses the fmt package to format strings.
// It uses the log package to log messages.
package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type ServerConfig struct {
	Host string
	Port int
}

type DatabaseConfig struct {
	Type     string
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type LoggingConfig struct {
	Level string
	File  string
}

type APIConfig struct {
	Key string
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logging  LoggingConfig
	API      APIConfig
}

func LoadConfig() *Config {
	var cfg Config

	// Load server configuration
	cfg.Server.Host = getEnv("SERVER_HOST", "localhost")
	cfg.Server.Port = getEnvAsInt("SERVER_PORT", 8080)
	log.Printf("Loaded server configuration: Host=%s, Port=%d", cfg.Server.Host, cfg.Server.Port)

	// Load database configuration
	cfg.Database.Type = getEnv("DB_TYPE", "postgres")
	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnvAsInt("DB_PORT", 5432)
	cfg.Database.User = getEnv("DB_USER", "user")
	cfg.Database.Password = getEnv("DB_PASSWORD", "password")
	cfg.Database.Name = getEnv("DB_NAME", "myanimeapi")
	log.Printf("Loaded database configuration: Type=%s, Host=%s, Port=%d, User=%s, Name=%s",
		cfg.Database.Type, cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Name)

	// Validate database password
	/*if cfg.Database.Password == "password" {
		log.Fatal("Database password must be set in environment variables (DB_PASSWORD)")
	}*/

	// Load logging configuration
	cfg.Logging.Level = getEnv("LOG_LEVEL", "info")
	cfg.Logging.File = getEnv("LOG_FILE", "/var/log/myapp.log")
	log.Printf("Loaded logging configuration: Level=%s, File=%s", cfg.Logging.Level, cfg.Logging.File)

	// Load API configuration
	cfg.API.Key = getEnv("API_KEY", "your_api_key")
	log.Printf("Loaded API configuration: Key=%s", cfg.API.Key)

	// Validate API key
	/*	if cfg.API.Key == "your_api_key" {
		log.Fatal("API key must be set in environment variables (API_KEY)")
	}*/

	return &cfg
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		boolValue, err := strconv.ParseBool(value)
		if err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		durationValue, err := time.ParseDuration(value)
		if err == nil {
			return durationValue
		}
	}
	return defaultValue
}
