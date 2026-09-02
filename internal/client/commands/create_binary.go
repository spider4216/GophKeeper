package commands

import (
	"context"
	"errors"
	"flag"
	"path/filepath"
)

func (c *Command) CreateBinary(ctx context.Context, args []string) (string, error) {
	fs := flag.NewFlagSet(CreateBinary.String(), flag.ExitOnError)

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

	meta := map[string]string{
		"Title":       *title,
		"ContentType": "binary",
		MetaFileName:  filepath.Base(*filePath),
	}

	if err := c.Service.SaveBinary(ctx, auth.UserID, meta, *filePath, c.Cfg.EncryptKey, c.Cfg.PartSize); err != nil {
		return "", err
	}

	return "Binary was successfully saved", nil
}
