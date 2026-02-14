package types

import "learn/internal/models"

type MonitorCreateRequest struct {
	Name            string `json:"name" validate:"required,min=1,max=100" example:"Google"`
	URL             string `json:"url" validate:"required,url" example:"https://google.com"`
	IntervalSeconds int    `json:"interval_seconds" validate:"omitempty,min=60,max=86400" example:"300"`
}

type MonitorResponseEnvelope struct {
	Success bool           `json:"success"`
	Status  int            `json:"status"`
	Message string         `json:"message"`
	Data    models.Monitor `json:"data"`
}

type MonitorListResponseEnvelope struct {
	Success bool             `json:"success"`
	Status  int              `json:"status"`
	Message string           `json:"message"`
	Data    []models.Monitor `json:"data"`
}

type MonitorWithLogsResponseEnvelope struct {
	Success bool                   `json:"success"`
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    models.MonitorWithLogs `json:"data"`
}

type MonitorDashboardResponse struct {
	TotalMonitors  int                       `json:"total_monitors"`
	ActiveMonitors int                       `json:"active_monitors"`
	UpMonitors     int                       `json:"up_monitors"`
	DownMonitors   int                       `json:"down_monitors"`
	RecentLogs     []models.RecentMonitorLog `json:"recent_logs"`
}

type MonitorDashboardResponseEnvelope struct {
	Success bool                     `json:"success"`
	Status  int                      `json:"status"`
	Message string                   `json:"message"`
	Data    MonitorDashboardResponse `json:"data"`
}
