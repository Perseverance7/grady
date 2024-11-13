package service

import (
	"github.com/Perseverance7/grady/internal/repository"
)

type ChatService struct {
	repo repository.Chat
}

func NewChatService(repo repository.Chat) *ChatService {
	return &ChatService{
		repo: repo,
	}
}
