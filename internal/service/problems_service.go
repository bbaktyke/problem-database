package service

import (
	"context"

	"git.01.alem.school/bbaktyke/test.project.git/internal/models"
	"git.01.alem.school/bbaktyke/test.project.git/internal/repository"
)

type ProblemService interface {
	CreateProblem(ctx context.Context, pwt models.ProblemWithTopics) (int, error)
	GetProblems(ctx context.Context, pageNum, pageSize int) ([]models.ProblemWithTopics, error)
	GetProblemByID(ctx context.Context, id int) (models.ProblemWithTopics, error)
	DeleteProblem(ctx context.Context, problemid int) error
	UpdateProblem(ctx context.Context, problemid int, upd models.ProblemUpdate) error
	AccessRight(ctx context.Context, userId, problemId int) error
	GetByParameter(ctx context.Context, topic, level string) ([]models.ProblemWithTopics, error)
	SearchProblem(ctx context.Context, title string) ([]models.ProblemWithTopics, error)
}

type ProblemServiceImpl struct {
	repo repository.Problem
}

func NewProblemService(repo repository.Problem) ProblemService {
	return &ProblemServiceImpl{
		repo: repo,
	}
}

func (p *ProblemServiceImpl) CreateProblem(ctx context.Context, merged models.ProblemWithTopics) (int, error) {
	if err := merged.Validate(); err != nil {
		return 0, nil
	}

	id, err := p.repo.CreateProblemRepo(ctx, merged)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (p *ProblemServiceImpl) GetProblems(ctx context.Context, pageNum, pageSize int) ([]models.ProblemWithTopics, error) {
	merged, err := p.repo.GetAllProblemsRepo(ctx, pageNum, pageSize)
	if err != nil {
		return nil, err
	}

	return merged, nil
}

func (p *ProblemServiceImpl) GetProblemByID(ctx context.Context, id int) (models.ProblemWithTopics, error) {
	merged, err := p.repo.GetByIDProblemsRepo(ctx, id)
	if err != nil {
		return merged, err
	}

	return merged, nil
}

func (p *ProblemServiceImpl) DeleteProblem(ctx context.Context, problemid int) error {
	return p.repo.DeleteRepo(ctx, problemid)
}

func (p *ProblemServiceImpl) UpdateProblem(ctx context.Context, problemid int, upd models.ProblemUpdate) error {
	if err := upd.Validate(); err != nil {
		return err
	}

	return p.repo.UpdateRepo(ctx, problemid, upd)
}

func (p *ProblemServiceImpl) AccessRight(ctx context.Context, userId, problemId int) error {
	return p.repo.AccessRight(ctx, userId, problemId)
}

func (p *ProblemServiceImpl) GetByParameter(ctx context.Context, topic, level string) ([]models.ProblemWithTopics, error) {
	return p.repo.GetByParameter(ctx, topic, level)
}

func (p *ProblemServiceImpl) SearchProblem(ctx context.Context, title string) ([]models.ProblemWithTopics, error) {
	return p.repo.SearchProblemRepo(ctx, title)
}
