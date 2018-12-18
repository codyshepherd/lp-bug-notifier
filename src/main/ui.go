package main

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
)

type Side struct {
	tracker *Tracker
}

func NewSide() *Side {
	s := new(Side)
	s.tracker = NewTracker()
	s.tracker.Read()
	return s
}

func (s *Side) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("side", -1, -1, maxX/3, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

		var noargs []string
		bugs := list_bugs(s.tracker, noargs)

		for _, line := range bugs {
			fmt.Fprintln(v, line)
		}
	}
	return nil
}

func (s *Side) FgSide(g *gocui.Gui, v *gocui.View) error {
	vs := g.Views()
	for _, view := range vs {
		if view.Name() == "AddTextBar" {
			contents := strings.Split(view.Buffer(), " ")
			view.Clear()

			for _, item := range contents {
				if _, ok := s.tracker.list[item]; !ok {
					s.tracker.Add(item)
				}
			}
		}
	}

	nv, err := g.SetCurrentView("side")
	nv.Clear()

	var noargs []string
	bugs := list_bugs(s.tracker, noargs)

	for _, line := range bugs {
		fmt.Fprintln(nv, line)
	}

	return err
}

type AddTextBar struct {
	buffer string
}

func NewAddTextBar() *AddTextBar {
	tb := new(AddTextBar)
	return tb
}

func (tb *AddTextBar) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	v, err := g.SetView("AddTextBar", maxX/2-10, maxY/2, maxX/2+10, maxY/2+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "")
		v.Editable = true
	}
	return nil
}

func (tb *AddTextBar) Edit(g *gocui.Gui, v *gocui.View) error {
	nv, err := g.SetCurrentView("AddTextBar")
	nv.SetCursor(0, 0)
	if err != nil {
		return err
	}

	nv.Editor = gocui.EditorFunc(simpleEditor)
	return nil
}

func list_bugs(t *Tracker, args []string) []string {
	buffer := []string{}
	t.lock.Lock()
	for k, v := range t.list {
		tm, _ := time.Parse(time.RFC3339, v.BugStruct.Date_last_message)

		// check for updated-ness and prepend string if appropriate
		// note that once "Updated!" is displayed, the Changed flag on the
		// bug will be turned off

		log.Debug("k: " + k)
		log.Debug("v: " + v.BugStruct.Title)
		buffer = append(buffer, fmt.Sprintf("%s [%s]: %s", k,
			time.Since(tm).Truncate(time.Minute).String(), v.BugStruct.Title))
	}
	t.lock.Unlock()

	// sort ascending by time since last update
	sort.Sort(byTime(buffer))

	newbuffer := []string{}
	for _, item := range buffer {
		k := strings.TrimRight(item, " ")
		if v, ok := t.list[k]; ok {
			updated := ""
			if v.Changed {
				updated = "**Updated!** "
				v.Changed = false
			}
			newbuffer = append(newbuffer, updated+item)
		} else {
			newbuffer = append(newbuffer, item)
		}

	}

	return newbuffer
}

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
	log.Debug(s)
	r := regexp.MustCompile(`\[[0-9]*`)
	string_i := strings.Trim(r.FindString(s[i]), "[")
	string_j := strings.Trim(r.FindString(s[j]), "[")
	int_i, erri := strconv.Atoi(string_i)
	int_j, errj := strconv.Atoi(string_j)
	if erri != nil || errj != nil {
		log.Panic("Problem with converting string to int")
	}
	return int_i < int_j
}
