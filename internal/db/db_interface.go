// internal/db/db.go
package db

import "gorm.io/gorm"

// DBInterface defines the methods required for database interactions
type DBInterface interface {
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Find(dest interface{}, conds ...interface{}) *gorm.DB
	Where(query interface{}, args ...interface{}) *gorm.DB
	Create(value interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	Delete(value interface{}, conds ...interface{}) *gorm.DB
	Preload(column string, conditions ...interface{}) *gorm.DB
	Offset(offset int) *gorm.DB
}
