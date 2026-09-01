package commands

import (
	"context"
	"errors"
	"flag"
)

func (c *Command) CreateText(ctx context.Context, args []string) (string, error) {
	fs := flag.NewFlagSet(CreateText.String(), flag.ExitOnError)

	title := fs.String("title", "", "Title")
	userID := fs.Int64("user-id", 0, "User ID")
	filePath := fs.String("path", "", "File path")

	if err := fs.Parse(args); err != nil {
		return "", err
	}

	if *title == "" || *userID == 0 || *filePath == "" {
		return "", errors.New("title, user-id and file path are required")
	}

	auth, err := c.Service.GetToken(ctx, *userID)
	if err != nil {
		return "", err
	}

	if err := c.Service.SaveHugeText(ctx, auth.UserID, *title, *filePath, c.Cfg.EncryptKey); err != nil {
		return "", err
	}

	return "Text was successfully saved", nil
}
