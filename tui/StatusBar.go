package tui

import (
	"fmt"
	"strings"
	"github.com/gdamore/tcell/v2"
	"github.com/Rye123/notepad--/util"
)

// (Cursor XY Position)        | 100% | (Line End Mode) | (Encoding)
// StatusBar: Bottom bar that shows detail about the file
type StatusBar struct {
	hidden bool
	textbox *Textbox
	drawn bool
	appstate *util.AppState
}

func NewStatusBar(appstate *util.AppState, textbox *Textbox) *StatusBar {
	return &StatusBar{false, textbox, false, appstate}
}

func (elem *StatusBar) Draw() {
	if elem.hidden {
		return
	}

	// Don't update if textbox is not active and if this is already drawn
	if !elem.textbox.active && elem.drawn {
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
	otherText := fmt.Sprintf("| 100%% | %v | %v ", appstate.Options.LineEndModeString(), appstate.Options.Encoding)

	// Generate full string
	spaceBetween := scr_w - len(otherText) - len(cursorText)
	fullText := cursorText
	// if not enough space between, we only show the cursor text
	if spaceBetween > 0 {
		// otherwise, we can show both
		fullText = cursorText + strings.Repeat(" ", spaceBetween) + otherText
	}
	drawText(appstate.Screen, 0, scr_h, scr_w, scr_h, appstate.BarStyle, fullText)

	elem.drawn = true
}

func (elem *StatusBar) IsActive() bool {
	return false
}

func (elem *StatusBar) Focus() {
	panic("not implemented")
}

func (elem *StatusBar) Unfocus() {
	return
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
	elem.drawn = false
}

func (elem *StatusBar) Show() {
	elem.hidden = false
}

func (elem *StatusBar) HandleKey(keyEvent *tcell.EventKey) {
	return
}
