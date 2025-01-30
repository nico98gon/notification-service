package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"nilus-challenge-backend/internal/domain/notification"
	"nilus-challenge-backend/internal/infrastructure/messaging"
	"nilus-challenge-backend/internal/infrastructure/middleware"
)

type NotificationHandler struct {
	notificationService *notification.Service
}

func NewNotificationHandler(service *notification.Service) *NotificationHandler {
	return &NotificationHandler{notificationService: service}
}

func (h *NotificationHandler) GetAllNotifications(w http.ResponseWriter, r *http.Request) {
	notifications, err := h.notificationService.GetAllNotifications()
	if err != nil {
		middleware.ErrorResponse(w, http.StatusInternalServerError, "Error al obtener las notificaciones: "+err.Error())
		return
	}

	middleware.SuccessResponse(w, notifications)
}

func (h *NotificationHandler) GetNotificationByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
			middleware.ErrorResponse(w, http.StatusBadRequest, "El parámetro 'id' es obligatorio")
			return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, "ID no válido")
			return
	}

	notification, err := h.notificationService.GetNotificationByID(id)
	if err != nil {
			middleware.ErrorResponse(w, http.StatusInternalServerError, "Error al obtener la notificación: "+err.Error())
			return
	}

	if notification == nil {
			middleware.ErrorResponse(w, http.StatusNotFound, "Notificación no encontrada")
			return
	}

	middleware.SuccessResponse(w, notification)
}


func (h *NotificationHandler) ScheduleNotification(w http.ResponseWriter, r *http.Request) {
	var n notification.Notification
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		middleware.ErrorResponse(w, http.StatusBadRequest, "Error al procesar la solicitud: "+err.Error())
		return
	}

	if err := h.notificationService.ScheduleNotification(n.UserID, &n); err != nil {
		middleware.ErrorResponse(w, http.StatusInternalServerError, "Error al programar la notificación: "+err.Error())
		return
	}

	middleware.SuccessResponse(w, n)
}

func (h *NotificationHandler) UpdateNotification(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		middleware.ErrorResponse(w, http.StatusBadRequest, "El parámetro 'id' es obligatorio")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.ErrorResponse(w, http.StatusBadRequest, "ID no válido")
		return
	}

	var n notification.Notification
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		middleware.ErrorResponse(w, http.StatusBadRequest, "Error al procesar la solicitud: "+err.Error())
		return
	}

	if err := h.notificationService.UpdateNotification(id, &n); err != nil {
		middleware.ErrorResponse(w, http.StatusInternalServerError, "Error al actualizar la notificación: "+err.Error())
		return
	}

	middleware.SuccessResponse(w, n)
}

func (h *NotificationHandler) DeleteNotification(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		middleware.ErrorResponse(w, http.StatusBadRequest, "El parámetro 'id' es obligatorio")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.ErrorResponse(w, http.StatusBadRequest, "ID no válido")
		return
	}

	if err := h.notificationService.DeleteNotification(id); err != nil {
		middleware.ErrorResponse(w, http.StatusInternalServerError, "Error al eliminar la notificación: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *NotificationHandler) HandleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := messaging.UpgradeConnection(w, r)
	if err != nil {
		middleware.ErrorResponse(w, http.StatusInternalServerError, "Error al actualizar a WebSocket: "+err.Error())
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
