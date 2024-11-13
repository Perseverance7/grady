package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/Perseverance7/grady/internal/models"
	"github.com/Perseverance7/grady/internal/repository"
)

const (
	tokenTTL = 12 * time.Hour
)

var (
	signingKey = os.Getenv("SIGNING_KEY")
)

type TokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (a *AuthService) CreateUser(user models.UserRegister) (int, error) {
	var err error

	user.Salt, err = GenerateSalt()
	if err != nil {
		return 0, err
	}

	user.Password = HashPassword(user.Password, user.Salt)
	id, err := a.repo.CreateUser(user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"users_username_key\"") {
			return 0, errors.New("user already exists")
		}
		return 0, errors.New(err.Error())
	}
	return id, nil
}

func (a *AuthService) GenerateToken(email, password string) (string, error) {
	salt, err := a.repo.GetUserSalt(email)

	if err != nil {
		return "", errors.New("invalid login or password")
	} else {
		id, err := a.repo.GetUser(email, HashPassword(password, salt))
		if err != nil {
			return "", errors.New("invalid login or password")
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(tokenTTL).Unix(),
				IssuedAt:  time.Now().Unix(),
			},
			id,
		})

		return token.SignedString([]byte(signingKey))
	}

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
