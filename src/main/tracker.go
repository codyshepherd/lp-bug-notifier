package main

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Tracker struct {
	list map[string]string
}

func NewTracker() *Tracker {
	log.Debug("New tracker created")
	var t *Tracker = new(Tracker)

	t.list = make(map[string]string)

	return t
}

func (t *Tracker) Add(args []string) {
	if len(args) < 1 {
		log.Error(fmt.Sprintf("Usage: add <bug #>",
			strings.Join(args, " ")))
		return
	}
	log.Debug(fmt.Sprintf("add: %s", args[0]))
	t.list[args[0]] = "url here"
}

func (t *Tracker) Drop(args []string) {
	if len(args) < 1 {
		log.Error(fmt.Sprintf("Usage: drop <bug #>",
			strings.Join(args, " ")))
		return
	}
	if _, ok := t.list[args[0]]; ok {
		log.Debug(fmt.Sprintf("drop: %s :: %s", args[0]))
		delete(t.list, args[0])
	} else {
		log.Error(fmt.Sprintf("Cannot drop: %s \nDoes not exist!", args[0]))
	}
}
