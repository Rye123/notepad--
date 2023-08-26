package tui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/Rye123/notepad--/util"
)

const TITLEBAR_STARTROW = 0
const TITLEBAR_ENDROW = 1

type TitleBar struct {
	hidden bool
	textbox *Textbox
	drawn bool
	appstate *util.AppState
}

func NewTitleBar(appstate *util.AppState, textbox *Textbox) *TitleBar {
	return &TitleBar{false, textbox, false, appstate}
}

func (elem *TitleBar) Draw() {
	if elem.hidden {
		return
	}

	// Don't update if textbox is not active and if this is alread drawn
	if !elem.textbox.active && elem.drawn {
		return
	}

	scr_w, scr_h := elem.appstate.Screen.Size()
	scr_w--; scr_h--

	filename := elem.appstate.Filename
	
	if len(filename) == 0 {
		// If filename not set, set temporary title
		filename = util.GetTemporaryTitle(elem.appstate.TextBuffer.String())
		if len(filename) == 0 {
			filename = "Untitled"
			elem.appstate.FileModified = false
		}
	}
	
	if elem.appstate.FileModified {
		filename = "*" + filename
	}

	titleText := "🗒 " + filename + " - " + elem.appstate.AppName

	// Pad with spaces
	titleText = fmt.Sprintf("%-*s", scr_w, titleText)
	
	drawText(elem.appstate.Screen, 0, TITLEBAR_STARTROW, len(titleText), TITLEBAR_STARTROW, elem.appstate.BarStyle, titleText)
	drawHorizontalLine(elem.appstate.Screen, 0, scr_w, TITLEBAR_ENDROW, elem.appstate.BarStyle)

	elem.drawn = true
}

func (elem *TitleBar) IsActive() bool {
	return false
}

func (elem *TitleBar) Focus() {
	panic("not implemented")
}

func (elem *TitleBar) Unfocus() {
	return
}

func (elem *TitleBar) GetCursorIndex() int {
	panic("not implemented")
}

func (elem *TitleBar) SetCursorIndex(newIndex int) {
	panic("not implemented")
}

func (elem *TitleBar) IsHidden() bool {
	return elem.hidden
}

func (elem *TitleBar) Hide() {
	elem.hidden = true
	elem.drawn = false
}

func (elem *TitleBar) Show() {
	elem.hidden = false
}

func (elem *TitleBar) HandleKey(keyEvent *tcell.EventKey) {
	return
}
