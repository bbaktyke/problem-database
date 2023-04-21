package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	usersTable           = "users"
	problemTable         = "problem"
	topicTable           = "topic"
	relationTable        = "relation"
	createUserQuery      = "INSERT INTO " + usersTable + " (name, username, password) VALUES ($1, $2, $3) RETURNING id"
	selectUserQuery      = "SELECT id FROM " + usersTable + " WHERE username=$1 AND password=$2"
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

type Confiq struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Confiq) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
