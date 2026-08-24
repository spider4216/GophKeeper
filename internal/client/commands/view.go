package commands

import (
	"context"
	"flag"
	"fmt"
	"os"
)

// todo return err
func (c *Command) View(args []string) {
	// todo  command name to constant
	fs := flag.NewFlagSet("view", flag.ExitOnError)

	// validation
	itemID := fs.Int64("item_id", 0, "Item ID")
	token := fs.String("token", "", "JWT from server")

	fs.Parse(args)

	if *itemID == 0 && *token == "" {
		// todo instead fmt use something with out source
		fmt.Println("jwt and item_id are required")
		os.Exit(1)
	}

	claims, err := c.GetClaims(*token)
	if err != nil {
		fmt.Printf("cannot parse jwt: %s\n", err)
		os.Exit(1)
	}

	// todo подумать как сдесь релизовать middleware
	ctx := context.WithValue(context.Background(), "userID", claims.UserID)

	if err != nil {
		fmt.Printf("cannot convert item ID to int: %s\n", err)
		os.Exit(1)
	}

	item, err := c.Service.GetUserItemByID(ctx, *itemID, claims.UserID)

	if err != nil {
		fmt.Printf("cannot get user item by id: %s\n", err)
		os.Exit(1)
	}

	decrypted, err := c.Service.DecryptData(item.Ciphertext, []byte(c.Cfg.EncryptKey))

	if err != nil {
		fmt.Printf("cannot decrypt data: %s\n", err)
		os.Exit(1)
	}

	// todo в зависимости от типа нужно делать unmarshall в конкретную структуру

	// var decData models.LoginPassFmt

	// if err := json.Unmarshal(decrypted, &decData); err != nil {
	// 	fmt.Printf("cannot unmarshal decrypted data: %s\n", err)
	// 	os.Exit(1)
	// }

	fmt.Println(string(decrypted))
}
