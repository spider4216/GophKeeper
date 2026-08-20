package commands

import (
	"flag"
	"fmt"
	"os"
)

func (c *Command) Register(args []string) {
	// todo  command name to constant
	fs := flag.NewFlagSet("register", flag.ExitOnError)

	email := fs.String("email", "", "user email")
	pass := fs.String("password", "", "user password")

	fs.Parse(args)

	if *email == "" || *pass == "" {
		fmt.Println("email and password are required")
		os.Exit(1)
	}

	fmt.Println("Register logic")
}
