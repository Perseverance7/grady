package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Perseverance7/grady/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) WebSocketEndpoint(c *gin.Context) {
	groupID := c.Query("group_id")
	if groupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no group id"})
		return
	}

	// Получаем историю чата
	limitParam := c.DefaultQuery("history_limit", "20")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid history_limit"})
		return
	}

	history, err := h.services.GetChatHistory(groupID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Проверка и извлечение токена из заголовков
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
		return
	}

	// Разделяем заголовок на части
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
		return
	}
	accessToken := parts[1]

	// Проверка токена
	payload, err := h.services.VerifyAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("invalid token: %v", err)})
		return
	}

	// Извлекаем user_id из payload
	userID := payload.ID

	// Устанавливаем WebSocket-соединение
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer func() {
		h.unregisterClient(conn)
		if err := conn.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}()

	// Отправляем клиенту историю чата
	if err := conn.WriteJSON(history); err != nil {
		log.Printf("Failed to send chat history: %v", err)
		return
	}

	// Регистрируем клиента
	clientChan := make(chan models.Message, 100) // Увеличен размер буфера
	h.registerClient(conn, clientChan)

	// Горутинa для отправки сообщений клиенту
	go h.sendMessages(conn, clientChan)

	// Читаем сообщения от клиента
	h.handleClientMessages(conn, groupID, userID)
}

func (h *Handler) registerClient(conn *websocket.Conn, clientChan chan models.Message) {
	h.mu.Lock()
	h.clients[conn] = clientChan
	h.mu.Unlock()
}

func (h *Handler) unregisterClient(conn *websocket.Conn) {
	h.mu.Lock()
	if clientChan, ok := h.clients[conn]; ok {
		close(clientChan)
		delete(h.clients, conn)
	}
	h.mu.Unlock()
}

func (h *Handler) handleClientMessages(conn *websocket.Conn, groupID string, userID int64) {
	for {
		var msg models.Message
		if err := conn.ReadJSON(&msg); err != nil {
			// Если клиент отключился или произошла ошибка
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Println("Client disconnected")
				break
			}
			log.Printf("Error reading from client: %v", err)
			break
		}

		// Устанавливаем группу и отправляем сообщение в broadcast
		msg.GroupID = groupID
		msg.UserID = userID

		if err := h.services.SendMessage(msg); err != nil {
			log.Printf("Failed to save message: %v", err)
		} else {
			h.broadcast <- msg // Добавляем сообщение в общую рассылку
		}
	}
}

func (h *Handler) sendMessages(conn *websocket.Conn, clientChan chan models.Message) {
	for msg := range clientChan {
		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("Error sending to client: %v", err)
			break
		}
	}
	// Закрытие соединения после завершения работы
	h.unregisterClient(conn)
}

func (h *Handler) listenForMessages() {
	for msg := range h.broadcast {
		h.mu.Lock()
		for client, clientChan := range h.clients {
			select {
			case clientChan <- msg:
			default:
				// Если клиент не читает, закрываем соединение
				log.Printf("Client unresponsive, closing connection")
				close(clientChan)
				delete(h.clients, client)
				client.Close()
			}
		}
		h.mu.Unlock()
	}
}
