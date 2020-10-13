package common

const (
	DBNotFoundCode   = 404
	UnauthorizedCode = 401
)

// Sever Error
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

// Unauthorized Error
//
// swagger:response UnauthorizedResponse
type UnauthorizedResponse struct {
	Message string `json:"message"`
}
