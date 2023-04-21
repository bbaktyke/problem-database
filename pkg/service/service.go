package service

import (
	"git.01.alem.school/bbaktyke/test.project.git/pkg/models"
	"git.01.alem.school/bbaktyke/test.project.git/pkg/repository"
)

type Problem interface {
	CreateService(userid int, problem models.ProblemWithTopics) (int, error)
	GetAllService(pageNum, pageSize int) ([]models.ProblemWithTopics, error)
	GetByIDService(id int) (models.ProblemWithTopics, error)
	DeleteService(problemid int) error
	UpdateService(problemid int, upd models.ProblemUpdate) error
	AccessRightService(userId, problemId int) error
	GetByParameter(topic, level string) ([]models.ProblemWithTopics, error)
	SearchProblemService(title string) ([]models.ProblemWithTopics, error)
}

type Authorization interface {
	CreateUserService(user models.User) (int, error)
	GenerateTokenService(username, password string) (string, error)
	ParseTokenService(tokenString string) (int, error)
}

type Service struct {
	Problem
	Authorization
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
		Problem:       NewProblemService(repo.Problem),
	}
}
