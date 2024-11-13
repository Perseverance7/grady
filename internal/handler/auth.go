package handler

import (
	"net/http"

	"github.com/Perseverance7/grady/internal/models"
	"github.com/gin-gonic/gin"
)

func(h *Handler) register(c *gin.Context){
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

func(h *Handler) login(c *gin.Context){
	var input models.UserLogin
	if err := c.BindJSON(&input); err != nil {
		newErrorResponce(c, http.StatusBadRequest, "invalid input body")
		return
	}

	token, err := h.services.GenerateToken(input.Email, input.Password)
	if err != nil {
		newErrorResponce(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

func(h *Handler) logout(c *gin.Context){
	
}