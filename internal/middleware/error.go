package middleware

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

// AppError represents a structured error response
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// NotFoundError represents a resource not found error
type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %s not found", e.Resource, e.ID)
}

// ErrorResponse wraps the AppError for JSON response
type ErrorResponse struct {
	Error AppError `json:"error"`
}

// ErrorHandler middleware handles application errors and converts them to structured responses
func ErrorHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err == nil {
			return nil
		}

		var response ErrorResponse

		switch {
		case errors.Is(err, &NotFoundError{}):
			response = ErrorResponse{
				Error: AppError{
					Code:    http.StatusNotFound,
					Message: "Not Found",
					Details: err.Error(),
				},
			}
			return c.JSON(http.StatusNotFound, response)

		case errors.As(err, &ValidationError{}):
			response = ErrorResponse{
				Error: AppError{
					Code:    http.StatusBadRequest,
					Message: "Validation Error",
					Details: err.Error(),
				},
			}
			return c.JSON(http.StatusBadRequest, response)

		default:
			response = ErrorResponse{
				Error: AppError{
					Code:    http.StatusInternalServerError,
					Message: "Internal Server Error",
					Details: err.Error(),
				},
			}
			return c.JSON(http.StatusInternalServerError, response)
		}
	}
}
