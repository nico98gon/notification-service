package repository

import (
	"database/sql"
	"time"

	"nilus-challenge-backend/internal/domain/notification"
)

type PostgresNotificationRepository struct {
	db *sql.DB
}

func NewPostgresNotificationRepository(db *sql.DB) *PostgresNotificationRepository {
	return &PostgresNotificationRepository{db: db}
}

func (r *PostgresNotificationRepository) SaveNotification(n *notification.Notification) error {
	query := `INSERT INTO notifications (user_id, title, content, scheduled, status) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(query, n.UserID, n.Title, n.Content, n.Scheduled, "SCHEDULED")
	if err != nil {
		return err
	}
	return nil
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

func (r *PostgresNotificationRepository) UpdateNotificationStatus(id int, status string) error {
	_, err := r.db.Exec(`UPDATE notifications SET status = $1 WHERE id = $2`, status, id)
	return err
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

