package main

import (
	"crypto/tls"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	shrConfig "github.com/spider4216/GophKeeper/config"
	"github.com/spider4216/GophKeeper/internal/client/commands"
	"github.com/spider4216/GophKeeper/internal/client/config"
	"github.com/spider4216/GophKeeper/internal/client/repositories"
	"github.com/spider4216/GophKeeper/internal/logger"
	commonRep "github.com/spider4216/GophKeeper/internal/repository"
	migCli "github.com/spider4216/GophKeeper/migrations/client"
)

type app struct {
	logger  *slog.Logger
	repo    repositories.Repository
	cfg     *config.Config
	cli     *http.Client
	cmdName commands.CmdName
	args    []string
}

func newApp() *app {
	return &app{}
}

func (a *app) Run() error {
	_, err := shrConfig.NewBuilder(a).
		Step((*app).initConfig).
		Step((*app).initLogger).
		Step((*app).initRepo).
		Step((*app).initMigrations).
		Step((*app).initCli).
		Step((*app).InitArgs).
		Build()

	return err
}

func (a *app) initConfig() error {
	var cfg *config.Config

	cfg, err := config.New()
	if err != nil {
		return err
	}

	if cfg.DbDsn == "" {
		return errors.New("dsn didtn't passed")
	}

	a.cfg = cfg

	return nil
}

func (a *app) initRepo() error {
	common, err := commonRep.NewRepository(a.cfg.DbDsn, a.logger)
	if err != nil {
		return err
	}

	repo, err := repositories.NewRepository(a.cfg.DbDsn, a.logger, common)
	if err != nil {
		return err
	}

	a.repo = repo

	return nil
}

func (a *app) initMigrations() error {
	a.logger.Debug("Up migrations")

	repo, ok := a.repo.(*repositories.ClientRepository)

	if !ok {
		return fmt.Errorf("cannot cast to pgx store type in init migration")
	}

	src := repo.Source().(*sql.DB)

	if err := migCli.MigrateClient(src); err != nil {
		return err
	}

	a.logger.Debug("Migration done")

	return nil
}

func (a *app) initLogger() error {
	logger := logger.Init(a.cfg.LogLvl)

	a.logger = logger

	return nil
}

func (a *app) initCli() error {
	dialer := &net.Dialer{
		Timeout: a.cfg.DialerTimeout,
	}

	tlsCfg := tls.Config{
		InsecureSkipVerify: true,
	}

	trans := &http.Transport{
		DialContext:           dialer.DialContext,
		TLSHandshakeTimeout:   a.cfg.TLSHandshakeTimeout,
		ResponseHeaderTimeout: a.cfg.RespHeaderTimeout,
		TLSClientConfig:       &tlsCfg,
	}

	client := &http.Client{
		Transport: trans,
		Timeout:   a.cfg.CliTimeout,
	}

	a.cli = client

	return nil
}

func (a *app) InitArgs() error {
	if len(os.Args) < 2 {
		return errors.New("too few arguments")
	}

	a.cmdName = commands.CmdName(os.Args[1])
	a.args = os.Args[2:]

	return nil
}
