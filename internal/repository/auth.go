package repository

import (
	"fmt"

	"github.com/Perseverance7/grady/internal/models"
	"github.com/jmoiron/sqlx"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) CreateUser(user models.UserRegister) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, surname, patronymic, email, password_hash, password_salt) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", tableUsers)
	row := r.db.QueryRow(query, user.Name, user.Surname, user.Patronymic, user.Email, user.Password, user.Salt)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthRepository) GetUserSalt(email string) (string, error) {
	var salt string
	query := fmt.Sprintf("SELECT password_salt FROM %s WHERE email=$1", tableUsers)
	row := r.db.QueryRow(query, email)
	if err := row.Scan(&salt); err != nil {
		return "", err
	}

	return salt, nil
}

func (r *AuthRepository) GetUser(email, password string) (int, error) {
	var id int
	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1 AND password_hash=$2", tableUsers)
	err := r.db.Get(&id, query, email, password)

	return id, err
}

