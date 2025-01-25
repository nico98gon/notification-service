package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "¡Servicio de notificaciones en funcionamiento!")
	})

	http.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("http://weather-service:8083/")
		if err != nil {
			http.Error(w, "Error al conectar con weather-service: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(w, "Respuesta de weather-service: %s", body)
	})

	fmt.Println("Servicio de notificaciones escuchando en el puerto 8082...")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
