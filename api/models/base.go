package models

import "time"

// BaseModel represents the common fields used across all models.
// This is a Swagger-friendly version of gorm.Model.
type BaseModel struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`                                 // Unique identifier
	CreatedAt time.Time `json:"created_at" example:"2025-02-20T19:27:00Z"`                        // Creation timestamp
	UpdatedAt time.Time `json:"updated_at" example:"2025-02-20T19:27:00Z"`                        // Update timestamp
	DeletedAt time.Time `json:"deleted_at,omitempty" gorm:"index" example:"2025-02-20T19:27:00Z"` // Soft delete timestamp
}
