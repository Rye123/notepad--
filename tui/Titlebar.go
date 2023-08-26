package tui

import (
	"github.com/Rye123/notepad--/util"
)

const TITLEBAR_STARTROW = 0
const TITLEBAR_ENDROW = 1

type TitleBar struct {
	hidden bool
	appstate *util.AppState
}

func NewTitleBar(appstate *util.AppState) *TitleBar {
	return &TitleBar{false, appstate}
}

func (el *TitleBar) Draw() {
	if el.hidden {
		return
	}

	scr_w, scr_h := el.appstate.Screen.Size()
	scr_w--; scr_h--

	filename := el.appstate.Filename
	
	if len(filename) == 0 {
		// If filename not set, set temporary title
		filename = util.GetTemporaryTitle(el.appstate.TextBuffer.String())
		if len(filename) == 0 {
			filename = "Untitled"
			el.appstate.FileModified = false
		}
	}
	
	if el.appstate.FileModified {
		filename = "*" + filename
	}

	titleText := "🗒 " + filename + " - " + el.appstate.AppName
	drawText(el.appstate.Screen, 0, TITLEBAR_STARTROW, len(titleText), TITLEBAR_STARTROW, el.appstate.BarStyle, titleText)
	drawHorizontalLine(el.appstate.Screen, 0, scr_w, TITLEBAR_ENDROW, el.appstate.BarStyle)
}

func (el *TitleBar) IsActive() bool {
	return false
}

func (el *TitleBar) Focus() {
	panic("not implemented")
}

func (el *TitleBar) Unfocus() {
	panic("not implemented")
}

func (el *TitleBar) GetCursorIndex() int {
	panic("not implemented")
}

func (el *TitleBar) SetCursorIndex(newIndex int) {
	panic("not implemented")
}

func (el *TitleBar) IsHidden() bool {
	return el.hidden
}

func (el *TitleBar) Hide() {
	el.hidden = true
}

func (el *TitleBar) Show() {
	el.hidden = false
}

