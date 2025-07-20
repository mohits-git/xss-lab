package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mohits-git/xss-lab/internal/database"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func main() {
	fmt.Println("Starting the web server")

	apiCfg, err := initializeAPIConfig()
	if err != nil {
		log.Fatalf("Error initializing API config: %v\n", err)
	}
	defer func() {
		if err := database.CloseDB(); err != nil {
			fmt.Printf("Error closing database: %v\n", err)
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	// HTML endpoints
	mux.HandleFunc("/", apiCfg.landingPageHandler)
	mux.HandleFunc("GET /login", apiCfg.loginPageHandler)
	mux.HandleFunc("GET /register", apiCfg.registerPageHandler)
	mux.HandleFunc("GET /blogs", apiCfg.blogsPageHandler)
	mux.HandleFunc("GET /blogs/{id}", apiCfg.blogPageHandler)
	// API endpoints
	mux.HandleFunc("POST /api/login", apiCfg.loginHandler)
	mux.HandleFunc("POST /api/register", apiCfg.registerHandler)
	mux.HandleFunc("GET /api/blogs", apiCfg.getBlogsHandler)
	mux.HandleFunc("GET /api/blogs/count", apiCfg.getBlogsCountHandler)
	mux.HandleFunc("GET /api/users/{id}/blogs", apiCfg.getUserBlogsHandler)
	mux.Handle("POST /api/blogs", apiCfg.loggedInMiddleware(apiCfg.createBlogHandler))
	mux.Handle("PUT /api/blogs/{id}", apiCfg.loggedInMiddleware(apiCfg.updateBlogHandler))
	mux.Handle("DELETE /api/blogs/{id}", apiCfg.loggedInMiddleware(apiCfg.deleteBlogHandler))
	mux.Handle("POST /api/comments/{blog_id}", apiCfg.loggedInMiddleware(apiCfg.createCommentHandler));

	fmt.Println("Web server listening on :" + apiCfg.port)
	log.Fatal(http.ListenAndServe(":"+apiCfg.port, mux))
}
