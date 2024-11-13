package models

type UserRegister struct {
	Id         int
	Name       string `json:"name" binding:"required"`
	Surname    string `json:"surname" binding:"required"`
	Patronymic string `json:"patronymic"`
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Salt       string
}

type UserLogin struct {
	Id       int
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
