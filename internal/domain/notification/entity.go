package notification

import (
	"errors"
	"time"
)

type Notification struct {
	Id        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Status    string    `json:"status"` // "SCHEDULED", "SENT", "ERROR", "CANCELLED"
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Scheduled time.Time `json:"scheduled"`
	CreatedAt time.Time `json:"created_at"`
}	

func (u *Notification) Validate() error {
	if u.UserID <= 0 {
		return errors.New("user_id es requerido")
	}

	if len(u.Title) == 0 {
		return errors.New("titulo es requerido")
	}
	if len(u.Title) > 255 {
		return errors.New("titulo no debe exceder 255 caracteres")
	}

	if len(u.Content) > 1000 {
		return errors.New("contenido no debe exceder 1000 caracteres")
	}

	if u.Scheduled.Before(time.Now()) {
		return errors.New("scheduled time debe ser posterior a la fecha actual")
	}

	return nil
}
