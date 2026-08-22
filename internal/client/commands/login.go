package commands

import (
	"context"
	"database/sql"
	"errors"
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

	// todo thenk aboun ctx in commands
	ctx := context.Background()

	resp, err := c.Service.Login(req)

	if err != nil {
		fmt.Printf("cannot login: %s", err)
		os.Exit(1)
	}

	// todo логин можен осуществлять с другого клиента у которого пустая БД
	// Если это первый вход клиента с нового устройства, то нужно
	// Внести ему последнюю ревизию как 0
	if _, err := c.Service.GetLatestUserRev(ctx, resp.UserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// todo create latest rev as 0
			// todo logger
			c.Service.CreateLastUserRev(ctx, resp.UserID, 0)
		} else {
			fmt.Printf("cannot get latest revision: %s", err)
			os.Exit(1)
		}
	}

	fmt.Printf("JWT token\n\n%s\n", resp.Token)
	// todo вывести когда истекает в человекопонятном виде
}
