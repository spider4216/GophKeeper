package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/spider4216/GophKeeper/internal/client/commands"
	"github.com/spider4216/GophKeeper/internal/client/services"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("command is required")
		os.Exit(1)
	}

	// todo move to app
	// todo use config
	dialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}

	trans := &http.Transport{
		DialContext:           dialer.DialContext,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
	}

	client := &http.Client{
		Transport: trans,
		Timeout:   10 * time.Second,
	}

	// todo host to config
	service := services.New(client, "http://127.0.0.1:8080")

	cmd := commands.New(service)

	args := os.Args[2:]

	switch os.Args[1] {
	// todo команду в константу
	case "register":
		cmd.Register(args)
	case "login":
		cmd.Login(args)
	default:
		fmt.Println("command not found")
		os.Exit(1)
	}
}
