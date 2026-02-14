package repository

import (
	"context"
	"database/sql"
	"time"

	"learn/internal/models"
)

type SQLiteMonitorRepository struct {
	db *sql.DB
}

func NewSQLiteMonitorRepository(db *sql.DB) *SQLiteMonitorRepository {
	return &SQLiteMonitorRepository{db: db}
}

func (r *SQLiteMonitorRepository) Create(ctx context.Context, monitor models.Monitor) (models.Monitor, error) {
	_, err := r.db.ExecContext(ctx, `
INSERT INTO monitors (id, user_id, name, url, interval_seconds, is_active)
VALUES (?, ?, ?, ?, ?, ?)
`, monitor.ID, monitor.UserID, monitor.Name, monitor.URL, monitor.IntervalSeconds, boolToInt(monitor.IsActive))
	if err != nil {
		return models.Monitor{}, err
	}

	return r.GetByID(ctx, monitor.UserID, monitor.ID)
}

func (r *SQLiteMonitorRepository) ListByUser(ctx context.Context, userID string) ([]models.Monitor, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT m.id, m.user_id, m.name, m.url, m.interval_seconds, m.is_active, m.created_at,
	COALESCE((SELECT status FROM monitor_logs WHERE monitor_id = m.id ORDER BY checked_at DESC LIMIT 1), 'pending') as last_status
FROM monitors m
WHERE m.user_id = ?
ORDER BY m.created_at DESC
`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var monitors []models.Monitor
	for rows.Next() {
		var monitor models.Monitor
		var isActive int
		if err := rows.Scan(&monitor.ID, &monitor.UserID, &monitor.Name, &monitor.URL, &monitor.IntervalSeconds, &isActive, &monitor.CreatedAt, &monitor.LastStatus); err != nil {
			return nil, err
		}
		monitor.IsActive = isActive == 1
		monitors = append(monitors, monitor)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return monitors, nil
}

func (r *SQLiteMonitorRepository) GetByID(ctx context.Context, userID, id string) (models.Monitor, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT id, user_id, name, url, interval_seconds, is_active, created_at
FROM monitors
WHERE id = ? AND user_id = ?
`, id, userID)

	var monitor models.Monitor
	var isActive int
	if err := row.Scan(&monitor.ID, &monitor.UserID, &monitor.Name, &monitor.URL, &monitor.IntervalSeconds, &isActive, &monitor.CreatedAt); err != nil {
		return models.Monitor{}, err
	}
	monitor.IsActive = isActive == 1

	return monitor, nil
}

func (r *SQLiteMonitorRepository) Delete(ctx context.Context, userID, id string) (bool, error) {
	result, err := r.db.ExecContext(ctx, `
DELETE FROM monitors WHERE id = ? AND user_id = ?
`, id, userID)
	if err != nil {
		return false, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}

func (r *SQLiteMonitorRepository) Toggle(ctx context.Context, userID, id string) (bool, error) {
	result, err := r.db.ExecContext(ctx, `
UPDATE monitors SET is_active = NOT is_active WHERE id = ? AND user_id = ?
`, id, userID)
	if err != nil {
		return false, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}

func (r *SQLiteMonitorRepository) ListLogs(ctx context.Context, monitorID string, limit int) ([]models.MonitorLog, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT id, monitor_id, status, status_code, response_time_ms, COALESCE(error_message, ''), checked_at
FROM monitor_logs
WHERE monitor_id = ?
ORDER BY checked_at DESC
LIMIT ?
`, monitorID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.MonitorLog
	for rows.Next() {
		var log models.MonitorLog
		if err := rows.Scan(&log.ID, &log.MonitorID, &log.Status, &log.StatusCode, &log.ResponseTimeMs, &log.ErrorMessage, &log.CheckedAt); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

func (r *SQLiteMonitorRepository) ListRecentLogs(ctx context.Context, userID string, limit int) ([]models.RecentMonitorLog, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT ml.id, ml.monitor_id, ml.status, ml.status_code, ml.response_time_ms,
	COALESCE(ml.error_message, ''), ml.checked_at, m.name, m.url
FROM monitor_logs ml
JOIN monitors m ON ml.monitor_id = m.id
WHERE m.user_id = ?
ORDER BY ml.checked_at DESC
LIMIT ?
`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.RecentMonitorLog
	for rows.Next() {
		var log models.RecentMonitorLog
		if err := rows.Scan(&log.ID, &log.MonitorID, &log.Status, &log.StatusCode, &log.ResponseTimeMs, &log.ErrorMessage, &log.CheckedAt, &log.MonitorName, &log.MonitorURL); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

func (r *SQLiteMonitorRepository) CountStats(ctx context.Context, userID string) (models.MonitorStats, error) {
	var stats models.MonitorStats

	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM monitors WHERE user_id = ?", userID).Scan(&stats.Total); err != nil {
		return models.MonitorStats{}, err
	}

	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM monitors WHERE user_id = ? AND is_active = 1", userID).Scan(&stats.Active); err != nil {
		return models.MonitorStats{}, err
	}

	if err := r.db.QueryRowContext(ctx, `
SELECT COUNT(DISTINCT m.id) FROM monitors m
JOIN monitor_logs ml ON m.id = ml.monitor_id
WHERE m.user_id = ? AND ml.status = 'up'
AND ml.checked_at = (SELECT MAX(checked_at) FROM monitor_logs WHERE monitor_id = m.id)
`, userID).Scan(&stats.Up); err != nil {
		return models.MonitorStats{}, err
	}

	if err := r.db.QueryRowContext(ctx, `
SELECT COUNT(DISTINCT m.id) FROM monitors m
JOIN monitor_logs ml ON m.id = ml.monitor_id
WHERE m.user_id = ? AND ml.status = 'down'
AND ml.checked_at = (SELECT MAX(checked_at) FROM monitor_logs WHERE monitor_id = m.id)
`, userID).Scan(&stats.Down); err != nil {
		return models.MonitorStats{}, err
	}

	return stats, nil
}

func (r *SQLiteMonitorRepository) ListActive(ctx context.Context) ([]models.MonitorCheck, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT id, url, interval_seconds
FROM monitors
WHERE is_active = 1
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var monitors []models.MonitorCheck
	for rows.Next() {
		var monitor models.MonitorCheck
		if err := rows.Scan(&monitor.ID, &monitor.URL, &monitor.IntervalSeconds); err != nil {
			return nil, err
		}
		monitors = append(monitors, monitor)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return monitors, nil
}

func (r *SQLiteMonitorRepository) GetLastCheck(ctx context.Context, monitorID string) (time.Time, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT checked_at FROM monitor_logs WHERE monitor_id = ? ORDER BY checked_at DESC LIMIT 1
`, monitorID)
	var checkedAt time.Time
	if err := row.Scan(&checkedAt); err != nil {
		return time.Time{}, err
	}
	return checkedAt, nil
}

func (r *SQLiteMonitorRepository) CreateLog(ctx context.Context, log models.MonitorLog) error {
	_, err := r.db.ExecContext(ctx, `
INSERT INTO monitor_logs (id, monitor_id, status, status_code, response_time_ms, error_message)
VALUES (?, ?, ?, ?, ?, ?)
`, log.ID, log.MonitorID, log.Status, log.StatusCode, log.ResponseTimeMs, log.ErrorMessage)
	return err
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
