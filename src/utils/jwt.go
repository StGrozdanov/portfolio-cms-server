package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type TokenClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

var tokenKey string

// GenerateJWT generates a new JWT access token
func GenerateJWT() (string, error) {
	claims := TokenClaims{
		"administrator",
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return accessToken.SignedString([]byte(tokenKey))
}

// ParseJWT parses the JWT token, checking for correct signing method. Use token.Valid method after for
// final token validation
func ParseJWT(accessToken string) (token *jwt.Token, err error) {
	token, err = jwt.ParseWithClaims(accessToken, TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenKey), nil
	})
	return
}

// ExtractJWTClaims extracts the data from the token
func ExtractJWTClaims(token *jwt.Token) (*TokenClaims, error) {
	if claims, ok := token.Claims.(*TokenClaims); ok {
		return claims, nil
	}
	return nil, errors.New("could not extract token claims")
}

// GetJWTKey gets the JWT secret and uses it for new token signing and old tokens verification
func GetJWTKey(key string) {
	tokenKey = key
}
