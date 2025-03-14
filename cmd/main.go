// cmd/main.go
// This file is the entry point for the application. It loads the configuration, initializes the database connection, and starts the server.
// It imports the necessary packages and registers the routes.
// It uses the gorilla/mux package to create a new router.
// It uses the gorm package to open a connection to the database.
// It uses the database package to setup the database schema.
// It uses the handlers package to initialize the global DB variable.
// It uses the routes package to register the routes.
// It uses the config package to load the configuration.
// It uses the log package to log messages.
// It uses the fmt package to format strings.
// It uses the net/http package to start the server.
// It uses the gorm.io/driver/postgres package to open a connection to the PostgreSQL database.
// It uses the gorm.io/gorm package to interact with the database.
// It uses the myanimeapi/internal/config package to load the configuration.
// It uses the myanimeapi/internal/routes package to register the routes.
// It uses the myanimeapi/pkg/database package to setup the database schema.
// It uses the myanimeapi/pkg/handlers package to initialize the global DB variable.
// It uses the myanimeapi/pkg/handlers package to handle the requests.
// It uses the os package to listen for interrupt signals.
// It uses the os/signal package to listen for interrupt signals.
// It uses the syscall package to listen for interrupt signals.
// It uses the time package to create a context with a timeout for graceful shutdown.
// It uses the context package to create a context with a timeout for graceful shutdown.
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

	"myanimeapi/internal/config"
	"myanimeapi/internal/db"
	"myanimeapi/internal/routes"
	"myanimeapi/pkg/database"
	myhandlers "myanimeapi/pkg/handlers" // Alias for your custom handlers package

	gorillahandlers "github.com/gorilla/handlers" // Alias for Gorilla's handlers package
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	// Register all routes
	swaggerURL := "http://localhost:8080/swagger/doc.json" // or fetch from config
	routes.RegisterRoutes(router, swaggerURL, animeHandler, userHandler, reviewHandler, authHandler)
	log.Println("Routes registered successfully")

	// Configure CORS
	corsHandler := gorillahandlers.CORS(
		gorillahandlers.AllowedOrigins([]string{"*"}), // Allow all origins
		gorillahandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}),
		gorillahandlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Start the server with graceful shutdown
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: corsHandler(router),
	}

	// Channel to listen for interrupt signals
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		log.Printf("Server started on :%d", cfg.Server.Port)
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
