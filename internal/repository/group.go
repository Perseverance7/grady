package repository

import (
	"fmt"

	"github.com/Perseverance7/grady/internal/models"
	"github.com/jmoiron/sqlx"
)

type GroupRepository struct {
	db *sqlx.DB
}

func NewGroupRepository(db *sqlx.DB) *GroupRepository {
	return &GroupRepository{
		db: db,
	}
}

func (gr *GroupRepository) CreateGroup(group *models.CreateGroupReq) error {
	tx, err := gr.db.Beginx()
	if err != nil {
		return err
	}

	// Вставка группы и получение её ID
	var groupId int
	query := fmt.Sprintf("INSERT INTO %s (name, description, created_by) VALUES ($1, $2, $3) RETURNING id", tableGroups)
	err = tx.QueryRow(query, group.Name, group.Description, group.CreatedBy).Scan(&groupId)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Добавление пользователя в качестве владельца группы
	query = fmt.Sprintf("INSERT INTO %s (group_id, user_id, role) VALUES ($1, $2, $3)", tableGroupMembers)
	_, err = tx.Exec(query, groupId, group.CreatedBy, "owner")
	if err != nil {
		tx.Rollback()
		return err
	}

	// Завершение транзакции
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
