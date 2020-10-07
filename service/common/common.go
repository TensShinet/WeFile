package common

// Error Message
//
// swagger:response ErrorResponse
type ErrorResponse struct {
	Message string `json:"message"`
}

// NotFound Error
//
// swagger:response NotFoundError
type Response404 struct {
	Message string `json:"message"`
}

// Conflict Error
//
// swagger:response ConflictError
type ConflictError struct {
	Message string `json:"message"`
}

// Forbidden Error
//
// swagger:response ForbiddenResponse
type ForbiddenResponse struct {
	Message string `json:"message"`
}
