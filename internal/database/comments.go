package database

import (
	"database/sql"
)

type Comment struct {
	ID        int64        `json:"id"`
	BlogID    int64        `json:"blog_id"`
	UserID    int64        `json:"user_id"`
	Author    string	   `json:"author,omitempty"`
	Content   string       `json:"content"`
	CreatedAt sql.NullTime `json:"created_at"`
}

func (q *Queries) CreateComment(blogID, userID int64, content string) (*Comment, error) {
	query := "INSERT INTO comments (blog_id, user_id, content) VALUES (?, ?, ?)"
	result, err := q.db.Exec(query, blogID, userID, content)
	if err != nil {
		return nil, err
	}
	commentID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &Comment{
		ID:      commentID,
		BlogID:  blogID,
		UserID:  userID,
		Content: content,
	}, nil
}

func (q *Queries) GetCommentByID(commentID int) (*Comment, error) {
	query := "SELECT comments.id, blog_id, user_id, content, comments.created_at, users.name FROM comments JOIN users ON comments.user_id = users.id WHERE comments.id = ?"
	row := q.db.QueryRow(query, commentID)

	var comment Comment
	err := row.Scan(&comment.ID, &comment.BlogID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.Author)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No comment found
		}
		return nil, err // Other error
	}
	return &comment, nil
}

func (q *Queries) GetCommentsByBlogID(blogID int64) ([]*Comment, error) {
	query := "SELECT comments.id, blog_id, user_id, content, comments.created_at, users.name FROM comments JOIN users ON comments.user_id = users.id WHERE blog_id = ?"
	rows, err := q.db.Query(query, blogID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.BlogID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.Author)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (q *Queries) DeleteComment(commentID int) error {
	query := "DELETE FROM comments WHERE id = ?"
	_, err := q.db.Exec(query, commentID)
	if err != nil {
		return err
	}
	return nil
}

func (q *Queries) UpdateComment(commentID int, content string) (*Comment, error) {
	query := "UPDATE comments SET content = ? WHERE id = ?"
	_, err := q.db.Exec(query, content, commentID)
	if err != nil {
		return nil, err
	}
	return q.GetCommentByID(commentID)
}

func (q *Queries) GetCommentsByUserID(userID int) ([]*Comment, error) {
	query := "SELECT id, blog_id, user_id, content, created_at FROM comments WHERE user_id = ?"
	rows, err := q.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.BlogID, &comment.UserID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}
