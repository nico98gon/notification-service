package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"nilus-challenge-backend/internal/domain/notification"
	"nilus-challenge-backend/internal/infrastructure"
	"nilus-challenge-backend/internal/infrastructure/messaging"
	"nilus-challenge-backend/internal/infrastructure/repository"
	"time"
)

func main() {
	db := infrastructure.NewDBConnection()
	defer db.Close()

	notificationPostgresRepo := repository.NewPostgresNotificationRepository(db)
	notificationService := notification.NewService(notificationPostgresRepo)
	fmt.Println(notificationPostgresRepo)

	webSocketSender := messaging.NewWebSocketSender()
	infrastructure.NewRouter(notificationService)

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			err := notificationService.CheckAndSendNotifications(webSocketSender)
			if err != nil {
				log.Printf("Error al enviar notificaciones: %v", err)
			}
		}
	}()

	// consumer := messaging.NewConsumer("localhost:9092", "notifications", "notification-group")
	// err := consumer.StartProcessing(notificationService.ProcessEvent)
	// if err != nil {
	// 	log.Fatalf("Error al iniciar el consumidor: %v", err)
	// }

	http.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("http://weather-service:8083/")
		if err != nil {
			http.Error(w, "Error al conectar con weather-service: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(w, "Respuesta de weather-service: %s", body)
	})

	fmt.Println("Servicio de notificaciones escuchando en el puerto 8082...")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
