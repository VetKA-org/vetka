package entity

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/satori/go.uuid"
)

type JWTToken string

// Issue new token using HS256 signing method.
func NewJWTToken(user User, roles []Role, secret Secret, duration time.Duration) (JWTToken, error) {
	now := time.Now()

	assignedRoles := make([]string, len(roles))
	for i, role := range roles {
		assignedRoles[i] = role.Name
	}

	claims := jwt.MapClaims{}

	claims["iss"] = "Vetka"
	claims["jti"] = uuid.NewV4()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()
	claims["exp"] = now.Add(duration).Unix()

	// User info
	claims["sub"] = user.ID
	claims["login"] = user.Login
	claims["roles"] = assignedRoles

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := rawToken.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("jwt - Token.New - rawToken.SignedString: %w", err)
	}

	return JWTToken(signedToken), nil
}

type DecodedToken struct {
	jwt.RegisteredClaims
	Roles []string `json:"roles"`
}

// Decode token, verify it's signature and return claims if the token is valid.
func DecodeToken(rawToken string, secret Secret) (*DecodedToken, error) {
	claims := new(DecodedToken)
	if _, err := jwt.ParseWithClaims(rawToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}); err != nil {
		return nil, err
	}

	return claims, nil
}

// Decode token without verification.
// (!) Warning: shouldn't be used unless we 100% sure that token is valid and signed by our service.
func DecodeTokenUnverified(rawToken string) (*DecodedToken, error) {
	claims := new(DecodedToken)
	if _, _, err := new(jwt.Parser).ParseUnverified(rawToken, claims); err != nil {
		return nil, err
	}

	return claims, nil
}
