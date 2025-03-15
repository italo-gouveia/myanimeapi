// internal/errors/errors.go
// Package errors provides custom error types used in the application.
// It defines two custom error types: NotFoundError and ValidationError.
// These errors are used to handle specific error scenarios, such as resource not found and validation failures.
package errors

import "fmt"

// NotFoundError represents a "not found" error.
// It is used when a requested resource cannot be found in the system.
type NotFoundError struct {
	Resource string // Name of the resource that was not found
	ID       string // ID of the resource that was not found
	Code     string // Error code for categorization
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
	Errors map[string]string // Map of field names to error messages
	Code   string            // Error code for categorization
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
func NewValidationError(field string, message string, code string) *ValidationError {
	return &ValidationError{
		Errors: map[string]string{field: message},
		Code:   code,
	}
}

// NewValidationErrors creates a new ValidationError for multiple fields.
// It initializes the ValidationError with the provided map of field errors and error code.
func NewValidationErrors(errors map[string]string, code string) *ValidationError {
	return &ValidationError{
		Errors: errors,
		Code:   code,
	}
}
