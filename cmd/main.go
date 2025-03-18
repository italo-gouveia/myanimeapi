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

	"myanimeapi/api/database"
	myhandlers "myanimeapi/api/handlers" // Alias for your custom handlers package
	"myanimeapi/api/routes"
	"myanimeapi/internal/config"
	"myanimeapi/internal/db"

	gorillahandlers "github.com/gorilla/handlers" // Alias for Gorilla's handlers package
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	VERSION = "1.0.0" // Current version of the API
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

	// Initialize handlers with the database instance
	animeHandler := myhandlers.NewAnimeHandler(dbWrapper)
	userHandler := myhandlers.NewUserHandler(dbWrapper)
	reviewHandler := myhandlers.NewReviewHandler(dbWrapper)
	authHandler := myhandlers.NewAuthHandler(dbWrapper)
	log.Println("Handlers initialized successfully")

	// Create a new router
	router := mux.NewRouter()

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
	routes.RegisterRoutes(router, swaggerURL, animeHandler, userHandler, reviewHandler, authHandler, VERSION)
	log.Println("Routes registered successfully")

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
