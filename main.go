package main

import (
	"fmt"
	"net/http"

	"github.com/mohits-git/xss-lab/internal/database"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func main() {
	fmt.Println("Starting the web server")

	// initialize sqlite database
	err := database.InitializeDB("file:sqlite3.db")
	if err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		return
	}
	defer func() {
		if err := database.CloseDB(); err != nil {
			fmt.Printf("Error closing database: %v\n", err)
		}
	}()

	mux := http.NewServeMux()
	// for css, js, images and other static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", rootHandler)
	http.ListenAndServe(":8080", mux)
}
