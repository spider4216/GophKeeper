package commands

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"

	shrModel "github.com/spider4216/GophKeeper/internal/model"
)

// todo return err
func (c *Command) Login(args []string) (string, error) {
	// todo  command name to constant
	fs := flag.NewFlagSet(Login.String(), flag.ExitOnError)

	email := fs.String("email", "", "user email")
	pass := fs.String("password", "", "user password")

	fs.Parse(args)

	if *email == "" || *pass == "" {
		return "", errors.New("email and password are required")
	}

	req := shrModel.LoginReq{
		Email: *email,
		Pass:  *pass,
	}

	// todo thenk aboun ctx in commands
	ctx := context.Background()

	resp, err := c.Service.Login(req)

	if err != nil {
		return "", fmt.Errorf("cannot login: %s", err)
	}

	// todo логин можен осуществлять с другого клиента у которого пустая БД
	// Если это первый вход клиента с нового устройства, то нужно
	// Внести ему последнюю ревизию как 0
	if _, err := c.Service.GetLatestUserRev(ctx, resp.UserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// todo logger
			c.Service.CreateLastUserRev(ctx, resp.UserID, 0)
		} else {
			return "", fmt.Errorf("cannot get latest revision: %s", err)
		}
	}

	return fmt.Sprintf("JWT token\n\n%s\n", resp.Token), nil
	// todo вывести когда истекает в человекопонятном виде
}
