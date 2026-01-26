package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"learn/internal/config"
	"learn/internal/database"
	"learn/internal/routes"
	"learn/internal/worker"
)

func main() {
	// Load configuration
	config.Load()

	// Initialize database
	database.Init()
	defer database.Close()

	// Start background worker
	monitorWorker := worker.NewMonitorWorker()
	monitorWorker.Start()

	// Setup routes
	handler := routes.Setup()

	// Create server
	server := &http.Server{
		Addr:    ":" + config.AppConfig.Port,
		Handler: handler,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down server...")
		monitorWorker.Stop()
		server.Close()
	}()

	log.Printf("Server starting on http://localhost:%s", config.AppConfig.Port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal("Server error:", err)
	}
}
