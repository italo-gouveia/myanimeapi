// internal/errors/errors.go
// This package defines custom error types used in the application.
// It is used by the application to define custom error types.
// It defines two custom error types: NotFoundError and ValidationError.
// NotFoundError represents a "not found" error.
// ValidationError represents a validation error.
// It defines an Error method for each custom error type that returns a formatted error message.
// It imports the fmt package to format error messages.
package errors

import "fmt"

// NotFoundError represents a "not found" error.
type NotFoundError struct {
	Resource string
	ID       string // ID of the resource that was not found
	Code     string // Error code
}

func (e *NotFoundError) Error() string {
	if e.ID != "" {
		return fmt.Sprintf("%s with ID '%s' not found", e.Resource, e.ID)
	}
	return fmt.Sprintf("%s not found", e.Resource)
}

// NewNotFoundError creates a new NotFoundError
func NewNotFoundError(resource string, id string, code string) *NotFoundError {
	return &NotFoundError{
		Resource: resource,
		ID:       id,
		Code:     code,
	}
}

// ValidationError represents a validation error.
type ValidationError struct {
	Errors map[string]string // Map of field names to error messages
	Code   string            // Error code
}

func (e *ValidationError) Error() string {
	message := "Validation errors:"
	for field, msg := range e.Errors {
		message += fmt.Sprintf("\n- Field '%s': %s", field, msg)
	}
	return message
}

// NewValidationError creates a new ValidationError for a single field
func NewValidationError(field string, message string, code string) *ValidationError {
	return &ValidationError{
		Errors: map[string]string{field: message},
		Code:   code,
	}
}

// NewValidationErrors creates a new ValidationError for multiple fields
func NewValidationErrors(errors map[string]string, code string) *ValidationError {
	return &ValidationError{
		Errors: errors,
		Code:   code,
	}
}
