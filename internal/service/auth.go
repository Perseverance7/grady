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

var (
	accessSecret  = []byte(os.Getenv("ACCESS_SECRET"))
	refreshSecret = []byte(os.Getenv("REFRESH_SECRET"))

	AccessExpiry  = time.Minute * 15   // Access token живет 15 минут
	RefreshExpiry = time.Hour * 24 * 7 // Refresh token живет 7 дней
)

type TokenClaims struct {
	jwt.StandardClaims
	UserID int `json:"user_id"`
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

func (a *AuthService) GenerateTokens(email, password string) (string, string, error) {
	salt, err := a.repo.GetUserSalt(email)
	if err != nil {
		return "", "", errors.New("invalid login or password")
	}

	id, err := a.repo.GetUser(email, HashPassword(password, salt))
	if err != nil {
		return "", "", errors.New("invalid login or password")
	}

	// Создание access token с коротким сроком действия
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(AccessExpiry).Unix(), // Access token истекает быстро
			IssuedAt:  time.Now().Unix(),
		},
		UserID: id,
	})

	accessTokenString, err := accessToken.SignedString(accessSecret)
	if err != nil {
		return "", "", err
	}

	// Создание refresh token с более длинным сроком действия
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(RefreshExpiry).Unix(), // Refresh token живет дольше
			IssuedAt:  time.Now().Unix(),
		},
		UserID: id,
	})

	refreshTokenString, err := refreshToken.SignedString(refreshSecret)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
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
