package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/ilmakiage/argocd-bitbucket-proxy/internal/config"
	"github.com/ilmakiage/argocd-bitbucket-proxy/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	level := slog.LevelInfo
	switch cfg.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})))

	srv := server.New(cfg, cfg.ArgocdWebhookURL)
	handler := srv.Handler()

	addr := ":" + cfg.Port
	slog.Info("starting argocd-bitbucket-proxy", "addr", addr, "target", cfg.ArgocdWebhookURL)

	if err := http.ListenAndServe(addr, handler); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
