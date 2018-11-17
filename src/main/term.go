package main

import (
    "bufio"
    "os"
)

type Term struct{
    reader  *bufio.Reader
    tracker *Tracker
}

func NewTerm(tracker Tracker){
    var t *Term
    t.reader = bufio.NewReader(os.Stdin)
    t.tracker = newTracker()
}

func (t *Term) Run(){
    for {
        input,_  := t.reader.ReadString('\n')
        t.recv <-input
    }
}
