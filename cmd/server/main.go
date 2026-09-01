package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/spider4216/GophKeeper/internal/server/config"
	"github.com/spider4216/GophKeeper/internal/server/handlers"
	"github.com/spider4216/GophKeeper/internal/server/middlewares"
	"github.com/spider4216/GophKeeper/internal/server/services"
)

const (
	serverTimeout time.Duration = 5 * time.Second
)

func main() {
	app := newApp()

	if err := app.Run(); err != nil {
		log.Fatal("Cannot run app", err)
	}

	service := services.New(app.repo, app.logger)
	handler := handlers.New(app.cfg, app.logger, service)
	middleware := middlewares.New(app.logger, app.cfg, service)

	// ----
	mux := http.NewServeMux()
	mux.Handle("POST /auth/register", middleware.WithLogging(http.HandlerFunc(handler.CreateUser)))
	mux.Handle("POST /auth/login", middleware.WithLogging(http.HandlerFunc(handler.Login)))
	mux.Handle("POST /sync", middleware.WithJwt(middleware.WithLogging(http.HandlerFunc(handler.SyncIn))))
	mux.Handle("GET /sync", middleware.WithJwt(middleware.WithLogging(http.HandlerFunc(handler.SyncOut))))
	mux.Handle("PUT /items/{itemID}/chunks/{num}", middleware.WithJwt(middleware.WithLogging(http.HandlerFunc(handler.SaveChunk))))

	srv := &http.Server{
		Addr:         app.cfg.ServerAddress,
		Handler:      mux,
		ReadTimeout:  app.cfg.ReadTimeout,
		WriteTimeout: app.cfg.WriteTimeout,
		IdleTimeout:  app.cfg.IdleTimeout,
	}

	var wg sync.WaitGroup
	wg.Add(1)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	go func() {
		defer wg.Done()

		app.logger.Debug("Graceful shutdown mode on")
		<-ctx.Done()

		// Тут тоже останавляваем перехват сигналов
		stop()

		app.logger.Debug("Shutdown server...")

		ctxShutdown, cancel := context.WithTimeout(context.Background(), serverTimeout)
		defer cancel()

		if err := srv.Shutdown(ctxShutdown); err != nil {
			app.logger.Warn("Cannot shutdown main server", "error", err)
		}
	}()

	app.logger.Debug("Listen server", "address", app.cfg.ServerAddress)

	if err := runServer(srv, app.cfg, app.logger); err != nil && !errors.Is(err, http.ErrServerClosed) {
		app.logger.Error("Server error", "error", err)
		os.Exit(1)
	}

	wg.Wait()
}

func runServer(srv *http.Server, cfg *config.Config, logger *slog.Logger) error {
	if cfg.Https {
		logger.Info("Run HTTPS mode")
		return srv.ListenAndServeTLS(cfg.CrtPath, cfg.PKPath)
	}

	logger.Info("Run HTTP mode")
	return srv.ListenAndServe()
}
