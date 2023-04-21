package repository

import (
	"git.01.alem.school/bbaktyke/test.project.git/pkg/models"
	"github.com/jmoiron/sqlx"
)

type Problem interface {
	CreateProblemRepo(merged models.ProblemWithTopics) (int, error)
	GetAllProblemsRepo(pageNum, pageSize int) ([]models.ProblemWithTopics, error)
	GetByIDProblemsRepo(id int) (models.ProblemWithTopics, error)
	DeleteRepo(id int) error
	UpdateRepo(id int, upd models.ProblemUpdate) error
	GetByParameter(topic, level string) ([]models.ProblemWithTopics, error)
	AccessRight(userId, problemId int) error
	SearchProblemRepo(title string) ([]models.ProblemWithTopics, error)
}

type Authorization interface {
	CreateUserRepo(user models.User) (int, error)
	GetUserRepo(username, password string) (models.User, error)
}

type Repository struct {
	Problem
	Authorization
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthRepository(db),
		Problem:       NewProblemRepository(db),
	}
}
