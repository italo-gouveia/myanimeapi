package httphandler

import (
	"encoding/json"
	"net/http"

	"myanimeapi/internal/errors"
	"myanimeapi/internal/logger"
)

const (
	contentTypeJSON = "application/json"
	contentTypeKey  = "Content-Type"

	msgInvalidPayload      = "Invalid payload"
	msgPayloadNotRetrieved = "The request payload could not be retrieved."
	msgNotAuthenticated    = "User not authenticated"
	msgAuthRequired        = "Authentication required to access this resource"
	msgNotAuthReview       = "You must be logged in to create a review."
	msgNotAuthUpdateReview = "You must be logged in to update a review."
	msgNotAuthDeleteReview = "You must be logged in to delete a review."
	msgUnauthorizedUpdate  = "You are not authorized to update this review."
	msgUnauthorizedDelete  = "You are not authorized to delete this review."
	msgUserNotAuth         = "The user is not authenticated."

	msgFailedEncodeDetails = "An internal server error occurred while encoding the response."
)

// writeJSON sets Content-Type, writes status, and encodes v as JSON. On encode
// failure it logs (when log is non-nil), writes an error response, and returns
// false so the caller can return immediately.
func writeJSON(w http.ResponseWriter, status int, v interface{}, log *logger.Logger) bool {
	w.Header().Set(contentTypeKey, contentTypeJSON)
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		if log != nil {
			log.WithField("error", err.Error()).Error("Failed to encode response")
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", msgFailedEncodeDetails, nil)
		return false
	}
	return true
}

// respondInvalidPayload writes the standard "invalid payload" error and logs a warning.
func respondInvalidPayload(w http.ResponseWriter, log *logger.Logger) {
	if log != nil {
		log.Warning(msgInvalidPayload)
	}
	errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, msgInvalidPayload, msgPayloadNotRetrieved, nil)
}

// writeAppError checks whether err is an *errors.AppError and writes the
// appropriate HTTP response. Returns true if it was an AppError (caller should
// return), false otherwise.
func writeAppError(w http.ResponseWriter, err error) bool {
	if appErr, ok := err.(*errors.AppError); ok {
		errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
		return true
	}
	return false
}
