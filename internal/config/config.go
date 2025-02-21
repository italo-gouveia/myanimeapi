// internal/config/config.go
package config

import (
	"os"
	"strconv"
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

	// Load database configuration
	cfg.Database.Type = getEnv("DB_TYPE", "postgres")
	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnvAsInt("DB_PORT", 5432)
	cfg.Database.User = getEnv("DB_USER", "user")
	cfg.Database.Password = getEnv("DB_PASSWORD", "password")
	cfg.Database.Name = getEnv("DB_NAME", "myanimeapi")

	// Load logging configuration
	cfg.Logging.Level = getEnv("LOG_LEVEL", "info")
	cfg.Logging.File = getEnv("LOG_FILE", "/var/log/myapp.log")

	// Load API configuration
	cfg.API.Key = getEnv("API_KEY", "your_api_key")

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
