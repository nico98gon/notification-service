package notification

import (
	"log"
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

func (s *Service) ScheduleNotification(userID int, notification *Notification) error {
	return s.repo.SaveNotification(notification)
}

func (s *Service) ProcessEvent(event NotificationEvent) error {
	notification := &Notification{
		UserID:    event.UserID,
		Title:     event.Title,
		Content:   event.Content,
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

func (s *Service) CheckAndSendNotifications(sender NotificationSender) error {
	now := time.Now()
	notifications, err := s.repo.FindScheduledNotifications(now)
	if err != nil {
		return err
	}

	for _, n := range notifications {
		log.Printf("Enviando notificación: %+v\n", n)

		// Intenta enviar la notificación por WebSocket
		err := sender.Send(n)
		if err != nil {
			log.Printf("Error al enviar notificación por WebSocket (Usuario %d): %v", n.UserID, err)

			// Si falla, podría quedarse en estado "SCHEDULED" o actualizarse a "PENDING"
			continue
		}

		// Si el envío fue exitoso, actualiza el estado a "SENT"
		err = s.repo.UpdateNotificationStatus(n.Id, "SENT")
		if err != nil {
			log.Printf("Error al actualizar el estado de la notificación ID %d: %v", n.Id, err)
		}
	}

	return nil
}