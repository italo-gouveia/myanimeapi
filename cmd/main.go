package main

import (
	"fmt"
	"log"
	"net/http"

	"myanimeapi/internal/config"
	"myanimeapi/internal/routes"
	"myanimeapi/pkg/handlers" // Ensure you import the handlers package

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig() // Expecting one return value of type *Config
	if cfg == nil {
		log.Fatal("Error loading config")
	}

	// Build the connection string for PostgreSQL
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.Database.Host, cfg.Database.Port)

	// Open a connection to the database
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	// Initialize the global DB variable in handlers
	handlers.InitializeDB(db)

	// Create a new router
	router := mux.NewRouter()

	// Register all routes
	routes.RegisterRoutes(router)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", router))
}
