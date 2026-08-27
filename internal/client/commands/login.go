package commands

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"time"

	shrModel "github.com/spider4216/GophKeeper/internal/model"
)

func (c *Command) Login(ctx context.Context, args []string) (string, error) {
	fs := flag.NewFlagSet(Login.String(), flag.ExitOnError)

	email := fs.String("email", "", "user email")
	pass := fs.String("password", "", "user password")

	if err := fs.Parse(args); err != nil {
		return "", err
	}

	if *email == "" || *pass == "" {
		return "", errors.New("email and password are required")
	}

	req := shrModel.LoginReq{
		Email: *email,
		Pass:  *pass,
	}

	resp, err := c.Service.Login(ctx, req)
	if err != nil {
		return "", fmt.Errorf("cannot login: %w", err)
	}

	// логин можен осуществлять с другого клиента у которого пустая БД
	// Если это первый вход клиента с нового устройства, то нужно
	// Внести ему последнюю ревизию как 0
	if _, err := c.Service.GetLatestUserRev(ctx, resp.UserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.logger.Debug("no last rev. Create zero one")
			if err := c.Service.CreateLastUserRev(ctx, resp.UserID, 0); err != nil {
				return "", fmt.Errorf("cannot create last rev: %w", err)
			}
		} else {
			return "", fmt.Errorf("cannot get latest revision: %w", err)
		}
	}

	expObj := time.Unix(resp.ExpiredAt, 0)
	expStr := expObj.Format("2006-01-02 15:04:05")

	return fmt.Sprintf("JWT token\n\n%s\n\nExpired: %s\n", resp.Token, expStr), nil
}
