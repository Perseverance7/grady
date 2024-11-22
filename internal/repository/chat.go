package repository

import (
	"fmt"

	"github.com/Perseverance7/grady/internal/models"
	"github.com/jmoiron/sqlx"
)

type ChatRepository struct {
	db *sqlx.DB
}

func NewChatRepository(db *sqlx.DB) *ChatRepository {
	return &ChatRepository{
		db: db,
	}
}

func (cr *ChatRepository) SaveMessage(msg models.Message) error {
	query := fmt.Sprintf("INSERT INTO %s (group_id, user_id, content) VALUES ($1, $2, $3)", tableMessages)
	_, err := cr.db.Exec(query, msg.GroupID, msg.UserID, msg.Content)

	return err
}

func (cr *ChatRepository) GetMessagesByGroupId(groupID string, limit int) ([]models.Message, error) {
	var messages []models.Message

	query := fmt.Sprintf("SELECT group_id, user_id, content, sent_at FROM %s WHERE group_id = $1 ORDER BY sent_at DESC LIMIT $2", tableMessages)
	err := cr.db.Select(&messages, query, groupID, limit)
	if err != nil{
		return nil, err
	}

	return messages, nil
}
