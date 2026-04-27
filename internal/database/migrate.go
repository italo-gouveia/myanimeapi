// Package database provides database migration support via golang-migrate.
// SQL migration files are embedded directly into the binary so the application
// works without any external file system dependency.
package database

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"gorm.io/gorm"

	"myanimeapi/internal/logger"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// RunMigrations applies all pending SQL migrations to the given database.
// It is idempotent: already-applied migrations are skipped, and
// migrate.ErrNoChange is not treated as an error.
//
// The migrations are embedded in the binary at compile time from the
// internal/database/migrations/ directory.
func RunMigrations(db *gorm.DB) error {
	log := logger.New()

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("getting underlying sql.DB: %w", err)
	}

	// Build the source driver from embedded files.
	src, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("creating iofs migration source: %w", err)
	}

	// Build the database driver.
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("creating postgres migrate driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		return fmt.Errorf("initialising migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("applying migrations: %w", err)
	}

	version, _, _ := m.Version()
	log.WithField("version", version).Info("Database migrations applied successfully")
	return nil
}

// MigrateDown rolls back all applied migrations (useful in tests).
func MigrateDown(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("getting underlying sql.DB: %w", err)
	}

	src, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("creating iofs migration source: %w", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("creating postgres migrate driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		return fmt.Errorf("initialising migrate instance: %w", err)
	}

	if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("rolling back migrations: %w", err)
	}

	return nil
}
