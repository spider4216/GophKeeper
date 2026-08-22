package main

import (
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

	switch os.Args[1] {
	// todo команду в константу
	case "register":
		cmd.Register(args)
	case "login":
		cmd.Login(args)
	case "insert-loginpass":
		cmd.CreateItem(args)
	case "sync-send":
		cmd.SyncSend(args)
	case "list":
		cmd.UserList(args)
	case "view":
		cmd.View(args)
	case "sync-get":
		cmd.SyncGet(args)
	default:
		fmt.Println("command not found")
		os.Exit(1)
	}
}
