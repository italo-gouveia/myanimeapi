package database

import (
	"myanimeapi/pkg/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// SetupDatabase initializes the database schema
func SetupDatabase(db *gorm.DB) {
	// AutoMigrate to create/update schema
	db.AutoMigrate(&models.User{}, &models.Anime{}, &models.Review{})
}
