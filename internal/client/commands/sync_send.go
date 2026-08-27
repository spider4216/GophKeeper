package commands

import (
	"context"
	"errors"
	"flag"
	"fmt"
)

func (c *Command) SyncSend(ctx context.Context, args []string) (string, error) {
	fs := flag.NewFlagSet(SyncSend.String(), flag.ExitOnError)

	token := fs.String("token", "", "Token")

	fs.Parse(args)

	if *token == "" {
		return "", errors.New("Token is required")
	}

	claims, err := c.getClaims(*token)
	if err != nil {
		return "", fmt.Errorf("cannot parse jwt: %s", err)
	}

	if err := c.Service.SyncSend(ctx, claims.UserID, *token); err != nil {
		return "", fmt.Errorf("cannot sync: %s", err)
	}

	return "Sync was successfully done", nil
}
