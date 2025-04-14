package repositories

import (
	"database/sql"
	"financierGo/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) Create(user *models.User) error {
	query := `INSERT INTO users (username, email, password_hash, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id`
	return r.DB.QueryRow(query, user.Username, user.Email, user.Password).Scan(&user.ID)
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, created_at FROM users WHERE email = $1`
	row := r.DB.QueryRow(query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
