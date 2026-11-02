package httphandler

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
)
