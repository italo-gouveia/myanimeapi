// internal/db/db_interface.go
package db

import (
	"context"

	"gorm.io/gorm"
)

// DBInterface defines the methods required for database interactions
type DBInterface interface {
	// Basic CRUD operations
	First(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB
	Find(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB
	Where(ctx context.Context, query interface{}, args ...interface{}) *gorm.DB
	Create(ctx context.Context, value interface{}) *gorm.DB
	Save(ctx context.Context, value interface{}) *gorm.DB
	Delete(ctx context.Context, value interface{}, conds ...interface{}) *gorm.DB
	Preload(column string, ctx context.Context, conditions ...interface{}) *gorm.DB
	Offset(offset int) *gorm.DB
	Limit(limit int) *gorm.DB
	Unscoped(ctx context.Context) *gorm.DB

	// Transaction support
	Begin(ctx context.Context) *gorm.DB
	Commit(ctx context.Context) *gorm.DB
	Rollback(ctx context.Context) *gorm.DB

	// Context support
	WithContext(ctx context.Context) *gorm.DB // Add this method

	// Error handling
	GetError() error
}
