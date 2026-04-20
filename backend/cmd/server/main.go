package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JustJoeYo/trophy-collector/internal/api"
	"github.com/JustJoeYo/trophy-collector/internal/cache"
	"github.com/JustJoeYo/trophy-collector/internal/clients"
	"github.com/JustJoeYo/trophy-collector/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    slog.SetDefault(logger)

    cfg := config.Load()

    deadlockClient := clients.NewDeadlockClient(cfg.DeadlockAPIURL, cfg.AssetsURL)
    redisCache := cache.NewRedisCache(cfg.RedisAddr)
    handler := api.NewHandler(deadlockClient, redisCache)

    r := chi.NewRouter()
    r.Use(middleware.RequestID)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
    r.Use(corsMiddleware)

    handler.RegisterRoutes(r)

        srv := &http.Server{
        Addr:    ":" + cfg.Port,
        Handler: r,
    }

    go func() {
        slog.Info("server starting", "port", cfg.Port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            slog.Error("server error", "error", err)
            os.Exit(1)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    slog.Info("shutting down server")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        slog.Error("forced shutdown", "error", err)
    }

    slog.Info("server stopped")
}


func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}