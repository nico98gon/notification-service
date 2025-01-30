package middleware

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Status  string      `json:"status"` // "success", "error"
	Data    interface{} `json:"data"`   // Datos de respuesta (puede ser nil)
	Message string      `json:"message,omitempty"` // Mensaje de error o información adicional (opcional)
}

func RespondJSON(w http.ResponseWriter, statusCode int, status string, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := APIResponse{
		Status:  status,
		Data:    data,
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}

func SuccessResponse(w http.ResponseWriter, data interface{}) {
	RespondJSON(w, http.StatusOK, "success", data, "")
}

func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	RespondJSON(w, statusCode, "error", nil, message)
}
