package utils

import (
	"errors"
	"html"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// Claims holds jwt claims.
type Claims struct {
	ID uuid.UUID `json:"id"`
	jwt.StandardClaims
}

// GenerateJWTToken generates new JWT Token with claims, which includes the id and expiry time.
func GenerateJWTToken(id uuid.UUID, expires time.Duration, JwtSecretKey string) (string, error) {
	// Register the JWT claims, which includes the id and expiry time
	claims := &Claims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expires).Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Register the JWT string
	tokenString, err := token.SignedString([]byte(JwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ExtractJWTFromRequest extracts JWT claims form request.
func ExtractJWTFromRequest(r *http.Request, jwtSecretKey string) (*Claims, error) {
	// Get the JWT string
	tokenString := ExtractBearerToken(r)

	// Parse the JWT string
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, errors.New("invalid token signature")
		}
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ExtractBearerToken extracts bearer token from request Authorization header.
func ExtractBearerToken(r *http.Request) string {
	headerAuthorization := r.Header.Get("Authorization")
	bearerToken := strings.Split(headerAuthorization, " ")
	return html.EscapeString(bearerToken[1])
}
