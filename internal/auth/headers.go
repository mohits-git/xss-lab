package auth

import (
	"errors"
	"net/http"
)

func GetAuthHeader(h http.Header) (string, error) {
	authHeader := h.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header not found")
	}

	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", errors.New("invalid authorization header format")
	}

	token := authHeader[7:]
	if token == "" {
		return "", errors.New("token is empty")
	}

	return token, nil
}

func SetAuthHeader(h http.Header, token string) {
	if token == "" {
		h.Del("Authorization")
		return
	}
	h.Set("Authorization", "Bearer "+token)
}
