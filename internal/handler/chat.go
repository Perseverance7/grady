package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Perseverance7/grady/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Допускаем все запросы, в реальном приложении лучше ограничить.
	},
}

func (h *Handler) webSocketHandler(c *gin.Context) {
	groupID := c.Query("group_id")
	if groupID == "" {
		log.Println("Group ID is required")
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		log.Println("Invalid limit, using default value 20")
		limit = 20
	}

	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		log.Println("Invalid offset, using default value 0")
		offset = 0
	}

	userInfo, exists := c.Get(ctxUserKey)
	if !exists {
		log.Println("User info not found in context")
		return
	}

	user, ok := userInfo.(*models.UserInfo)
	if !ok {
		log.Println("Invalid user info type")
		return
	}

	isInGroup, err := h.services.Chat.IsUserInGroup(user.ID, groupID)
	if err != nil {
		log.Println("Error checking group membership:", err)
		return
	}

	if !isInGroup {
		log.Println("User not authorized for this group")
		newErrorResponce(c, http.StatusForbidden, errors.New("forbidden"))
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	if h.connections[groupID] == nil {
		h.connections[groupID] = make(map[*websocket.Conn]bool)
	}
	h.connections[groupID][conn] = true

	defer func() {
		conn.Close()
		delete(h.connections[groupID], conn)
		if len(h.connections[groupID]) == 0 {
			delete(h.connections, groupID)
		}
	}()

	messages, err := h.services.Chat.GetChatHistory(groupID, limit, offset)
	if err != nil {
		log.Println("Failed to load chat history:", err)
		return
	}

	for _, msg := range messages {
		err := conn.WriteJSON(msg)
		if err != nil {
			log.Println("Failed to send message to client:", err)
			return
		}
	}

	for {
		_, messageContent, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var incomingMsg models.IncomingMessage
		err = json.Unmarshal(messageContent, &incomingMsg)
		if err != nil {
			log.Println("JSON Unmarshal error:", err)
			break
		}

		userData, err := h.services.GetUserData(user.ID)
		if err != nil {
			log.Println("Getting user data error: %w", err)
			break
		}

		message := models.Message{
			UserID:     user.ID,
			GroupID:    groupID,
			Name:       userData.Name,
			Surname:    userData.Surname,
			Patronymic: userData.Patronymic,
			Content:    incomingMsg.Content,
			SentAt:     time.Now(),
		}

		err = h.services.Chat.SendMessage(&message)
		if err != nil {
			log.Println("SendMessage error:", err)
			break
		}

		h.broadcastMessageToGroup(groupID, &message)
	}
}

func (h *Handler) broadcastMessageToGroup(groupID string, message *models.Message) {
	if h.connections[groupID] == nil {
		return
	}

	for conn := range h.connections[groupID] {
		err := conn.WriteJSON(message)
		if err != nil {
			log.Printf("Error sending message to group %s: %v\n", groupID, err)
			conn.Close()
			delete(h.connections[groupID], conn)
		}
	}
}
