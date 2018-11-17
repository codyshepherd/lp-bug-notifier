package main

import (
    "fmt"
    "strings"
)

type Tracker struct {
    recv chan []byte
    send chan []byte

    list map[string]string
}

var ops = map[string]func(t *Tracker, string, string){
            "add": func(t *Tracker, s string, url string){
                return t.add(s, url)
             },
            "drop": func(t *Tracker, s string, url string){
                return t.drop(s, url)
            },
}

func NewTracker() *Tracker{
    log.Debug("New tracker created")
    var t *Tracker
    t.list = make(map[string]string)
    t.recv = make(chan []byte, 1024)
    return &t
}

func (t *Tracker) Add(s string, url string){
    log.Debug(fmt.Sprintf("add: %s :: %s", s, url))
    t.list[s] = s
}

func (t *Tracker) Drop(s string, url string){
    if u, ok := t.list[s]; ok {
        log.Debug(fmt.Sprintf("drop: %s :: %s", s, u))
        del t.list[s]
    }
    else {
        log.Error(fmt.Sprintf("Cannot drop: %s \nDoes not exist!", s))
    }
}

func (t *Tracker) GetList() map[string]string {
    return t.list
}

func (t *Tracker) Run(){
    for {
        select {
            case msg := <-t.recv
                words := strings.Split(msg, ' ')
                if len(words) != 3 {
                    log.Error(fmt.Sprintf(`Tracker recieved msg with incorrect
                    number of args: %s`, msg))
                    continue
                }
                op := words[0]
                s := words[1]
                url := words[2]
                ops(t, s, url)
            case msg := <-t.send

        }
    }
}

