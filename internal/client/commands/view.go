package commands

import (
	"context"
	"errors"
	"flag"
	"fmt"
)

func (c *Command) View(args []string) (string, error) {
	fs := flag.NewFlagSet(View.String(), flag.ExitOnError)

	// validation
	itemID := fs.Int64("item_id", 0, "Item ID")
	token := fs.String("token", "", "JWT from server")

	fs.Parse(args)

	if *itemID == 0 && *token == "" {
		return "", errors.New("jwt and item_id are required")
	}

	claims, err := c.getClaims(*token)
	if err != nil {
		return "", fmt.Errorf("cannot parse jwt: %s", err)
	}

	// todo подумать как сдесь релизовать middleware
	ctx := context.WithValue(context.Background(), "userID", claims.UserID)

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

	// todo в зависимости от типа нужно делать unmarshall в конкретную структуру

	// var decData models.LoginPassFmt

	// if err := json.Unmarshal(decrypted, &decData); err != nil {
	// 	fmt.Printf("cannot unmarshal decrypted data: %s\n", err)
	// 	os.Exit(1)
	// }

	return string(decrypted), nil
}
