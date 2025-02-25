package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorHandler(t *testing.T) {
	tests := []struct {
		name           string
		error          error
		expectedStatus int
		expectedBody   ErrorResponse
	}{
		{
			name:           "Record Not Found Error",
			error:          echo.NewHTTPError(http.StatusNotFound, "record not found"),
			expectedStatus: http.StatusNotFound,
			expectedBody: ErrorResponse{
				Error: AppError{
					Code:    http.StatusNotFound,
					Message: "Not Found",
					Details: "record not found",
				},
			},
		},
		{
			name: "Validation Error",
			error: &ValidationError{
				Field:   "test",
				Message: "validation failed",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: ErrorResponse{
				Error: AppError{
					Code:    http.StatusBadRequest,
					Message: "Validation failed",
					Details: "test: validation failed",
				},
			},
		},
		{
			name:           "Unknown Error",
			error:          errors.New("unknown error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: ErrorResponse{
				Error: AppError{
					Code:    http.StatusInternalServerError,
					Message: "Internal server error",
					Details: "An unexpected error occurred",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(ErrorHandler())

			router.GET("/test", func(c *gin.Context) {
				c.Error(tt.error)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)
		})
	}
}

type TestStruct struct {
	Name  string `json:"name" validate:"required,min=3"`
	Email string `json:"email" validate:"required,email"`
}

func TestValidationMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		contentType    string
		expectedStatus int
	}{
		{
			name:           "Valid Content Type",
			contentType:    "application/json",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Content Type",
			contentType:    "text/plain",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(ValidationMiddleware())

			router.POST("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/test", bytes.NewBuffer([]byte("{}")))
			req.Header.Set("Content-Type", tt.contentType)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestRequestValidator(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		payload        interface{}
		expectedStatus int
	}{
		{
			name: "Valid Request",
			payload: TestStruct{
				Name:  "John Doe",
				Email: "john@example.com",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid Request - Missing Email",
			payload: TestStruct{
				Name: "John Doe",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid Request - Short Name",
			payload: TestStruct{
				Name:  "Jo",
				Email: "john@example.com",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(ErrorHandler())

			router.POST("/test", RequestValidator(&TestStruct{}), func(c *gin.Context) {
				validated := c.MustGet("validated").(*TestStruct)
				assert.NotNil(t, validated)
				c.Status(http.StatusOK)
			})

			body, _ := json.Marshal(tt.payload)
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
