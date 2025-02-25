// Package middleware provides HTTP middleware components for the application.
//
// The validation middleware package implements request validation using go-playground/validator.
// It provides a custom validator that can be used with Echo framework to validate
// incoming HTTP requests against struct tags.
//
// Key features:
// - Custom validator implementation for Echo framework
// - Support for struct field validation using tags
// - Extensible validation rules
// - Integration with go-playground/validator
//
// Usage:
//  e := echo.New()
//  e.Validator = middleware.NewCustomValidator()
//
// Example struct with validation tags:
//  type CreateUserRequest struct {
//      Name  string `json:"name" validate:"required,min=2"`
//      Email string `json:"email" validate:"required,email"`
//      Age   int    `json:"age" validate:"required,gte=18"`
//  }
package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// GetValidator returns the validator instance
func GetValidator() *validator.Validate {
	return validate
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", v.Field, v.Message)
}

// CustomValidator is a custom validator for Echo
type CustomValidator struct {
	Validator *validator.Validate
}

// Validate implements echo.Validator interface
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errorMessages []string
			for _, e := range validationErrors {
				message := fmt.Sprintf("%s: failed validation for '%s'", e.Field(), e.Tag())
				errorMessages = append(errorMessages, message)
			}
			return &ValidationError{
				Field:   "multiple",
				Message: strings.Join(errorMessages, "; "),
			}
		}
		return err
	}
	return nil
}

// ValidationMiddleware creates a validator middleware for Echo
func ValidationMiddleware(e *echo.Echo) {
	e.Validator = &CustomValidator{Validator: validate}

	// Add request validator middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method == http.MethodPost || c.Request().Method == http.MethodPut {
				if c.Request().Header.Get(echo.HeaderContentType) != echo.MIMEApplicationJSON {
					return echo.NewHTTPError(http.StatusBadRequest, "Invalid Content-Type. Expected application/json")
				}
			}
			return next(c)
		}
	})
}

