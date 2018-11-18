package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	//log "github.com/sirupsen/logrus"
)

type Term struct {
	recv    chan string
	reader  *bufio.Reader
	tracker *Tracker
}

var ops = map[string]func(*Tracker, []string){
	"add": func(t *Tracker, args []string) {
		t.Add(args)
	},
	"drop": func(t *Tracker, args []string) {
		t.Drop(args)
	},
	"ls": func(t *Tracker, args []string) {
		for k, v := range t.list {
			tm, err := time.Parse(time.RFC3339, v.Date_last_message)

			if err != nil {
				fmt.Println(fmt.Sprintf("%s: %s", k, v.Title))
			} else {
				fmt.Println(fmt.Sprintf("%s: %s [%s]", k, v.Title,
					time.Since(tm).Truncate(time.Minute).String()))
			}
		}
	},
	"refresh": func(t *Tracker, args []string) {
		for k := range t.list {
			t.Add([]string{k})
		}
	},
	"quit": func(t *Tracker, args []string) {
		fmt.Println("quitting...")
	},
	"exit": func(t *Tracker, args []string) {
		fmt.Println("quitting...")
	},
}

func NewTerm(tracker *Tracker) *Term {
	var t *Term = new(Term)
	t.recv = make(chan string, 1024)
	t.reader = bufio.NewReader(os.Stdin)
	t.tracker = NewTracker()
	t.tracker.Read()
	return t
}

func (t *Term) Run() {
	go t.handle()

	lastCheck := time.Now()

	for {
		if time.Since(lastCheck).Minutes() > 5.0 {
			for k := range t.tracker.list {
				t.tracker.Add([]string{k})
			}
			lastCheck = time.Now()
		}

		fmt.Print("> ")
		input, _ := t.reader.ReadString('\n')
		t.recv <- input
		time.Sleep(100 * time.Millisecond)
		if strings.HasPrefix(input, "quit") {
			t.tracker.Save()
			break
		}
	}
}

func (t *Term) handle() {
	for {
		select {
		case msg := <-t.recv:
			words := strings.Split(msg, " ")
			op := strings.TrimRight(words[0], " \n")
			args := words[1:]

			if f, ok := ops[op]; ok {
				f(t.tracker, args)

				if strings.HasPrefix(op, "quit") {
					break
				}
			} else {
				fmt.Println("Available commands:")
				for k := range ops {
					fmt.Println(k)
				}
			}
		}
	}

}
