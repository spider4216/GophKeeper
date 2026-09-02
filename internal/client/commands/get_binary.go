package commands

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func (c *Command) DownloadBinary(ctx context.Context, args []string) (string, error) {
	fs := flag.NewFlagSet(DownloadBinary.String(), flag.ExitOnError)

	userID := fs.Int64("user-id", 0, "User ID")
	itemID := fs.String("item-id", "", "Item ID")
	out := fs.String("out", "", "Path out")

	if err := fs.Parse(args); err != nil {
		return "", err
	}

	if *itemID == "" || *userID == 0 || *out == "" {
		return "", errors.New("user-id, item-id and out are required")
	}

	auth, err := c.Service.GetToken(ctx, *userID)
	if err != nil {
		return "", err
	}

	metadata, err := c.Service.GetMetadataByItemID(ctx, *itemID)

	if err != nil {
		return "", fmt.Errorf("cannot get item id: %w", err)
	}

	metamap := make(map[string]string)

	for _, m := range metadata {
		metamap[m.Key] = m.Value
	}

	metaName, ok := metamap[MetaFileName]

	if !ok {
		return "", errors.New("cannot get filename from metadata")
	}

	extPath := filepath.Join(*out, metaName)

	file, err := os.Create(extPath)

	if err != nil {
		return "", fmt.Errorf("cannot create out file: %w", err)
	}

	if err := c.Service.GetBinaryData(ctx, auth.UserID, *itemID, c.Cfg.EncryptKey, file); err != nil {
		return "", fmt.Errorf("cannot get huge text: %w", err)
	}

	return fmt.Sprintf("Saved to %s", extPath), nil
}
