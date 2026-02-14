package routes

import (
	"net/http"

	"learn/internal/api/handlers"
)

func RegisterMonitorRoutes(mux *http.ServeMux, handler *handlers.MonitorHandler, auth func(http.Handler) http.Handler) {
	mux.Handle("GET /monitors", auth(http.HandlerFunc(handler.GetMonitors)))
	mux.Handle("POST /monitors", auth(http.HandlerFunc(handler.CreateMonitor)))
	mux.Handle("GET /monitors/{id}", auth(http.HandlerFunc(handler.GetMonitor)))
	mux.Handle("DELETE /monitors/{id}", auth(http.HandlerFunc(handler.DeleteMonitor)))
	mux.Handle("PATCH /monitors/{id}/toggle", auth(http.HandlerFunc(handler.ToggleMonitor)))
	mux.Handle("GET /dashboard", auth(http.HandlerFunc(handler.GetDashboard)))
}
