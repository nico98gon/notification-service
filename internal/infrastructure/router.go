package infrastructure

import (
	"net/http"
	"nilus-challenge-backend/internal/domain/notification"
	handler "nilus-challenge-backend/internal/infrastructure/http"
)

func NewRouter(notificationService *notification.Service) {
	notificationHandler := handler.NewNotificationHandler(notificationService)

	http.HandleFunc("POST /notify", notificationHandler.ScheduleNotification)
	http.HandleFunc("GET /ws", notificationHandler.HandleWebSocketConnection)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
}