package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/mohits-git/xss-lab/internal/auth"
)

type protectedHandlerFunc func(w http.ResponseWriter, r *http.Request, userID int64)

func (cfg *apiConfig) loggedInMiddleware(next protectedHandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check auth header
		token, err := auth.GetAuthHeader(r.Header)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unauthorized: %v", err), http.StatusUnauthorized)
			return
		}

		// verify token
		uID, err := auth.ValidateJWT(token, cfg.jwtSecret)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unauthorized: %v", err), http.StatusUnauthorized)
			return
		}

		userID, err := strconv.ParseInt(uID, 10, 64)
		if err != nil {
			fmt.Println("Error parsing user ID:", err)
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		next(w, r, userID)
	})
}
