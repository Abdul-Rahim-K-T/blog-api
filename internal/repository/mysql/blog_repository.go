package mysql

import (
	"blog-api/internal/entity"
	"database/sql"
)

type BlogRepository struct {
	DB *sql.DB
}

func NewBlogRepository(db *sql.DB) *BlogRepository {
	return &BlogRepository{DB: db}
}

func (r *BlogRepository) Create(blog *entity.Blog) error {
	_, err := r.DB.Exec("INSERT INTO blogs (title, content, user_id, thumbnail) VALUES (?, ?, ?, ?)",
		blog.Title, blog.Content, blog.UserID, blog.Thumbnail)
	return err
}

func (r *BlogRepository) GetAll() ([]*entity.Blog, error) {
	rows, err := r.DB.Query("SELECT id, title, content, user_id, thumbnail FROM blogs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []*entity.Blog
	for rows.Next() {
		var blog entity.Blog
		if err := rows.Scan(&blog.ID, &blog.Title, &blog.Content, &blog.UserID, &blog.Thumbnail); err != nil {
			return nil, err
		}
		blogs = append(blogs, &blog)
	}

	return blogs, nil
}

func (r *BlogRepository) GetByID(id int) (*entity.Blog, error) {
	row := r.DB.QueryRow("SELECT id, title, content, user_id, thumbnail FROM blogs WHERE id = ?", id)

	var blog entity.Blog
	if err := row.Scan(&blog.ID, &blog.Title, &blog.Content, &blog.UserID, &blog.Thumbnail); err != nil {
		return nil, err
	}
	return &blog, nil
}

func (r *BlogRepository) Update(blog *entity.Blog) error {
	_, err := r.DB.Exec("UPDATE blogs SET title = ?, content = ?, thumbnail = ? WHERE id = ?",
		blog.Title, blog.Content, blog.Thumbnail, blog.ID)
	return err
}

func (r *BlogRepository) Delete(id int) error {
	_, err := r.DB.Exec("DELETE FROM blogs WHERE id = ?", id)
	return err
}

func (r *BlogRepository) CreateComment(comment *entity.Comment) error {
	query := `INSERT INTO comments (content, user_id, blog_id) VALUES (?, ?, ?)`
	_, err := r.DB.Exec(query, comment.Content, comment.UserID, comment.BlogID)
	return err
}
