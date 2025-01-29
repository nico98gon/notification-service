package notification

import "time"

type Repository interface {
	SaveNotification(notification *Notification) error
	FindScheduledNotifications(now time.Time) ([]Notification, error)
	UpdateNotificationStatus(id int, status string) error
}