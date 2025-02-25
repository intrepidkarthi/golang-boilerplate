package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
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
			e := echo.New()
			e.Use(ErrorHandler())

			e.GET("/test", func(c echo.Context) error {
				return tt.error
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			var response ErrorResponse
			err := json.Unmarshal(rec.Body.Bytes(), &response)
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
			e := echo.New()
			e.Use(ValidationMiddleware())

			e.POST("/test", func(c echo.Context) error {
				return c.NoContent(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer([]byte("{}")))
			req.Header.Set("Content-Type", tt.contentType)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestRequestValidator(t *testing.T) {
	tests := []struct {
		name           string
		payload        interface{}
		expectedStatus int
	}{
		{
			name: "Valid Request",
			payload: TestStruct{
				Name:  "Super star",
				Email: "superstar@example.com",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid Request - Missing Email",
			payload: TestStruct{
				Name: "Super star",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid Request - Short Name",
			payload: TestStruct{
				Name:  "Star",
				Email: "superstar@example.com",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.Use(ErrorHandler())

			e.POST("/test", func(c echo.Context) error {
				var validated TestStruct
				if err := c.Bind(&validated); err != nil {
					return err
				}
				if err := c.Validate(&validated); err != nil {
					return err
				}
				return c.NoContent(http.StatusOK)
			})

			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
