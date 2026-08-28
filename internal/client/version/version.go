package version

import "fmt"

var (
	Version   = "dev"
	BuildDate = "unknown"
)

func String() string {
	return fmt.Sprintf("version %s, build date %s", Version, BuildDate)
}
