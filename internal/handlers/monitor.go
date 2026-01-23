package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"learn/internal/database"
	"learn/internal/middleware"
	"learn/internal/models"
	"learn/internal/response"
	"learn/internal/validator"
)

func CreateMonitor(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		response.Unauthorized(w, "User not found in context")
		return
	}

	var req models.CreateMonitorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	if err := validator.Validate(req); err != nil {
		response.BadRequest(w, validator.FormatErrorsString(err))
		return
	}

	// Default interval to 5 minutes if not specified
	if req.IntervalSeconds == 0 {
		req.IntervalSeconds = 300
	}

	result, err := database.DB.Exec(
		"INSERT INTO monitors (user_id, name, url, interval_seconds) VALUES (?, ?, ?, ?)",
		user.ID, req.Name, req.URL, req.IntervalSeconds,
	)
	if err != nil {
		response.InternalError(w, "Failed to create monitor")
		return
	}

	id, _ := result.LastInsertId()

	monitor := models.Monitor{
		ID:              id,
		UserID:          user.ID,
		Name:            req.Name,
		URL:             req.URL,
		IntervalSeconds: req.IntervalSeconds,
		IsActive:        true,
	}

	response.Created(w, "Monitor created successfully", monitor)
}

func GetMonitors(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		response.Unauthorized(w, "User not found in context")
		return
	}

	rows, err := database.DB.Query(`
		SELECT m.id, m.user_id, m.name, m.url, m.interval_seconds, m.is_active, m.created_at,
			COALESCE((SELECT status FROM monitor_logs WHERE monitor_id = m.id ORDER BY checked_at DESC LIMIT 1), 'pending') as last_status
		FROM monitors m
		WHERE m.user_id = ?
		ORDER BY m.created_at DESC
	`, user.ID)
	if err != nil {
		response.InternalError(w, "Database error")
		return
	}
	defer rows.Close()

	var monitors []models.Monitor
	for rows.Next() {
		var m models.Monitor
		if err := rows.Scan(&m.ID, &m.UserID, &m.Name, &m.URL, &m.IntervalSeconds, &m.IsActive, &m.CreatedAt, &m.LastStatus); err != nil {
			continue
		}
		monitors = append(monitors, m)
	}

	if monitors == nil {
		monitors = []models.Monitor{}
	}

	response.Success(w, "Monitors retrieved successfully", monitors)
}

func GetMonitor(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		response.Unauthorized(w, "User not found in context")
		return
	}

	id := r.PathValue("id")

	var monitor models.Monitor
	err := database.DB.QueryRow(`
		SELECT id, user_id, name, url, interval_seconds, is_active, created_at
		FROM monitors WHERE id = ? AND user_id = ?
	`, id, user.ID).Scan(&monitor.ID, &monitor.UserID, &monitor.Name, &monitor.URL, &monitor.IntervalSeconds, &monitor.IsActive, &monitor.CreatedAt)

	if err == sql.ErrNoRows {
		response.NotFound(w, "Monitor not found")
		return
	}
	if err != nil {
		response.InternalError(w, "Database error")
		return
	}

	// Get logs for this monitor
	rows, err := database.DB.Query(`
		SELECT id, monitor_id, status, status_code, response_time_ms, COALESCE(error_message, ''), checked_at
		FROM monitor_logs
		WHERE monitor_id = ?
		ORDER BY checked_at DESC
		LIMIT 50
	`, id)
	if err != nil {
		response.InternalError(w, "Database error")
		return
	}
	defer rows.Close()

	var logs []models.MonitorLog
	var upCount, totalCount int
	for rows.Next() {
		var log models.MonitorLog
		if err := rows.Scan(&log.ID, &log.MonitorID, &log.Status, &log.StatusCode, &log.ResponseTimeMs, &log.ErrorMessage, &log.CheckedAt); err != nil {
			continue
		}
		logs = append(logs, log)
		totalCount++
		if log.Status == "up" {
			upCount++
		}
	}

	if logs == nil {
		logs = []models.MonitorLog{}
	}

	uptime := float64(0)
	if totalCount > 0 {
		uptime = float64(upCount) / float64(totalCount) * 100
	}

	response.Success(w, "Monitor retrieved successfully", models.MonitorWithLogs{
		Monitor: monitor,
		Logs:    logs,
		Uptime:  uptime,
	})
}

func DeleteMonitor(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		response.Unauthorized(w, "User not found in context")
		return
	}

	id := r.PathValue("id")

	result, err := database.DB.Exec("DELETE FROM monitors WHERE id = ? AND user_id = ?", id, user.ID)
	if err != nil {
		response.InternalError(w, "Database error")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		response.NotFound(w, "Monitor not found")
		return
	}

	response.Success(w, "Monitor deleted successfully", nil)
}

func ToggleMonitor(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		response.Unauthorized(w, "User not found in context")
		return
	}

	id := r.PathValue("id")

	result, err := database.DB.Exec(`
		UPDATE monitors SET is_active = NOT is_active WHERE id = ? AND user_id = ?
	`, id, user.ID)
	if err != nil {
		response.InternalError(w, "Database error")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		response.NotFound(w, "Monitor not found")
		return
	}

	response.Success(w, "Monitor toggled successfully", nil)
}

func GetDashboard(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		response.Unauthorized(w, "User not found in context")
		return
	}

	// Get summary stats
	var totalMonitors, activeMonitors, upMonitors, downMonitors int

	database.DB.QueryRow("SELECT COUNT(*) FROM monitors WHERE user_id = ?", user.ID).Scan(&totalMonitors)
	database.DB.QueryRow("SELECT COUNT(*) FROM monitors WHERE user_id = ? AND is_active = 1", user.ID).Scan(&activeMonitors)

	// Count up/down based on latest status
	database.DB.QueryRow(`
		SELECT COUNT(DISTINCT m.id) FROM monitors m
		JOIN monitor_logs ml ON m.id = ml.monitor_id
		WHERE m.user_id = ? AND ml.status = 'up'
		AND ml.checked_at = (SELECT MAX(checked_at) FROM monitor_logs WHERE monitor_id = m.id)
	`, user.ID).Scan(&upMonitors)

	database.DB.QueryRow(`
		SELECT COUNT(DISTINCT m.id) FROM monitors m
		JOIN monitor_logs ml ON m.id = ml.monitor_id
		WHERE m.user_id = ? AND ml.status = 'down'
		AND ml.checked_at = (SELECT MAX(checked_at) FROM monitor_logs WHERE monitor_id = m.id)
	`, user.ID).Scan(&downMonitors)

	// Get recent logs across all monitors
	rows, err := database.DB.Query(`
		SELECT ml.id, ml.monitor_id, ml.status, ml.status_code, ml.response_time_ms,
			COALESCE(ml.error_message, ''), ml.checked_at, m.name, m.url
		FROM monitor_logs ml
		JOIN monitors m ON ml.monitor_id = m.id
		WHERE m.user_id = ?
		ORDER BY ml.checked_at DESC
		LIMIT 20
	`, user.ID)
	if err != nil {
		response.InternalError(w, "Database error")
		return
	}
	defer rows.Close()

	type RecentLog struct {
		models.MonitorLog
		MonitorName string `json:"monitor_name"`
		MonitorURL  string `json:"monitor_url"`
	}

	var recentLogs []RecentLog
	for rows.Next() {
		var log RecentLog
		if err := rows.Scan(&log.ID, &log.MonitorID, &log.Status, &log.StatusCode, &log.ResponseTimeMs, &log.ErrorMessage, &log.CheckedAt, &log.MonitorName, &log.MonitorURL); err != nil {
			continue
		}
		recentLogs = append(recentLogs, log)
	}

	if recentLogs == nil {
		recentLogs = []RecentLog{}
	}

	response.Success(w, "Dashboard retrieved successfully", map[string]any{
		"total_monitors":  totalMonitors,
		"active_monitors": activeMonitors,
		"up_monitors":     upMonitors,
		"down_monitors":   downMonitors,
		"recent_logs":     recentLogs,
	})
}
