package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/jroimartin/gocui"
	//log "github.com/sirupsen/logrus"
)

type Ui struct {
	// gui     *gocui.Gui
	tracker *Tracker
}

func NewUi() *Ui {
	u := new(Ui)
	u.tracker = NewTracker()
	u.tracker.Read()
	/*
		var err error
		u.gui, err = gocui.NewGui(gocui.OutputNormal)
		if err != nil {
			log.Panic(err)
		}

		u.gui.SetManagerFunc(u.MkLayout())

		if err := u.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
			log.Panic(err)
		}
	*/
	return u
}

/*
// "curried" function for passing term contents into layout
func (u *Ui) MkLayout() func(*gocui.Gui) error {
	return func(g *gocui.Gui) error {
		_, maxY := g.Size()
		if v, err := g.SetView("side", -1, -1, 30, maxY); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Highlight = true
			v.SelBgColor = gocui.ColorGreen
			v.SelFgColor = gocui.ColorBlack
			var noargs []string
			bugs := list_bugs(u.tracker, noargs)

			for _, line := range bugs {
				fmt.Fprintln(v, line)
			}
		}
		return nil
	}
}
*/

func (u *Ui) Layout(g *gocui.Gui) error {
	_, maxY := g.Size()
	if v, err := g.SetView("side", -1, -1, 30, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		var noargs []string
		bugs := list_bugs(u.tracker, noargs)

		for _, line := range bugs {
			fmt.Fprintln(v, line)
		}
	}
	return nil

}

type TextBar struct {
	buffer []string
}

func NewTextBar() *TextBar {
	tb := new(TextBar)
	return tb
}

func (tb *TextBar) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	v, err := g.SetView("TextBar", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Input Bug")
		cursor_keybindings(g, v)
	}
	return nil
}

func (tb *TextBar) Edit(g *gocui.Gui, v *gocui.View) error {
	_, err := g.SetCurrentView("TextBar")
	if err != nil {
		return err
	}

	return nil
}

func add_bug(t *Tracker, args []string) {

}

func list_bugs(t *Tracker, args []string) []string {
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

	// sort ascending by time since last update
	sort.Sort(byTime(buffer))
	return buffer
}
