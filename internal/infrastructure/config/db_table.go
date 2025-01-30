package config

import (
	"database/sql"
	"fmt"
	"log"
)

func CreateNotificationsTable(db *sql.DB) {
	query := `
		CREATE TABLE IF NOT EXISTS notifications (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		locality_id VARCHAR(255) NOT NULL,
		status VARCHAR(50) NOT NULL,
		title VARCHAR(255),
		content TEXT,
		scheduled TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Error creando la tabla de notificaciones: %v", err)
	} else {
		fmt.Println("Tabla `notifications` creada con éxito.")
	}

	// Verificar si hay datos en la tabla
	checkQuery := `SELECT COUNT(*) FROM notifications`
	var count int
	err = db.QueryRow(checkQuery).Scan(&count)
	if err != nil {
		log.Fatalf("Error al verificar la tabla de notificaciones: %v", err)
	}

	// Si no hay notificaciones, insertar algunas por defecto
	if count == 0 {
		fmt.Println("No se encontró ninguna notificación. Insertando algunas por defecto...")

		insertQuery := `
		INSERT INTO notifications (user_id, locality_id, status, scheduled)
		VALUES
			(1, '123', 'SCHEDULED', NOW() + INTERVAL '1 hour'),
			(2, '456', 'PENDING', NOW() + INTERVAL '2 hour'),
			(3, '789', 'SCHEDULED', NOW() + INTERVAL '3 hour')
		`

		_, err = db.Exec(insertQuery)
		if err != nil {
			log.Fatalf("Error insertando notificaciones por defecto: %v", err)
		} else {
			fmt.Println("Notificaciones por defecto insertadas con éxito.")
		}
	} else {
		fmt.Printf("La tabla `notifications` tiene %d notificaciones. No se agregaron notificaciones por defecto.\n", count)
	}
}