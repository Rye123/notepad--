package tui

import (
	"fmt"
	"strings"
	"github.com/Rye123/notepad--/util"
	"github.com/Rye123/notepad--/textbuffer"
	"github.com/gdamore/tcell/v2"
)

type TUIElem interface {
	Draw()  // Draws the element on screen
	IsActive() bool // Returns if the current element is focused on
	Focus() // Change cursor focus to this element.
	Unfocus() // Change cursor focus off this element.
	GetCursorIndex() int // Returns the cursor index on this element
	SetCursorIndex(newCursorIndex int) // Moves the cursor index on this element
	Hide() // Hides the elemnent
	Show() // Shows the element
}

// TitleBar: The top bar that shows the application name and the file name
type TitleBar struct {
	hidden bool
	row1, row2 int // Rows occupied by the titlebar
	appstate *util.AppState
}

func NewTitleBar(row1, row2 int, appstate *util.AppState) *TitleBar {
	return &TitleBar{false, row1, row2, appstate}
}

func (elem *TitleBar) Draw() {
	if elem.hidden {
		return
	}
	
	appstate := elem.appstate
	
	scr_w, scr_h := appstate.Screen.Size()
	scr_w--; scr_h--

	filename := appstate.Filename
	if appstate.FileModified {
		filename = "*" + filename
	}
	titleText := "🗒 " + filename + " - " + appstate.AppName

	drawText(appstate.Screen, 0, elem.row1, len(titleText), elem.row2, appstate.BarStyle, titleText)
	drawHorizontalLine(appstate.Screen, 0, scr_w, elem.row2, appstate.BarStyle)
}

func (elem *TitleBar) IsActive() bool{
	return false
}

func (elem *TitleBar) Focus() {
	panic("not implemented")
}

func (elem *TitleBar) Unfocus() {
	panic("not implemented")
}

func (elem *TitleBar) GetCursorIndex() int {
	panic("not implemented")
}

func (elem *TitleBar) SetCursorIndex(newCursorIndex int) {
	panic("not implemented")
}

func (elem *TitleBar) Hide() {
	elem.hidden = true
}

func (elem *TitleBar) Show() {
	elem.hidden = false
}

// MenuBar: A bar that shows Alt-functions (e.g. File, Edit, etc)
type MenuBar struct {
	hidden bool
	active bool // True if this element is focused on
	row1, row2 int // Rows occupied by the menubar
	cursorIndex int // Current-selected menu button
	appstate *util.AppState
}

func NewMenuBar(row1, row2 int, appstate *util.AppState) *MenuBar {
	return &MenuBar{false, false, row1, row2, 0, appstate}
}

const MENU_BUTTON_COUNT = 5

type menuButton struct {
	active bool
	x int
	text string
	hotkeyIndex int // Which index in the text is the hotkey
}

func (elem *MenuBar) Draw() {
	// File | Edit | Format | View | Help
	if elem.hidden {
		return
	}

	appstate := elem.appstate
	
	scr_w, scr_h := appstate.Screen.Size()
	scr_w--; scr_h--

	// Define Buttons
	buttons := [MENU_BUTTON_COUNT]menuButton{
		{false, 0, "File", 0},
		{false, 6, "Edit", 0},
		{false, 12, "Format", 1},
		{false, 20, "View", 0},
		{false, 26, "Help", 0},
	}

	// Set Active Button
	if elem.active {
		buttons[elem.cursorIndex].active = true
	}

	// Draw Buttons
	for _, button := range(buttons) {
		button.drawAtRow(appstate, elem.row1)
	}

	// Draw Divider
	drawHorizontalLine(appstate.Screen, 0, scr_w, elem.row2, appstate.BarStyle)
}

func (elem *MenuBar) IsActive() bool {
	return elem.active
}

func (button *menuButton) drawAtRow(appstate *util.AppState, row int) {
	x1, x2 := button.x, button.x + len(button.text)
	hotkeyX := button.x + button.hotkeyIndex
	style := appstate.ButtonStyle
	if button.active {
		style = appstate.ButtonActiveStyle
	}
	
	drawText(appstate.Screen, x1, row, x2, row, style, button.text)
	drawText(appstate.Screen, hotkeyX, row, hotkeyX, row, style.Underline(true), string(button.text[button.hotkeyIndex]))
}

func (elem *MenuBar) Focus() {
	elem.active = true
	elem.cursorIndex = 0
}

func (elem *MenuBar) Unfocus() {
	elem.active = false
	elem.cursorIndex = -1
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

func (elem *MenuBar) Hide() {
	elem.hidden = true
}

func (elem *MenuBar) Show() {
	elem.hidden = false
}

// Textbox: The main textbox for data entry
type Textbox struct {
	hidden bool
	active bool // True if this element is focused on
	startRow int // Starting row of the textbox
	cursorIndex int // Current position of the cursor index in the buffer
	leftIndex int // x-coordinate of leftmost index -- this is to allow 'scrolling' for non-wordwrapped text
	buf textbuffer.TextBuffer
	appstate *util.AppState
}

func (elem *Textbox) GetLeft() int {
	return elem.leftIndex
}

func NewTextbox(row int, initialText string, appstate *util.AppState) *Textbox {
	buf := textbuffer.NewGapBuffer()
	if len(initialText) > 0 {
		buf.Append(initialText)
	}
	
	textbox := Textbox{
		false,
		true,
		row,
		0, 0,
		buf,
		appstate,
	}

	textbox.SetCursorIndex(len(initialText))
	return &textbox
}

func (elem *Textbox) GetCursorXY() (x int, y int) {
	// x is the number of characters in last line, y is number of lines
	linelength := 0
	linecount := 0
	textTillCursor := elem.buf.String()[:elem.cursorIndex]
	for _, ch := range(textTillCursor) {
		if ch == '\n' {
			linelength = 0
			linecount++
		} else {
			linelength++
		}
	}
	x = linelength
	y = linecount

	return x, y
}

func (elem *Textbox) Draw() {
	if elem.hidden {
		return
	}

	appstate := elem.appstate
	
	scr_w, scr_h := appstate.Screen.Size()
	scr_w--; scr_h--

	startRow, endRow := elem.startRow, scr_h-2
	rows := make([]string, endRow-startRow + 1)

	cursor_x, cursor_y := 0, 0

	if appstate.Options.WordWrap {
		// Convert text into a series of lines short enough to fit within scr_w
		curLine := 0
		for _, str := range(strings.Split(elem.buf.StringBeforeIndex(), "\n")) {
			if (startRow + curLine) > endRow {
				break
			}

			// Split string
			if len(str) <= scr_w {
				rows[curLine] = str
			} else {
				fullString := str
				for len(fullString) > scr_w {
					if (startRow + curLine) >= endRow {
						break
					}
					rows[curLine] = fullString[:scr_w]
					curLine++
					fullString = fullString[scr_w:]
				}
				// Add the final line
				if len(fullString) <= scr_w {
					rows[curLine] = fullString
				} else {
					rows[curLine] = fullString[:scr_w]
				}
			}
			curLine++
		}
		curLine--
		cursor_x = len(rows[curLine])
		cursor_y = curLine + startRow

		// Process the other lines
		for _, str := range(strings.Split(elem.buf.StringAfterInclIndex(), "\n")) {
			if (startRow + curLine) > endRow {
				break
			}

			// Split string
			fullString := rows[curLine] + str
			if len(fullString) <= scr_w {
				rows[curLine] = fullString
			} else {
				for len(fullString) > scr_w {
					if (startRow + curLine) >= endRow {
						break
					}
					rows[curLine] = fullString[:scr_w]
					curLine++
					fullString = fullString[scr_w:]
				}
				// Add the final line
				if len(fullString) <= scr_w {
					rows[curLine] = fullString
				} else {
					rows[curLine] = fullString[:scr_w]
				}
			}
			curLine++
		}

		
	} else {
		if elem.active {
			cursor_x, cursor_y = elem.GetCursorXY()
		}
		// convert text into a series of truncated lines
		for i, str := range(strings.Split(elem.buf.String(), "\n")) {
			if (startRow + i) > endRow {
				break
			}

			// Truncate string based on leftIndex
			if len(str) < elem.leftIndex {
				rows[i] = ""
			} else if len(str[elem.leftIndex:]) <= scr_w {
				rows[i] = str[elem.leftIndex:]
			} else {
				rows[i] = str[elem.leftIndex:elem.leftIndex + scr_w+1]
			}
		}

		// set cursor
		cursor_x -= elem.leftIndex
		cursor_y += startRow
	}

	// Draw Text on screen
	for i, str := range(rows) {
		if startRow + i > endRow {
			break
		}
		drawText(appstate.Screen, 0, startRow + i, scr_w, startRow + i, appstate.TextboxStyle, str)
	}

	// Show Cursor if active
	if elem.active {
		appstate.Screen.ShowCursor(cursor_x, cursor_y)
	}
}

func (elem *Textbox) IsActive() bool {
	return elem.active
}

func (elem *Textbox) Focus() {
	elem.active = true
}

func (elem *Textbox) Unfocus() {
	elem.active = false
}

func (elem *Textbox) GetCursorIndex() int {
	return elem.cursorIndex
}

func (elem *Textbox) SetCursorIndex(newCursorIndex int) {
	if newCursorIndex < 0 {
		newCursorIndex = 0
	} else if newCursorIndex > elem.buf.Length() {
		newCursorIndex = elem.buf.Length()
	}
	
	elem.buf.MoveIndex(newCursorIndex)
	elem.cursorIndex = newCursorIndex
	
	// Set left accordingly
	cursor_x, _ := elem.GetCursorXY()
	
	scr_w, _ := elem.appstate.Screen.Size()
	scr_w--
		
	for (cursor_x - elem.leftIndex) > scr_w {
		elem.leftIndex++
	}

	for (cursor_x - elem.leftIndex) < 0 {
		elem.leftIndex--
	}
	
}

func (elem *Textbox) Hide() {
	elem.hidden = true
}

func (elem *Textbox) Show() {
	elem.hidden = false
}

// StatusBar: Bottom bar that shows details about the file
type StatusBar struct {
	hidden bool
	textbox *Textbox // Textbox the statusbar is responsible for
	appstate *util.AppState
}

func NewStatusBar(textbox *Textbox, appstate *util.AppState) *StatusBar {
	return &StatusBar{false, textbox, appstate}
}

func (elem *StatusBar) Draw() {
	// (Cursor Position)           | 100% | (Line End Mode) | (Encoding)
	if elem.hidden {
		return
	}

	appstate := elem.appstate
	
	scr_w, scr_h := appstate.Screen.Size()
	scr_w--; scr_h--;

	drawHorizontalLine(appstate.Screen, 0, scr_w, scr_h-1, appstate.BarStyle)

	// Status Data
	cursor_x, cursor_y := elem.textbox.GetCursorXY()
	//TODO: Remove debugging cursorIndex
	cursorText := fmt.Sprintf("Ln %d, Col %d (%d)", cursor_y+1, cursor_x+1, elem.textbox.cursorIndex)
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

func (elem *StatusBar) Hide() {
	elem.hidden = true
}

func (elem *StatusBar) Show() {
	elem.hidden = false
}

/* HELPER FUNCTIONS */

// Draw Text
// Code adapted from https://github.com/gdamore/tcell/blob/main/TUTORIAL.md
func drawText(screen tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		screen.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}


// Draw Horizontal Line
func drawHorizontalLine(screen tcell.Screen, x1, x2, y int, style tcell.Style) {
	if x2 < x1 {
		x1, x2 = x2, x1
	}
	for col := x1; col <= x2; col ++ {
		screen.SetContent(col, y, tcell.RuneHLine, nil, style)
	}
}
