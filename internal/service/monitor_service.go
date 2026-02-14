package service

import (
	"context"

	"github.com/google/uuid"
	"learn/internal/models"
	"learn/internal/repository"
)

type MonitorService struct {
	monitors repository.MonitorRepository
}

func NewMonitorService(monitors repository.MonitorRepository) *MonitorService {
	return &MonitorService{monitors: monitors}
}

func (s *MonitorService) Create(ctx context.Context, userID string, name, url string, intervalSeconds int) (models.Monitor, error) {
	if intervalSeconds == 0 {
		intervalSeconds = 300
	}

	monitor := models.Monitor{
		ID:              uuid.NewString(),
		UserID:          userID,
		Name:            name,
		URL:             url,
		IntervalSeconds: intervalSeconds,
		IsActive:        true,
	}

	return s.monitors.Create(ctx, monitor)
}

func (s *MonitorService) ListByUser(ctx context.Context, userID string) ([]models.Monitor, error) {
	return s.monitors.ListByUser(ctx, userID)
}

func (s *MonitorService) GetWithLogs(ctx context.Context, userID, id string) (models.MonitorWithLogs, error) {
	monitor, err := s.monitors.GetByID(ctx, userID, id)
	if err != nil {
		return models.MonitorWithLogs{}, err
	}

	logs, err := s.monitors.ListLogs(ctx, id, 50)
	if err != nil {
		return models.MonitorWithLogs{}, err
	}
	if logs == nil {
		logs = []models.MonitorLog{}
	}

	var upCount int
	for _, log := range logs {
		if log.Status == "up" {
			upCount++
		}
	}

	uptime := float64(0)
	if len(logs) > 0 {
		uptime = float64(upCount) / float64(len(logs)) * 100
	}

	return models.MonitorWithLogs{
		Monitor: monitor,
		Logs:    logs,
		Uptime:  uptime,
	}, nil
}

func (s *MonitorService) Delete(ctx context.Context, userID, id string) (bool, error) {
	return s.monitors.Delete(ctx, userID, id)
}

func (s *MonitorService) Toggle(ctx context.Context, userID, id string) (bool, error) {
	return s.monitors.Toggle(ctx, userID, id)
}

func (s *MonitorService) Dashboard(ctx context.Context, userID string) (models.MonitorStats, []models.RecentMonitorLog, error) {
	stats, err := s.monitors.CountStats(ctx, userID)
	if err != nil {
		return models.MonitorStats{}, nil, err
	}

	logs, err := s.monitors.ListRecentLogs(ctx, userID, 20)
	if err != nil {
		return models.MonitorStats{}, nil, err
	}

	return stats, logs, nil
}
