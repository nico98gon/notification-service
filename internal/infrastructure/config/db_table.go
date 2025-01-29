package config

import (
	"database/sql"
	"fmt"
	"log"
)

func CreateUsersTable(db *sql.DB) {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			opt_out BOOLEAN DEFAULT FALSE,
			locality_id INT UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Error creando tabla de usuarios: %v", err)
	} else {
		fmt.Println("Tabla `users` creada con éxito.")
	}

	checkQuery := `SELECT COUNT(*) FROM users`
	var count int
	err = db.QueryRow(checkQuery).Scan(&count)
	if err != nil {
		log.Fatalf("Error checkeando la tabla de usuarios: %v", err)
	}

	if count == 0 {
		fmt.Println("No se encontró usuarios en la tabla `users`. Insertando usuarios por defecto...")

		insertQuery := `
			INSERT INTO users (name, email, opt_out)
			VALUES
				('Gonzalo Darín', 'gonza.darin@example.com', FALSE),
				('Julieta Castaño', 'juli.castaño@example.com', TRUE),
				('Fernando Martínez', 'fernando.Martinez@example.com', TRUE),
				('Carla López', 'carla.lopez@example.com', FALSE),
				('Lucas Fernández', 'lucas.fernandez@example.com', TRUE),
				('María González', 'maria.gonzalez@example.com', FALSE),
				('Santiago Ruiz', 'santiago.ruiz@example.com', TRUE),
				('Camila Vega', 'camila.vega@example.com', FALSE),
				('Martín Silva', 'martin.silva@example.com', TRUE)
		`
		_, err = db.Exec(insertQuery)
		if err != nil {
			log.Fatalf("Error insertando usuarios por defecto: %v", err)
		} else {
			fmt.Println("Usuarios por defecto insertados con éxito.")
		}
	} else {
		fmt.Printf("La tabla `users` tiene %d usuarios. No se agregaron usuarios por defecto.\n", count)
	}
}
