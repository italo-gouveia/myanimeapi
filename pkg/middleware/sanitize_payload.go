// pkg/middleware/sanitize_payload.go
package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/microcosm-cc/bluemonday"
)

var validate = validator.New()
var p = bluemonday.UGCPolicy() // Sanitization policy

// Define a custom type for context keys
type contextKeyValidation string

// Define the key for the validated payload
const ValidatedPayloadKey contextKeyValidation = "validatedPayload"

// Generic middleware for validation and sanitization
func ValidateAndSanitizePayload(next http.Handler, payloadType interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a new instance of the payload type
		payload := reflect.New(reflect.TypeOf(payloadType)).Interface()

		// Decode the request body into the payload
		if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

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

			// Return the custom error messages as JSON
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"errors": errorMessages})
			return
		}

		// Sanitize string fields in the payload
		sanitizePayload(payload)

		// Store the validated and sanitized payload in the context using the custom key
		ctx := context.WithValue(r.Context(), ValidatedPayloadKey, payload)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// sanitizePayload sanitizes all string fields in the payload
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

// sanitizePayload sanitizes all string fields in the payload, including nested structs
/*func sanitizePayload(payload interface{}) {
	v := reflect.ValueOf(payload).Elem() // Get the underlying value of the pointer
	sanitizeValue(v)
}

// sanitizeValue recursively sanitizes all string fields in the value
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
