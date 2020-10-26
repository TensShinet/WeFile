package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetBadRequestResponse(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, BadRequestResponse{Message: message})
}

func SetSimpleResponse(c *gin.Context, code int, message string) {
	switch code {
	case http.StatusNotFound:
		c.JSON(code, Response404{Message: message})
	case http.StatusBadRequest:
		c.JSON(code, BadRequestResponse{Message: message})
	case http.StatusAccepted:
		c.JSON(code, AcceptedResponse{Message: message})
	case http.StatusConflict:
		c.JSON(code, ConflictError{Message: message})
	default:
		c.JSON(code, ErrorResponse{Message: message})
	}
}
