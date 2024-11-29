package mysql

import (
	"blog-api/internal/entity"
	"database/sql"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *entity.User) error {
	_, err := r.DB.Exec("INSERT INTO users (username, password, email, role) VALUES(?, ?, ?, ?)",
		user.Username, user.Password, user.Email, user.Role)
	return err
}

func (r *UserRepository) GetByUsernameAndPassword(username, password string) (*entity.User, error) {
	row := r.DB.QueryRow("SELECT id, username, email, role FROM users WHERE username = ? AND password = ?",
		username, password)

	var user entity.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUsernameOrEmail(username, email string) (*entity.User, error) {
	row := r.DB.QueryRow("SELECT id, username, email, role FROM users WHERE username = ? OR email = ?",
		username, email)

	var user entity.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*entity.User, error) {
	row := r.DB.QueryRow("SELECT id, username, password, email, role FROM users WHERE username = ?", username)

	var user entity.User
	if err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Role); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
