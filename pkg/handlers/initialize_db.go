// pkg/handlers/initialize_db.go
// Package handlers provides functionality to initialize and manage the global database instance for the MyAnimeAPI application.
// It defines a function to set the global database instance, which is used by the handlers to interact with the database.
// The package ensures that the database instance is accessible globally within the application.
//
// Example usage:
//
//	db := // initialize your database connection
//	handlers.InitializeDB(db)
//
//	// Later, retrieve the database instance
//	dbInstance := handlers.GetDB(context.Background())
package handlers

import (
	"context"
	"log"
	"myanimeapi/internal/db"
)

// database is the global database instance used by the handlers.
var database db.DBInterface

// InitializeDB sets the global database instance to the provided dbInstance.
// This function should be called during application startup to ensure the database instance is available globally.
//
// Example:
//
//	db := // initialize your database connection
//	handlers.InitializeDB(db)
func InitializeDB(dbInstance db.DBInterface) {
	log.Println("Initializing global database instance")
	database = dbInstance
}

// GetDB returns the global database instance.
// It accepts a context and returns the database instance set by InitializeDB.
// This function is used by handlers to interact with the database.
//
// Example:
//
//	dbInstance := handlers.GetDB(context.Background())
func GetDB(ctx context.Context) db.DBInterface {
	return database
}
