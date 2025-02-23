// pkg/handlers/initialize_db.go
// This package defines a function to set the global database instance.
// It is used by the server to set the global database instance.
// It provides a function to set the global database instance to be used by the handlers.
package handlers

import (
	"context"
	"log"
	"myanimeapi/internal/db"
)

var database db.DBInterface

// InitializeDB sets the global database instance
func InitializeDB(dbInstance db.DBInterface) {
	log.Println("Initializing global database instance")
	database = dbInstance
}

// GetDB returns the global database instance
func GetDB(ctx context.Context) db.DBInterface {
	return database
}
