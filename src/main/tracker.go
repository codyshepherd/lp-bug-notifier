package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

const (
	saveFileName = ".bugsave"
)

type Tracker struct {
	client *http.Client
	list   map[string]*BugDevel
	lock   *sync.Mutex
}

func NewTracker() *Tracker {
	log.Debug("New tracker created")
	var t *Tracker = new(Tracker)

	t.client = &http.Client{}
	t.list = make(map[string]*BugDevel)
	t.lock = &sync.Mutex{}

	return t
}

// Add a bug to be tracked and store its BugDevel response object
func (t *Tracker) Add(args []string) {
	if len(args) < 1 {
		log.Error("Empty arg list for Add()")
		fmt.Println("Usage: add <bug #>")
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
		log.Debug(fmt.Sprintf("Found: %s", obj.Title))

		log.Info(fmt.Sprintf("added: %s", bug))
		t.list[bug] = &obj
		t.Save()
	} else {
		log.Error("GET error: ", err)
		fmt.Println("Error: There was a problem adding that bug.")
	}
}

// Drop a bug from the tracker
func (t *Tracker) Drop(args []string) {
	if len(args) < 1 {
		log.Error("Empty arg list for Drop()")
		fmt.Println("Usage: drop <bug #>")
		return
	}
	if _, ok := t.list[args[0]]; ok {
		log.Info(fmt.Sprintf("dropped: %s :: %s", args[0]))
		delete(t.list, args[0])
		t.Save()
	} else {
		log.Error(fmt.Sprintf("Cannot drop: %s \nHas not been added!", args[0]))
		fmt.Println("Error: That bug has not been added")
	}
}

// Load save file into in-memory map
func (t *Tracker) Read() {
	bytes, err := ioutil.ReadFile(saveFileName)

	if err != nil {
		log.Error("Problem reading save file")
		log.Error(err)
		return
	}

	var obj map[string]*BugDevel
	err = json.Unmarshal(bytes, &obj)

	if err != nil {
		log.Error("Problem unmarshalling savefile data.")
		log.Error(err)
		return
	}

	t.list = obj
}

// Save off contents of map to a file
func (t *Tracker) Save() {
	out, err := json.Marshal(t.list)
	if err != nil {
		log.Error("Failed to marshal.")
		return
	}

	ioutil.WriteFile(saveFileName, out, 0644)
}
