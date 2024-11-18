package service

import (
	"github.com/Perseverance7/grady/internal/models"
	"github.com/Perseverance7/grady/internal/repository"
)

type Authorization interface {
	CreateUser(input models.UserRegisterReq) (models.UserRegisterRes, error)
	GetUser(email, password string) (models.UserLogin, error)
	UpdateUser(user *models.User) (*models.User, error)
	CreateSession(session *models.Session) (*models.Session, error)
	GetSession(id string) (*models.Session, error)
	RevokeSession(id string) error
	DeleteSession(id string) error 
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
