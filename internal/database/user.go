package database

import (
	"database/sql"
)

type User struct {
	ID    int64    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (q *Queries) GetUserByID(userID int) (*User, error) {
	query := "SELECT id, name, email FROM users WHERE id = ?"
	row := q.db.QueryRow(query, userID)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found
		}
		return nil, err // Other error
	}
	return &user, nil
}

func (q *Queries) CreateUser(name, email, password_hash string) (*User, error) {
	query := "INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)"
	result, err := q.db.Exec(query, name, email, password_hash)
	if err != nil {
		return nil, err
	}
	userID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &User{
		ID:    userID,
		Name:  name,
		Email: email,
	}, nil
}

func (q *Queries) UpdateUser(userID int, name, email string) (*User, error) {
	query := "UPDATE users SET name = ?, email = ? WHERE id = ?"
	_, err := q.db.Exec(query, name, email, userID)
	if err != nil {
		return nil, err
	}
	return q.GetUserByID(userID)
}
