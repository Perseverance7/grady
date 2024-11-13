package repository

import "github.com/jmoiron/sqlx"

type NotificationRepository struct {
	db *sqlx.DB
}

func NewNotificationRepository(db *sqlx.DB) *NotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}
