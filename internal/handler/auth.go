package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Perseverance7/grady/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) register(c *gin.Context) {
	var userReq models.UserRegisterReq

	if err := c.BindJSON(&userReq); err != nil {
		newErrorResponce(c, http.StatusBadRequest, "invalid input body")
		return
	}

	userRes, err := h.services.Authorization.CreateUser(userReq)
	if err != nil {
		newErrorResponce(c, http.StatusConflict, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"user": userRes,
	})
}

func (h *Handler) login(c *gin.Context) {
	var userReq models.UserLoginReq
	if err := c.BindJSON(&userReq); err != nil {
		newErrorResponce(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.services.GetUser(userReq.Email, userReq.Password)
	if err != nil {
		newErrorResponce(c, http.StatusInternalServerError, err.Error())
		return
	}

	accessToken, accessClaims, err := h.tokenMaker.CreateToken(user.ID, user.Email, user.IsAdmin, 15*time.Minute)
	if err != nil {
		newErrorResponce(c, http.StatusInternalServerError, err.Error())
		return
	}

	refreshToken, refreshClaims, err := h.tokenMaker.CreateToken(user.ID, user.Email, user.IsAdmin, 24*time.Hour)
	if err != nil {
		newErrorResponce(c, http.StatusInternalServerError, err.Error())
		return
	}

	session, err := h.services.CreateSession(&models.Session{
		ID:           refreshClaims.RegisteredClaims.ID,
		UserEmail:    user.Email,
		RefreshToken: refreshToken,
		IsRevoked:    false,
		ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
	})

	if err != nil {
		newErrorResponce(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.SetCookie("refresh_token", refreshToken, int(24*time.Hour.Seconds()), "/api/v1/auth", "", false, true)

	res := models.UserLoginRes{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessClaims.RegisteredClaims.ExpiresAt.Time,
		User: models.UserLogin{
			Name:       user.Name,
			Surname:    user.Surname,
			Patronymic: user.Patronymic,
			Email:      user.Email,
			IsAdmin:    user.IsAdmin,
		},
	}

	c.JSON(http.StatusOK, res)

}

func (h *Handler) logout(c *gin.Context) {
	var req models.LogoutRequest
	if err := c.BindJSON(&req); err != nil {
		newErrorResponce(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.SessionID == "" {
		newErrorResponce(c, http.StatusBadRequest, "missing session ID")
		return
	}

	err := h.services.DeleteSession(req.SessionID)
	if err != nil {
		newErrorResponce(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "successful logout"})
}

func (h *Handler) renewAccessToken(c *gin.Context) {
	var req models.RenewAccessTokenReq
	if err := c.BindJSON(&req); err != nil {
		newErrorResponce(c, http.StatusBadRequest, err.Error())
		return
	}

	refreshClaims, err := h.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		newErrorResponce(c, http.StatusUnauthorized, "token verification failed")
		return
	}

	session, err := h.services.GetSession(refreshClaims.RegisteredClaims.ID)
	if err != nil {
		newErrorResponce(c, http.StatusInternalServerError, err.Error())
		return
	}

	if session.IsRevoked {
		newErrorResponce(c, http.StatusUnauthorized, "session is revoked")
		return
	}

	if session.UserEmail != refreshClaims.Email {
		newErrorResponce(c, http.StatusUnauthorized, "inappropriate email address")
		return
	}

	accessToken, accessClaims, err := h.tokenMaker.CreateToken(refreshClaims.ID, refreshClaims.Email, refreshClaims.IsAdmin, 15*time.Minute)
	if err != nil{
		newErrorResponce(c, http.StatusInternalServerError, err.Error())
		return 
	}

	res := models.RenewAccessTokenRes{
		AccessToken: accessToken,
		AccessTokenExpiresAt: accessClaims.RegisteredClaims.ExpiresAt.Time,
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) revokeSession(c *gin.Context) {
	var req models.LogoutRequest
	if err := c.BindJSON(&req); err != nil {
		newErrorResponce(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.SessionID == "" {
		newErrorResponce(c, http.StatusBadRequest, "missing session ID")
		return
	}

	err := h.services.RevokeSession(req.SessionID)
	if err != nil {
		newErrorResponce(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "succesfuly revoking session"})
}

func (h *Handler) authMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            newErrorResponce(c, http.StatusUnauthorized, "missing authorization header")
            c.Abort()
            return
        }

        // Разделяем заголовок на части
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            newErrorResponce(c, http.StatusUnauthorized, "invalid authorization header format")
            c.Abort()
            return
        }
        accessToken := parts[1]

        // Проверяем валидность токена
        payload, err := h.tokenMaker.VerifyToken(accessToken)
        if err != nil {
            newErrorResponce(c, http.StatusUnauthorized, fmt.Sprintf("invalid token: %v", err))
            c.Abort()
            return
        }

        // Сохраняем информацию о пользователе в контексте
        c.Set("user_id", payload.ID)
        c.Set("email", payload.Email)
        c.Set("is_admin", payload.IsAdmin)

        // Передаем управление следующему обработчику
        c.Next()
    }
}
