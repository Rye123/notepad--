package tui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
)

// Draw Title Bar
func DrawTitleBar(screen tcell.Screen, appname string, filename string) {
	scr_w, scr_h := screen.Size()
	barStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	scr_w--
	scr_h--

	titleText := "🗒 " + filename + " - " + appname
	drawText(screen, 0, 0, len(titleText), 0, barStyle, titleText)
	drawHorizontalLine(screen, 0, scr_w, 1, barStyle)
}

// Draw Menu Bar
// File | Edit | Format | View | Help
func DrawMenuBar(screen tcell.Screen) {
	scr_w, scr_h := screen.Size()
	barStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	scr_w--
	scr_h--

	// Draw Buttons
	drawMenuButton(screen, 0, 4, 2, barStyle, "File", 0)
	drawMenuButton(screen, 6, 4, 2, barStyle, "Edit", 0)
	drawMenuButton(screen, 12, 6, 2, barStyle, "Format", 1)
	drawMenuButton(screen, 20, 4, 2, barStyle, "View", 0)
	drawMenuButton(screen, 26, 4, 2, barStyle, "Help", 0)

	// Draw Divider
	drawHorizontalLine(screen, 0, scr_w, 3, barStyle)
}

func drawMenuButton(screen tcell.Screen, x1, width, y int, style tcell.Style, text string, hotkeyIndex int) {
	drawText(screen, x1, y, x1+width, y, style, text)
	drawText(screen, x1+hotkeyIndex, y, x1+hotkeyIndex, y, style.Underline(true), string(text[hotkeyIndex]))
}

// Draw the Status Bar
// (Cursor Position)           | 100% | (Line End Mode) | (Encoding)
func DrawStatusBar(screen tcell.Screen, cursor_x, cursor_y int, line_end_mode string, encoding string) {
	scr_w, scr_h := screen.Size()
	barStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	scr_w--
	scr_h--
	
	// Status Bar
	drawHorizontalLine(screen, 0, scr_w, scr_h-1, barStyle)
	//// Cursor Data
	cursorText := fmt.Sprintf("Ln %d, Col %d", cursor_x+1, cursor_y+1)
	drawText(screen, 0, scr_h, len(cursorText), scr_h, barStyle, cursorText)
	
	//// Other Data
	line_end_mode_text := ""
	switch line_end_mode {
	case "LF":
		line_end_mode_text = "Unix (LF)"
	case "CRLF":
		line_end_mode_text = "Windows (CRLF)"
	default:
		line_end_mode_text = "Windows (CRLF)"
	}
	otherText := fmt.Sprintf("| 100%% | %v   | %v   ", line_end_mode_text, encoding)
	drawText(screen, scr_w - len(otherText), scr_h, scr_w, scr_h, barStyle, otherText)
}

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
