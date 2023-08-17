package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"chord-drawer/app/internal/domain/service/chord"
	"chord-drawer/app/internal/session"
	v1 "chord-drawer/app/internal/transport/http/v1"
	"chord-drawer/app/pkg/config"
)

const network = "tcp"

type Application struct {
	cfg    *config.Config
	router *http.ServeMux
	log    *slog.Logger
	mu     *sync.Mutex
}

func ApplicationInit(cfg *config.Config, router *http.ServeMux, log *slog.Logger, mu *sync.Mutex) *Application {
	return &Application{
		cfg:    cfg,
		router: router,
		log:    log,
		mu:     mu,
	}
}

func (a *Application) Start(ctx context.Context) error {
	listener, err := net.Listen(network, fmt.Sprintf("%s:%s", a.cfg.Server.BindIP, a.cfg.Server.Port))
	if err != nil {
		return err
	}
	service := chord.NewService(a.cfg, a.mu)
	handler := v1.NewHandler(a.log, a.cfg, service)
	manager, err := session.Init(a.log, a.cfg, a.mu)
	if err != nil {
		return err
	}
	handler.Register(a.router, manager)
	server := http.Server{
		Handler:      a.router,
		WriteTimeout: a.cfg.Server.Timeout * time.Second,
		ReadTimeout:  a.cfg.Server.Timeout * time.Second,
	}
	go func() {
		if err = server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.ErrorContext(ctx, "failed to start server", slog.Any("details", err))
			panic(err)
		}
	}()
	a.log.Info("listening on", slog.String("IP", a.cfg.Server.BindIP), slog.String("port", a.cfg.Server.Port))
	<-ctx.Done()
	a.log.Info("server is shutting down...")
	sdCtx, cancel := context.WithTimeout(context.Background(), a.cfg.Server.ShutdownDuration*time.Second)
	defer cancel()
	if err = server.Shutdown(sdCtx); err != nil {
		return err
	}
	longShutDown := make(chan struct{}, 1)
	go func() {
		manager.Reset()
		longShutDown <- struct{}{}
	}()
	select {
	case <-sdCtx.Done():
		return sdCtx.Err()
	case <-longShutDown:
		a.log.Info("shutdown finished")
	}
	return nil
}
