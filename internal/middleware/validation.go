package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strings"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", v.Field, v.Message)
}

// ValidateStruct validates a struct using validator tags
func ValidateStruct(obj interface{}) error {
	if err := validate.Struct(obj); err != nil {
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

// ValidationMiddleware validates request bodies against their struct definitions
func ValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			if c.ContentType() != "application/json" {
				c.AbortWithStatusJSON(400, ErrorResponse{
					Error: AppError{
						Code:    400,
						Message: "Invalid Content-Type",
						Details: "Expected application/json",
					},
				})
				return
			}
		}
		c.Next()
	}
}

// RequestValidator middleware factory for validating specific request structs
func RequestValidator(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(obj); err != nil {
			c.Error(&ValidationError{
				Field:   "request",
				Message: err.Error(),
			})
			c.Abort()
			return
		}

		if err := ValidateStruct(obj); err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		c.Set("validated", obj)
		c.Next()
	}
}
