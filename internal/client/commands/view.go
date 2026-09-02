package commands

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/spider4216/GophKeeper/internal/client/models"
	"github.com/spider4216/GophKeeper/internal/enum"
	shrModel "github.com/spider4216/GophKeeper/internal/model"
)

func (c *Command) View(ctx context.Context, args []string) (string, error) {
	fs := flag.NewFlagSet(View.String(), flag.ExitOnError)

	// todo validation
	itemID := fs.String("item_id", "", "Item ID")
	userID := fs.Int64("user-id", 0, "User ID from server")

	if err := fs.Parse(args); err != nil {
		return "", err
	}

	if *itemID == "" || *userID == 0 {
		return "", errors.New("jwt and item_id are required")
	}

	auth, err := c.Service.GetToken(ctx, *userID)
	if err != nil {
		return "", err
	}

	item, err := c.Service.GetUserItemByID(ctx, *itemID, auth.UserID)
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

	switch item.Type {
	case enum.LoginPass:
		res, errCmd = c.outLoginPass(decrypted, metas)
	case enum.Card:
		res, errCmd = c.outCard(decrypted, metas)
	default:
		res, errCmd = "", errors.New("undefined type")
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

	for _, v := range meta {
		fmt.Fprintf(&builder, "%d: %s: %s\n", v.ID, v.Key, v.Value)
	}

	return builder.String(), nil
}

func (c *Command) outCard(decrypted []byte, meta []shrModel.MetadataRepo) (string, error) {
	var data models.CardFmt

	if err := json.Unmarshal(decrypted, &data); err != nil {
		return "", fmt.Errorf("cannot unmaeshall card: %s", err)
	}

	var builder strings.Builder

	fmt.Fprintf(&builder, "PAN: %s\n", data.Pan)
	fmt.Fprintf(&builder, "CVC: %s\n", data.Cvc)
	fmt.Fprintf(&builder, "Holder: %s\n", data.Holder)
	fmt.Fprintf(&builder, "Date: %s\n", data.Date.Format("02.01.2006"))
	fmt.Fprint(&builder, "Meta:\n")

	for _, v := range meta {
		fmt.Fprintf(&builder, "%d: %s: %s\n", v.ID, v.Key, v.Value)
	}

	return builder.String(), nil
}
