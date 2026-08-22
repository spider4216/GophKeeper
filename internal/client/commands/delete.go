package commands

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
)

// todo return err
func (c *Command) DeleteItem(args []string) {
	// todo  command name to constant
	fs := flag.NewFlagSet("delete-item", flag.ExitOnError)

	itemID := fs.String("id", "", "Item id")
	token := fs.String("token", "", "JWT from server")

	fs.Parse(args)

	if *token == "" || *itemID == "" {
		// todo instead fmt use something with out source
		fmt.Println("login, password, title and jwt are required")
		os.Exit(1)
	}

	claims := &claims{}
	_, err := jwt.ParseWithClaims(*token, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(c.Cfg.JWTKey), nil
		})

	if err != nil {
		fmt.Printf("cannot parse jwt: %s", err)
		os.Exit(1)
	}

	// todo подумать как сдесь релизовать middleware
	ctx := context.WithValue(context.Background(), "userID", claims.UserID)

	itemIDInt, err := strconv.ParseInt(*itemID, 10, 64)

	if err != nil {
		fmt.Printf("cannot parse itemID: %s", err)
		os.Exit(1)
	}

	if err := c.Service.DeleteUserItem(ctx, itemIDInt, claims.UserID); err != nil {
		fmt.Printf("cannot delete item: %s", err)
		os.Exit(1)
	}

	fmt.Println("Item successfully deleted")
}
