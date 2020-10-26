package common

import "fmt"

const (
	DBNotFoundCode   = 404
	UnauthorizedCode = 401
	DBConflictCode   = 409
	DBServiceError   = 500

	// role id
	GeneralUserRoleID = 10000
	AdminRoleID       = 1

	GeneralUserRoleName = "普通用户"
	AdminRoleName       = "超级用户"
)

var (
	ErrConflict = fmt.Errorf("conflict error")
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

// BadRequest
//
// swagger:response BadRequestResponse
type BadRequestResponse struct {
	Message string `json:"message"`
}

// Accepted Response
//
// swagger:response AcceptedResponse
type AcceptedResponse struct {
	Message string `json:"message"`
}
