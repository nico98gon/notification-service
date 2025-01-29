package notification

import "time"

type NotificationValidator interface {
	Validate() error
}

type NotificationEvent struct {
	EventType string 		`json:"event_type"` // Ej.: "USER_NOTIFIED"
	UserID    int    		`json:"user_id"`
	Title     string 		`json:"title"`
	Content   string 		`json:"content"`
	Scheduled time.Time `json:"scheduled"` // ISO-8601
}
