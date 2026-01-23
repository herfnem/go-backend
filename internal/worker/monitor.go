package worker

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"learn/internal/database"
	"learn/internal/handlers"
)

type MonitorWorker struct {
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	httpClient *http.Client
}

func NewMonitorWorker() *MonitorWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &MonitorWorker{
		ctx:    ctx,
		cancel: cancel,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (w *MonitorWorker) Start() {
	log.Println("Monitor worker started")

	// Run monitor checks every minute
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		// Run immediately on start
		w.checkAllMonitors()

		for {
			select {
			case <-w.ctx.Done():
				log.Println("Monitor worker stopping...")
				return
			case <-ticker.C:
				w.checkAllMonitors()
			}
		}
	}()

	// Run snippet cleanup every hour
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-w.ctx.Done():
				return
			case <-ticker.C:
				deleted := handlers.DeleteExpiredSnippets()
				if deleted > 0 {
					log.Printf("Cleaned up %d expired snippets", deleted)
				}
			}
		}
	}()
}

func (w *MonitorWorker) Stop() {
	w.cancel()
	w.wg.Wait()
	log.Println("Monitor worker stopped")
}

func (w *MonitorWorker) checkAllMonitors() {
	rows, err := database.DB.Query(`
		SELECT id, url, interval_seconds
		FROM monitors
		WHERE is_active = 1
	`)
	if err != nil {
		log.Printf("Error fetching monitors: %v", err)
		return
	}
	defer rows.Close()

	type monitorInfo struct {
		ID       int64
		URL      string
		Interval int
	}

	var monitors []monitorInfo
	for rows.Next() {
		var m monitorInfo
		if err := rows.Scan(&m.ID, &m.URL, &m.Interval); err != nil {
			continue
		}
		monitors = append(monitors, m)
	}

	if len(monitors) == 0 {
		return
	}

	log.Printf("Checking %d monitors...", len(monitors))

	// Check monitors concurrently using goroutines
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10) // Limit concurrent checks to 10

	for _, m := range monitors {
		// Check if this monitor needs to be checked based on interval
		var lastCheck time.Time
		err := database.DB.QueryRow(
			"SELECT checked_at FROM monitor_logs WHERE monitor_id = ? ORDER BY checked_at DESC LIMIT 1",
			m.ID,
		).Scan(&lastCheck)

		if err == nil && time.Since(lastCheck) < time.Duration(m.Interval)*time.Second {
			continue // Skip, not time yet
		}

		wg.Add(1)
		go func(monitor monitorInfo) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			w.checkMonitor(monitor.ID, monitor.URL)
		}(m)
	}

	wg.Wait()
}

func (w *MonitorWorker) checkMonitor(monitorID int64, url string) {
	start := time.Now()

	req, err := http.NewRequestWithContext(w.ctx, "GET", url, nil)
	if err != nil {
		w.logResult(monitorID, "down", 0, 0, err.Error())
		return
	}

	req.Header.Set("User-Agent", "UptimeNinja/1.0")

	resp, err := w.httpClient.Do(req)
	responseTime := time.Since(start).Milliseconds()

	if err != nil {
		w.logResult(monitorID, "down", 0, responseTime, err.Error())
		return
	}
	defer resp.Body.Close()

	status := "up"
	if resp.StatusCode >= 400 {
		status = "down"
	}

	w.logResult(monitorID, status, resp.StatusCode, responseTime, "")
}

func (w *MonitorWorker) logResult(monitorID int64, status string, statusCode int, responseTimeMs int64, errorMsg string) {
	_, err := database.DB.Exec(`
		INSERT INTO monitor_logs (monitor_id, status, status_code, response_time_ms, error_message)
		VALUES (?, ?, ?, ?, ?)
	`, monitorID, status, statusCode, responseTimeMs, errorMsg)

	if err != nil {
		log.Printf("Error logging monitor result: %v", err)
	}
}
