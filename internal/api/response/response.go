package response

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Envelope struct {
	Success bool        `json:"success"`
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(payload)
}

func WriteSuccess(w http.ResponseWriter, status int, data interface{}, message ...string) {
	msg := http.StatusText(status)
	if len(message) > 0 {
		custom := strings.TrimSpace(message[0])
		if custom != "" {
			msg = custom
		}
	}

	WriteJSON(w, status, Envelope{
		Success: true,
		Status:  status,
		Message: msg,
		Data:    data,
	})
}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, Envelope{
		Success: false,
		Status:  status,
		Message: message,
	})
}
