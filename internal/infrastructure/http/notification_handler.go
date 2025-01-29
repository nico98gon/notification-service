package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"nilus-challenge-backend/internal/domain/notification"
	"nilus-challenge-backend/internal/infrastructure/messaging"
)

type NotificationHandler struct {
	notificationService *notification.Service
}

func (h *NotificationHandler) GetAllNotifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	notifications, err := h.notificationService.GetAllNotifications()
	if err != nil {
		http.Error(w, "Error al obtener las notificaciones: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(notifications); err != nil {
		http.Error(w, "Error al serializar las notificaciones: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *NotificationHandler) GetNotificationByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "El parámetro 'id' es obligatorio", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID no válido", http.StatusBadRequest)
		return
	}

	notification, err := h.notificationService.GetNotificationByID(id)
	if err != nil {
		http.Error(w, "Error al obtener la notificación: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if notification == nil {
		http.Error(w, "Notificación no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(notification); err != nil {
		http.Error(w, "Error al serializar la notificación: "+err.Error(), http.StatusInternalServerError)
	}
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

func (h *NotificationHandler) UpdateNotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "El parámetro 'id' es obligatorio", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID no válido", http.StatusBadRequest)
		return
	}

	var n notification.Notification
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		http.Error(w, "Error al procesar el cuerpo de la solicitud: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.notificationService.UpdateNotification(id, &n); err != nil {
		http.Error(w, "Error al actualizar la notificación: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(n); err != nil {
		http.Error(w, "Error al serializar la notificación: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *NotificationHandler) DeleteNotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "El parámetro 'id' es obligatorio", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID no válido", http.StatusBadRequest)
		return
	}

	if err := h.notificationService.DeleteNotification(id); err != nil {
		http.Error(w, "Error al eliminar la notificación: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
