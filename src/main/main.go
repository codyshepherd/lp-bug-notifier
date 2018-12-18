package main

import (
	"fmt"
	"os"

	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
)

func main() {
	f, err := os.OpenFile(".buglog", os.O_WRONLY|os.O_CREATE,
		0666)
	if err != nil {
		panic(fmt.Sprintf("Failed to open file for logging: %s", err))
	}

	log.SetOutput(f)
	log.SetLevel(log.DebugLevel)
	log.Info("LP Bug Notifier started")

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panic(err)
	}
	defer g.Close()

	g.Cursor = true

	u := NewUi()
	tb := NewTextBar()

	g.SetManager(tb, u)

	if err := g.SetKeybinding("", gocui.KeyCtrlA, gocui.ModNone, tb.Edit); err != nil {
		log.Panic(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, fgUi); err != nil {
		log.Panic(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panic(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panic(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func fgUi(g *gocui.Gui, v *gocui.View) error {
	g.SetCurrentView("side")
	return nil
}

func fgTextBar(g *gocui.Gui, v *gocui.View) error {
	g.SetCurrentView("TextBar")
	return nil
}

func cursor_keybindings(g *gocui.Gui, v *gocui.View) {
	if err := g.SetKeybinding("TextBar", gocui.KeyArrowLeft, gocui.ModNone,
		cursorLeft); err != nil {
		log.Panic(err)
	}

	if err := g.SetKeybinding("TextBar", gocui.KeyArrowRight, gocui.ModNone,
		cursorRight); err != nil {
		log.Panic(err)
	}
}

func cursorRight(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, _ := v.Origin()
		cx, _ := v.Cursor()
		if cx+1 < ox+10 {
			v.MoveCursor(1, 0, false)
		}
	}
	return nil
}

func cursorLeft(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, _ := v.Origin()
		cx, _ := v.Cursor()
		if cx-1 >= ox {
			v.MoveCursor(-1, 0, false)
		}
	}
	return nil
}
