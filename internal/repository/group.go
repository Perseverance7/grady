package repository

import "github.com/jmoiron/sqlx"

type GroupRepository struct {
	db *sqlx.DB
}

func NewGroupRepository(db *sqlx.DB) *GroupRepository {
	return &GroupRepository{
		db: db,
	}
}
