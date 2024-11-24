package service

import (
	"time"

	"github.com/Perseverance7/grady/internal/models"
	"github.com/Perseverance7/grady/internal/repository"
)

type Authorization interface {
	CreateUser(input models.UserRegisterReq) (models.UserRegisterRes, error)
	GetUser(email, password string) (models.UserLogin, error)
	UpdateUser(user *models.User) (*models.User, error)
	CreateToken(id int64, email string, isAdmin bool, duration time.Duration) (string, *models.UserClaims, error)
	VerifyAccessToken(tokenStr string) (*models.UserClaims, error)
	VerifyRefreshToken(refreshTokenUUID string) (*models.UserClaims, error)
	CreateSession(session *models.Session) (*models.Session, error)
	GetSession(id string) (*models.Session, error)
	RevokeSession(id string) error
	DeleteSession(id string) error 
}

type Task interface {
}

type Group interface {
	CreateGroup(group *models.CreateGroupReq) error
}

type Notification interface {
}

type Chat interface {
	SendMessage(message *models.Message) error
	GetChatHistory(groupID string) ([]models.Message, error)
	GetUserData(userID int64) (models.Message, error)
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

func NewService(repo *repository.Repository, secretKey []byte) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization, secretKey),
		Task:          NewTaskService(repo.Task),
		Group:         NewGroupService(repo.Group),
		Notification:  NewNotificationService(repo.Notification),
		Chat:          NewChatService(repo.Chat),
		Statistics:    NewStatisticsService(repo.Statistics),
	}
}
