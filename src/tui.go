package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"lookup"
	"strings"
)

func enterWord(tui *gocui.Gui, view *gocui.View) error {
	word := strings.TrimSuffix(view.Buffer(), "\n")
	definition, err := lookup.GetDefinition(word)
	if err != nil {
		return err
	}

	outputView, _ := tui.View("definitionview")
	outputView.Clear()
	fmt.Fprint(outputView, definition)

	//Clears view and moves cursor back to start of view
	view.Clear()
	cx, cy := view.Cursor()
	view.MoveCursor(-cx, -cy, false)

	return nil
}

func drawTui(tui *gocui.Gui) error {
	width, height := tui.Size()

	wordView, err := tui.SetView("wordview", width/3, 0, 2*width/3, 2)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	wordView.Editable = true
	wordView.Title = "Enter Word"

	definitionView, err := tui.SetView("definitionview", width/10, 5, 9*width/10, height-1)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	definitionView.Wrap = true
	definitionView.Title = "Definition"

	tui.SetCurrentView("wordview")

	return nil
}

func quit(tui *gocui.Gui, view *gocui.View) error {
	return gocui.ErrQuit
}

func run() error {
	//Initialize TUI
	tui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}
	defer tui.Close()

	tui.SetManagerFunc(drawTui)

	//Create keybindings
	if err := tui.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := tui.SetKeybinding("wordview", gocui.KeyEnter, gocui.ModNone, enterWord); err != nil {
		return err
	}

	//Run mainloop
	if err := tui.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
}
