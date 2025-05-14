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
	"log"
	"net/http"
	"os"
	"os/signal"
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

	// Alias for Gorilla's handlers package
	gorillahandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	VERSION = "1.4.0" // Current version of the API
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
	// Load configuration
	cfg := config.LoadConfig()
	if cfg == nil {
		log.Fatal("Error loading config")
	}
	log.Println("Configuration loaded successfully")

	// Build the connection string for PostgreSQL
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.Database.Host, cfg.Database.Port)

	// Open a connection to the database
	gormDB, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}
	log.Println("Database connection established successfully")

	// Wrap the *gorm.DB instance in the GormDB struct
	dbWrapper := db.NewGormDB(gormDB)

	// AutoMigrate the database schema
	err = database.SetupDatabase(gormDB)
	if err != nil {
		log.Fatalf("Error migrating database schema: %v", err)
	}
	log.Println("Database schema migrated successfully")

	// Initialize storage service
	var storageSvc *services.StorageService
	if os.Getenv("ENV") == "production" {
		// Use S3 storage in production
		s3Strategy, err := services.NewS3StorageStrategy(
			os.Getenv("AWS_REGION"),
			os.Getenv("AWS_S3_BUCKET"),
			os.Getenv("AWS_S3_BASE_URL"),
			"uploads",
		)
		if err != nil {
			log.Fatalf("Failed to initialize S3 storage: %v", err)
		}
		storageSvc = services.NewStorageService(s3Strategy)
	} else {
		// Use local storage in development
		localStrategy, err := services.NewLocalStorageStrategy(
			".",
			"/media",
		)
		if err != nil {
			log.Fatalf("Failed to initialize local storage: %v", err)
		}
		storageSvc = services.NewStorageService(localStrategy)
	}
	log.Println("Storage service initialized successfully")

	// Create a new router
	router := mux.NewRouter()

	// Add request ID middleware
	router.Use(middleware.RequestIDMiddleware())

	// Determine environment
	env := os.Getenv("ENVIRONMENT")
	if env == "production" {
		log.Println("Running in production mode")
		// Apply HTTPS redirection middleware
		//router.Use(middleware.HTTPSRedirectMiddleware)
	}

	// Serve Swagger UI
	swaggerURL := os.Getenv("SWAGGER_URL")
	// Register all routes
	routes.RegisterRoutes(router, swaggerURL, dbWrapper, storageSvc, VERSION)
	log.Println("Routes registered successfully")

	// Serve static files for media
	fs := http.FileServer(http.Dir("uploads"))
	router.PathPrefix("/api/media/").Handler(http.StripPrefix("/api/media/", fs))
	log.Println("Static file serving configured")

	// Configure CORS
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	corsHandler := gorillahandlers.CORS(
		gorillahandlers.AllowedOrigins([]string{allowedOrigins}),
		gorillahandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}),
		gorillahandlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Start the server with graceful shutdown
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:           corsHandler(router),
		ReadHeaderTimeout: 10 * time.Second, // Add a timeout for reading headers
	}

	// Channel to listen for interrupt signals
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		log.Printf("Starting MyAnimeAPI version %s on :%d", VERSION, cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-done
	log.Println("Server is shutting down...")

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}
	log.Println("Server shut down gracefully")
}
