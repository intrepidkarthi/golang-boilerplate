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
