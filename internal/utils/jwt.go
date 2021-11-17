package utils

import (
	"errors"
	"html"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// GenerateJWTToken generates new JWT Token with claims, which includes the id and expiry time.
func GenerateJWTToken(id uuid.UUID, expires time.Duration, jwtSecretKey string) (string, error) {
	// Register the JWT claims, which includes the id and expiry time
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expires)),
		ID:        id.String(),
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Register the JWT string
	tokenString, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ExtractJWTFromRequest extracts JWT claims form request.
func ExtractJWTFromRequest(r *http.Request, jwtSecretKey string) (*jwt.RegisteredClaims, error) {
	// Get the JWT string
	tokenString, err := ExtractBearerToken(r)
	if err != nil {
		return nil, err
	}

	// Parse the JWT string
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ExtractBearerToken extracts bearer token from request Authorization header.
func ExtractBearerToken(r *http.Request) (string, error) {
	headerAuthorization := r.Header.Get("Authorization")
	bearerToken := strings.Split(headerAuthorization, " ")
	if len(bearerToken) != 2 {
		return "", errors.New("invalid header format")
	}

	return html.EscapeString(bearerToken[1]), nil
}
