package models

import (
	"errors"
	"strings"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Problem struct {
	ID          int       `json:"-" db:"id"`
	UserID      int       `json:"-" db:"user_id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Level       string    `json:"level" db:"level"`
	Samples     string    `json:"samples" db:"samples"`
	CreatedAt   time.Time `json:"-" db:"created_at"`
	UpdatedAt   time.Time `json:"-" db:"updated_at"`
}

type Topic struct {
	ID     int    `json:"id" db:"id"`
	Topics string `json:"name" db:"name"`
}

type Relation struct {
	ID         int
	ProblemsID int
	TopicID    int
}

type ProblemWithTopics struct {
	Problem Problem `json:"problem"`
	Topics  []Topic `json:"topics"`
}

type ProblemUpdate struct {
	Title       *string   `json:"title"`
	Description *string   `json:"description" `
	Level       *string   `json:"level" `
	Samples     *string   `json:"samples"`
	UpdatedAt   time.Time `json:"-" `
}

func (pu *ProblemUpdate) Validate() error {
	if pu.Title == nil && pu.Description == nil && pu.Level == nil && pu.Samples == nil {
		return errors.New("update struct has no values")
	}
	return nil
}

func (pwt *ProblemWithTopics) Validate() error {
	p := pwt.Problem
	if strings.TrimSpace(p.Title) == "" {
		return errors.New("title cannot be empty or contain only whitespaces")
	}
	if strings.TrimSpace(p.Description) == "" {
		return errors.New("description cannot be empty or contain only whitespaces")
	}

	if strings.TrimSpace(p.Level) == "" {
		return errors.New("level cannot be empty or contain only whitespaces")
	}
	if p.Level != "easy" && p.Level != "middle" && p.Level != "hard" {
		return errors.New("Not Valid values for level")
	}
	if strings.TrimSpace(p.Samples) == "" {
		return errors.New("samples cannot be empty or contain only whitespaces")
	}
	t := pwt.Topics
	for _, v := range t {
		if v.ID < 0 || v.ID > 4 {
			return errors.New("not valid topic id")
		}
	}
	return nil
}
