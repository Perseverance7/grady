package repository

import "github.com/jmoiron/sqlx"

type Authorization interface {
}

type Task interface {
}

type Group interface {
}

type Notification interface {
}

type Chat interface {
}

type Statistics interface {
}

type Repository struct {
	Authorization
	Task
	Group
	Notification
	Chat
	Statistics
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthRepository(db),
		Task:          NewTaskRepository(db),
		Group:         NewGroupRepository(db),
		Notification:  NewNotificationRepository(db),
		Chat:          NewChatRepository(db),
		Statistics:    NewStatisticsRepository(db),
	}
}
