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
	cfg := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "user"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "dbname"),
		},
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
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
