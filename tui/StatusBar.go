package tui

import (
	"fmt"
	"github.com/Rye123/notepad--/util"
)

// (Cursor XY Position)        | 100% | (Line End Mode) | (Encoding)
// StatusBar: Bottom bar that shows detail about the file
type StatusBar struct {
	hidden bool
	textbox *Textbox
	appstate *util.AppState
}

func NewStatusBar(appstate *util.AppState, textbox *Textbox) *StatusBar {
	return &StatusBar{false, textbox, appstate}
}

func (elem *StatusBar) Draw() {
	if elem.hidden {
		return
	}

	appstate := elem.appstate
	scr_w, scr_h := appstate.Screen.Size()
	scr_w--; scr_h--

	drawHorizontalLine(appstate.Screen, 0, scr_w, scr_h - 1, appstate.BarStyle)

	// Status Data
	cursorX, cursorY := elem.textbox.GetCursorXY()
	//TODO: Remove debugging cursorIndex
	cursorText := fmt.Sprintf("Ln %d, Col %d (%d)", cursorY+1, cursorX+1, elem.textbox.cursorIndex)
	drawText(appstate.Screen, 0, scr_h, len(cursorText), scr_h, appstate.BarStyle, cursorText)

	otherText := fmt.Sprintf("| 100%% | %v | %v ", appstate.Options.LineEndModeString(), appstate.Options.Encoding)
	drawText(appstate.Screen, scr_w - len(otherText), scr_h, scr_w, scr_h, appstate.BarStyle, otherText)
}

func (elem *StatusBar) IsActive() bool {
	return false
}

func (elem *StatusBar) Focus() {
	panic("not implemented")
}

func (elem *StatusBar) Unfocus() {
	panic("not implemented")
}

func (elem *StatusBar) GetCursorIndex() int {
	panic("not implemented")
}

func (elem *StatusBar) SetCursorIndex(newCursorIndex int) {
	panic("not implemented")
}

func (elem *StatusBar) IsHidden() bool {
	return elem.hidden
}

func (elem *StatusBar) Hide() {
	elem.hidden = true
}

func (elem *StatusBar) Show() {
	elem.hidden = false
}
