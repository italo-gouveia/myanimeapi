// api/adapters/http/initialize_db.go
// Package httphandler provides functionality to initialize and manage the global database instance.
package httphandler

import (
	"log"
	"myanimeapi/internal/db"
)

// database is the global database instance used by the handlers.
var database db.DBInterface

// InitializeDB sets the global database instance to the provided dbInstance.
func InitializeDB(dbInstance db.DBInterface) {
	log.Println("Initializing global database instance")
	database = dbInstance
}

// GetDB returns the global database instance.
func GetDB() db.DBInterface {
	return database
}
