package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Perseverance7/grady/internal/models"
	"github.com/gin-gonic/gin"
)

const ctxUserKey = "user"

func (h *Handler) register(c *gin.Context) {
	var userReq models.UserRegisterReq

	if err := c.BindJSON(&userReq); err != nil {
		newErrorResponce(c, http.StatusBadRequest, errors.New("invalid input body"))
		return
	}

	userRes, err := h.services.Authorization.CreateUser(userReq)
	if err != nil {
		newErrorResponce(c, http.StatusConflict, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"user": userRes,
	})
}

func (h *Handler) login(c *gin.Context) {
	var userReq models.UserLoginReq
	if err := c.BindJSON(&userReq); err != nil {
		newErrorResponce(c, http.StatusBadRequest, err)
		return
	}

	user, err := h.services.GetUser(userReq.Email, userReq.Password)
	if err != nil {
		newErrorResponce(c, http.StatusInternalServerError, err)
		return
	}

	accessToken, accessClaims, err := h.services.CreateToken(user.ID,
		user.Email,
		user.IsAdmin,
		24*time.Hour)
	if err != nil {
		newErrorResponce(c, http.StatusInternalServerError, err)
		return
	}

	refreshToken, refreshClaims, err := h.services.CreateToken(user.ID,
		user.Email,
		user.IsAdmin,
		7*24*time.Hour)
	if err != nil {
		newErrorResponce(c, http.StatusInternalServerError, err)
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
		newErrorResponce(c, http.StatusInternalServerError, err)
		return
	}

	c.SetCookie("refresh_token_uuid",
		session.ID,
		int(24*time.Hour.Seconds()),
		"/api/v1/auth",
		"",
		false,
		true)
	// secure поменять при переходе на https

	c.SetCookie("session_id",
		session.ID,
		int(24*time.Hour.Seconds()),
		"/",
		"",
		false,
		true)
	// secure поменять при переходе на https

	res := models.UserLoginRes{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessClaims.RegisteredClaims.ExpiresAt.Time,
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
		newErrorResponce(c, http.StatusBadRequest, err)
		return
	}

	if req.SessionID == "" {
		newErrorResponce(c, http.StatusBadRequest, errors.New("missing session ID"))
		return
	}

	err := h.services.DeleteSession(req.SessionID)
	if err != nil {
		newErrorResponce(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "successful logout"})
}

func (h *Handler) renewAccessToken(c *gin.Context) {
	var req models.RenewAccessTokenReq
	if err := c.BindJSON(&req); err != nil {
		newErrorResponce(c, http.StatusBadRequest, err)
		return
	}

	refreshClaims, err := h.services.VerifyRefreshToken(req.RefreshTokenUUID)
	if err != nil {
		newErrorResponce(c, http.StatusUnauthorized, errors.New("token verification failed"))
		return
	}

	session, err := h.services.GetSession(refreshClaims.RegisteredClaims.ID)
	if err != nil {
		newErrorResponce(c, http.StatusInternalServerError, err)
		return
	}

	if session.IsRevoked {
		newErrorResponce(c, http.StatusUnauthorized, errors.New("session is revoked"))
		return
	}

	if session.UserEmail != refreshClaims.Email {
		newErrorResponce(c, http.StatusUnauthorized, errors.New("inappropriate email address"))
		return
	}

	accessToken, accessClaims, err := h.services.CreateToken(refreshClaims.ID, refreshClaims.Email, refreshClaims.IsAdmin, 15*time.Minute)
	if err != nil {
		newErrorResponce(c, http.StatusInternalServerError, err)
		return
	}

	res := models.RenewAccessTokenRes{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessClaims.RegisteredClaims.ExpiresAt.Time,
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) revokeSession(c *gin.Context) {
	var req models.LogoutRequest
	if err := c.BindJSON(&req); err != nil {
		newErrorResponce(c, http.StatusBadRequest, err)
		return
	}

	if req.SessionID == "" {
		newErrorResponce(c, http.StatusBadRequest, errors.New("missing session ID"))
		return
	}

	err := h.services.RevokeSession(req.SessionID)
	if err != nil {
		newErrorResponce(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "succesfuly revoking session"})
}

func (h *Handler) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			newErrorResponce(c, http.StatusUnauthorized, errors.New("missing authorization header"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			newErrorResponce(c, http.StatusUnauthorized, errors.New("invalid authorization header format"))
			c.Abort()
			return
		}
		accessToken := parts[1]

		payload, err := h.services.VerifyAccessToken(accessToken)
		if err != nil {
			newErrorResponce(c, http.StatusUnauthorized, err)
			c.Abort()
			return
		}

		c.Set(ctxUserKey, &models.UserInfo{
			ID:      payload.ID,
			Email:   payload.Email,
			IsAdmin: payload.IsAdmin,
		})

		c.Next()
	}
}

func getUserInfo(c *gin.Context) (*models.UserInfo, error) {
	userInfo, exists := c.Get(ctxUserKey)
	if !exists {
		return nil, fmt.Errorf("user not found in context")
	}

	user, ok := userInfo.(*models.UserInfo)
	if !ok {
		return nil, fmt.Errorf("failed to cast user info")
	}

	return user, nil
}
