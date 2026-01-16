package repository

import (
	"database/sql"
	"example/web-service-gin/app/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user models.User) error {
	_, err := r.db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, user.Password)
	return err
}

func (r *UserRepository) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	row := r.db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username)
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return user, err
	}
	return user, nil
}
