package commands

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/spider4216/GophKeeper/internal/client/models"

	"github.com/golang-jwt/jwt/v4"
)

// todo return err
func (c *Command) CreateItem(args []string) {
	// todo  command name to constant
	fs := flag.NewFlagSet("create-item", flag.ExitOnError)

	login := fs.String("login", "", "Login")
	pass := fs.String("password", "", "Password")
	token := fs.String("token", "", "JWT from server")
	title := fs.String("title", "", "Title")

	fs.Parse(args)

	if *login == "" || *pass == "" || *token == "" || *title == "" {
		// todo instead fmt use something with out source
		fmt.Println("login, password, title and jwt are required")
		os.Exit(1)
	}

	req := models.LoginPassReq{
		Login: *login,
		Pass:  *pass,
		Title: *title,
		JWT:   *token,
	}

	// todo извлечь из токена user id

	claims := &claims{}
	_, err := jwt.ParseWithClaims(req.JWT, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(c.Cfg.JWTKey), nil
		})

	if err != nil {
		fmt.Printf("cannot parse jwt: %s", err)
		os.Exit(1)
	}

	// todo подумать как сдесь релизовать middleware
	ctx := context.WithValue(context.Background(), "userID", claims.UserID)

	// todo transaction
	// todo подумать что сделать с контекстом
	// todo тип в коснтанту
	itemID, err := c.Service.CreateItem(ctx, "login_pass", req, c.Cfg.EncryptKey, claims.UserID)

	if err != nil {
		fmt.Printf("cannot create item: %s", err)
		os.Exit(1)
	}

	_, err = c.Service.CreateMeta(ctx, itemID, "Title", req.Title)

	if err != nil {
		fmt.Printf("cannot create meta for item: %s", err)
		os.Exit(1)
	}

	// op to const and custom type
	err = c.Service.CreatePendingChange(ctx, itemID, "CREATE")

	if err != nil {
		fmt.Printf("cannot create pending for item: %s", err)
		os.Exit(1)
	}

	fmt.Println("Item successfully created")
}
