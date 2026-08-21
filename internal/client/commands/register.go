package commands

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/spider4216/GophKeeper/internal/client/models"
)

// todo return err
func (c *Command) Register(args []string) {
	// todo  command name to constant
	// todo transaction
	fs := flag.NewFlagSet("register", flag.ExitOnError)

	email := fs.String("email", "", "user email")
	pass := fs.String("password", "", "user password")

	fs.Parse(args)

	if *email == "" || *pass == "" {
		// todo instead fmt use something with out source
		fmt.Println("email and password are required")
		os.Exit(1)
	}

	// todo подумать что то с контекстом
	ctx := context.Background()

	req := models.RegisterReq{
		Email: *email,
		Pass:  *pass,
	}

	// todo condext for register user
	resp, err := c.Service.Register(req)

	if err != nil {
		fmt.Printf("cannot create user: %s", err)
		os.Exit(1)
	}

	// todo command logger

	if err := c.Service.CreateLastUserRev(ctx, resp.UserID, 0); err != nil {
		fmt.Printf("cannot create latest revision for user: %s", err)
		os.Exit(1)
	}

	fmt.Println("Register is OK")
}
