package database

import (
	"database/sql"
)

type Blog struct {
	ID          int64        `json:"id"`
	Title       string       `json:"title"`
	Content     string       `json:"content"`
	UserID      int          `json:"user_id"`
	PublishedAt sql.NullTime `json:"published_at"`
}

func (q *Queries) CreateBlog(title, content string, userID int) (*Blog, error) {
	query := "INSERT INTO blogs (title, content, user_id) VALUES (?, ?, ?)"
	result, err := q.db.Exec(query, title, content, userID)
	if err != nil {
		return nil, err
	}
	blogID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &Blog{
		ID:     blogID,
		Title:  title,
		Content: content,
		UserID: userID,
	}, nil
}

func (q *Queries) UpdateBlog(blogID int, title, content string) (*Blog, error) {
	query := "UPDATE blogs SET title = ?, content = ? WHERE id = ?"
	_, err := q.db.Exec(query, title, content, blogID)
	if err != nil {
		return nil, err
	}
	return q.GetBlogByID(blogID)
}

func (q *Queries) GetBlogByID(blogID int) (*Blog, error) {
	query := "SELECT id, title, content, user_id, created_at FROM blogs WHERE id = ?"
	row := q.db.QueryRow(query, blogID)

	var blog Blog
	err := row.Scan(&blog.ID, &blog.Title, &blog.Content, &blog.UserID, &blog.PublishedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No blog found
		}
		return nil, err // Other error
	}
	return &blog, nil
}

func (q *Queries) DeleteBlog(blogID int) error {
	query := "DELETE FROM blogs WHERE id = ?"
	_, err := q.db.Exec(query, blogID)
	if err != nil {
		return err
	}
	return nil
}

func (q *Queries) GetBlogs(string, limit, offset int) ([]*Blog, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}
	query := "SELECT * FROM blogs LIMIT ? OFFSET ?"
	rows, err := q.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []*Blog
	for rows.Next() {
		var blog Blog
		if err := rows.Scan(&blog.ID, &blog.Title, &blog.Content, &blog.UserID, &blog.PublishedAt); err != nil {
			return nil, err
		}
		blogs = append(blogs, &blog)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return blogs, nil
}

func (q *Queries) GetBlogCount() (int, error) {
	query := "SELECT COUNT(*) FROM blogs"
	row := q.db.QueryRow(query)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (q *Queries) GetBlogsByTitle(title string, limit, offset int) ([]*Blog, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}
	query := "SELECT * FROM blogs WHERE title LIKE ? LIMIT ? OFFSET ?"
	rows, err := q.db.Query(query, "%"+title+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []*Blog
	for rows.Next() {
		var blog Blog
		if err := rows.Scan(&blog.ID, &blog.Title, &blog.Content, &blog.UserID, &blog.PublishedAt); err != nil {
			return nil, err
		}
		blogs = append(blogs, &blog)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return blogs, nil
}

func (q *Queries) GetBlogsByUserID(userID int) ([]*Blog, error) {
	query := "SELECT * FROM blogs WHERE user_id = ?"
	rows, err := q.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []*Blog
	for rows.Next() {
		var blog Blog
		if err := rows.Scan(&blog.ID, &blog.Title, &blog.Content, &blog.UserID); err != nil {
			return nil, err
		}
		blogs = append(blogs, &blog)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return blogs, nil
}
