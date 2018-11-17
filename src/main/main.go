package main

import (
	//"net/http"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("LP Bug Notifier started")

	var term *Term = NewTerm(NewTracker())
	term.Run()
}
