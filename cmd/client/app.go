package main

import (
	"crypto/tls"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/spider4216/GophKeeper/internal/client/commands"
	"github.com/spider4216/GophKeeper/internal/client/config"
	"github.com/spider4216/GophKeeper/internal/client/repositories"
	"github.com/spider4216/GophKeeper/internal/logger"
	commonRep "github.com/spider4216/GophKeeper/internal/repository"
	migCli "github.com/spider4216/GophKeeper/migrations/client"
	"go.uber.org/zap"
)

type app struct {
	logger  *zap.SugaredLogger
	repo    repositories.Repository
	cfg     *config.Config
	cli     *http.Client
	cmdName commands.CmdName
	args    []string
}

func newApp() *app {
	return &app{}
}

func (app *app) Run() error {
	if err := app.initConfig(); err != nil {
		return err
	}

	if err := app.initLogger(); err != nil {
		return err
	}

	if err := app.initRepo(); err != nil {
		return err
	}

	if err := app.initMigrations(); err != nil {
		return err
	}

	if err := app.initCli(); err != nil {
		return err
	}

	if err := app.InitArgs(); err != nil {
		return err
	}

	return nil
}

func (app *app) initConfig() error {
	var cfg *config.Config

	cfg, err := config.New()
	if err != nil {
		return err
	}

	if cfg.DbDsn == "" {
		return errors.New("dsn didtn't passed")
	}

	app.cfg = cfg

	return nil
}

func (app *app) initRepo() error {
	common, err := commonRep.NewRepository(app.cfg.DbDsn, app.logger)
	if err != nil {
		return err
	}

	repo, err := repositories.NewRepository(app.cfg.DbDsn, app.logger, common)
	if err != nil {
		return err
	}

	app.repo = repo

	return nil
}

func (app *app) initMigrations() error {
	app.logger.Debug("Up migrations")

	repo, ok := app.repo.(*repositories.ClientRepository)

	if !ok {
		return fmt.Errorf("cannot cast to pgx store type in init migration")
	}

	src := repo.Source().(*sql.DB)

	if err := migCli.MigrateClient(src); err != nil {
		return err
	}

	app.logger.Debug("Migration done")

	return nil
}

func (app *app) initLogger() error {
	logger, err := logger.InitZap(app.cfg.LogLvl)
	if err != nil {
		return err
	}

	app.logger = logger

	return nil
}

func (app *app) initCli() error {
	dialer := &net.Dialer{
		Timeout: app.cfg.DialerTimeout,
	}

	tlsCfg := tls.Config{
		InsecureSkipVerify: true,
	}

	trans := &http.Transport{
		DialContext:           dialer.DialContext,
		TLSHandshakeTimeout:   app.cfg.TLSHandshakeTimeout,
		ResponseHeaderTimeout: app.cfg.RespHeaderTimeout,
		TLSClientConfig:       &tlsCfg,
	}

	client := &http.Client{
		Transport: trans,
		Timeout:   app.cfg.CliTimeout,
	}

	app.cli = client

	return nil
}

func (app *app) InitArgs() error {
	if len(os.Args) < 2 {
		return errors.New("too few arguments")
	}

	app.cmdName = commands.CmdName(os.Args[1])
	app.args = os.Args[2:]

	return nil
}
