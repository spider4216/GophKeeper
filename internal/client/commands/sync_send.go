package commands

import (
	"context"
	"errors"
	"flag"
	"fmt"
)

func (c *Command) SyncSend(ctx context.Context, args []string) (string, error) {
	fs := flag.NewFlagSet(SyncSend.String(), flag.ExitOnError)

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

	if err := c.Service.SyncSend(ctx, auth.UserID, auth.Token, c.Cfg.SyncChankSize); err != nil {
		return "", fmt.Errorf("cannot sync: %s", err)
	}

	return "Sync was successfully done", nil
}
