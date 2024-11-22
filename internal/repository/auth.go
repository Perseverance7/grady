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

func (r *AuthRepository) CreateUser(userReq models.UserRegisterReq) (models.UserRegisterRes, error) {
	var userRes models.UserRegisterRes
	query := fmt.Sprintf(`INSERT INTO %s (name, surname, patronymic, email, password_hash, password_salt) 
						  VALUES ($1, $2, $3, $4, $5, $6) 
						  RETURNING name, surname, patronymic, email, is_admin`, tableUsers)

	row := r.db.QueryRow(query,
		userReq.Name,
		userReq.Surname,
		userReq.Patronymic,
		userReq.Email,
		userReq.Password,
		userReq.Salt)

	if err := row.Scan(&userRes.Name, &userRes.Surname, &userRes.Patronymic, &userRes.Email, &userRes.IsAdmin); err != nil {
		return models.UserRegisterRes{}, err
	}
	return userRes, nil
}

func (r *AuthRepository) GetUser(email, password string) (models.UserLogin, error) {
	var user models.UserLogin
	query := fmt.Sprintf("SELECT id, name, surname, patronymic, email, is_admin FROM %s WHERE email=$1 AND password_hash=$2", tableUsers)

	row := r.db.QueryRow(query, email, password)

	if err := row.Scan(&user.ID, &user.Name, &user.Surname, &user.Patronymic, &user.Email, &user.IsAdmin); err != nil {
		return models.UserLogin{}, err
	}

	return user, nil
}

func (r *AuthRepository) UpdateUser(user *models.User) (*models.User, error) {
	query := fmt.Sprintf("UPDATE %s SET name=$1, surname=$2, patronymic=$3, password_hash=$4 WHERE id=$5", tableUsers)
	_, err := r.db.Exec(query, user.Name, user.Surname, user.Patronymic, user.Password, user.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *AuthRepository) DeleteUser(id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", tableUsers)
	_, err := r.db.Exec(query, id)

	return err
}

func (a *AuthRepository) GetUserSalt(email string) (string, error) {
	var salt string
	query := fmt.Sprintf("SELECT password_salt FROM %s WHERE email=$1", tableUsers)
	row := a.db.QueryRow(query, email)
	if err := row.Scan(&salt); err != nil {
		return "", err
	}

	return salt, nil
}

func (a *AuthRepository) GetRefreshToken(id string) (string, error) {
	var refreshToken string
	query := fmt.Sprintf("SELECT refresh_token FROM %s WHERE id=$1", tableSessions)

	row := a.db.QueryRow(query, id)
	if err := row.Scan(&refreshToken); err != nil {
		return "", fmt.Errorf("error getting session %w", err)
	}

	return refreshToken, nil
}

func (a *AuthRepository) CreateSession(session *models.Session) (*models.Session, error) {
	query := fmt.Sprintf("INSERT INTO %s (id, user_email, refresh_token, is_revoked, expires_at) VALUES ($1, $2, $3, $4, $5)", tableSessions)
	_, err := a.db.Exec(query, session.ID, session.UserEmail, session.RefreshToken, session.IsRevoked, session.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("error inserting session %w", err)
	}

	return session, nil
}

func (a *AuthRepository) GetSession(id string) (*models.Session, error) {
	var s models.Session
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", tableSessions)

	row := a.db.QueryRow(query, id)
	if err := row.Scan(&s.ID, &s.UserEmail, &s.RefreshToken, &s.IsRevoked, &s.CreatedAt, &s.ExpiresAt); err != nil {
		return nil, fmt.Errorf("error getting session %w", err)
	}

	return &s, nil
}

func (a *AuthRepository) RevokeSession(id string) error {
	query := fmt.Sprintf("UPDATE %s SET is_revoked=true WHERE id=$1", tableSessions)
	_, err := a.db.Exec(query, id)

	return err
}

func (a *AuthRepository) DeleteSession(id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", tableSessions)
	_, err := a.db.Exec(query, id)

	return err
}
