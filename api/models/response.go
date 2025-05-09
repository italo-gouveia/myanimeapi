package models

// SuccessResponse represents a generic success response.
// It includes a success message and optional data.
//
// Example:
//
//	{
//	  "message": "Operation completed successfully",
//	  "data": {
//	    "id": 1,
//	    "name": "Example"
//	  }
//	}
type SuccessResponse struct {
	Message string      `json:"message" example:"Operation completed successfully"` // Success message
	Data    interface{} `json:"data,omitempty"`                                     // Optional response data
}
