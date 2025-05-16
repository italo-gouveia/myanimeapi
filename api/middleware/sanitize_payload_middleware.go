// api/middleware/sanitize_payload_middleware.go
// Package middleware provides HTTP middleware utilities for handling requests.
// This file defines a middleware for validating and sanitizing request payloads.
// It uses the `validator` package for validation and `bluemonday` for sanitization.
package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"

	"myanimeapi/internal/errors"
	"myanimeapi/internal/logger"

	"github.com/go-playground/validator/v10"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sirupsen/logrus"
)

var validate = validator.New() // Global validator instance
var p = bluemonday.UGCPolicy() // Sanitization policy for user-generated content

// contextKeyValidation is a custom type for context keys to avoid key collisions.
type contextKeyValidation string

// ValidatedPayloadKey is the context key for storing the validated and sanitized payload.
const ValidatedPayloadKey contextKeyValidation = "validatedPayload"

// ValidateAndSanitizePayload is a middleware that validates and sanitizes the request payload.
// It decodes the request body into the provided payload type, validates it using the `validator` package,
// sanitizes string fields using `bluemonday`, and stores the validated payload in the request context.
//
// Example usage:
//
//	type MyPayload struct {
//	    Name  string `json:"name" validate:"required,min=3,max=50"`
//	    Email string `json:"email" validate:"required,email"`
//	}
//
//	http.Handle("/path", ValidateAndSanitizePayload(myHandler, MyPayload{}))
//
// This will validate and sanitize the request payload for "/path".
func ValidateAndSanitizePayload(next http.Handler, payloadType interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.Get()
		log.WithFields(logrus.Fields{
			"path":       r.URL.Path,
			"method":     r.Method,
			"request_id": r.Context().Value(RequestIDContextKey),
		}).Debug("Starting payload validation and sanitization")

		// Create a new instance of the payload type
		payload := reflect.New(reflect.TypeOf(payloadType)).Interface()

		// Decode the request body into the payload
		if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
			log.WithError(err).WithFields(logrus.Fields{
				"path":       r.URL.Path,
				"method":     r.Method,
				"request_id": r.Context().Value(RequestIDContextKey),
			}).Warn("Failed to decode request body")
			errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid input", "The request body could not be decoded.")
			return
		}

		// Log the decoded payload
		log.WithFields(logrus.Fields{
			"payload":    payload,
			"path":       r.URL.Path,
			"method":     r.Method,
			"request_id": r.Context().Value(RequestIDContextKey),
		}).Debug("Request body decoded successfully")

		// Validate the payload
		if err := validate.Struct(payload); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			errorMessages := make(map[string]string)

			for _, e := range validationErrors {
				// Customize error messages based on the field and tag
				switch e.Tag() {
				case "required":
					errorMessages[e.Field()] = e.Field() + " is a required field."
				case "min":
					errorMessages[e.Field()] = e.Field() + " must be at least " + e.Param() + " characters long."
				case "max":
					errorMessages[e.Field()] = e.Field() + " cannot exceed " + e.Param() + " characters."
				case "gte":
					errorMessages[e.Field()] = e.Field() + " must be greater than or equal to " + e.Param() + "."
				case "lte":
					errorMessages[e.Field()] = e.Field() + " must be less than or equal to " + e.Param() + "."
				default:
					errorMessages[e.Field()] = "Validation failed for " + e.Field() + "."
				}
			}

			// Log the validation errors
			log.WithFields(logrus.Fields{
				"errors":     errorMessages,
				"path":       r.URL.Path,
				"method":     r.Method,
				"request_id": r.Context().Value(RequestIDContextKey),
			}).Warn("Payload validation failed")

			// Return the custom error messages as JSON
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(map[string]interface{}{"errors": errorMessages}); err != nil {
				log.WithError(err).WithFields(logrus.Fields{
					"path":       r.URL.Path,
					"method":     r.Method,
					"request_id": r.Context().Value(RequestIDContextKey),
				}).Error("Failed to encode error response")
				errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode error response", "An internal server error occurred while encoding the response.")
				return
			}
			return
		}

		// Sanitize string fields in the payload
		sanitizePayload(payload)

		// Log the sanitized payload
		log.WithFields(logrus.Fields{
			"payload":    payload,
			"path":       r.URL.Path,
			"method":     r.Method,
			"request_id": r.Context().Value(RequestIDContextKey),
		}).Debug("Payload sanitized successfully")

		// Store the validated and sanitized payload in the context using the custom key
		ctx := context.WithValue(r.Context(), ValidatedPayloadKey, payload)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// sanitizePayload sanitizes all string fields in the payload.
// It uses the `bluemonday` policy to sanitize user-generated content.
func sanitizePayload(payload interface{}) {
	v := reflect.ValueOf(payload).Elem() // Get the underlying value of the pointer
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.String {
			sanitizedValue := p.Sanitize(field.String()) // Sanitize the string
			field.SetString(sanitizedValue)              // Update the field with the sanitized value
		}
	}
}

// sanitizePayload (commented out) is an alternative implementation that supports nested structs and pointers.
// It recursively sanitizes all string fields in the payload, including those in nested structs and pointers.
/*
func sanitizePayload(payload interface{}) {
	v := reflect.ValueOf(payload).Elem() // Get the underlying value of the pointer
	sanitizeValue(v)
}

// sanitizeValue recursively sanitizes all string fields in the value.
func sanitizeValue(v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		switch field.Kind() {
		case reflect.String:
			sanitizedValue := p.Sanitize(field.String()) // Sanitize the string
			field.SetString(sanitizedValue)              // Update the field with the sanitized value
		case reflect.Struct:
			sanitizeValue(field) // Recursively sanitize nested structs
		case reflect.Ptr:
			if !field.IsNil() {
				sanitizeValue(field.Elem()) // Recursively sanitize nested pointers
			}
		}
	}
}
*/
