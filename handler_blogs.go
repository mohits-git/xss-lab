package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/mohits-git/xss-lab/internal/database"
)

func (cfg *apiConfig) blogsPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("templates/blogs.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (cfg *apiConfig) blogPageHandler(w http.ResponseWriter, r *http.Request) {
	blogID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || blogID <= 0 {
		fmt.Println("Error parsing blog ID:", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	blog, err := cfg.db.GetBlogByID(blogID)
	if err != nil {
		fmt.Println("Error fetching blog:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if blog == nil {
		http.NotFound(w, r)
		return
	}

	comments, err := cfg.db.GetCommentsByBlogID(blogID)
	if err != nil {
		fmt.Println("Error fetching comments:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Blog     *database.Blog
		Comments []*database.Comment
	}{
		Blog:    blog,
		Comments: comments,
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl, err := template.ParseFiles("templates/blog.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (cfg *apiConfig) getBlogsHandler(w http.ResponseWriter, r *http.Request) {
	// pagination
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1" // default to page 1 if not specified
	}
	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	limit := 10
	offset := (pageNum - 1) * limit
	// get blogs from the database
	blogs, err := cfg.db.GetBlogs(limit, offset)
	if err != nil {
		fmt.Println("Error fetching blogs:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(blogs); err != nil {
		fmt.Println("Error encoding blogs to JSON:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (cfg *apiConfig) createBlogHandler(w http.ResponseWriter, r *http.Request, userID int64) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	newBlog, err := cfg.db.CreateBlog(title, content, userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newBlog); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (cfg *apiConfig) getUserBlogsHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("id")
	if userIDStr == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil || userID <= 0 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	blogs, err := cfg.db.GetBlogsByUserID(userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(blogs); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (cfg *apiConfig) updateBlogHandler(w http.ResponseWriter, r *http.Request, userID int64) {
	blogIDStr := r.PathValue("id")
	if blogIDStr == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	blogID, err := strconv.ParseInt(blogIDStr, 10, 64)
	if err != nil || blogID <= 0 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	existingBlog, err := cfg.db.GetBlogByID(blogID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if existingBlog == nil {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}
	// authorization check
	if existingBlog.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	updatedBlog, err := cfg.db.UpdateBlog(blogID, title, content)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedBlog); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (cfg *apiConfig) deleteBlogHandler(w http.ResponseWriter, r *http.Request, userID int64) {
	blogIDStr := r.PathValue("id")
	if blogIDStr == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	blogID, err := strconv.ParseInt(blogIDStr, 10, 64)
	if err != nil || blogID <= 0 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	existingBlog, err := cfg.db.GetBlogByID(blogID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if existingBlog == nil {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}
	// authorization check
	if existingBlog.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := cfg.db.DeleteBlog(blogID); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
