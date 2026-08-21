package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/spider4216/GophKeeper/internal/server/handlers"
	"github.com/spider4216/GophKeeper/internal/server/middlewares"
	"github.com/spider4216/GophKeeper/internal/server/services"
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
	// todo gracefull shutdown

	app.logger.Debugf("Listen server on %s", app.cfg.ServerAddress)

	if err := srv.ListenAndServe(); err != nil {
		app.logger.Fatalf("Server error: %s", err)
	}

	app.logger.Debug("Everything is ok")
}
