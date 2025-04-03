// api/database/database.go
// Package database provides functionality to initialize and manage the database schema for the MyAnimeAPI application.
// It uses GORM (Go Object-Relational Mapping) to automatically create or update the database schema based on the defined models.
// The package is responsible for ensuring the database is properly set up when the application starts.
//
// Example usage:
//
//	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
//	if err != nil {
//	    log.Fatalf("Error opening database connection: %v", err)
//	}
//
//	err = database.SetupDatabase(db)
//	if err != nil {
//	    log.Fatalf("Error setting up database schema: %v", err)
//	}
package database

import (
	"log"
	"myanimeapi/api/models"

	"gorm.io/gorm"
)

// SetupDatabase initializes the database schema by automatically migrating the defined models.
// It uses GORM's AutoMigrate function to create or update the database tables for the User, Anime, Review, and Favorite models.
// If the migration fails, an error is returned.
//
// Example:
//
//	err := database.SetupDatabase(db)
//	if err != nil {
//	    log.Fatalf("Error setting up database schema: %v", err)
//	}
func SetupDatabase(db *gorm.DB) error {
	// AutoMigrate to create/update schema
	err := db.AutoMigrate(&models.User{}, &models.Anime{}, &models.Review{}, &models.Favorite{})
	if err != nil {
		log.Printf("Error setting up database schema: %v", err)
		return err
	}
	log.Println("Database schema initialized successfully")
	return nil
}
