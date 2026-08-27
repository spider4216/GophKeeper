package commands

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/spider4216/GophKeeper/internal/client/models"
	"github.com/spider4216/GophKeeper/internal/enum"
	shrModel "github.com/spider4216/GophKeeper/internal/model"
)

func (c *Command) View(ctx context.Context, args []string) (string, error) {
	fs := flag.NewFlagSet(View.String(), flag.ExitOnError)

	// todo validation
	itemID := fs.Int64("item_id", 0, "Item ID")
	token := fs.String("token", "", "JWT from server")

	if err := fs.Parse(args); err != nil {
		return "", err
	}

	if *itemID == 0 && *token == "" {
		return "", errors.New("jwt and item_id are required")
	}

	claims, err := c.getClaims(*token)
	if err != nil {
		return "", fmt.Errorf("cannot parse jwt: %s", err)
	}

	if err != nil {
		return "", fmt.Errorf("cannot convert item ID to int: %s", err)
	}

	item, err := c.Service.GetUserItemByID(ctx, *itemID, claims.UserID)
	if err != nil {
		return "", fmt.Errorf("cannot get user item by id: %s", err)
	}

	decrypted, err := c.Service.DecryptData(item.Ciphertext, []byte(c.Cfg.EncryptKey))
	if err != nil {
		return "", fmt.Errorf("cannot decrypt data: %s", err)
	}

	metas, err := c.Service.GetMetadataByItemID(ctx, item.ID)
	if err != nil {
		return "", fmt.Errorf("cannot get metadata: %s", err)
	}

	var res string
	var errCmd error

	if item.Type == enum.LoginPass {
		res, errCmd = c.outLoginPass(decrypted, metas)
	}

	if errCmd != nil {
		return "", errCmd
	}

	return res, nil
}

func (c *Command) outLoginPass(decrypted []byte, meta []shrModel.MetadataRepo) (string, error) {
	var data models.LoginPassFmt

	if err := json.Unmarshal(decrypted, &data); err != nil {
		return "", fmt.Errorf("cannot unmaeshall loginpass: %s", err)
	}

	var builder strings.Builder

	fmt.Fprintf(&builder, "Login: %s\n", data.Login)
	fmt.Fprintf(&builder, "Password: %s\n", data.Pass)
	fmt.Fprint(&builder, "Meta:\n")

	log.Println(meta)

	for _, v := range meta {
		fmt.Fprintf(&builder, "%d: %s: %s\n", v.ID, v.Key, v.Value)
	}

	return builder.String(), nil
}
