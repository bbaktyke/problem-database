package service

import (
	"git.01.alem.school/bbaktyke/test.project.git/pkg/models"
)

type ProblemService interface {
	CreateProblem(userid int, problem models.ProblemWithTopics) (int, error)
	GetProblems(pageNum, pageSize int) ([]models.ProblemWithTopics, error)
	GetProblemByID(id int) (models.ProblemWithTopics, error)
	DeleteProblem(problemid int) error
	UpdateProblem(problemid int, upd models.ProblemUpdate) error
	AccessRight(userId, problemId int) error
	GetByParameter(topic, level string) ([]models.ProblemWithTopics, error)
	SearchProblem(title string) ([]models.ProblemWithTopics, error)
}

type AuthorizationService interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(tokenString string) (int, error)
}
