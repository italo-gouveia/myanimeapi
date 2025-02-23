// pkg/database/database.go
// This package defines the database schema and initializes the database.
// It is used by the server to setup the database schema.
// It initializes the database schema by calling the AutoMigrate function on the database instance.
// It returns an error if the AutoMigrate function fails.
package database

import (
	"log"
	"myanimeapi/pkg/models"

	"gorm.io/gorm"
)

// SetupDatabase initializes the database schema
func SetupDatabase(db *gorm.DB) error {
	// AutoMigrate to create/update schema
	err := db.AutoMigrate(&models.User{}, &models.Anime{}, &models.Review{})
	if err != nil {
		log.Printf("Error setting up database schema: %v", err)
		return err
	}
	log.Println("Database schema initialized successfully")
	return nil
}
