// pkg/database/database.go
package database

import (
	"myanimeapi/pkg/models"

	"gorm.io/gorm"
)

// SetupDatabase initializes the database schema
func SetupDatabase(db *gorm.DB) error {
	// AutoMigrate to create/update schema
	err := db.AutoMigrate(&models.User{}, &models.Anime{}, &models.Review{})
	if err != nil {
		return err
	}
	return nil
}
