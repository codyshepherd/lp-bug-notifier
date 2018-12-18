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

	s := NewSide()
	tb := NewAddTextBar()

	g.SetManager(tb, s)

	if err := g.SetKeybinding("", gocui.KeyCtrlA, gocui.ModNone, tb.Edit); err != nil {
		log.Panic(err)
	}

	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, s.FgSide); err != nil {
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

func fgSide(g *gocui.Gui, v *gocui.View) error {
	g.SetCurrentView("side")
	return nil
}

func fgAddTextBar(g *gocui.Gui, v *gocui.View) error {
	g.SetCurrentView("AddTextBar")
	return nil
}

func simpleEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyDelete:
		v.EditDelete(false)
	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	case key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case key == gocui.KeyArrowRight:
		ox, _ := v.Origin()
		cx, _ := v.Cursor()
		if cx+1 < ox+100 {
			v.MoveCursor(1, 0, false)
		}
	}
}

func cursor_keybindings(g *gocui.Gui, v *gocui.View) {
	if err := g.SetKeybinding("AddTextBar", gocui.KeyArrowLeft, gocui.ModNone,
		cursorLeft); err != nil {
		log.Panic(err)
	}

	if err := g.SetKeybinding("AddTextBar", gocui.KeyArrowRight, gocui.ModNone,
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
