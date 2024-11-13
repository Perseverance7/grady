package service

import (
	"github.com/Perseverance7/grady/internal/repository"
)

type NotificationService struct {
	repo repository.Notification
}

func NewNotificationService(repo repository.Notification) *NotificationService {
	return &NotificationService{
		repo: repo,
	}
}
