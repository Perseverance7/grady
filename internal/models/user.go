package models

import "time"

type User struct {
	ID         int64
	Name       string `json:"name" binding:"required" db:"name"`
	Surname    string `json:"surname" binding:"required" db:"surname"`
	Patronymic string `json:"patronymic" db:"surname"`
	Email      string `json:"email" binding:"required" db:"surname"`
	IsAdmin    bool   `json:"is_admin" binding:"required" db:"is_admin"`
	Password   string `json:"password" binding:"required"`
	Salt       string
}

type UserRegisterReq struct {
	Name       string `json:"name" binding:"required"`
	Surname    string `json:"surname" binding:"required"`
	Patronymic string `json:"patronymic"`
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Salt       string
}

type UserRegisterRes struct {
	Name       string `json:"name" binding:"required" db:"name"`
	Surname    string `json:"surname" binding:"required" db:"surname"`
	Patronymic string `json:"patronymic" db:"surname"`
	Email      string `json:"email" binding:"required" db:"surname"`
	IsAdmin    bool   `json:"is_admin" binding:"required" db:"is_admin"`
}

type UserLogin struct {
	ID         int64  `json:"id" db:"id"`
	Name       string `json:"name" binding:"required" db:"name"`
	Surname    string `json:"surname" binding:"required" db:"surname"`
	Patronymic string `json:"patronymic" db:"surname"`
	Email      string `json:"email" binding:"required" db:"surname"`
	IsAdmin    bool   `json:"is_admin" binding:"required" db:"is_admin"`
}

type UserLoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginRes struct {
	SessionID             string    `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	User                  UserLogin `json:"user"`
}

type RenewAccessTokenReq struct {
	RefreshToken string `json:"refresh_token"`
}

type RenewAccessTokenRes struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

