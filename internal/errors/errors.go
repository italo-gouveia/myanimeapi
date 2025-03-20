// internal/errors/errors.go
// Package errors provides custom error types and structured error responses for the MyAnimeAPI application.
// It defines error types such as NotFoundError and ValidationError, and provides utilities to return
// structured error responses with error codes, messages, and additional context.
//
// Example Usage:
//
//	// Writing a structured error response
//	errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid input", "The 'title' field is missing.")
//
//	// Creating a custom error
//	notFoundErr := errors.NewNotFoundError("Anime", "123", errors.ErrResourceNotFound)
//	errors.WriteErrorResponse(w, http.StatusNotFound, notFoundErr.Code, notFoundErr.Error(), "")
//
// The package also includes constants for common error codes to ensure consistency across the application.
package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrorResponse represents a structured error response.
// It includes an error code, a message, and optional details for additional context.
//
// Example:
//
//	{
//	  "error": {
//	    "code": "ERR-001",
//	    "message": "Invalid input: title is required",
//	    "details": "The 'title' field is missing in the request body."
//	  }
//	}
type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`              // A unique error code for categorization.
		Message string `json:"message"`           // A human-readable error message.
		Details string `json:"details,omitempty"` // Additional context or details about the error.
	} `json:"error"`
}

// NewErrorResponse creates a new ErrorResponse instance.
// It initializes the ErrorResponse with the provided code, message, and details.
//
// Parameters:
//   - code: A unique error code (e.g., errors.ErrInvalidInput).
//   - message: A human-readable error message.
//   - details: Additional context or details about the error (optional).
//
// Returns:
//   - *ErrorResponse: A pointer to the newly created ErrorResponse.
func NewErrorResponse(code, message, details string) *ErrorResponse {
	errResp := &ErrorResponse{}
	errResp.Error.Code = code
	errResp.Error.Message = message
	errResp.Error.Details = details
	return errResp
}

// WriteErrorResponse writes an error response to the HTTP response writer.
// It sets the appropriate HTTP status code and encodes the error response as JSON.
//
// Parameters:
//   - w: The HTTP response writer.
//   - statusCode: The HTTP status code to return (e.g., http.StatusBadRequest).
//   - code: A unique error code (e.g., errors.ErrInvalidInput).
//   - message: A human-readable error message.
//   - details: Additional context or details about the error (optional).
func WriteErrorResponse(w http.ResponseWriter, statusCode int, code, message, details string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errResp := NewErrorResponse(code, message, details)
	if err := json.NewEncoder(w).Encode(errResp); err != nil {
		http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
	}
}

// NotFoundError represents a "not found" error.
// It is used when a requested resource cannot be found in the system.
type NotFoundError struct {
	Resource string // Name of the resource that was not found.
	ID       string // ID of the resource that was not found.
	Code     string // Error code for categorization.
}

// Error returns a formatted error message for the NotFoundError.
// If an ID is provided, it includes the ID in the message.
func (e *NotFoundError) Error() string {
	if e.ID != "" {
		return fmt.Sprintf("%s with ID '%s' not found", e.Resource, e.ID)
	}
	return fmt.Sprintf("%s not found", e.Resource)
}

// NewNotFoundError creates a new NotFoundError instance.
// It initializes the NotFoundError with the provided resource, ID, and error code.
//
// Parameters:
//   - resource: The name of the resource that was not found.
//   - id: The ID of the resource that was not found.
//   - code: A unique error code for categorization (e.g., errors.ErrResourceNotFound).
//
// Returns:
//   - *NotFoundError: A pointer to the newly created NotFoundError.
func NewNotFoundError(resource string, id string, code string) *NotFoundError {
	return &NotFoundError{
		Resource: resource,
		ID:       id,
		Code:     code,
	}
}

// ValidationError represents a validation error.
// It is used when input data fails validation checks.
type ValidationError struct {
	Errors map[string]string // Map of field names to error messages.
	Code   string            // Error code for categorization.
}

// Error returns a formatted error message for the ValidationError.
// It lists all validation errors by field name and their corresponding messages.
func (e *ValidationError) Error() string {
	message := "Validation errors:"
	for field, msg := range e.Errors {
		message += fmt.Sprintf("\n- Field '%s': %s", field, msg)
	}
	return message
}

// NewValidationError creates a new ValidationError for a single field.
// It initializes the ValidationError with the provided field, message, and error code.
//
// Parameters:
//   - field: The name of the field that failed validation.
//   - message: The validation error message for the field.
//   - code: A unique error code for categorization (e.g., errors.ErrValidationFailed).
//
// Returns:
//   - *ValidationError: A pointer to the newly created ValidationError.
func NewValidationError(field string, message string, code string) *ValidationError {
	return &ValidationError{
		Errors: map[string]string{field: message},
		Code:   code,
	}
}

// NewValidationErrors creates a new ValidationError for multiple fields.
// It initializes the ValidationError with the provided map of field errors and error code.
//
// Parameters:
//   - errors: A map of field names to their corresponding error messages.
//   - code: A unique error code for categorization (e.g., errors.ErrValidationFailed).
//
// Returns:
//   - *ValidationError: A pointer to the newly created ValidationError.
func NewValidationErrors(errors map[string]string, code string) *ValidationError {
	return &ValidationError{
		Errors: errors,
		Code:   code,
	}
}

// Constants for error codes.
const (
	ErrInvalidInput       = "ERR-001" // Invalid input provided by the client.
	ErrResourceNotFound   = "ERR-002" // Requested resource not found.
	ErrValidationFailed   = "ERR-003" // Input validation failed.
	ErrInternalServer     = "ERR-004" // Internal server error.
	ErrUnauthorized       = "ERR-005" // Unauthorized access.
	ErrRateLimitExceeded  = "ERR-006" // Rate limit exceeded.
	ErrDatabaseConnection = "ERR-007" // Database connection error.
	ErrConflict           = "ERR-008" // Conflict (e.g., duplicate resource).
	ErrForbidden          = "ERR-009" // Forbidden access (e.g., insufficient permissions).
	ErrServiceUnavailable = "ERR-010" // Service unavailable (e.g., third-party service down).
)
