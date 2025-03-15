// internal/db/gorm_db.go
// Package db provides database interaction utilities.
// This file defines the GormDB struct, which wraps a *gorm.DB instance to implement the DBInterface interface.
// It provides methods to interact with the database using the GORM library.
// The GormDB struct is used by the server to perform database operations in a consistent and abstracted manner.
package db

import (
	"context"

	"gorm.io/gorm"
)

// GormDB wraps a *gorm.DB instance to implement the DBInterface interface.
// It provides methods for CRUD operations, transactions, and context handling.
type GormDB struct {
	db *gorm.DB // Underlying GORM database instance
}

// NewGormDB creates a new GormDB instance with the provided *gorm.DB instance.
// It returns a pointer to the GormDB struct.
func NewGormDB(db *gorm.DB) *GormDB {
	return &GormDB{db: db}
}

// WithContext attaches the given context to the underlying *gorm.DB instance.
// It returns a new *gorm.DB instance with the context applied.
func (g *GormDB) WithContext(ctx context.Context) *gorm.DB {
	return g.db.WithContext(ctx)
}

// First retrieves the first record matching the given conditions.
// It stores the result in `dest` and returns a *gorm.DB instance.
func (g *GormDB) First(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB {
	return g.db.WithContext(ctx).First(dest, conds...)
}

// Find retrieves all records matching the given conditions.
// It stores the results in `dest` and returns a *gorm.DB instance.
func (g *GormDB) Find(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB {
	return g.db.WithContext(ctx).Find(dest, conds...)
}

// Where filters records based on the given query and arguments.
// It returns a *gorm.DB instance for chaining.
func (g *GormDB) Where(ctx context.Context, query interface{}, args ...interface{}) *gorm.DB {
	return g.db.WithContext(ctx).Where(query, args...)
}

// Create inserts a new record into the database.
// It returns a *gorm.DB instance for chaining.
func (g *GormDB) Create(ctx context.Context, value interface{}) *gorm.DB {
	return g.db.WithContext(ctx).Create(value)
}

// Save updates or inserts a record in the database.
// It returns a *gorm.DB instance for chaining.
func (g *GormDB) Save(ctx context.Context, value interface{}) *gorm.DB {
	return g.db.WithContext(ctx).Save(value)
}

// Delete removes a record from the database.
// It returns a *gorm.DB instance for chaining.
func (g *GormDB) Delete(ctx context.Context, value interface{}, conds ...interface{}) *gorm.DB {
	return g.db.WithContext(ctx).Delete(value, conds...)
}

// Preload eagerly loads related records using the specified column and conditions.
// It returns a *gorm.DB instance for chaining.
func (g *GormDB) Preload(column string, ctx context.Context, conditions ...interface{}) *gorm.DB {
	return g.db.WithContext(ctx).Preload(column, conditions...)
}

// Offset skips a specified number of records in the result set.
// It returns a *gorm.DB instance for chaining.
func (g *GormDB) Offset(offset int) *gorm.DB {
	return g.db.Offset(offset)
}

// Limit restricts the number of records returned in the result set.
// It returns a *gorm.DB instance for chaining.
func (g *GormDB) Limit(limit int) *gorm.DB {
	return g.db.Limit(limit)
}

// Unscoped includes soft-deleted records in the query results.
// It returns a *gorm.DB instance for chaining.
func (g *GormDB) Unscoped(ctx context.Context) *gorm.DB {
	return g.db.WithContext(ctx).Unscoped()
}

// Begin starts a new database transaction.
// It returns a *gorm.DB instance representing the transaction.
func (g *GormDB) Begin(ctx context.Context) *gorm.DB {
	return g.db.WithContext(ctx).Begin()
}

// Commit commits the current transaction.
// It returns a *gorm.DB instance for chaining.
func (g *GormDB) Commit(ctx context.Context) *gorm.DB {
	return g.db.WithContext(ctx).Commit()
}

// Rollback rolls back the current transaction.
// It returns a *gorm.DB instance for chaining.
func (g *GormDB) Rollback(ctx context.Context) *gorm.DB {
	return g.db.WithContext(ctx).Rollback()
}

// GetError returns the last error encountered during database operations.
func (g *GormDB) GetError() error {
	return g.db.Error
}
