package models

import "time"

type Monitor struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	Name            string    `json:"name"`
	URL             string    `json:"url"`
	IntervalSeconds int       `json:"interval_seconds"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	LastStatus      string    `json:"last_status,omitempty"`
}

type MonitorLog struct {
	ID             string    `json:"id"`
	MonitorID      string    `json:"monitor_id"`
	Status         string    `json:"status"`
	StatusCode     int       `json:"status_code"`
	ResponseTimeMs int64     `json:"response_time_ms"`
	ErrorMessage   string    `json:"error_message,omitempty"`
	CheckedAt      time.Time `json:"checked_at"`
}

type MonitorWithLogs struct {
	Monitor Monitor      `json:"monitor"`
	Logs    []MonitorLog `json:"logs"`
	Uptime  float64      `json:"uptime_percentage"`
}

type MonitorStats struct {
	Total  int
	Active int
	Up     int
	Down   int
}

type RecentMonitorLog struct {
	MonitorLog
	MonitorName string `json:"monitor_name"`
	MonitorURL  string `json:"monitor_url"`
}

type MonitorCheck struct {
	ID              string
	URL             string
	IntervalSeconds int
}
