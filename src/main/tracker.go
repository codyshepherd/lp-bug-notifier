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
	list   map[string]string
}

func NewTracker() *Tracker {
	log.Debug("New tracker created")
	var t *Tracker = new(Tracker)

	t.client = &http.Client{}
	t.list = make(map[string]string)

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
		"https://api.launchpad.net/devel/bugs", nil)
	q := req.URL.Query()
	q.Add("ws.op", "getBugData")
	q.Add("bug_id", bug)
	req.URL.RawQuery = q.Encode()
	//req.Header.Set("bug_id", bug)
	resp, err := t.client.Do(req)

	log.Info("resp content length: ", resp.ContentLength)

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	var obj []BugDevel

	json.Unmarshal(bodyBytes, &obj)

	log.Info("len obj: ", len(obj))
	log.Info(obj[0].title)

	if err == nil {
		log.Info(fmt.Sprintf("added: %s", bug))
		t.list[bug] = "info here"
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
