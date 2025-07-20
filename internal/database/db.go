package database

import (
	"database/sql"
)

var db *sql.DB

type Queries struct {
	db *sql.DB
}

func InitializeDB(dataSourceName string) error {
	var err error
	db, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return err
	}
	// Check if the database is reachable
	if err = db.Ping(); err != nil {
		return err
	}
	// Enable foreign keys
	if _, err = db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return err
	}
	return nil
}

func GetQueries() *Queries {
	if db == nil {
		panic("Database not initialized. Call InitializeDB first.")
	}
	return &Queries{db: db}
}

func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	db = nil
	return nil
}
