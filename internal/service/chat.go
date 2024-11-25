package service

import (
	"github.com/Perseverance7/grady/internal/models"
	"github.com/Perseverance7/grady/internal/repository"
)

type ChatService struct {
	repo repository.Chat
}

func NewChatService(repo repository.Chat) *ChatService {
	return &ChatService{repo: repo}
}

func (s *ChatService) SendMessage(message *models.Message) error {
	return s.repo.SaveMessage(message)
}

func (s *ChatService) GetUserData(userID int64) (models.Message, error) {
	return s.repo.GetUserData(userID)
}

func (s *ChatService) GetChatHistory(groupID string, limit, offset int) ([]models.Message, error) {
	return s.repo.GetMessagesByGroup(groupID, limit, offset)
}

func (s *ChatService) IsUserInGroup(userID int64, groupID string) (bool, error) {
    return s.repo.IsUserInGroup(userID, groupID)
}

