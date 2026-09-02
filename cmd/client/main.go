package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spider4216/GophKeeper/internal/client/commands"
	"github.com/spider4216/GophKeeper/internal/client/services"
)

func main() {
	app := newApp()

	if err := app.Run(); err != nil {
		log.Fatal("Cannot run app", err)
	}

	service, errService := services.New(
		services.WithHTTPClient(app.cli),
		services.WithHost(app.cfg.SrvHost),
		services.WithRepo(app.repo),
		services.WithLogger(app.logger),
	)

	if errService != nil {
		fmt.Println(errService)
		os.Exit(1)
	}

	cmd := commands.New(service, app.cfg, app.logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var err error
	var msg string

	switch app.cmdName {
	case commands.Register:
		msg, err = cmd.Register(ctx, app.args)
	case commands.Login:
		msg, err = cmd.Login(ctx, app.args)
	case commands.InsertLoginPass:
		msg, err = cmd.CreateLoginpass(ctx, app.args)
	case commands.SyncSend:
		msg, err = cmd.SyncSend(ctx, app.args)
	case commands.List:
		msg, err = cmd.UserList(ctx, app.args)
	case commands.View:
		msg, err = cmd.View(ctx, app.args)
	case commands.SyncGet:
		msg, err = cmd.SyncGet(ctx, app.args)
	case commands.Delete:
		msg, err = cmd.DeleteItem(ctx, app.args)
	case commands.UpdateLoginPass:
		msg, err = cmd.UpdateLoginPass(ctx, app.args)
	case commands.Version:
		msg, err = cmd.PrintVersion()
	case commands.CreateText:
		msg, err = cmd.CreateText(ctx, app.args)
	case commands.GetText:
		msg, err = cmd.GetText(ctx, app.args)
	case commands.CreateBinary:
		msg, err = cmd.CreateBinary(ctx, app.args)
	case commands.DownloadBinary:
		msg, err = cmd.DownloadBinary(ctx, app.args)
	default:
		err = errors.New("command not found")
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(msg)
}
