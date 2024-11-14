package handler

import (
	"net/http"

	"github.com/Perseverance7/grady/internal/models"
	"github.com/Perseverance7/grady/internal/service"
	"github.com/gin-gonic/gin"
)

func (h *Handler) register(c *gin.Context) {
	var input models.UserRegister

	if err := c.BindJSON(&input); err != nil {
		newErrorResponce(c, http.StatusBadRequest, "invalid input body")
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponce(c, http.StatusConflict, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) login(c *gin.Context) {
	var input models.UserLogin

	// Парсинг JSON-запроса с данными пользователя
	if err := c.BindJSON(&input); err != nil {
		newErrorResponce(c, http.StatusBadRequest, "invalid input body")
		return
	}

	// Генерация access и refresh токенов
	accessToken, refreshToken, err := h.services.GenerateTokens(input.Email, input.Password)
	if err != nil {
		newErrorResponce(c, http.StatusUnauthorized, "invalid email or password")
		return
	}

	// Опционально: установка токенов в HttpOnly cookies
	c.SetCookie("access_token", accessToken, int(service.AccessExpiry.Seconds()), "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, int(service.RefreshExpiry.Seconds()), "/", "", false, true)

	// Возвращаем токены в JSON-ответе
	c.JSON(http.StatusOK, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}


func (h *Handler) logout(c *gin.Context) {

}

