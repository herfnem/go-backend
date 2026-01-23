package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool   `json:"success"`
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func JSON(w http.ResponseWriter, status int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := Response{
		Success: status >= 200 && status < 300,
		Status:  status,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(resp)
}

func Success(w http.ResponseWriter, message string, data any) {
	JSON(w, http.StatusOK, message, data)
}

func Created(w http.ResponseWriter, message string, data any) {
	JSON(w, http.StatusCreated, message, data)
}

func BadRequest(w http.ResponseWriter, message string) {
	JSON(w, http.StatusBadRequest, message, nil)
}

func Unauthorized(w http.ResponseWriter, message string) {
	JSON(w, http.StatusUnauthorized, message, nil)
}

func NotFound(w http.ResponseWriter, message string) {
	JSON(w, http.StatusNotFound, message, nil)
}

func Conflict(w http.ResponseWriter, message string) {
	JSON(w, http.StatusConflict, message, nil)
}

func InternalError(w http.ResponseWriter, message string) {
	JSON(w, http.StatusInternalServerError, message, nil)
}
