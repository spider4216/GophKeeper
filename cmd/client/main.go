package main

import (
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

	args := os.Args[2:]

	var err error
	var msg string

	switch os.Args[1] {
	// todo команду в константу
	case "register":
		msg, err = cmd.Register(args)
	case "login":
		msg, err = cmd.Login(args)
	case "insert-loginpass":
		msg, err = cmd.CreateLoginpass(args)
	case "sync-send":
		msg, err = cmd.SyncSend(args)
	case "list":
		msg, err = cmd.UserList(args)
	case "view":
		msg, err = cmd.View(args)
	case "sync-get":
		msg, err = cmd.SyncGet(args)
	case "delete-item":
		msg, err = cmd.DeleteItem(args)
	case "update-loginpass":
		msg, err = cmd.UpdateLoginPass(args)
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
