package repository

import (
	"git.01.alem.school/bbaktyke/test.project.git/internal/models"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

const (
	createUserQuery = "INSERT INTO " + usersTable + " (name, username, password) VALUES ($1, $2, $3) RETURNING id"
	selectUserQuery = "SELECT id FROM " + usersTable + " WHERE username=$1 AND password=$2"
)

func NewAuthRepository(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (auth *AuthPostgres) CreateUserRepo(user models.User) (int, error) {
	var id int
	row := auth.db.QueryRow(createUserQuery, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (auth *AuthPostgres) GetUserRepo(username, password string) (models.User, error) {
	var user models.User
	err := auth.db.Get(&user, selectUserQuery, username, password)
	return user, err
}
