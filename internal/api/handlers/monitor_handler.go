package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"learn/internal/api/middleware"
	"learn/internal/api/response"
	"learn/internal/api/validator"
	"learn/internal/models"
	"learn/internal/service"
	"learn/internal/types"
)

type MonitorHandler struct {
	monitors *service.MonitorService
}

func NewMonitorHandler(monitors *service.MonitorService) *MonitorHandler {
	return &MonitorHandler{monitors: monitors}
}

// CreateMonitor godoc
// @Summary Create a monitor
// @Tags monitors
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body types.MonitorCreateRequest true "Create monitor"
// @Success 201 {object} types.MonitorResponseEnvelope
// @Failure 400 {object} types.ErrorResponseEnvelope
// @Failure 401 {object} types.ErrorResponseEnvelope
// @Failure 500 {object} types.ErrorResponseEnvelope
// @Router /monitors [post]
func (h *MonitorHandler) CreateMonitor(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	var req types.MonitorCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validator.Validate(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, validator.FormatErrorsString(err))
		return
	}

	monitor, err := h.monitors.Create(r.Context(), user.ID, req.Name, req.URL, req.IntervalSeconds)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to create monitor")
		return
	}

	response.WriteSuccess(w, http.StatusCreated, monitor, "Monitor created successfully")
}

// GetMonitors godoc
// @Summary List monitors
// @Tags monitors
// @Security BearerAuth
// @Produce json
// @Success 200 {object} types.MonitorListResponseEnvelope
// @Failure 401 {object} types.ErrorResponseEnvelope
// @Failure 500 {object} types.ErrorResponseEnvelope
// @Router /monitors [get]
func (h *MonitorHandler) GetMonitors(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	monitors, err := h.monitors.ListByUser(r.Context(), user.ID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Database error")
		return
	}

	if monitors == nil {
		monitors = []models.Monitor{}
	}

	response.WriteSuccess(w, http.StatusOK, monitors, "Monitors retrieved successfully")
}

// GetMonitor godoc
// @Summary Get monitor with logs
// @Tags monitors
// @Security BearerAuth
// @Produce json
// @Param id path string true "Monitor ID"
// @Success 200 {object} types.MonitorWithLogsResponseEnvelope
// @Failure 400 {object} types.ErrorResponseEnvelope
// @Failure 401 {object} types.ErrorResponseEnvelope
// @Failure 404 {object} types.ErrorResponseEnvelope
// @Failure 500 {object} types.ErrorResponseEnvelope
// @Router /monitors/{id} [get]
func (h *MonitorHandler) GetMonitor(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		response.WriteError(w, http.StatusBadRequest, "Invalid monitor ID")
		return
	}

	result, err := h.monitors.GetWithLogs(r.Context(), user.ID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.WriteError(w, http.StatusNotFound, "Monitor not found")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Database error")
		return
	}

	response.WriteSuccess(w, http.StatusOK, result, "Monitor retrieved successfully")
}

// DeleteMonitor godoc
// @Summary Delete a monitor
// @Tags monitors
// @Security BearerAuth
// @Produce json
// @Param id path string true "Monitor ID"
// @Success 200 {object} types.EmptyResponseEnvelope
// @Failure 401 {object} types.ErrorResponseEnvelope
// @Failure 404 {object} types.ErrorResponseEnvelope
// @Failure 500 {object} types.ErrorResponseEnvelope
// @Router /monitors/{id} [delete]
func (h *MonitorHandler) DeleteMonitor(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		response.WriteError(w, http.StatusBadRequest, "Invalid monitor ID")
		return
	}

	deleted, err := h.monitors.Delete(r.Context(), user.ID, id)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if !deleted {
		response.WriteError(w, http.StatusNotFound, "Monitor not found")
		return
	}

	response.WriteSuccess(w, http.StatusOK, nil, "Monitor deleted successfully")
}

// ToggleMonitor godoc
// @Summary Toggle monitor
// @Tags monitors
// @Security BearerAuth
// @Produce json
// @Param id path string true "Monitor ID"
// @Success 200 {object} types.EmptyResponseEnvelope
// @Failure 401 {object} types.ErrorResponseEnvelope
// @Failure 404 {object} types.ErrorResponseEnvelope
// @Failure 500 {object} types.ErrorResponseEnvelope
// @Router /monitors/{id}/toggle [patch]
func (h *MonitorHandler) ToggleMonitor(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		response.WriteError(w, http.StatusBadRequest, "Invalid monitor ID")
		return
	}

	updated, err := h.monitors.Toggle(r.Context(), user.ID, id)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if !updated {
		response.WriteError(w, http.StatusNotFound, "Monitor not found")
		return
	}

	response.WriteSuccess(w, http.StatusOK, nil, "Monitor toggled successfully")
}

// GetDashboard godoc
// @Summary Get monitor dashboard
// @Tags monitors
// @Security BearerAuth
// @Produce json
// @Success 200 {object} types.MonitorDashboardResponseEnvelope
// @Failure 401 {object} types.ErrorResponseEnvelope
// @Failure 500 {object} types.ErrorResponseEnvelope
// @Router /dashboard [get]
func (h *MonitorHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	stats, logs, err := h.monitors.Dashboard(r.Context(), user.ID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if logs == nil {
		logs = []models.RecentMonitorLog{}
	}

	response.WriteSuccess(w, http.StatusOK, types.MonitorDashboardResponse{
		TotalMonitors:  stats.Total,
		ActiveMonitors: stats.Active,
		UpMonitors:     stats.Up,
		DownMonitors:   stats.Down,
		RecentLogs:     logs,
	}, "Dashboard retrieved successfully")
}
