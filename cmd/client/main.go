package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/spider4216/GophKeeper/internal/client/commands"
	"github.com/spider4216/GophKeeper/internal/client/services"
)

func main() {
	app := newApp()

	if err := app.Run(); err != nil {
		log.Fatal("Cannot run app", err)
	}

	// todo перенести в init app
	if len(os.Args) < 2 {
		fmt.Println("command is required")
		os.Exit(1)
	}

	// todo move to app
	// todo use config
	dialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}

	// todo move to app and use cfg
	trans := &http.Transport{
		DialContext:           dialer.DialContext,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
	}

	// todo move to app and use cfg
	client := &http.Client{
		Transport: trans,
		Timeout:   10 * time.Second,
	}

	// todo host to config
	service := services.New(client, "http://127.0.0.1:8080", app.repo, app.logger)

	cmd := commands.New(service, app.cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := os.Args[2:]

	var err error
	var msg string

	switch commands.CmdName(os.Args[1]) {
	// todo команду в константу
	case commands.Register:
		msg, err = cmd.Register(ctx, args)
	case commands.Login:
		msg, err = cmd.Login(ctx, args)
	case commands.InsertLoginPass:
		msg, err = cmd.CreateLoginpass(ctx, args)
	case commands.SyncSend:
		msg, err = cmd.SyncSend(ctx, args)
	case commands.List:
		msg, err = cmd.UserList(ctx, args)
	case commands.View:
		msg, err = cmd.View(ctx, args)
	case commands.SyncGet:
		msg, err = cmd.SyncGet(ctx, args)
	case commands.Delete:
		msg, err = cmd.DeleteItem(ctx, args)
	case commands.UpdateLoginPass:
		msg, err = cmd.UpdateLoginPass(ctx, args)
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
