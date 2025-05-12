package models

import "time"

// BaseModel represents the common fields used across all models.
// It provides a standardized structure for database entities with:
//   - Unique identifier (ID)
//   - Timestamps for creation and updates
//   - Soft delete functionality
//
// This is a Swagger-friendly version of gorm.Model that includes
// proper JSON tags and example values for API documentation.
//
// Example:
//
//	{
//	  "id": 1,
//	  "created_at": "2025-02-20T19:27:00Z",
//	  "updated_at": "2025-02-20T19:27:00Z",
//	  "deleted_at": "2025-02-21T10:00:00Z"  // Only present when soft deleted
//	}
type BaseModel struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`                                 // Unique identifier
	CreatedAt time.Time `json:"created_at" example:"2025-02-20T19:27:00Z"`                        // Creation timestamp
	UpdatedAt time.Time `json:"updated_at" example:"2025-02-20T19:27:00Z"`                        // Update timestamp
	DeletedAt time.Time `json:"deleted_at,omitempty" gorm:"index" example:"2025-02-20T19:27:00Z"` // Soft delete timestamp
}
