package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strings"


	"github.com/Perseverance7/grady/internal/models"
	"github.com/Perseverance7/grady/internal/repository"
)


type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (a *AuthService) CreateUser(userReq models.UserRegisterReq) (models.UserRegisterRes, error) {
	var err error

	userReq.Salt, err = GenerateSalt()
	if err != nil {
		return models.UserRegisterRes{}, err
	}

	userReq.Password = HashPassword(userReq.Password, userReq.Salt)
	userRes, err := a.repo.CreateUser(userReq)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"users_username_key\"") {
			return models.UserRegisterRes{}, errors.New("user already exists")
		}
		return models.UserRegisterRes{}, errors.New(err.Error())
	}
	return userRes, nil
}

func (a *AuthService) GetUser(email, password string) (models.UserLogin, error) {
	salt, err := a.repo.GetUserSalt(email)
	if err != nil{
		return models.UserLogin{}, err
	}

	user, err := a.repo.GetUser(email, HashPassword(password, salt))
	if err != nil{
		return models.UserLogin{}, err
	}

	return user, nil
}

func (a *AuthService) UpdateUser(user *models.User) (*models.User, error) {
	return a.repo.UpdateUser(user)
}

func (a *AuthService) DeleteUser(id int64) error{
	return a.repo.DeleteUser(id)
}

func (a *AuthService) CreateSession(session *models.Session) (*models.Session, error) {
	return a.repo.CreateSession(session)
}

func (a *AuthService) GetSession(id string) (*models.Session, error) {
	return a.repo.GetSession(id)
}

func (a *AuthService) RevokeSession(id string) error {
	return a.repo.RevokeSession(id)
}

func (a *AuthService) DeleteSession(id string) error {
	return a.repo.DeleteSession(id)
}

func GenerateSalt() (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

// HashPassword - функция для хеширования пароля с использованием соли
func HashPassword(password, salt string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password + salt))
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}
