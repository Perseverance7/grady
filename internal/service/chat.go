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

// Отправка сообщения
func (s *ChatService) SendMessage(message *models.Message) error {
	return s.repo.SaveMessage(message)
}

func (s *ChatService) GetUserData(userID int64) (models.Message, error) {
	return s.repo.GetUserData(userID)
}

// Получение истории чата
func (s *ChatService) GetChatHistory(groupID string, limit, offset int) ([]models.Message, error) {
	return s.repo.GetMessagesByGroup(groupID, limit, offset)
}

// Проверяет, принадлежит ли пользователь указанной группе
func (s *ChatService) IsUserInGroup(userID int64, groupID string) (bool, error) {
    // Обратитесь в репозиторий для проверки
    return s.repo.IsUserInGroup(userID, groupID)
}

