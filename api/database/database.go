// api/database/database.go
// Package database provides post-migration database initialisation for the
// MyAnimeAPI application.  Schema management is handled separately by
// internal/database.RunMigrations, which must be called before SetupDatabase.
//
// Example usage:
//
//	// 1. Run SQL migrations first
//	if err := internaldb.RunMigrations(gormDB); err != nil {
//	    log.Fatalf("migrations: %v", err)
//	}
//	// 2. Seed the admin user (if env vars are set)
//	if err := database.SetupDatabase(gormDB); err != nil {
//	    log.Fatalf("setup: %v", err)
//	}
package database

import (
	"myanimeapi/api/auth"
	"myanimeapi/api/models"
	"os"

	"myanimeapi/internal/logger"

	"gorm.io/gorm"
	gormErrors "gorm.io/gorm/logger"
)

// SetupDatabase creates the admin user when the ADMIN_USERNAME, ADMIN_EMAIL,
// and ADMIN_PASSWORD environment variables are all set and the user does not
// already exist.  It is idempotent — running it multiple times is safe.
func SetupDatabase(db *gorm.DB) error {
	log := logger.New()

	// Create admin user if needed
	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if adminUsername == "" || adminEmail == "" || adminPassword == "" {
		log.Info("ADMIN_USERNAME, ADMIN_EMAIL, or ADMIN_PASSWORD not set. Skipping admin user creation.")
		return nil
	}

	// Check if admin user with email already exists
	var existingUserByEmail models.User
	err := db.Where("email = ?", adminEmail).First(&existingUserByEmail).Error
	if err == nil {
		log.WithFields(map[string]interface{}{
			"email":  adminEmail,
			"userID": existingUserByEmail.ID,
		}).Info("Admin user with this email already exists. Skipping creation.")
		return nil
	}
	if err != gorm.ErrRecordNotFound && err != gormErrors.ErrRecordNotFound {
		log.WithFields(map[string]interface{}{
			"email": adminEmail,
			"error": err.Error(),
		}).Error("Error checking for admin user by email. Skipping admin creation.")
		return err
	}

	// Check if admin user with username already exists
	var existingUserByUsername models.User
	err = db.Where("username = ?", adminUsername).First(&existingUserByUsername).Error
	if err == nil {
		log.WithFields(map[string]interface{}{
			"username": adminUsername,
			"userID":   existingUserByUsername.ID,
		}).Info("Admin user with this username already exists. Skipping creation.")
		return nil
	}
	if err != gorm.ErrRecordNotFound && err != gormErrors.ErrRecordNotFound {
		log.WithFields(map[string]interface{}{
			"username": adminUsername,
			"error":    err.Error(),
		}).Error("Error checking for admin user by username. Skipping admin creation.")
		return err
	}

	log.WithFields(map[string]interface{}{
		"username": adminUsername,
		"email":    adminEmail,
	}).Info("Creating admin user.")

	hashedPassword, err := auth.HashPassword(adminPassword)
	if err != nil {
		log.WithField("error", err.Error()).Error("Error hashing admin password")
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
		log.WithFields(map[string]interface{}{
			"username": adminUsername,
			"error":    err.Error(),
		}).Error("Error creating admin user")
		return err
	}

	log.WithFields(map[string]interface{}{
		"username": adminUsername,
		"userID":   adminUser.ID,
	}).Info("Admin user created successfully.")
	return nil
}
