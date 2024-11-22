package service

import (
	"github.com/Perseverance7/grady/internal/models"
	"github.com/Perseverance7/grady/internal/repository"
)

type ChatService struct {
	repo      repository.Chat
	broadcast chan models.Message
}

func NewChatService(repo repository.Chat) *ChatService {
	return &ChatService{
		repo:      repo,
		broadcast: make(chan models.Message),
	}
}

func (cs *ChatService) SendMessage(msg models.Message) error {
	err := cs.repo.SaveMessage(msg)
	if err != nil {
		return err
	}

	cs.broadcast <- msg
	return nil
}

func (cs *ChatService) GetChatHistory(groupID string, limit int) ([]models.Message, error) {
	return cs.repo.GetMessagesByGroupId(groupID, limit)
}


