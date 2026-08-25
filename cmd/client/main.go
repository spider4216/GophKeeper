package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spider4216/GophKeeper/internal/client/commands"
	"github.com/spider4216/GophKeeper/internal/client/services"
)

func main() {
	app := newApp()

	if err := app.Run(); err != nil {
		log.Fatal("Cannot run app", err)
	}

	service := services.New(app.cli, app.cfg.SrvHost, app.repo, app.logger)

	cmd := commands.New(service, app.cfg)

	// timeout to config
	ctx, cancel := context.WithTimeout(context.Background(), app.cfg.CtxTimeout)
	defer cancel()

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
	default:
		err = errors.New("command not found")
		fmt.Println("command not found")
		os.Exit(1)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(msg)
}
