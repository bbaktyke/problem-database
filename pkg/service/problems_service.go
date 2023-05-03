package service

import (
	"time"

	"git.01.alem.school/bbaktyke/test.project.git/pkg/models"
	"git.01.alem.school/bbaktyke/test.project.git/pkg/repository"
)

type ProblemServiceImpl struct {
	repo repository.Problem
}

func NewProblemService(repo repository.Problem) ProblemService {
	return &ProblemServiceImpl{
		repo: repo,
	}
}

func (p *ProblemServiceImpl) CreateProblem(userid int, merged models.ProblemWithTopics) (int, error) {
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

func (p *ProblemServiceImpl) GetProblems(pageNum, pageSize int) ([]models.ProblemWithTopics, error) {
	merged, err := p.repo.GetAllProblemsRepo(pageNum, pageSize)
	if err != nil {
		return nil, err
	}
	return merged, nil
}

func (p *ProblemServiceImpl) GetProblemByID(id int) (models.ProblemWithTopics, error) {
	var merged models.ProblemWithTopics
	merged, err := p.repo.GetByIDProblemsRepo(id)
	if err != nil {
		return merged, err
	}
	return merged, nil
}

func (p *ProblemServiceImpl) DeleteProblem(problemid int) error {
	return p.repo.DeleteRepo(problemid)
}

func (p *ProblemServiceImpl) UpdateProblem(problemid int, upd models.ProblemUpdate) error {
	if err := upd.Validate(); err != nil {
		return err
	}
	upd.UpdatedAt = *GetTime()
	return p.repo.UpdateRepo(problemid, upd)
}

func (p *ProblemServiceImpl) AccessRight(userId, problemId int) error {
	return p.repo.AccessRight(userId, problemId)
}

func (p *ProblemServiceImpl) GetByParameter(topic, level string) ([]models.ProblemWithTopics, error) {
	return p.repo.GetByParameter(topic, level)
}

func (p *ProblemServiceImpl) SearchProblem(title string) ([]models.ProblemWithTopics, error) {
	return p.repo.SearchProblemRepo(title)
}

func GetTime() *time.Time {
	time, err := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	if err != nil {
		return nil
	}
	return &time
}
