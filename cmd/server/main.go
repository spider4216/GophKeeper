package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi"
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

	r := chi.NewRouter()

	r.Post("/auth/register", http.HandlerFunc(handler.CreateUser))
	r.Post("/auth/login", http.HandlerFunc(handler.Login))
	r.With(middleware.WithJwt).Post("/sync", http.HandlerFunc(handler.SyncIn))
	r.With(middleware.WithJwt).Get("/sync", http.HandlerFunc(handler.SyncOut))

	srv := &http.Server{
		Addr:         app.cfg.ServerAddress,
		Handler:      r,
		ReadTimeout:  app.cfg.ReadTimeout,
		WriteTimeout: app.cfg.WriteTimeout,
		IdleTimeout:  app.cfg.IdleTimeout,
	}

	// todo https

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
			app.logger.Warnf("Cannot shutdown main server: %s", err)
		}
	}()

	app.logger.Debugf("Listen server on %s", app.cfg.ServerAddress)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		app.logger.Fatalf("Server error: %s", err)
	}

	wg.Wait()
}
