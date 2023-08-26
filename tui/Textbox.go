package tui

import (
	"strings"
	"github.com/Rye123/notepad--/util"
	"github.com/Rye123/notepad--/textbuffer"
)

const TEXTBOX_STARTROW = 4

type Textbox struct {
	hidden bool
	active bool // True if this element is focused on
	cursorIndex int // Current position of the cursorIndex in the buffer
	leftIndex int // x-coordinate of leftmost index, to allow horizontal scrolling for non-wordwrapped text
	buf textbuffer.TextBuffer
	appstate *util.AppState
}

func NewTextbox(appstate *util.AppState) *Textbox {
	textbox := Textbox{
		false,
		true,
		0, 0,
		appstate.TextBuffer,
		appstate,
	}

	textbox.SetCursorIndex(appstate.TextBuffer.Length())
	return &textbox
}

func (elem *Textbox) Draw() {
	if elem.hidden {
		return
	}

	appstate := elem.appstate
	scr_w, scr_h := appstate.Screen.Size()
	scr_w--; scr_h--

	startRow, endRow := TEXTBOX_STARTROW, scr_h - 2
	rows := make([]string, endRow - startRow + 1)

	// While GetCursorXY returns the true X and Y coordinates of the cursor, if word-wrapped is enabled, we need to calculate the view X and Y coordinates of the cursor.
	cursorX, cursorY := 0, 0
	if appstate.Options.WordWrap {
		// Convert text into lines short enough to fit in [0, scr_w]
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

		// Set view X and Y coordinates
		curLine--;
		cursorX = len(rows[curLine])
		cursorY = curLine + startRow

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
		cursorX, cursorY = elem.GetCursorXY()

		// Convert text into a series of truncated lines starting from leftIndex
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

		// Update cursor with textbox positioning
		cursorX -= elem.leftIndex
		cursorY += TEXTBOX_STARTROW
	}

	// Draw text
	for i, str := range(rows) {
		if startRow + i > endRow {
			break
		}
		drawText(appstate.Screen, 0, TEXTBOX_STARTROW + i, scr_w, TEXTBOX_STARTROW + i, appstate.TextboxStyle, str)
	}

	// Show Cursor
	if elem.active {
		appstate.Screen.ShowCursor(cursorX, cursorY)
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
	// Limit cursor movement at limits
	if newCursorIndex < 0 {
		newCursorIndex = 0
	} else if newCursorIndex > elem.buf.Length() {
		newCursorIndex = elem.buf.Length()
	}

	elem.buf.MoveIndex(newCursorIndex)
	elem.cursorIndex = newCursorIndex

	// Update left accordingly
	cursorX, _ := elem.GetCursorXY()
	scr_w, _ := elem.appstate.Screen.Size()
	scr_w --

	for (cursorX - elem.leftIndex) > scr_w {
		elem.leftIndex++
	}

	for (cursorX - elem.leftIndex) < 0 {
		elem.leftIndex--
	}
}

func (elem *Textbox) IsHidden() bool {
	return elem.hidden
}

func (elem *Textbox) Hide() {
	elem.hidden = true
}

func (elem *Textbox) Show() {
	elem.hidden = false
}


func (elem *Textbox) Content() string {
	return elem.buf.String()
}

// Returns the left offset of the displayed text
func (elem *Textbox) GetLeft() int {
	return elem.leftIndex
}

// Returns the true X and Y coordinates of the cursor
func (elem *Textbox) GetCursorXY() (x int, y int) {
	// Don't compute all that if not active
	if !elem.active {
		return 0, 0
	}
	
	// x is number of characters in last line, y is number of lines
	lines := strings.Split(elem.buf.StringBeforeIndex(), "\n")
	x = len(lines[len(lines) - 1])
	y = len(lines) - 1
	return x, y
}

func (elem *Textbox) Insert(key rune) {
	elem.buf.Insert(elem.cursorIndex, key)
	elem.SetCursorIndex(elem.cursorIndex + 1)
}

// Deletes the character directly after the cursor.
func (elem *Textbox) Delete() {
	elem.buf.Delete(elem.cursorIndex)
}

// Deletes the character directly before the cursor, and shifts the cursor backward
func (elem *Textbox) Backspace() {
	elem.SetCursorIndex(elem.cursorIndex - 1)
	elem.buf.Delete(elem.cursorIndex)
}

