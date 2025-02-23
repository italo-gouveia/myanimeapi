// internal/db/gorm_db.go
// This package defines the GormDB struct that wraps a *gorm.DB instance to implement the DBInterface interface.
// It provides methods to interact with the database using the GORM library.
// It is used by the server to perform database operations.
// It is a simple wrapper around the GORM library to provide a consistent interface for database operations.
// It implements the DBInterface interface defined in the db_interface.go file.
// It is used by the server to interact with the database using the GORM library.
package db

import (
	"context"

	"gorm.io/gorm"
)

// GormDB wraps a *gorm.DB instance to implement DBInterface
type GormDB struct {
	db *gorm.DB
}

// NewGormDB creates a new GormDB instance
func NewGormDB(db *gorm.DB) *GormDB {
	return &GormDB{db: db}
}

func (g *GormDB) First(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB {
	return g.db.WithContext(ctx).First(dest, conds...)
}

func (g *GormDB) Find(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB {
	return g.db.WithContext(ctx).Find(dest, conds...)
}

func (g *GormDB) Where(ctx context.Context, query interface{}, args ...interface{}) *gorm.DB {
	return g.db.WithContext(ctx).Where(query, args...)
}

func (g *GormDB) Create(ctx context.Context, value interface{}) *gorm.DB {
	return g.db.WithContext(ctx).Create(value)
}

func (g *GormDB) Save(ctx context.Context, value interface{}) *gorm.DB {
	return g.db.WithContext(ctx).Save(value)
}

func (g *GormDB) Delete(ctx context.Context, value interface{}, conds ...interface{}) *gorm.DB {
	return g.db.WithContext(ctx).Delete(value, conds...)
}

func (g *GormDB) Preload(column string, ctx context.Context, conditions ...interface{}) *gorm.DB {
	return g.db.WithContext(ctx).Preload(column, conditions...)
}

func (g *GormDB) Offset(offset int) *gorm.DB {
	return g.db.Offset(offset)
}

func (g *GormDB) Limit(limit int) *gorm.DB {
	return g.db.Limit(limit)
}

func (g *GormDB) Unscoped(ctx context.Context) *gorm.DB {
	return g.db.WithContext(ctx).Unscoped()
}

func (g *GormDB) Begin(ctx context.Context) *gorm.DB {
	return g.db.WithContext(ctx).Begin()
}

func (g *GormDB) Commit(ctx context.Context) *gorm.DB {
	return g.db.WithContext(ctx).Commit()
}

func (g *GormDB) Rollback(ctx context.Context) *gorm.DB {
	return g.db.WithContext(ctx).Rollback()
}

func (g *GormDB) GetError() error {
	return g.db.Error
}
