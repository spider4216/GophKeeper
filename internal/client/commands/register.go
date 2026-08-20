package commands

import (
	"flag"
	"fmt"
	"os"

	"github.com/spider4216/GophKeeper/internal/client/models"
)

// todo return err
func (c *Command) Register(args []string) {
	// todo  command name to constant
	fs := flag.NewFlagSet("register", flag.ExitOnError)

	email := fs.String("email", "", "user email")
	pass := fs.String("password", "", "user password")

	fs.Parse(args)

	if *email == "" || *pass == "" {
		// todo instead fmt use something with out source
		fmt.Println("email and password are required")
		os.Exit(1)
	}

	req := models.RegisterReq{
		Email: *email,
		Pass:  *pass,
	}

	if err := c.Service.Register(req); err != nil {
		fmt.Printf("cannot create user: %s", err)
		os.Exit(1)
	}

	fmt.Println("Register is OK")
}
