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
