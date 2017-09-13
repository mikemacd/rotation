package main

import (
	"github.com/andlabs/ui"
)

var CTRLPressed bool = false

type WindowHandler struct {
	Window *ui.Window
}

// Gets called when "something" changes in the environment and the Window system needs to tell the application to "redraw" itself.
func (self WindowHandler) Draw(a *ui.Area, db *ui.AreaDrawParams) {
		R.Render(a,db)
		R.Rotate()
}

// Stub to match ui.WindowHandler interface
func (self WindowHandler) MouseEvent(a *ui.Area, me *ui.AreaMouseEvent) {}

// Stub to match ui.WindowHandler interface
func (self WindowHandler) MouseCrossed(a *ui.Area, left bool) {}

// Stub to match ui.WindowHandler interface
func (self WindowHandler) DragBroken(a *ui.Area) {}

func (self WindowHandler) KeyEvent(a *ui.Area, ke *ui.AreaKeyEvent) (handled bool) {
	if ke.Modifier == ui.Ctrl {
		CTRLPressed = true
	} else if CTRLPressed && ke.Key == 'c' {
		ui.Quit()
	} else {
		CTRLPressed = false
	}
	return true
}
