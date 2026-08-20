package main

import (
	"fmt"
	"os"

	"github.com/spider4216/GophKeeper/internal/client/commands"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("command is required")
		os.Exit(1)
	}

	cmd := commands.New()

	args := os.Args[2:]

	switch os.Args[1] {
	case "register":
		cmd.Register(args)
	default:
		fmt.Println("command not found")
		os.Exit(1)
	}
}
