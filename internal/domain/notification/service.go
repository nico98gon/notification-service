package notification

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"nilus-challenge-backend/internal/infrastructure/services"
	"strconv"
	"time"
)

type NotificationSender interface {
	Send(n Notification) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAllNotifications() ([]Notification, error) {
	return s.repo.FindAllNotifications()
}

func (s *Service) GetNotificationByID(id int) (*Notification, error) {
	return s.repo.FindNotificationByID(id)
}

func (s *Service) ScheduleNotification(userID int, notification *Notification) error {
	return s.repo.SaveNotification(notification)
}

func (s *Service) ProcessEvent(event NotificationEvent) error {
	notification := &Notification{
		UserID:    event.UserID,
		Title:     &event.Title,
		Content:   &event.Content,
		Scheduled: event.Scheduled,
	}

	if err := notification.Validate(); err != nil {
		return err
	}

	if err := s.repo.SaveNotification(notification); err != nil {
		return err
	}

	log.Printf("Notificación programada para el usuario %d", event.UserID)
	return nil
}

func (s *Service) UpdateNotification(id int, notification *Notification) error {
	if err := notification.Validate(); err != nil {
		return err
	}

	return s.repo.UpdateNotification(id, notification)
}

func (s *Service) DeleteNotification(id int) error {
	return s.repo.DeleteNotification(id)
}

func (s *Service) CheckAndSendNotifications(sender NotificationSender) error {
	now := time.Now()

	err := s.repo.MarkNotificationsAsPending(now)
	if err != nil {
		return fmt.Errorf("error al actualizar notificaciones pendientes: %v", err)
	}

	notifications, err := s.repo.FindPendingNotifications(now) 
	if err != nil {
		return err
	}

	for _, n := range notifications {
		log.Printf("Preparando notificación para el usuario %d\n", n.UserID)

		cityID, err := strconv.Atoi(n.LocalityID)
		if err != nil {
			log.Printf("Error al convertir cityID a int para el usuario %d: %v", n.UserID, err)
			continue
		}

		isCoastal := false
		waveURL := fmt.Sprintf("http://weather-service:8083/wave-forecast?city_id=%s&day=0", cityID)
		resp, err := http.Get(waveURL)
		if err == nil {
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err == nil {
				var waveData struct {
					Status string `json:"status"`
				}
				err = json.Unmarshal(body, &waveData)
				if err == nil && waveData.Status == "success" {
					isCoastal = true
				}
			}
		}

		title, content, err := services.GetWeatherForecast(cityID, isCoastal)
		if err != nil {
			log.Printf("Error al obtener el clima para el usuario %d: %v", n.UserID, err)
			continue 
		}

		n.Title = &title
		n.Content = &content

		err = sender.Send(n)
		if err != nil {
			log.Printf("Error al enviar notificación (Usuario %d): %v", n.UserID, err)
			continue
		}

		err = s.repo.SendAndMarkAsSent(n.Id) 
		if err != nil {
			log.Printf("Error al actualizar estado de la notificación ID %d: %v", n.Id, err)
		}
	}

	return nil
}
