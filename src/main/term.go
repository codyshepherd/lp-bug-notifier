package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

//////////////////////////////////////////////////////////////////////////////
// byTime allows us to sort strings by regex
type byTime []string
func (s byTime) Len() int {
	return len(s)
}
func (s byTime) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byTime) Less(i, j int) bool {
	r := regexp.MustCompile(`\[[0-9]*`)
	string_i := strings.Trim(r.FindString(s[i]), "[")
	string_j := strings.Trim(r.FindString(s[j]), "[")
	int_i, erri := strconv.Atoi(string_i)
	int_j, errj := strconv.Atoi(string_j)
	if erri != nil || errj != nil {
		panic("Problem with converting string to int")
	}
	return int_i < int_j
}
//////////////////////////////////////////////////////////////////////////////

// ops defines the behavior of available terminal commands
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
		buffer := []string{}
		t.lock.Lock()
		for k, v := range t.list {
			tm, err := time.Parse(time.RFC3339, v.BugStruct.Date_last_message)

            // check for updated-ness and prepend string if appropriate
            // note that once "Updated!" is displayed, the Changed flag on the
            // bug will be turned off
            updated := ""
            if v.Changed {
                updated = "**Updated!** "
                v.Changed = false
            }

            // check for time conversion failure
			if err != nil {
				buffer = append(buffer, fmt.Sprintf("%s%s: %s", updated, k,
                v.BugStruct.Title))
			} else {
				buffer = append(buffer, fmt.Sprintf("%s%s [%s]: %s", updated,
                k, time.Since(tm).Truncate(time.Minute).String(), v.BugStruct.Title))
			}
		}
		t.lock.Unlock()

		fmt.Println()

        // sort ascending by time since last update
		sort.Sort(byTime(buffer))
		for i := range buffer {
			fmt.Println(buffer[i])
		}

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

// Term handles input from user and display
type Term struct {
	recv    chan string
	reader  *bufio.Reader
	tracker *Tracker
}

// initializes a new Term object
func NewTerm() *Term {
	var t *Term = new(Term)
	t.recv = make(chan string, 1024)     // maybe i'll add more channels eventually
	t.reader = bufio.NewReader(os.Stdin) // get input from user
	t.tracker = NewTracker()             // in-memory map for storing info
	t.tracker.Read()                     // load any info saved to file
	return t
}

// The runtime loop for the terminal
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
