package tui

import "github.com/Rye123/notepad--/util"

const MENUBAR_STARTROW = 2
const MENUBAR_ENDROW = 3

// File | Edit | Format | View | Help
const MENU_BUTTON_COUNT = 5

// MenuBar: A bar that shows Alt-functions (e.g. File, Edit, etc)
type MenuBar struct {
	hidden bool
	active bool
	cursorIndex int
	appstate *util.AppState
}

func NewMenuBar(appstate *util.AppState) *MenuBar {
	return &MenuBar{false, false, 0, appstate}
}

func (elem *MenuBar) Draw() {
	if elem.hidden {
		return
	}

	appstate := elem.appstate
	scr_w, scr_h := appstate.Screen.Size()
	scr_w--; scr_h--

	// Define Buttons
	buttons := [MENU_BUTTON_COUNT]menuButton{
		{false, 0, "File", 0, appstate},
		{false, 6, "Edit", 0, appstate},
		{false, 12, "Format", 1, appstate},
		{false, 20, "View", 0, appstate},
		{false, 26, "Help", 0, appstate},
	}

	// Set Active Button
	if elem.active {
		buttons[elem.cursorIndex].active = true
	}

	// Draw Buttons
	for _, button := range(buttons) {
		button.drawAtRow(MENUBAR_STARTROW)
	}

	// Draw Divider
	drawHorizontalLine(appstate.Screen, 0, scr_w, MENUBAR_ENDROW, appstate.BarStyle)
}

func (elem *MenuBar) IsActive() bool {
	return elem.active
}

func (elem *MenuBar) Focus() {
	elem.active = true
}

func (elem *MenuBar) Unfocus() {
	elem.active = false
}

func (elem *MenuBar) GetCursorIndex() int {
	return elem.cursorIndex
}

func (elem *MenuBar) SetCursorIndex(newCursorIndex int) {
	// Set cursor to wrap around where necessary
	if newCursorIndex < 0 {
		newCursorIndex = MENU_BUTTON_COUNT-1
	}

	if newCursorIndex >= MENU_BUTTON_COUNT {
		newCursorIndex = 0
	}

	elem.cursorIndex = newCursorIndex
}

func (elem *MenuBar) IsHidden() bool {
	return elem.hidden
}

func (elem *MenuBar) Hide() {
	elem.hidden = true
}

func (elem *MenuBar) Show() {
	elem.hidden = false
}


// Defines a button in the menu bar.
type menuButton struct {
	active bool
	x int
	text string
	hotkeyIndex int // Which index in the text is the hotkey
	appstate *util.AppState
}

func (but *menuButton) drawAtRow(row int) {
	appstate := but.appstate
	x1, x2 := but.x, but.x + len(but.text)
	hotkeyX := but.x + but.hotkeyIndex
	style := appstate.ButtonStyle
	if but.active {
		style = appstate.ButtonActiveStyle
	}

	drawText(appstate.Screen, x1, row, x2, row, style, but.text)
	drawText(appstate.Screen, hotkeyX, row, hotkeyX, row, style.Underline(true), string(but.text[but.hotkeyIndex]))
}
