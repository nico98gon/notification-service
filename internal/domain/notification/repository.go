package notification

import "time"

type Repository interface {
	FindAllNotifications() ([]Notification, error)
	FindNotificationByID(id int) (*Notification, error)
	FindScheduledNotifications(now time.Time) ([]Notification, error)
	FindPendingNotifications(now time.Time) ([]Notification, error)
	SaveNotification(notification *Notification) error
	UpdateNotification(id int, notification *Notification) error
	UpdateNotificationStatus(id int, status string) error
	DeleteNotification(id int) error
	MarkNotificationsAsPending(now time.Time) error
	SendAndMarkAsSent(id int) error
}