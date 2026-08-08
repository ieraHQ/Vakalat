package websocket

import (
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/ieraHQ/Vakalat/backend/api/logger"
	"go.uber.org/zap"
)

// Client represents a WebSocket client.
type Client struct {
	ID   string
	Conn *websocket.Conn
}

// Hub manages WebSocket clients and broadcasts messages.
type Hub struct {
	clients    map[string]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

// NewHub creates a new Hub.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the Hub's event loop.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.ID] = client
			h.mu.Unlock()
			logger.GetLogger().Info("Client registered", zap.String("id", client.ID))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				client.Conn.Close()
			}
			h.mu.Unlock()
			logger.GetLogger().Info("Client unregistered", zap.String("id", client.ID))

		case message := <-h.broadcast:
			h.mu.Lock()
			for id, client := range h.clients {
				if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
					logger.GetLogger().Error("Failed to write message", zap.String("id", id), zap.Error(err))
					delete(h.clients, id)
					client.Conn.Close()
				}
			}
			h.mu.Unlock()
		}
	}
}

// WebSocketHandler handles WebSocket connections.
func WebSocketHandler(hub *Hub) fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		client := &Client{
			ID:   c.LocalAddr().String(),
			Conn: c,
		}

		hub.register <- client
		defer func() { hub.unregister <- client }()

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				logger.GetLogger().Error("Failed to read message", zap.Error(err))
				break
			}
			logger.GetLogger().Info("Received message", zap.String("message", string(msg)))
			hub.broadcast <- msg
		}
	})
}