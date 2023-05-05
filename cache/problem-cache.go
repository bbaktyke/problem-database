package cache

import "git.01.alem.school/bbaktyke/test.project.git/internal/models"

type ProblemCash interface {
	Set(key string, value *models.ProblemWithTopics) error
	Get(key string) *models.ProblemWithTopics
	SetArr(key string, value *[]models.ProblemWithTopics) error
	GetArr(key string) *[]models.ProblemWithTopics
}
