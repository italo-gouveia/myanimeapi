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
	"myanimeapi/api/auth"
	"myanimeapi/api/models"
	"os"

	"gorm.io/gorm"
	gormErrors "gorm.io/gorm/logger"
)

// SetupDatabase initializes the database schema by automatically migrating the defined models
// and creates an admin user if specified by environment variables and not already present.
// It uses GORM's AutoMigrate function to create or update the database tables.
// If the migration or admin creation fails, an error is returned.
//
// Example:
//
//	err := database.SetupDatabase(db)
//	if err != nil {
//	    log.Fatalf("Error setting up database schema: %v", err)
//	}
func SetupDatabase(db *gorm.DB) error {
	// AutoMigrate to create/update schema
	err := db.AutoMigrate(
		&models.User{},
		&models.Anime{},
		&models.Review{},
		&models.Favorite{},
		&models.Genre{},
		&models.Tag{},
		&models.PasswordResetToken{},
		&models.MediaAttachment{},
	)
	if err != nil {
		log.Printf("Error setting up database schema during AutoMigrate: %v", err)
		return err
	}
	log.Println("Database schema auto-migration completed successfully")

	// Create admin user if needed
	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if adminUsername == "" || adminEmail == "" || adminPassword == "" {
		log.Println("ADMIN_USERNAME, ADMIN_EMAIL, or ADMIN_PASSWORD not set. Skipping admin user creation.")
		return nil // Not an error, just skipping
	}

	// Check if admin user with email already exists
	var existingUserByEmail models.User
	err = db.Where("email = ?", adminEmail).First(&existingUserByEmail).Error
	if err == nil {
		log.Printf("Admin user with email '%s' already exists (ID: %d). Skipping creation.", adminEmail, existingUserByEmail.ID)
		return nil
	}
	// Only proceed if error is RecordNotFound, otherwise it's an unexpected DB error
	if err != gorm.ErrRecordNotFound && err != gormErrors.ErrRecordNotFound {
		log.Printf("Error checking for admin user by email '%s': %v. Skipping admin creation.", adminEmail, err)
		return err // Return actual DB error
	}

	// Check if admin user with username already exists
	var existingUserByUsername models.User
	err = db.Where("username = ?", adminUsername).First(&existingUserByUsername).Error
	if err == nil {
		log.Printf("Admin user with username '%s' already exists (ID: %d). Skipping creation.", adminUsername, existingUserByUsername.ID)
		return nil
	}
	if err != gorm.ErrRecordNotFound && err != gormErrors.ErrRecordNotFound {
		log.Printf("Error checking for admin user by username '%s': %v. Skipping admin creation.", adminUsername, err)
		return err // Return actual DB error
	}

	log.Printf("Creating admin user: Username='%s', Email='%s'", adminUsername, adminEmail)

	hashedPassword, err := auth.HashPassword(adminPassword)
	if err != nil {
		log.Printf("Error hashing admin password: %v", err)
		return err
	}

	adminUser := models.User{
		Username: adminUsername,
		Email:    adminEmail,
		Password: hashedPassword,
		IsAdmin:  true,
		IsActive: true,
	}

	if err := db.Create(&adminUser).Error; err != nil {
		log.Printf("Error creating admin user '%s': %v", adminUsername, err)
		return err
	}

	log.Printf("Admin user '%s' created successfully with ID %d.", adminUsername, adminUser.ID)
	return nil
}
