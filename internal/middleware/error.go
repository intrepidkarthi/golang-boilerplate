package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// AppError represents a structured error response
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// ErrorResponse wraps the AppError for JSON response
type ErrorResponse struct {
	Error AppError `json:"error"`
}

// ErrorHandler middleware handles application errors and converts them to structured responses
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Only handle errors if there are any
		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		var response ErrorResponse

		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			response = ErrorResponse{
				Error: AppError{
					Code:    http.StatusNotFound,
					Message: "Resource not found",
					Details: err.Error(),
				},
			}
		case errors.Is(err, &ValidationError{}):
			validationErr := err.(*ValidationError)
			response = ErrorResponse{
				Error: AppError{
					Code:    http.StatusBadRequest,
					Message: "Validation failed",
					Details: validationErr.Error(),
				},
			}
		default:
			// Log unexpected errors here
			response = ErrorResponse{
				Error: AppError{
					Code:    http.StatusInternalServerError,
					Message: "Internal server error",
					Details: "An unexpected error occurred",
				},
			}
		}

		c.JSON(response.Error.Code, response)
		c.Abort()
	}
}
