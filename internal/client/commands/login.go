package commands

import (
	"flag"
	"fmt"
	"os"

	"github.com/spider4216/GophKeeper/internal/client/models"
)

// todo return err
func (c *Command) Login(args []string) {
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

	req := models.LoginReq{
		Email: *email,
		Pass:  *pass,
	}

	resp, err := c.Service.Login(req)

	if err != nil {
		fmt.Printf("cannot create user: %s", err)
		os.Exit(1)
	}

	fmt.Printf("JWT token\n\n%s\n", resp.Token)
	// todo вывести когда истекает в человекопонятном виде
}
