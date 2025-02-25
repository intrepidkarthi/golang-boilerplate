// Package middleware provides HTTP middleware components for the application.
//
// The RBAC (Role-Based Access Control) middleware implements authorization
// for the HTTP server. It manages user roles and permissions, controlling
// access to various endpoints based on the user's assigned roles.
//
// Key features:
// - Role-based authorization
// - Permission management
// - Resource-level access control
// - Role hierarchy support
// - Dynamic permission checking
//
// Supported Roles:
// - ADMIN: Full system access
// - MANAGER: Department-level access
// - USER: Basic access rights
//
// Example Permission Structure:
//  {
//      "users": ["read", "write"],
//      "reports": ["read"],
//      "settings": ["read", "write", "delete"]
//  }
//
// Usage:
//  e := echo.New()
//  e.Use(middleware.RBACMiddleware())
package middleware

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go-boilerplate/internal/auth"
	"net/http"
)

func RBAC(requiredService, requiredPermission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, ok := c.Get("jwtClaims").(auth.Claims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}

			hasAccess := false
			for _, perm := range claims.Permissions {
				if perm == fmt.Sprintf("%s:%s", requiredService, requiredPermission) {
					hasAccess = true
					break
				}
			}

			if !hasAccess {
				return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions")
			}

			return next(c)
		}
	}
}
