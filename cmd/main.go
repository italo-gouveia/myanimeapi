// cmd/main.go
// Package main is the entry point for the MyAnimeAPI application.
// It loads the configuration, initializes the database connection, and starts the server.
// The application uses the Gorilla Mux router, GORM for database interactions, and supports graceful shutdown.
// It also sets up CORS, registers routes, and initializes handlers for managing anime, users, reviews, and authentication.
//
// The application is designed to be configurable via environment variables or a configuration file.
// It supports PostgreSQL as the database and provides a RESTful API for managing anime and reviews.
//
// Example usage:
//   go run cmd/main.go
//
// The server listens on the port specified in the configuration and can be accessed via HTTP.

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"myanimeapi/api/database"
	"myanimeapi/api/middleware"
	"myanimeapi/api/routes"
	"myanimeapi/api/services"
	"myanimeapi/internal/config"
	"myanimeapi/internal/db"
	"myanimeapi/internal/logger"

	gorillahandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	version = "1.9.0"
)

// @title MyAnimeAPI
// @version 1.0
// @description This is a sample API for managing anime and reviews.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@myanimeapi.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /v1
// @schemes http
//
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Use the format "Bearer <JWT_TOKEN>". Example: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
// @security ApiKeyAuth
func main() {
	// Initialize logger with environment-based level
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	log := logger.New()
	logger.SetDefaultLogger(log)
	log.WithField("level", logLevel).Info("Logger initialized")

	// Validate required secrets before anything else
	if secret := os.Getenv("JWT_SECRET_KEY"); len(secret) < 32 {
		log.Error("JWT_SECRET_KEY must be set and at least 32 characters long")
		os.Exit(1)
	}

	// Load configuration
	cfg := config.LoadConfig()
	if cfg == nil {
		log.Error("Error loading config")
		os.Exit(1)
	}
	log.WithField("version", version).Info("Configuration loaded successfully")

	// Build the connection string for PostgreSQL
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.Database.Host, cfg.Database.Port)

	// Open a connection to the database
	gormDB, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.WithField("error", err.Error()).Error("Error opening database connection")
		os.Exit(1)
	}
	log.Info("Database connection established successfully")

	// Wrap the *gorm.DB instance in the GormDB struct
	dbWrapper := db.NewGormDB(gormDB)

	// AutoMigrate the database schema
	err = database.SetupDatabase(gormDB)
	if err != nil {
		log.WithField("error", err.Error()).Error("Error setting up database")
		os.Exit(1)
	}
	log.Info("Database setup completed successfully")

	// Initialize storage service
	var storageSvc services.StorageServiceInterface
	if os.Getenv("ENV") == "production" {
		// Use S3 storage in production
		s3Strategy, errS3 := services.NewS3StorageStrategy(
			os.Getenv("AWS_REGION"),
			os.Getenv("AWS_S3_BUCKET"),
			os.Getenv("AWS_S3_BASE_URL"),
			"uploads",
		)
		if errS3 != nil {
			log.WithField("error", errS3.Error()).Error("Failed to initialize S3 storage")
			os.Exit(1)
		}
		storageSvc = services.NewStorageService(s3Strategy)
	} else {
		// Use local storage in development
		localStrategy, errLocal := services.NewLocalStorageStrategy(
			".",
			"/media",
		)
		if errLocal != nil {
			log.WithField("error", errLocal.Error()).Error("Failed to initialize local storage")
			os.Exit(1)
		}
		storageSvc = services.NewStorageService(localStrategy)
	}
	log.Info("Storage service initialized successfully")

	// Create a new router
	router := mux.NewRouter()

	// Add request ID middleware
	router.Use(middleware.RequestIDMiddleware())

	// Determine environment and configure accordingly
	if os.Getenv("ENV") == "production" {
		log.Info("Running in production mode")
		router.Use(middleware.HTTPSMiddleware)
	}

	// Register all routes
	routes.RegisterRoutes(router, os.Getenv("SWAGGER_URL"), dbWrapper, storageSvc, version)
	log.Info("Routes registered successfully")

	// Serve static files for media
	fs := http.FileServer(http.Dir("uploads"))
	router.PathPrefix("/api/media/").Handler(http.StripPrefix("/api/media/", fs))
	log.Info("Static file serving configured")

	// Configure CORS — require explicit origins; never default to wildcard
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		log.Error("ALLOWED_ORIGINS must be set (e.g. http://localhost:3000). Refusing to start with wildcard CORS.")
		os.Exit(1)
	}

	// Split allowed origins by comma and trim spaces
	origins := strings.Split(allowedOrigins, ",")
	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}

	corsHandler := gorillahandlers.CORS(
		gorillahandlers.AllowedOrigins(origins),
		gorillahandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}),
		gorillahandlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Start the server with graceful shutdown
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:           corsHandler(router),
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Channel to listen for interrupt signals
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		log.WithFields(map[string]interface{}{
			"version": version,
			"port":    cfg.Server.Port,
		}).Info("Starting MyAnimeAPI")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithField("error", err.Error()).Error("Error starting server")
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	<-done
	log.Info("Server is shutting down...")

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.WithField("error", err.Error()).Error("Error shutting down server")
		os.Exit(1)
	}

	// Close database connections
	if sqlDB, err := gormDB.DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			log.WithField("error", err.Error()).Error("Error closing database connection")
		}
	}
	log.Info("Server shut down gracefully")
}
