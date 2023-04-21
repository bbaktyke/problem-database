package service

import (
	"time"

	"git.01.alem.school/bbaktyke/test.project.git/pkg/models"
	"git.01.alem.school/bbaktyke/test.project.git/pkg/repository"
)

type ProblemService struct {
	repo repository.Problem
}

func NewProblemService(repo repository.Problem) *ProblemService {
	return &ProblemService{
		repo: repo,
	}
}

func (p *ProblemService) CreateService(userid int, merged models.ProblemWithTopics) (int, error) {
	merged.Problem.UserID = userid
	var err error
	merged.Problem.CreatedAt = *GetTime()
	merged.Problem.UpdatedAt = *GetTime()
	id, err := p.repo.CreateProblemRepo(merged)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (p *ProblemService) GetAllService(pageNum, pageSize int) ([]models.ProblemWithTopics, error) {
	merged, err := p.repo.GetAllProblemsRepo(pageNum, pageSize)
	if err != nil {
		return nil, err
	}
	return merged, nil
}

func (p *ProblemService) GetByIDService(id int) (models.ProblemWithTopics, error) {
	var merged models.ProblemWithTopics
	merged, err := p.repo.GetByIDProblemsRepo(id)
	if err != nil {
		return merged, err
	}
	return merged, nil
}

func (p *ProblemService) DeleteService(problemid int) error {
	return p.repo.DeleteRepo(problemid)
}

func (p *ProblemService) UpdateService(problemid int, upd models.ProblemUpdate) error {
	if err := upd.Validate(); err != nil {
		return err
	}
	upd.UpdatedAt = *GetTime()
	return p.repo.UpdateRepo(problemid, upd)
}

func (p *ProblemService) AccessRightService(userId, problemId int) error {
	return p.repo.AccessRight(userId, problemId)
}

func (p *ProblemService) GetByParameter(topic, level string) ([]models.ProblemWithTopics, error) {
	return p.repo.GetByParameter(topic, level)
}

func (p *ProblemService) SearchProblemService(title string) ([]models.ProblemWithTopics, error) {
	return p.repo.SearchProblemRepo(title)
}

func GetTime() *time.Time {
	time, err := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	if err != nil {
		return nil
	}
	return &time
}
