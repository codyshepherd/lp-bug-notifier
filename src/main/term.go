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
			fmt.Println(fmt.Sprintf("%s: %s", k, v.Title))
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
	return t
}

func (t *Term) Run() {
	go t.handle()

	for {
		fmt.Print("> ")
		input, _ := t.reader.ReadString('\n')
		t.recv <- input
		time.Sleep(100 * time.Millisecond)
		if strings.HasPrefix(input, "quit") {
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
