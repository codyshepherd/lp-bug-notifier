package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	f, err := os.OpenFile(".buglog", os.O_WRONLY|os.O_CREATE,
		0666)
	if err != nil {
		panic(fmt.Sprintf("Failed to open file for logging: %s", err))
	}

	log.SetOutput(f)
	log.SetLevel(log.DebugLevel)
	log.Info("LP Bug Notifier started")

	// var term *Term = NewTerm()
	// term.Run()

	var ui *Ui = NewUi()
	ui.Run()
}
