// pkg/handlers/handlers.go
package handlers

import (
	"myanimeapi/internal/db"
)

var database db.DBInterface

// InitializeDB sets the global database instance
func InitializeDB(dbInstance db.DBInterface) {
	database = dbInstance
}
