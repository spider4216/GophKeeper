package commands

import (
	"context"
	"errors"
	"flag"
	"fmt"
)

func (c *Command) SyncGet(ctx context.Context, args []string) (string, error) {
	fs := flag.NewFlagSet(SyncGet.String(), flag.ExitOnError)

	userID := fs.Int64("user-id", 0, "User ID")

	if err := fs.Parse(args); err != nil {
		return "", err
	}

	if *userID == 0 {
		return "", errors.New("token is required")
	}

	auth, err := c.Service.GetToken(ctx, *userID)

	if err != nil {
		return "", err
	}

	if err := c.Service.SyncGet(ctx, auth.UserID, auth.Token, c.Cfg.SyncLimit); err != nil {
		return "", fmt.Errorf("cannot sync: %s", err)
	}

	return "Sync was successfully done", nil
}
