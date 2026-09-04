package commands

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
)

func (c *Command) GetText(ctx context.Context, args []string) (string, error) {
	fs := flag.NewFlagSet(GetText.String(), flag.ExitOnError)

	userID := fs.Int64("user-id", 0, "User ID")
	itemID := fs.String("item-id", "", "Item ID")

	if err := fs.Parse(args); err != nil {
		return "", err
	}

	if *itemID == "" || *userID == 0 {
		return "", errors.New("user-id and item-id are required")
	}

	auth, err := c.Service.GetToken(ctx, *userID)
	if err != nil {
		return "", err
	}

	if err := c.Service.GetBinaryData(ctx, auth.UserID, *itemID, c.Cfg.EncryptKey, os.Stdout); err != nil {
		return "", fmt.Errorf("cannot get huge text: %w", err)
	}

	return "Success", nil
}
