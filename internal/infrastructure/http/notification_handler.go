package http

import (
	"encoding/json"
	"log"
	"net/http"

	"nilus-challenge-backend/internal/domain/notification"
	"nilus-challenge-backend/internal/infrastructure/messaging"
)

type NotificationHandler struct {
	notificationService *notification.Service
}

func NewNotificationHandler(service *notification.Service) *NotificationHandler {
	return &NotificationHandler{notificationService: service}
}

func (h *NotificationHandler) ScheduleNotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var n notification.Notification
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		http.Error(w, "Error al procesar el cuerpo de la solicitud: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.notificationService.ScheduleNotification(n.UserID, &n); err != nil {
		http.Error(w, "Error al programar la notificación: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(n)
}

func (h *NotificationHandler) HandleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := messaging.UpgradeConnection(w, r)
	if err != nil {
		http.Error(w, "Error al actualizar a WebSocket: "+err.Error(), http.StatusInternalServerError)
		return
	}

	defer conn.Close()
	log.Println("Nueva conexión WebSocket establecida")

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Conexión WebSocket cerrada")
			break
		}
	}
}
