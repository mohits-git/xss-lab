package main

import (
	"fmt"
	"os"

	"github.com/mohits-git/xss-lab/internal/database"
)

type apiConfig struct {
	db        *database.Queries
	jwtSecret string
	port      string
}

func initializeAPIConfig() (*apiConfig, error) {
	// Load environment variables
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		dbUrl = "file:sqlite3.db"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default_secret_key"
	}

	// initialize sqlite database
	err := database.InitializeDB(dbUrl)
	if err != nil {
		return nil, fmt.Errorf("Error initializing database: %v\n", err)
	}

	return &apiConfig{
		db:        database.GetQueries(),
		jwtSecret: jwtSecret,
		port:      port,
	}, nil
}
