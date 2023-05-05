package repository

import (
	"context"

	"git.01.alem.school/bbaktyke/test.project.git/internal/models"
	"github.com/jmoiron/sqlx"
)

type Problem interface {
	CreateProblemRepo(ctx context.Context, merged models.ProblemWithTopics) (int, error)
	GetAllProblemsRepo(ctx context.Context, pageNum, pageSize int) ([]models.ProblemWithTopics, error)
	GetByIDProblemsRepo(ctx context.Context, id int) (models.ProblemWithTopics, error)
	DeleteRepo(ctx context.Context, id int) error
	UpdateRepo(ctx context.Context, id int, upd models.ProblemUpdate) error
	GetByParameter(ctx context.Context, topic, level string) ([]models.ProblemWithTopics, error)
	AccessRight(ctx context.Context, userId, problemId int) error
	SearchProblemRepo(ctx context.Context, title string) ([]models.ProblemWithTopics, error)
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
