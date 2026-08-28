package commands

import (
	"fmt"

	"github.com/spider4216/GophKeeper/internal/client/version"
)

func (c *Command) PrintVersion() (string, error) {
	return fmt.Sprintf("version %s, build date %s", version.Version, version.BuildDate), nil
}
