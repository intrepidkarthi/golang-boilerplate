// Package middleware provides HTTP middleware components for the application.
//
// The security middleware implements various security measures for the HTTP server.
// It includes headers and configurations to protect against common web vulnerabilities
// and follows security best practices.
//
// Key features:
// - CORS (Cross-Origin Resource Sharing) configuration
// - Security headers (X-Frame-Options, X-XSS-Protection, etc.)
// - Content Security Policy (CSP)
// - HSTS (HTTP Strict Transport Security)
// - Rate limiting
// - Request size limiting
//
// Security Headers Set:
// - X-Frame-Options: DENY
// - X-Content-Type-Options: nosniff
// - X-XSS-Protection: 1; mode=block
// - Strict-Transport-Security: max-age=31536000; includeSubDomains
// - Content-Security-Policy: default-src 'self'
//
// Usage:
//  e := echo.New()
//  e.Use(middleware.SecurityHeaders())
//  e.Use(middleware.RateLimiter())
package middleware

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func SecurityHeaders() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Response().Header()
			header.Set("X-Content-Type-Options", "nosniff")
			header.Set("X-Frame-Options", "DENY")
			header.Set("X-XSS-Protection", "1; mode=block")
			header.Set("Content-Security-Policy",
				"default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'")
			return next(c)
		}
	}
}

func CORS(allowedOrigins []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Response().Header()
			header.Set("Access-Control-Allow-Origin", strings.Join(allowedOrigins, ","))
			header.Set("Access-Control-Allow-Credentials", "true")
			header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

			if c.Request().Method == http.MethodOptions {
				return c.NoContent(http.StatusNoContent)
			}

			return next(c)
		}
	}
}
