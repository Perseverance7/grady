package handler

import (
	"github.com/Perseverance7/grady/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	services    *service.Service
	connections map[string]map[*websocket.Conn]bool
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services:    services,
		connections: make(map[string]map[*websocket.Conn]bool),
	}

}

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/register", h.register)
		auth.POST("/login", h.login)
		auth.POST("/logout", h.logout)
		auth.POST("/token/renew", h.renewAccessToken)
		auth.POST("/session/revoke", h.revokeSession)
	}

	api.Use(h.authMiddleware())

	user := api.Group("/users/me")
	{
		user.GET("/", h.getUserProfile)
		user.PUT("/", h.updateUserProfile)
	}

	group := api.Group("/groups")
	{
		group.POST("/", h.createGroup)
		group.GET("/", h.listGroups)
		group.GET("/:group_id", h.getGroupDetails)
		group.POST("/:group_id/add_member", h.joinMember)
		group.DELETE("/:group_id/remove_member", h.removeMember)

		tasks := group.Group("/:group_id/tasks")
		{
			tasks.POST("/", h.createTask)
			tasks.GET("/", h.listTasks)
			tasks.GET("/:task_id", h.getTask)
			tasks.POST("/:task_id/submit", h.submitTask)
			tasks.GET("/:task_id/results", h.getTaskResults)
		}

		chat := group.Group("/chat")
		{
			chat.GET("/", h.webSocketHandler)
		}

		stats := group.Group("/:group_id/stats")
		{
			stats.GET("/", h.getGroupStats)
		}
	}

	notifications := api.Group("/notifications")
	{
		notifications.GET("/", h.listNotifications)
		notifications.POST("/", h.sendNotification)
		notifications.POST("/read", h.markNotificationsAsRead)
	}

	return r
}
