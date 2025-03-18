// internal/db/db_interface.go
// Package db provides database interaction utilities.
// This file defines the DBInterface, which abstracts the methods required for database operations.
// It uses the `gorm.io/gorm` package for database interactions.
package db

import (
	"context"

	"gorm.io/gorm"
)

// DBInterface defines the methods required for database interactions.
// It abstracts common database operations such as CRUD, transactions, and context handling.
// Implementations of this interface can be used to interact with different database backends.
type DBInterface interface {
	// First retrieves the first record matching the given conditions.
	// It stores the result in `dest` and returns a *gorm.DB instance.
	First(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB

	// Find retrieves all records matching the given conditions.
	// It stores the results in `dest` and returns a *gorm.DB instance.
	Find(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB

	// Where filters records based on the given query and arguments.
	// It returns a *gorm.DB instance for chaining.
	Where(ctx context.Context, query interface{}, args ...interface{}) *gorm.DB

	// Create inserts a new record into the database.
	// It returns a *gorm.DB instance for chaining.
	Create(ctx context.Context, value interface{}) *gorm.DB

	// Save updates or inserts a record in the database.
	// It returns a *gorm.DB instance for chaining.
	Save(ctx context.Context, value interface{}) *gorm.DB

	// Delete removes a record from the database.
	// It returns a *gorm.DB instance for chaining.
	Delete(ctx context.Context, value interface{}, conds ...interface{}) *gorm.DB

	// Preload eagerly loads related records using the specified column and conditions.
	// It returns a *gorm.DB instance for chaining.
	Preload(column string, ctx context.Context, conditions ...interface{}) *gorm.DB

	// Offset skips a specified number of records in the result set.
	// It returns a *gorm.DB instance for chaining.
	Offset(offset int) *gorm.DB

	// Limit restricts the number of records returned in the result set.
	// It returns a *gorm.DB instance for chaining.
	Limit(limit int) *gorm.DB

	// Unscoped includes soft-deleted records in the query results.
	// It returns a *gorm.DB instance for chaining.
	Unscoped(ctx context.Context) *gorm.DB

	// Begin starts a new database transaction.
	// It returns a *gorm.DB instance representing the transaction.
	Begin(ctx context.Context) *gorm.DB

	// Commit commits the current transaction.
	// It returns a *gorm.DB instance for chaining.
	Commit(ctx context.Context) *gorm.DB

	// Rollback rolls back the current transaction.
	// It returns a *gorm.DB instance for chaining.
	Rollback(ctx context.Context) *gorm.DB

	// WithContext attaches the given context to the database instance.
	// It returns a *gorm.DB instance for chaining.
	WithContext(ctx context.Context) *gorm.DB

	// GetError returns the last error encountered during database operations.
	GetError() error

	// IsHealthy checks if the database connection is healthy.
	// It pings the database to verify connectivity and returns true if successful.
	// If the connection is unhealthy, it logs the error and returns false.
	IsHealthy() bool
}
