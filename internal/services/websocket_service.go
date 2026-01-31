package services

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketService interface {
	RegisterClient(conn *websocket.Conn)
	UnregisterClient(conn *websocket.Conn)
	BroadcastMessage(message string) error
}

type webSocketService struct {
	clients map[*websocket.Conn]bool
	lock    sync.Mutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Implement your origin checking logic here for production
		return true
	},
}

func NewWebSocketService() WebSocketService {
	return &webSocketService{
		clients: make(map[*websocket.Conn]bool),
		lock:    sync.Mutex{}}
}

func (ws *webSocketService) RegisterClient(conn *websocket.Conn) {
	ws.lock.Lock()
	defer ws.lock.Unlock()
	ws.clients[conn] = true
}

func (ws *webSocketService) UnregisterClient(conn *websocket.Conn) {
	ws.lock.Lock()
	defer ws.lock.Unlock()
	ws.unregisterClientInternal(conn)
}

func (ws *webSocketService) unregisterClientInternal(conn *websocket.Conn) {
	if _, exists := ws.clients[conn]; exists {
		delete(ws.clients, conn)
		conn.Close()
	}
}

func (ws *webSocketService) BroadcastMessage(message string) error {
	ws.lock.Lock()
	defer ws.lock.Unlock()
	for conn := range ws.clients {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			ws.unregisterClientInternal(conn)
		}
	}
	return nil
}
