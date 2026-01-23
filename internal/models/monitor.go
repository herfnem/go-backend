package models

import "time"

type Monitor struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	Name            string    `json:"name"`
	URL             string    `json:"url"`
	IntervalSeconds int       `json:"interval_seconds"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	LastStatus      string    `json:"last_status,omitempty"`
}

type MonitorLog struct {
	ID             int64     `json:"id"`
	MonitorID      int64     `json:"monitor_id"`
	Status         string    `json:"status"`
	StatusCode     int       `json:"status_code"`
	ResponseTimeMs int64     `json:"response_time_ms"`
	ErrorMessage   string    `json:"error_message,omitempty"`
	CheckedAt      time.Time `json:"checked_at"`
}

type CreateMonitorRequest struct {
	Name            string `json:"name" validate:"required,min=1,max=100"`
	URL             string `json:"url" validate:"required,url"`
	IntervalSeconds int    `json:"interval_seconds" validate:"omitempty,min=60,max=86400"`
}

type MonitorWithLogs struct {
	Monitor Monitor      `json:"monitor"`
	Logs    []MonitorLog `json:"logs"`
	Uptime  float64      `json:"uptime_percentage"`
}
