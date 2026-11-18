package httphandler

const (
	msgInvalidIDFormat      = "Invalid ID format"
	msgInvalidIDDetail      = "The provided ID is not a valid unsigned integer."
	msgInvalidPagination    = "Invalid pagination parameters"
	msgDatabaseError        = "Database error"
	msgInvalidGenreID       = "Invalid genre ID"
	msgFailedGetUserID      = "Failed to get user ID from context"
	msgFailedDecodeBody     = "Failed to decode request body"
	msgInvalidRequestBody   = "Invalid request body"
	msgInvalidJSONBody      = "The request body is not a valid JSON object."
	msgFailedGetUser        = "Failed to get user"
	msgInternalServerGetUser = "An internal server error occurred while retrieving the user."
	queryUserID             = "user_id = ?"
)
