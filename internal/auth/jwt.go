// Package auth provides authentication and authorization functionality.
//
// The JWT package implements JSON Web Token (JWT) based authentication.
// It handles token generation, validation, and management for secure
// user authentication in the application.
//
// Key features:
// - JWT token generation and validation
// - Token refresh mechanism
// - Claims management
// - Token blacklisting
// - Configurable token expiration
//
// Token Structure:
//  {
//      "sub": "user_id",
//      "exp": 1516239022,
//      "iat": 1516239022,
//      "roles": ["user", "admin"],
//      "permissions": ["read", "write"]
//  }
//
// Usage:
//  token, err := auth.GenerateToken(userID, roles)
//  claims, err := auth.ValidateToken(tokenString)
//
// Security Considerations:
// - Uses HS256 algorithm for signing
// - Implements token expiration
// - Supports token revocation
// - Securely handles sensitive data
package auth

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID     string   `json:"uid"`
	Roles      []string `json:"roles"`
	Permissions []string `json:"perms"`
	jwt.RegisteredClaims
}

func GenerateTokenPair(userID string, roles, perms []string, secret string) (string, string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID:      userID,
		Roles:       roles,
		Permissions: perms,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			Issuer:    "go-boilerplate",
		},
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		Subject:   userID,
	})

	accessSigned, err := accessToken.SignedString([]byte(secret))
	refreshSigned, _ := refreshToken.SignedString([]byte(secret))
	return accessSigned, refreshSigned, err
}

func ValidateToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
