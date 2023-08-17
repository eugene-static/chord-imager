package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"chord-drawer/app/internal/server"
	"chord-drawer/app/pkg/config"
	"chord-drawer/app/pkg/logger/slogger"
)

func main() {
	cfgPath := flag.String("cfg", "app/config/config.json", "path to config file")
	flag.Parse()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	cfg, err := config.GetConfig(cfgPath)
	if err != nil || cfg == nil {
		slog.Error("failed to get config:", "details", err)
		os.Exit(1)
	}
	log := slogger.New(cfg.System)
	log.Info("config initialized successfully")
	router := http.NewServeMux()
	var mu sync.Mutex
	app := server.ApplicationInit(cfg, router, log, &mu)
	if err = app.Start(ctx); err != nil {
		log.Error("failed to start server", slog.Any("details", err))
		os.Exit(1)
	}
}
