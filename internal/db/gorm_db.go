// internal/db/gorm_db.go
package db

import "gorm.io/gorm"

// GormDB wraps a *gorm.DB instance to implement DBInterface
type GormDB struct {
	db *gorm.DB
}

// NewGormDB creates a new GormDB instance
func NewGormDB(db *gorm.DB) *GormDB {
	return &GormDB{db: db}
}

func (g *GormDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	return g.db.First(dest, conds...)
}

func (g *GormDB) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	return g.db.Find(dest, conds...)
}

func (g *GormDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return g.db.Where(query, args...)
}

func (g *GormDB) Create(value interface{}) *gorm.DB {
	return g.db.Create(value)
}

func (g *GormDB) Save(value interface{}) *gorm.DB {
	return g.db.Save(value)
}

func (g *GormDB) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	return g.db.Delete(value, conds...)
}

func (g *GormDB) Preload(column string, conditions ...interface{}) *gorm.DB {
	return g.db.Preload(column, conditions...)
}

func (g *GormDB) Offset(offset int) *gorm.DB {
	return g.db.Offset(offset)
}
