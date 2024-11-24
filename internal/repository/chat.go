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
	return &ChatRepository{db: db}
}

// Сохранить сообщение в базе данных
func (r *ChatRepository) SaveMessage(message *models.Message) error {
	query := fmt.Sprintf(`INSERT INTO %s (group_id, user_id, content, sent_at) VALUES ($1, $2, $3, $4)`, tableMessages)
	_, err := r.db.Exec(query, message.GroupID, message.UserID, message.Content, message.SentAt)
	return err
}

// Получить сообщения для группы
func (r *ChatRepository) GetMessagesByGroup(groupID string, limit, offset int) ([]models.Message, error) {
	var messages []models.Message
	query := fmt.Sprintf(`SELECT m.user_id AS user_id,
	                             m.group_id AS group_id, 
								 u.name AS name,
								 u.surname AS surname,
								 u.patronymic AS patronymic,
								 m.content AS content, 
								 m.sent_at AS sent_at 
								 FROM %s m
								 LEFT JOIN %s u
								 ON u.id = m.user_id
								 WHERE m.group_id=$1
								 ORDER BY m.sent_at DESC
								 LIMIT $2 OFFSET $3`, tableMessages, tableUsers)
	err := r.db.Select(&messages, query, groupID, limit, offset)
	return messages, err
}

func (r *ChatRepository) GetUserData(userID int64) (models.Message, error) {
	var messageData models.Message
	query := fmt.Sprintf(`SELECT name, 
								 surname, 
								 patronymic 
								 FROM %s
								 WHERE id=$1`, tableUsers)
	err := r.db.Get(&messageData, query, userID)
	if err != nil {
		return models.Message{}, err
	}

	return messageData, nil
}

func (r *ChatRepository) IsUserInGroup(userID int64, groupID string) (bool, error) {
	var exists bool
	query := fmt.Sprintf(`SELECT EXISTS (
		SELECT 1 FROM %s WHERE user_id = $1 AND group_id = $2
	)`, tableGroupMembers)
	err := r.db.QueryRow(query, userID, groupID).Scan(&exists)
	return exists, err
}
