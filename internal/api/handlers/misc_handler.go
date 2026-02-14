package handlers

import (
	"net/http"

	"learn/internal/api/response"
	"learn/internal/types"
)

type MiscHandler struct{}

func NewMiscHandler() *MiscHandler {
	return &MiscHandler{}
}

// Home godoc
// @Summary API welcome
// @Tags misc
// @Produce json
// @Success 200 {object} types.HomeResponseEnvelope
// @Router / [get]
func (h *MiscHandler) Home(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, types.HomeResponse{
		Version: "1.0.0",
	}, "Welcome to the API")
}

// Health godoc
// @Summary Health check
// @Tags misc
// @Produce json
// @Success 200 {object} types.HealthResponseEnvelope
// @Router /health [get]
func (h *MiscHandler) Health(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, types.HealthResponse{
		Status: "healthy",
	}, "Service is healthy")
}
