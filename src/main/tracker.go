package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Tracker struct {
	client *http.Client
	list   map[string]*BugDevel
}

func NewTracker() *Tracker {
	log.Debug("New tracker created")
	var t *Tracker = new(Tracker)

	t.client = &http.Client{}
	t.list = make(map[string]*BugDevel)

	return t
}

func (t *Tracker) Add(args []string) {
	if len(args) < 1 {
		log.Error(fmt.Sprintf("Usage: add <bug #>",
			strings.Join(args, " ")))
		return
	}

	bug := strings.TrimRight(args[0], " \n")

	req, _ := http.NewRequest("GET",
		"https://api.launchpad.net/devel/bugs/"+bug, nil)
	resp, err := t.client.Do(req)

	if err == nil {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		var obj BugDevel
		json.Unmarshal(bodyBytes, &obj)
		log.Info(fmt.Sprintf("Found: %s", obj.Title))

		log.Info(fmt.Sprintf("added: %s", bug))
		t.list[bug] = &obj
	} else {
		log.Error("GET error: ", err)
	}
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
