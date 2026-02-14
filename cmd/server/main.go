package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "learn/docs"
	"learn/internal/api/handlers"
	"learn/internal/api/middleware"
	"learn/internal/api/routes"
	"learn/internal/config"
	"learn/internal/repository"
	"learn/internal/service"
)

// @title Backend Misc API
// @version 1.0
// @description A Go REST API that bundles auth, posts, uptime monitoring, and snippets.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and your JWT.
// @BasePath /
func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
			switch attr.Key {
			case slog.TimeKey:
				return slog.String(slog.TimeKey, attr.Value.Time().Format(time.RFC3339))
			case slog.LevelKey:
				return slog.String(slog.LevelKey, strings.ToUpper(attr.Value.String()))
			default:
				return attr
			}
		},
	}))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	db, err := config.OpenDB(cfg.DBPath)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := config.Migrate(context.Background(), db); err != nil {
		logger.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	userRepo := repository.NewSQLiteUserRepository(db)
	monitorRepo := repository.NewSQLiteMonitorRepository(db)
	snippetRepo := repository.NewSQLiteSnippetRepository(db)

	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiry)
	userService := service.NewUserService(userRepo)
	monitorService := service.NewMonitorService(monitorRepo)
	snippetService := service.NewSnippetService(snippetRepo)
	postService := service.NewPostService()

	monitorWorker := service.NewMonitorWorker(monitorRepo, snippetService)
	monitorWorker.Start()

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	monitorHandler := handlers.NewMonitorHandler(monitorService)
	snippetHandler := handlers.NewSnippetHandler(snippetService)
	postHandler := handlers.NewPostHandler(postService)
	miscHandler := handlers.NewMiscHandler()

	mux := http.NewServeMux()
	routes.RegisterSwaggerRoutes(mux)
	routes.RegisterMiscRoutes(mux, miscHandler)
	routes.RegisterAuthRoutes(mux, authHandler)

	authMiddleware := middleware.Auth(userRepo, cfg.JWTSecret)
	routes.RegisterUserRoutes(mux, userHandler, authMiddleware)
	routes.RegisterPostRoutes(mux, postHandler, authMiddleware)
	routes.RegisterMonitorRoutes(mux, monitorHandler, authMiddleware)
	routes.RegisterSnippetRoutes(mux, snippetHandler)

	handler := middleware.Chain(mux,
		middleware.Recovery(logger),
		middleware.SecurityHeaders(),
		middleware.CORS(cfg.AllowedOrigins),
		middleware.Timeout(cfg.RequestTimeout),
		middleware.Logging(logger),
	)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		logger.Info("server started", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server stopped unexpectedly", "error", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("shutting down")
	monitorWorker.Stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown error", "error", err)
	}
}
