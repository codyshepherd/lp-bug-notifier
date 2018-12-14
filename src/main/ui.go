package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
)

type Ui struct {
	gui  *gocui.Gui
	term *Term
}

func NewUi() *Ui {
	u := new(Ui)
	var err error
	u.gui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panic(err)
	}

	u.gui.SetManagerFunc(layout)

	if err := u.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panic(err)
	}

	return u
}

func (u *Ui) Run() {
	defer u.gui.Close()

	if err := u.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panic(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Hello world!")
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
