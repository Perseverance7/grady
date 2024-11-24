package models

import "time"

type Message struct {
	UserID     int64     `json:"user_id" db:"user_id"`
	GroupID    string    `db:"group_id"`
	Name       string    `json:"name" db:"name"`
	Surname    string    `json:"surname" db:"surname"`
	Patronymic string    `json:"patronymic" db:"patronymic"`
	Content    string    `json:"content" db:"content"`
	SentAt     time.Time `json:"sent_at" db:"sent_at"`
}

type IncomingMessage struct {
	Content string `json:"content"`
}
