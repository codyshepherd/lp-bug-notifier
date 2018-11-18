package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type Term struct {
	recv    chan string
	reader  *bufio.Reader
	tracker *Tracker
}

var ops = map[string]func(*Tracker, []string){
	"add": func(t *Tracker, args []string) {
		t.lock.Lock()
		t.Add(args)
		t.lock.Unlock()
	},
	"drop": func(t *Tracker, args []string) {
		t.lock.Lock()
		t.Drop(args)
		t.lock.Unlock()
	},
	"ls": func(t *Tracker, args []string) {
		fmt.Println()
		t.lock.Lock()
		for k, v := range t.list {
			tm, err := time.Parse(time.RFC3339, v.Date_last_message)

			if err != nil {
				fmt.Println(fmt.Sprintf("%s: %s", k, v.Title))
			} else {
				fmt.Println(fmt.Sprintf("%s [%s]: %s", k,
					time.Since(tm).Truncate(time.Minute).String(), v.Title))
			}
		}
		t.lock.Unlock()
		fmt.Println()
	},
	"refresh": func(t *Tracker, args []string) {
		for k := range t.list {
			t.lock.Lock()
			t.Add([]string{k})
			t.lock.Unlock()
		}
	},
	"quit": func(t *Tracker, args []string) {
		fmt.Println("quitting...")
	},
	"exit": func(t *Tracker, args []string) {
		fmt.Println("quitting...")
	},
}

func NewTerm() *Term {
	var t *Term = new(Term)
	t.recv = make(chan string, 1024)     // maybe i'll add more channels eventually
	t.reader = bufio.NewReader(os.Stdin) // get input from user
	t.tracker = NewTracker()             // in-memory map for storing info
	t.tracker.Read()                     // load any info saved to file
	return t
}

func (t *Term) Run() {
	go t.handle()      // spin off serialized command handler
	go t.checkUpdate() // spin off periodic updater

	// here we hand commands to handler via channel
	// this allows the handler to take commands from elsewhere, and we get
	// the synchronization for free with channels
	for {
		fmt.Print("> ")
		input, _ := t.reader.ReadString('\n')
		t.recv <- input
		time.Sleep(100 * time.Millisecond) // QoL delay
		if strings.HasPrefix(input, "quit") ||
			strings.HasPrefix(input, "exit") {
			t.tracker.Save()
			break
		}
	}
}

// Refresh once every five minutes
func (t *Term) checkUpdate() {
	for {
		for k := range t.tracker.list {
			t.tracker.Add([]string{k})
		}
		fmt.Println()
		ops["ls"](t.tracker, []string{})
		fmt.Print("> ")

		time.Sleep(5 * time.Minute)
	}
}

// Take commands out of the channel and handle them
func (t *Term) handle() {
	for {
		select {
		case msg := <-t.recv:
			words := strings.Split(msg, " ")
			op := strings.TrimRight(words[0], " \n")
			args := words[1:]

			if f, ok := ops[op]; ok {
				f(t.tracker, args)

				if strings.HasPrefix(op, "quit") ||
					strings.HasPrefix(op, "exit") {
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
