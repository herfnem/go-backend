package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"learn/internal/models"
	"learn/internal/repository"
)

type MonitorWorker struct {
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	httpClient *http.Client
	monitors   repository.MonitorRepository
	snippets   *SnippetService
}

func NewMonitorWorker(monitors repository.MonitorRepository, snippets *SnippetService) *MonitorWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &MonitorWorker{
		ctx:    ctx,
		cancel: cancel,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		monitors: monitors,
		snippets: snippets,
	}
}

func (w *MonitorWorker) Start() {
	log.Println("Monitor worker started")

	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

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
				deleted, err := w.snippets.DeleteExpired(w.ctx)
				if err != nil {
					log.Printf("Error cleaning expired snippets: %v", err)
					continue
				}
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
	monitors, err := w.monitors.ListActive(w.ctx)
	if err != nil {
		log.Printf("Error fetching monitors: %v", err)
		return
	}

	if len(monitors) == 0 {
		return
	}

	log.Printf("Checking %d monitors...", len(monitors))

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10)

	for _, monitor := range monitors {
		lastCheck, err := w.monitors.GetLastCheck(w.ctx, monitor.ID)
		if err == nil && time.Since(lastCheck) < time.Duration(monitor.IntervalSeconds)*time.Second {
			continue
		}
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Printf("Error fetching last check for monitor %d: %v", monitor.ID, err)
			continue
		}

		wg.Add(1)
		go func(m models.MonitorCheck) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			w.checkMonitor(m.ID, m.URL)
		}(monitor)
	}

	wg.Wait()
}

func (w *MonitorWorker) checkMonitor(monitorID string, url string) {
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

func (w *MonitorWorker) logResult(monitorID string, status string, statusCode int, responseTimeMs int64, errorMsg string) {
	logEntry := models.MonitorLog{
		ID:             uuid.NewString(),
		MonitorID:      monitorID,
		Status:         status,
		StatusCode:     statusCode,
		ResponseTimeMs: responseTimeMs,
		ErrorMessage:   errorMsg,
	}

	if err := w.monitors.CreateLog(w.ctx, logEntry); err != nil {
		log.Printf("Error logging monitor result: %v", err)
	}
}
