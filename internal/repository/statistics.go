package repository

import "github.com/jmoiron/sqlx"

type StatisticsRepository struct {
	db *sqlx.DB
}

func NewStatisticsRepository(db *sqlx.DB) *StatisticsRepository {
	return &StatisticsRepository{
		db: db,
	}
}
