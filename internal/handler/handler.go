package handler

import (
	"github.com/Perseverance7/grady/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")

	// Authentication
	api.POST("/register", h.register)
	api.POST("/login", h.login)
	api.POST("/logout", h.logout)

	// User profile
	user := api.Group("/users/me")
	{
		user.GET("/", h.getUserProfile)
		user.PUT("/", h.updateUserProfile)
	}

	// Group routes (for authenticated users)
	group := api.Group("/groups")
	{
		group.POST("/", h.createGroup)
		group.GET("/", h.listGroups)
		group.GET("/:group_id", h.getGroupDetails)
		group.POST("/:group_id/join_member", h.joinMember)
		group.DELETE("/:group_id/remove_member", h.removeMember)

		// Test management within groups
		tasks := group.Group("/:group_id/tasks")
		{
			tasks.POST("/", h.createTask)
			tasks.GET("/", h.listTasks)
			tasks.GET("/:task_id", h.getTask)
			tasks.POST("/:task_id/submit", h.submitTask)
			tasks.GET("/:task_id/results", h.getTaskResults)
		}

		// Group chat and messaging
		chat := group.Group("/:group_id/chat")
		{
			chat.GET("/messages", h.getMessages)
			chat.POST("/messages", h.sendMessage)
		}

		// Group statistics
		stats := group.Group("/:group_id/stats")
		{
			stats.GET("/", h.getGroupStats)
		}
	}

	// Notifications
	notifications := api.Group("/notifications")
	{
		notifications.GET("/", h.listNotifications)
		notifications.POST("/", h.sendNotification) // Changed to POST "/" for consistency
		notifications.POST("/read", h.markNotificationsAsRead)
	}

	return r
}
