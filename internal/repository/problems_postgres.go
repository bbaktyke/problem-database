package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"git.01.alem.school/bbaktyke/test.project.git/internal/models"
	"github.com/jmoiron/sqlx"
)

const (
	createProblemQuery   = "INSERT INTO " + problemTable + " (user_id, title, description, level, samples, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	createdRelationQuery = "INSERT INTO " + relationTable + " (problem_id, topic_id) VALUES ($1, $2)"
	getProblemQuery      = "SELECT * FROM " + problemTable + " ORDER BY created_at DESC LIMIT $1 OFFSET $2"
	getTopicsQuery       = "SELECT t.id, t.topics FROM " + topicTable + " t LEFT JOIN relation r ON t.id = r.topic_id WHERE r.problem_id = $1 OR r.problem_id IS NULL"
	getProblemByID       = "SELECT * FROM " + problemTable + " WHERE id = $1 ORDER BY created_at DESC"
	deleteProblem        = "DELETE FROM " + problemTable + " WHERE id = $1"
	getUserIdProblem     = "SELECT user_id FROM " + problemTable + " WHERE id = $1"
	getTopicIDQuery      = "SELECT id FROM " + topicTable + " WHERE topics = $1"
	getParameterQuery    = "SELECT problem.id FROM problem INNER JOIN relation on relation.problem_id=problem.id where relation.topic_id = $1 and problem.level = $2"
)

type ProblemsPostgres struct {
	db *sqlx.DB
}

func NewProblemRepository(db *sqlx.DB) *ProblemsPostgres {
	return &ProblemsPostgres{db: db}
}

func (p *ProblemsPostgres) CreateProblemRepo(ctx context.Context, merged models.ProblemWithTopics) (int, error) {
	merged.Problem.CreatedAt = GetTime()
	merged.Problem.UpdatedAt = GetTime()
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	var problemID int

	row := tx.QueryRowContext(ctx, createProblemQuery, merged.Problem.UserID, merged.Problem.Title, merged.Problem.Description, merged.Problem.Level, merged.Problem.Samples, merged.Problem.CreatedAt, merged.Problem.UpdatedAt)
	if err := row.Scan(&problemID); err != nil {
		tx.Rollback()
		return 0, err
	}
	for _, v := range merged.Topics {
		_, err := tx.ExecContext(ctx, createdRelationQuery, problemID, v.ID)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return problemID, nil
}

func (p *ProblemsPostgres) GetAllProblemsRepo(ctx context.Context, pageNum, pageSize int) ([]models.ProblemWithTopics, error) {
	offset := (pageNum - 1) * pageSize
	limit := pageSize
	var mergedstruct []models.ProblemWithTopics

	row, err := p.db.QueryContext(ctx, getProblemQuery, limit, offset)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	for row.Next() {
		pwt := models.ProblemWithTopics{}

		err := row.Scan(&pwt.Problem.ID, &pwt.Problem.UserID, &pwt.Problem.Title, &pwt.Problem.Description, &pwt.Problem.Level, &pwt.Problem.Samples, &pwt.Problem.CreatedAt, &pwt.Problem.UpdatedAt)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, models.ErrNoRecord
			} else {
				return nil, err
			}
		}

		problemIDRow, err := p.db.QueryContext(ctx, getTopicsQuery, pwt.Problem.ID)
		if err != nil {
			return nil, err
		}
		defer problemIDRow.Close()

		var topics []models.Topic
		for problemIDRow.Next() {
			t := models.Topic{}
			err := problemIDRow.Scan(&t.ID, &t.Topics)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, models.ErrNoRecord
				} else {
					return nil, err
				}
			}
			topics = append(topics, t)
		}

		if err := problemIDRow.Err(); err != nil {
			return nil, err
		}

		pwt.Topics = topics
		mergedstruct = append(mergedstruct, pwt)
	}

	if err := row.Err(); err != nil {
		return nil, err
	}

	return mergedstruct, nil
}

func (p *ProblemsPostgres) GetByIDProblemsRepo(ctx context.Context, id int) (models.ProblemWithTopics, error) {
	var pwt models.ProblemWithTopics

	row := p.db.QueryRowContext(ctx, getProblemByID, id)
	if err := row.Scan(&pwt.Problem.ID, &pwt.Problem.UserID, &pwt.Problem.Title, &pwt.Problem.Description, &pwt.Problem.Level, &pwt.Problem.Samples, &pwt.Problem.CreatedAt, &pwt.Problem.UpdatedAt); err != nil {
		return pwt, err
	}
	problemIDRow, err := p.db.QueryContext(ctx, getTopicsQuery, id)
	if err != nil {
		return pwt, err
	}
	defer problemIDRow.Close()

	var topics []models.Topic
	for problemIDRow.Next() {
		t := models.Topic{}
		err := problemIDRow.Scan(&t.ID, &t.Topics)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return pwt, models.ErrNoRecord
			} else {
				return pwt, err
			}
		}
		topics = append(topics, t)
	}
	pwt.Topics = topics

	return pwt, nil
}

func (p *ProblemsPostgres) DeleteRepo(ctx context.Context, problemid int) error {
	_, err := p.db.ExecContext(ctx, deleteProblem, problemid)
	return err
}

func (p *ProblemsPostgres) UpdateRepo(ctx context.Context, problemid int, upd models.ProblemUpdate) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1
	upd.UpdatedAt = GetTime()

	if upd.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *upd.Title)
		argId++
	}
	if upd.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *upd.Description)
		argId++
	}
	if upd.Level != nil {
		setValues = append(setValues, fmt.Sprintf("level=$%d", argId))
		args = append(args, *upd.Level)
		argId++
	}
	if upd.Samples != nil {
		setValues = append(setValues, fmt.Sprintf("samples=$%d", argId))
		args = append(args, *upd.Samples)
		argId++
	}
	setValues = append(setValues, fmt.Sprintf("updated_at=$%d", argId))
	args = append(args, upd.UpdatedAt)
	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = %d", problemTable, setQuery, problemid)

	_, err := p.db.ExecContext(ctx, query, args...)
	return err
}

func (p *ProblemsPostgres) AccessRight(ctx context.Context, userId, problemId int) error {
	var authorId int
	row := p.db.QueryRowContext(ctx, getUserIdProblem, problemId)

	if err := row.Scan(&authorId); err != nil {
		return err
	}

	if authorId != userId {
		return errors.New("Not your post. U can not delete it")
	}
	return nil
}

func (p *ProblemsPostgres) GetByParameter(ctx context.Context, topic, level string) ([]models.ProblemWithTopics, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var topicID int
	row := tx.QueryRowContext(ctx, getTopicIDQuery, topic)
	if err := row.Scan(&topicID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	tx2, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx2.Rollback()

	problemIDRow, err := tx2.QueryContext(ctx, getParameterQuery, topicID, level)
	if err != nil {
		return nil, err
	}
	defer problemIDRow.Close()

	var idSlice []int
	for problemIDRow.Next() {
		var id int
		err := problemIDRow.Scan(&id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, models.ErrNoRecord
			} else {
				return nil, err
			}
		}
		idSlice = append(idSlice, id)
	}

	if err := tx2.Commit(); err != nil {
		return nil, err
	}

	var mergedstruct []models.ProblemWithTopics
	for _, v := range idSlice {
		unit, err := p.GetByIDProblemsRepo(ctx, v)
		if err != nil {
			return nil, err
		}
		mergedstruct = append(mergedstruct, unit)
	}
	return mergedstruct, nil
}

func (p *ProblemsPostgres) SearchProblemRepo(ctx context.Context, title string) ([]models.ProblemWithTopics, error) {
	stmt, err := p.db.Prepare("SELECT id FROM problem WHERE title LIKE $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, fmt.Sprintf("%%%s%%", title))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var idArr []int
	for rows.Next() {
		var id int

		err := rows.Scan(&id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, models.ErrNoRecord
			} else {
				return nil, err
			}
		}
		idArr = append(idArr, id)

	}
	var mergedstruct []models.ProblemWithTopics
	for _, v := range idArr {
		unit, err := p.GetByIDProblemsRepo(ctx, v)
		if err != nil {
			return nil, err
		}
		mergedstruct = append(mergedstruct, unit)
	}
	return mergedstruct, nil
}

func GetTime() time.Time {
	t, err := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	if err != nil {
		return time.Time{}
	}
	return t
}
