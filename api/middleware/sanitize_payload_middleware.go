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

	"github.com/go-playground/validator/v10"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sirupsen/logrus"

	"myanimeapi/internal/errors"
	"myanimeapi/internal/logger"
)

var validate = validator.New()
var p = bluemonday.UGCPolicy()

// ValidateAndSanitizePayload is a middleware that validates and sanitizes request payloads
func ValidateAndSanitizePayload(payloadType interface{}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.Get()
			logFields := logrus.Fields{
				"path":       r.URL.Path,
				"method":     r.Method,
				"request_id": r.Context().Value(RequestIDContextKey),
			}
			log.WithFields(logFields).Debug("Starting payload validation and sanitization")

			// Create a new instance of the payload type
			payload := reflect.New(reflect.TypeOf(payloadType)).Interface()

			// Decode the request body into the payload
			if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
				log.WithFields(logFields).WithField("error", err).Error("Failed to decode request body")
				errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid input", "The request body could not be decoded.", nil)
				return
			}

			// Log the decoded payload
			log.WithFields(logFields).WithField("payload", payload).Debug("Request body decoded successfully")

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
				log.WithFields(logFields).WithField("errors", errorMessages).Error("Payload validation failed")

				// Return the custom error messages as JSON
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				if err := json.NewEncoder(w).Encode(map[string]interface{}{"errors": errorMessages}); err != nil {
					log.WithFields(logFields).WithField("error", err).Error("Failed to encode error response")
					errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode error response", "An internal server error occurred while encoding the response.", nil)
					return
				}
				return
			}

			// Sanitize string fields in the payload
			if err := sanitizePayload(payload); err != nil {
				log.WithFields(logFields).WithField("error", err).Error("Failed to sanitize payload")
				errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to sanitize payload", "An internal server error occurred while sanitizing the payload.", nil)
				return
			}

			// Store the validated payload in the context
			ctx := r.Context()
			ctx = context.WithValue(ctx, ValidatedPayloadKey, payload)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// sanitizePayload recursively sanitizes string fields in a struct
func sanitizePayload(payload interface{}) error {
	val := reflect.ValueOf(payload)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.String {
			if field.CanSet() {
				sanitized := p.Sanitize(field.String())
				field.SetString(sanitized)
			}
		} else if field.Kind() == reflect.Struct {
			if err := sanitizePayload(field.Addr().Interface()); err != nil {
				return err
			}
		} else if field.Kind() == reflect.Ptr && field.Elem().Kind() == reflect.Struct {
			if err := sanitizePayload(field.Interface()); err != nil {
				return err
			}
		}
	}

	return nil
}
