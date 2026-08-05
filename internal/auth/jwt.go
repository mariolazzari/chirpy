package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const tokenIssuer = "chirpy-access"

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    tokenIssuer,
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn).UTC()),
			Subject:   userID.String(),
		},
	)

	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(token *jwt.Token) (any, error) {
			return []byte(tokenSecret), nil
		})

	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, errors.New("Token is not valid")
	}

	issuer, err := claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}

	if issuer != tokenIssuer {
		return uuid.Nil, errors.New("Token issuer mismatch")
	}

	userID, err := claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, err
	}

	return userUUID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authorization := headers.Get("Authorization")

	// Check if authorization is empty
	if authorization == "" {
		return "", errors.New("Missing Authorization header")
	}

	// Retrieve and clean Authorization token
	return strings.TrimPrefix(authorization, "Bearer "), nil
}
