package config

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/vault/api"
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
	JWTKey   string
}

func LoadConfig() *Config {
	// Initialize Vault client
	client, err := api.NewClient(&api.Config{
		Address: "http://vault:8200", // Use the service name "vault" in Docker
	})
	if err != nil {
		panic(fmt.Errorf("failed to create Vault client: %v", err))
	}

	// Set the Vault token (use the root token for development)
	client.SetToken("root")

	// Read the secrets from Vault
	secret, err := client.Logical().Read("secret/data/myapp")
	if err != nil {
		panic(fmt.Errorf("failed to read secret from Vault: %v", err))
	}

	// Extract the secret values
	var cfg Config
	if secret != nil && secret.Data != nil {
		data, ok := secret.Data["data"].(map[string]interface{})
		if !ok {
			panic("invalid secret data format")
		}

		// Server configuration
		cfg.Server.Host = getString(data, "SERVER_HOST", "localhost")
		cfg.Server.Port = getInt(data, "SERVER_PORT", 8080)

		// Database configuration
		cfg.Database.Type = getString(data, "DB_TYPE", "postgres")
		cfg.Database.Host = getString(data, "DB_HOST", "db")
		cfg.Database.Port = getInt(data, "DB_PORT", 5432)
		cfg.Database.User = getString(data, "DB_USER", "user")
		cfg.Database.Password = getString(data, "DB_PASSWORD", "password")
		cfg.Database.Name = getString(data, "DB_NAME", "myanimeapi")

		// Logging configuration
		cfg.Logging.Level = getString(data, "LOG_LEVEL", "info")
		cfg.Logging.File = getString(data, "LOG_FILE", "/var/log/myapp.log")

		// API configuration
		cfg.API.Key = getString(data, "API_KEY", "odagenius")

		// JWT secret key
		cfg.JWTKey = getString(data, "JWT_SECRET_KEY", "betweenearthandheaveniamthechosenone")
	} else {
		panic("secret not found")
	}

	return &cfg
}

func getString(data map[string]interface{}, key string, defaultValue string) string {
	if value, exists := data[key]; exists {
		return value.(string)
	}
	return defaultValue
}

func getInt(data map[string]interface{}, key string, defaultValue int) int {
	if value, exists := data[key]; exists {
		intValue, err := strconv.Atoi(value.(string))
		if err == nil {
			return intValue
		}
	}
	return defaultValue
}
