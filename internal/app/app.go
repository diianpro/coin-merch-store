package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"diianpro/coin-merch-store/config"
	"diianpro/coin-merch-store/internal/repo"
	"diianpro/coin-merch-store/internal/service"
	"diianpro/coin-merch-store/internal/transport"
	"diianpro/coin-merch-store/pkg/hasher"
	"diianpro/coin-merch-store/pkg/httpserver"
	"diianpro/coin-merch-store/pkg/postgres"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func Run() {
	// Configuration
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Repositories
	log.Info("Initializing postgres...")
	pg, err := postgres.New(ctx, &postgres.Config{
		URL:      cfg.Postgres.URL,
		MinConns: cfg.Postgres.MinConns,
		MaxConns: cfg.Postgres.MaxConns,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - pgdb.NewServices: %w", err))
	}
	defer pg.Close()

	// Repositories
	log.Info("Initializing repositories...")
	repositories := repo.NewRepositories(pg)

	// Services dependencies
	log.Info("Initializing services...")
	deps := service.ServicesDependencies{
		Repos:    repositories,
		Hasher:   hasher.New(cfg.Hasher.Salt),
		SignKey:  cfg.JWT.SignKey,
		TokenTTL: cfg.JWT.TokenTTL,
	}
	services := service.NewServices(deps)

	// Echo handler
	log.Info("Initializing handlers and routes...")
	handler := echo.New()
	transport.NewRouter(handler, services)

	// HTTP server
	log.Info("Starting http server...")
	log.Debugf("Server port: %s", cfg.HTTPPort)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTPPort))

	// Waiting signal
	log.Info("Configuring graceful shutdown...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Graceful shutdown
	log.Info("Shutting down...")
	err = httpServer.Shutdown()
	if err != nil {
		log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
