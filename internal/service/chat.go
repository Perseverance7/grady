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


