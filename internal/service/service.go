package service

import (
	"github.com/Perseverance7/grady/internal/models"
	"github.com/Perseverance7/grady/internal/repository"
)

type Authorization interface {
	CreateUser(input models.UserRegister) (int, error)
	GenerateTokens(email, password string) (string, string, error)
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

type Service struct {
	Authorization
	Task
	Group
	Notification
	Chat
	Statistics
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
		Task:          NewTaskService(repo.Task),
		Group:         NewGroupService(repo.Group),
		Notification:  NewNotificationService(repo.Notification),
		Chat:          NewChatService(repo.Chat),
		Statistics:    NewStatisticsService(repo.Statistics),
	}
}
