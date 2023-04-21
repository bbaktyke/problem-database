package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"git.01.alem.school/bbaktyke/test.project.git/pkg/models"
	"github.com/jmoiron/sqlx"
)

type ProblemsPostgres struct {
	db *sqlx.DB
}

func NewProblemRepository(db *sqlx.DB) *ProblemsPostgres {
	return &ProblemsPostgres{db: db}
}

func (p *ProblemsPostgres) CreateProblemRepo(merged models.ProblemWithTopics) (int, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return 0, err
	}
	var problemID int

	row := tx.QueryRow(createProblemQuery, merged.Problem.UserID, merged.Problem.Title, merged.Problem.Description, merged.Problem.Level, merged.Problem.Samples, merged.Problem.CreatedAt, merged.Problem.UpdatedAt)
	if err := row.Scan(&problemID); err != nil {
		tx.Rollback()
		return 0, err
	}
	for _, v := range merged.Topics {
		_, err := tx.Exec(createdRelationQuery, problemID, v.ID)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	return problemID, tx.Commit()
}

func (p *ProblemsPostgres) GetAllProblemsRepo(pageNum, pageSize int) ([]models.ProblemWithTopics, error) {
	offset := (pageNum - 1) * pageSize
	limit := pageSize
	var mergedstruct []models.ProblemWithTopics

	row, err := p.db.Query(getProblemQuery, limit, offset)
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

		problemIDRow, err := p.db.Query(getTopicsQuery, pwt.Problem.ID)
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

func (p *ProblemsPostgres) GetByIDProblemsRepo(id int) (models.ProblemWithTopics, error) {
	var pwt models.ProblemWithTopics
	row := p.db.QueryRow(getProblemByID, id)
	if err := row.Scan(&pwt.Problem.ID, &pwt.Problem.UserID, &pwt.Problem.Title, &pwt.Problem.Description, &pwt.Problem.Level, &pwt.Problem.Samples, &pwt.Problem.CreatedAt, &pwt.Problem.UpdatedAt); err != nil {
		return pwt, err
	}
	problemIDRow, err := p.db.Query(getTopicsQuery, id)
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

func (p *ProblemsPostgres) DeleteRepo(problemid int) error {
	_, err := p.db.Exec(deleteProblem, problemid)
	return err
}

func (p *ProblemsPostgres) UpdateRepo(problemid int, upd models.ProblemUpdate) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

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

	_, err := p.db.Exec(query, args...)
	return err
}

func (p *ProblemsPostgres) AccessRight(userId, problemId int) error {
	var authorId int
	row := p.db.QueryRow(getUserIdProblem, problemId)

	if err := row.Scan(&authorId); err != nil {
		return err
	}

	if authorId != userId {
		return errors.New("Not your post. U can not delete it")
	}
	return nil
}

func (p *ProblemsPostgres) GetByParameter(topic, level string) ([]models.ProblemWithTopics, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var topicID int
	row := tx.QueryRow(getTopicIDQuery, topic)
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

	tx2, err := p.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx2.Rollback()

	problemIDRow, err := tx2.Query(getParameterQuery, topicID, level)
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
		unit, err := p.GetByIDProblemsRepo(v)
		if err != nil {
			return nil, err
		}
		mergedstruct = append(mergedstruct, unit)
	}
	return mergedstruct, nil
}

func (p *ProblemsPostgres) SearchProblemRepo(title string) ([]models.ProblemWithTopics, error) {
	stmt, err := p.db.Prepare("SELECT id FROM problem WHERE title LIKE $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(fmt.Sprintf("%%%s%%", title))
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
		unit, err := p.GetByIDProblemsRepo(v)
		if err != nil {
			return nil, err
		}
		mergedstruct = append(mergedstruct, unit)
	}
	return mergedstruct, nil
}
