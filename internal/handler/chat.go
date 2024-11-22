package handler

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/Perseverance7/grady/internal/models"
	"github.com/Perseverance7/grady/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatHandler struct {
	services  *service.Service
	clients   map[*websocket.Conn]chan models.Message // Каналы для каждого клиента
	broadcast chan models.Message                     // Канал для общей рассылки
	mu        sync.Mutex                              // Мьютекс для защиты clients
}

func NewChatHandler(services *service.Service) *ChatHandler {
	handler := &ChatHandler{
		services:  services,
		clients:   make(map[*websocket.Conn]chan models.Message),
		broadcast: make(chan models.Message),
	}

	// Запускаем прослушивание канала broadcast в отдельной горутине
	go handler.listenForMessages()

	return handler
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (ch *ChatHandler) WebSocketEndpoint(c *gin.Context) {
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

	history, err := ch.services.GetChatHistory(groupID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Устанавливаем WebSocket-соединение
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer func() {
		ch.unregisterClient(conn)
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
	ch.registerClient(conn, clientChan)

	// Горутинa для отправки сообщений клиенту
	go ch.sendMessages(conn, clientChan)

	// Читаем сообщения от клиента
	ch.handleClientMessages(conn, groupID)
}

func (ch *ChatHandler) registerClient(conn *websocket.Conn, clientChan chan models.Message) {
	ch.mu.Lock()
	ch.clients[conn] = clientChan
	ch.mu.Unlock()
}

func (ch *ChatHandler) unregisterClient(conn *websocket.Conn) {
	ch.mu.Lock()
	if clientChan, ok := ch.clients[conn]; ok {
		close(clientChan)
		delete(ch.clients, conn)
	}
	ch.mu.Unlock()
}

func (ch *ChatHandler) handleClientMessages(conn *websocket.Conn, groupID string) {
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
		if err := ch.services.SendMessage(msg); err != nil {
			log.Printf("Failed to save message: %v", err)
		} else {
			ch.broadcast <- msg // Добавляем сообщение в общую рассылку
		}
	}
}

func (ch *ChatHandler) sendMessages(conn *websocket.Conn, clientChan chan models.Message) {
	for msg := range clientChan {
		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("Error sending to client: %v", err)
			break
		}
	}
	// Закрытие соединения после завершения работы
	ch.unregisterClient(conn)
}

func (ch *ChatHandler) listenForMessages() {
	for msg := range ch.broadcast {
		ch.mu.Lock()
		for client, clientChan := range ch.clients {
			select {
			case clientChan <- msg:
			default:
				// Если клиент не читает, закрываем соединение
				log.Printf("Client unresponsive, closing connection")
				close(clientChan)
				delete(ch.clients, client)
				client.Close()
			}
		}
		ch.mu.Unlock()
	}
}
