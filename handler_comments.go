package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) createCommentHandler(w http.ResponseWriter, r *http.Request, userID int64) {
	blogIDStr := r.PathValue("blog_id")
	blogID, err := strconv.ParseInt(blogIDStr, 10, 64)
	if err != nil || blogID <= 0 {
		fmt.Println("Invalid blog ID:", blogIDStr, err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")

	comment, err := cfg.db.CreateComment(blogID, userID, content)
	if err != nil {
		fmt.Println("Error creating comment:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(comment); err != nil {
		fmt.Println("Error encoding comment to JSON:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
