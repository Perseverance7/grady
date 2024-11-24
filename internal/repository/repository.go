package repository

import (
	"github.com/Perseverance7/grady/internal/models"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(input models.UserRegisterReq) (models.UserRegisterRes, error)
	GetUser(email, password string) (models.UserLogin, error)
	GetUserSalt(email string) (string, error)
	GetRefreshToken(id string) (string, error)
	UpdateUser(user *models.User) (*models.User, error)
	DeleteUser(id int64) error
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
	SaveMessage(message *models.Message) error
	GetMessagesByGroup(groupID string, limit, offset int) ([]models.Message, error)
	GetUserData(userID int64) (models.Message, error)
	IsUserInGroup(userID int64, groupID string) (bool, error)
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
