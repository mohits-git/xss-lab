package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func MakeJwt(userId string, tokenSecret string, expiresIn time.Duration) (string, error) {
	if userId == "" {
		return "", errors.New("userId cannot be empty")
	}

	if tokenSecret == "" {
		return "", errors.New("tokenSecret cannot be empty")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "xss-lab",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userId,
	})

	jwtStr, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return jwtStr, nil
}

func ValidateJWT(tokenString, tokenSecret string) (string, error) {
	if tokenSecret == "" {
		return "", errors.New("tokenSecret cannot be empty")
	}

	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	}, jwt.WithIssuer("xss-lab"))
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", jwt.ErrSignatureInvalid
	}

	userId := claims.Subject
	return userId, nil
}
