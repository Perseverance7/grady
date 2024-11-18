package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	tableUsers           = "users"
	tableGroups          = "groups"
	tableGroupMembers    = "group_members"
	tableTests           = "tests"
	tableQuestions       = "questions"
	tableAnswerOptions   = "answer_options"
	tableTestSubmissions = "test_submissions"
	tableStudentAnswers  = "student_answers"
	tableMessages        = "messages"
	tableSessions        = "sessions"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg *Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	return db, nil
}
