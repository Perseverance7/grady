package models

import "time"

type Message struct {
	GroupID  string `json:"group_id" db:"group_id"`
	UserID int64 `json:"user_id" db:"user_id"`
	Content  string `json:"content" db:"content"`
	SentAt time.Time `json:"sent_at" db:"sent_at"`
}
