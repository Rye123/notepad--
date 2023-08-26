package tui

import (
	"github.com/gdamore/tcell/v2"
)

type TUIElem interface {
	Draw()  // Draws the element on screen
	IsActive() bool // Returns if the current element is focused on
	Focus() // Change cursor focus to this element.
	Unfocus() // Change cursor focus off this element.
	GetCursorIndex() int // Returns the cursor index on this element
	SetCursorIndex(newCursorIndex int) // Moves the cursor index on this element
	IsHidden() bool
	Hide() // Hides the element
	Show() // Shows the element
	HandleKey(keyEvent *tcell.EventKey)
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
