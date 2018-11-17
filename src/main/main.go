package main

import (
    "net/http"
    "github.com/codyshepherd/lp-bug-notifier/src/term"
    "github.com/codyshepherd/lp-bug-notifier/src/tracker"

    log "github.com/sirupsen/logrus"
)

func main() {
    log.Info("LP Bug Notifier started");

    term := NewTerm(NewTracker())
    term.Run()
}
