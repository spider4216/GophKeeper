package main

import "log"

func main() {
	app := newApp()

	if err := app.Run(); err != nil {
		log.Fatal("Cannot run app", err)
	}

	app.logger.Debug("Everything is ok")
}
