package messaging

import (
	"fmt"
	"log"
	"net/http"
	"nilus-challenge-backend/internal/domain/notification"

	"github.com/gorilla/websocket"
)

type WebSocketSender struct {
	connections map[int]*websocket.Conn
}

func NewWebSocketSender() *WebSocketSender {
	return &WebSocketSender{connections: make(map[int]*websocket.Conn)}
}

func (s *WebSocketSender) RegisterConnection(userID int, conn *websocket.Conn) {
	s.connections[userID] = conn
}

func (s *WebSocketSender) Send(n notification.Notification) error {
	conn, ok := s.connections[n.UserID]
	if !ok {
		return fmt.Errorf("No hay conexión activa para el usuario %d", n.UserID)
	}

	return conn.WriteJSON(n)
}

func (s *WebSocketSender) UnregisterConnection(userID int) {
	delete(s.connections, userID)
	log.Printf("Conexión eliminada para el usuario %d", userID)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func UpgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	return upgrader.Upgrade(w, r, nil)
}

