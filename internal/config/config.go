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

	cfg.Database.User = getEnv("DB_USER", "defaultuser")
	cfg.Database.Password = getEnv("DB_PASSWORD", "defaultpassword")
	cfg.Database.Name = getEnv("DB_NAME", "defaultdb")
	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnvAsInt("DB_PORT", 5432)

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
