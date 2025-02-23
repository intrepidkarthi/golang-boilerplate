package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-boilerplate/internal/auth"
)

func RBAC(requiredService, requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("jwtClaims")
		if !exists {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		hasAccess := false
		for _, perm := range claims.(auth.Claims).Permissions {
			if perm == fmt.Sprintf("%s:%s", requiredService, requiredPermission) {
				hasAccess = true
				break
			}
		}

		if !hasAccess {
			c.AbortWithStatusJSON(403, gin.H{"error": "Insufficient permissions"})
		}
		c.Next()
	}
}
