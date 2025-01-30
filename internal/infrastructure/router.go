package infrastructure

import (
	"net/http"
	"nilus-challenge-backend/internal/domain/notification"
	handler "nilus-challenge-backend/internal/infrastructure/http"
)

func NewRouter(notificationService *notification.Service) {
	notificationHandler := handler.NewNotificationHandler(notificationService)

	api := "/api/v1"

	http.HandleFunc("GET "+api+"/notification", notificationHandler.GetAllNotifications)
	http.HandleFunc("GET "+api+"/notification/{id}", notificationHandler.GetNotificationByID)
	http.HandleFunc("POST "+api+"/notification", notificationHandler.ScheduleNotification)
	http.HandleFunc("PUT "+api+"/notification/{id}", notificationHandler.UpdateNotification)
	http.HandleFunc("DELETE "+api+"/notification/{id}", notificationHandler.DeleteNotification)
	http.HandleFunc("GET "+api+"/ws", notificationHandler.HandleWebSocketConnection)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
}