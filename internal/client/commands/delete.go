package commands

import (
	"context"
	"errors"
	"flag"
	"fmt"
)

func (c *Command) DeleteItem(ctx context.Context, args []string) (string, error) {
	fs := flag.NewFlagSet(Delete.String(), flag.ExitOnError)

	itemID := fs.String("id", "", "Item id")
	userID := fs.Int64("user-id", 0, "User ID")

	if err := fs.Parse(args); err != nil {
		return "", err
	}

	if *userID == 0 || *itemID == "" {
		return "", errors.New("userID and itemID are required")
	}

	auth, err := c.Service.GetToken(ctx, *userID)

	if err != nil {
		return "", err
	}

	if err := c.Service.DeleteUserItem(ctx, *itemID, auth.UserID); err != nil {
		return "", fmt.Errorf("cannot delete item: %w", err)
	}

	return "Item successfully deleted", nil
}
