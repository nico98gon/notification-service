package notification

import (
	"errors"
	"time"
)

type Notification struct {
	Id        	int       `json:"id"`
	UserID    	int       `json:"user_id"`
	LocalityID 	string    `json:"locality_id"`
	Status    	string    `json:"status"`
	Title     	*string   `json:"title,omitempty"`
	Content   	*string   `json:"content,omitempty"`
	Scheduled 	time.Time `json:"scheduled"`
	CreatedAt 	time.Time `json:"created_at"`
}	

func (u *Notification) Validate() error {
	if u.UserID <= 0 {
		return errors.New("user_id es requerido")
	}

	if u.LocalityID == "" {
		return errors.New("locality_id es requerido")
	}

	if u.Scheduled.Before(time.Now()) {
		return errors.New("scheduled time debe ser posterior a la fecha actual")
	}

	return nil
}
