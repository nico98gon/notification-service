package infrastructure

import (
	"net/http"
	"nilus-challenge-backend/internal/domain/notification"
	handler "nilus-challenge-backend/internal/infrastructure/http"
)

func NewRouter(notificationService *notification.Service) {
	notificationHandler := handler.NewNotificationHandler(notificationService)

	http.HandleFunc("GET /notifications", notificationHandler.GetAllNotifications)
	http.HandleFunc("GET /notification", notificationHandler.GetNotificationByID)
	http.HandleFunc("POST /notify", notificationHandler.ScheduleNotification)
	http.HandleFunc("PUT /notification", notificationHandler.UpdateNotification)
	http.HandleFunc("DELETE /notification", notificationHandler.DeleteNotification)
	http.HandleFunc("GET /ws", notificationHandler.HandleWebSocketConnection)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
}