package repository

import (
	"context"
	"time"

	"learn/internal/models"
)

type MonitorRepository interface {
	Create(ctx context.Context, monitor models.Monitor) (models.Monitor, error)
	ListByUser(ctx context.Context, userID string) ([]models.Monitor, error)
	GetByID(ctx context.Context, userID, id string) (models.Monitor, error)
	Delete(ctx context.Context, userID, id string) (bool, error)
	Toggle(ctx context.Context, userID, id string) (bool, error)
	ListLogs(ctx context.Context, monitorID string, limit int) ([]models.MonitorLog, error)
	ListRecentLogs(ctx context.Context, userID string, limit int) ([]models.RecentMonitorLog, error)
	CountStats(ctx context.Context, userID string) (models.MonitorStats, error)
	ListActive(ctx context.Context) ([]models.MonitorCheck, error)
	GetLastCheck(ctx context.Context, monitorID string) (time.Time, error)
	CreateLog(ctx context.Context, log models.MonitorLog) error
}
