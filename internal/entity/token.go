package entity

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/satori/go.uuid"
)

type JWTToken string

// Issue new token using HS256 signing method.
func NewJWTToken(user User, secret string, duration time.Duration) (JWTToken, error) {
	now := time.Now()

	claims := jwt.MapClaims{}

	claims["iss"] = "Vetka"
	claims["jti"] = uuid.NewV4()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()
	claims["exp"] = now.Add(duration).Unix()

	// User info
	claims["sub"] = user.ID
	claims["login"] = user.Login

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := rawToken.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("jwt - Token.New - rawToken.SignedString: %w", err)
	}

	return JWTToken(signedToken), nil
}
