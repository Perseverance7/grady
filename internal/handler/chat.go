package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/Perseverance7/grady/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatHandler struct {
	connections map[string]map[*websocket.Conn]bool
}

func NewChatHandler() *ChatHandler {
	return &ChatHandler{
		connections: make(map[string]map[*websocket.Conn]bool),
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Допускаем все запросы, в реальном приложении лучше ограничить.
	},
}

func (h *Handler) webSocketHandler(c *gin.Context) {

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	groupID := c.Query("group_id") // Получаем ID группы из маршрута
	if groupID == "" {
		log.Println("Group ID is required")
		return
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

	ch := NewChatHandler()

	defer func() {
		conn.Close()
		delete(ch.connections[groupID], conn)
		if len(ch.connections[groupID]) == 0 {
			delete(ch.connections, groupID)
		}
	}()

	if ch.connections[groupID] == nil {
		ch.connections[groupID] = make(map[*websocket.Conn]bool)
	}
	ch.connections[groupID][conn] = true

	// Загружаем историю сообщений для группы
	messages, err := h.services.Chat.GetChatHistory(groupID)
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
		// Чтение сообщения от клиента
		_, messageContent, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		
		userData, err := h.services.GetUserData(user.ID)
		if err != nil{
			log.Println("Getting user data error: %w", err)
			break
		}

		message := models.Message{
			UserID:     user.ID,
			GroupID:    groupID,
			Name:       userData.Name,
			Surname:    userData.Surname,
			Patronymic: userData.Patronymic,
			Content:    string(messageContent),
			SentAt:     time.Now(),
		}

		// Сохраняем сообщение и отправляем остальным
		err = h.services.Chat.SendMessage(&message)
		if err != nil {
			log.Println("SendMessage error:", err)
			break
		}

		// Здесь можно добавить логику для рассылки сообщения другим клиентам через вебсокет-соединения.
		h.broadcastMessageToGroup(ch, groupID, &message)
	}
}

func (h *Handler) broadcastMessageToGroup(ch *ChatHandler, groupID string, message *models.Message) {
	if ch.connections[groupID] == nil {
		return
	}

	for conn := range ch.connections[groupID] {
		err := conn.WriteJSON(message)
		if err != nil {
			log.Printf("Error sending message to group %s: %v\n", groupID, err)
			conn.Close()
			delete(ch.connections[groupID], conn)
		}
	}
}
