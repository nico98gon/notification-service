package repository

import (
	"database/sql"
	"fmt"
	"time"

	"nilus-challenge-backend/internal/domain/notification"
)

type PostgresNotificationRepository struct {
	db *sql.DB
}

func NewPostgresNotificationRepository(db *sql.DB) *PostgresNotificationRepository {
	return &PostgresNotificationRepository{db: db}
}

func (r *PostgresNotificationRepository) FindAllNotifications() ([]notification.Notification, error) {
	query := `SELECT id, user_id, title, content, scheduled, status FROM notifications`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []notification.Notification
	for rows.Next() {
		var n notification.Notification
		if err := rows.Scan(&n.Id, &n.UserID, &n.Title, &n.Content, &n.Scheduled, &n.Status); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

func (r *PostgresNotificationRepository) FindNotificationByID(id int) (*notification.Notification, error) {
	query := `SELECT id, user_id, title, content, scheduled, status FROM notifications WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var n notification.Notification
	err := row.Scan(&n.Id, &n.UserID, &n.Title, &n.Content, &n.Scheduled, &n.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("notificación con ID %d no encontrada", id)
		}
		return nil, err
	}

	return &n, nil
}

func (r *PostgresNotificationRepository) FindScheduledNotifications(now time.Time) ([]notification.Notification, error) {
	query := `SELECT id, user_id, title, content, scheduled FROM notifications WHERE scheduled <= $1 AND status = 'SCHEDULED'`
	rows, err := r.db.Query(query, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []notification.Notification
	for rows.Next() {
		var n notification.Notification
		if err := rows.Scan(&n.Id, &n.UserID, &n.Title, &n.Content, &n.Scheduled); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

func (r *PostgresNotificationRepository) FindPendingNotifications(now time.Time) ([]notification.Notification, error) {
	var notifications []notification.Notification

	query := `SELECT id, user_id, title, content, scheduled FROM notifications WHERE scheduled <= $1 AND status = 'PENDING'`
	rows, err := r.db.Query(query, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var n notification.Notification
		if err := rows.Scan(&n.Id, &n.UserID, &n.Title, &n.Content, &n.Scheduled); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

func (r *PostgresNotificationRepository) SaveNotification(n *notification.Notification) error {
	query := `INSERT INTO notifications (user_id, title, content, scheduled, status) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(query, n.UserID, n.Title, n.Content, n.Scheduled, "SCHEDULED")
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresNotificationRepository) UpdateNotification(id int, n *notification.Notification) error {
	query := `UPDATE notifications SET user_id = $1, title = $2, content = $3, scheduled = $4, status = $5 WHERE id = $6`
	_, err := r.db.Exec(query, n.UserID, n.Title, n.Content, n.Scheduled, n.Status, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresNotificationRepository) UpdateNotificationStatus(id int, status string) error {
	_, err := r.db.Exec(`UPDATE notifications SET status = $1 WHERE id = $2`, status, id)
	return err
}

func (r *PostgresNotificationRepository) DeleteNotification(id int) error {
	query := `DELETE FROM notifications WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
