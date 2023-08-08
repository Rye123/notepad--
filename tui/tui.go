package tui

import (
	"fmt"
	"strings"
	"github.com/Rye123/notepad--/util"
	"github.com/Rye123/notepad--/textbuffer"
	"github.com/gdamore/tcell/v2"
)

// Returns cursor's x, y coordinates based on a given cursor_i
func GetCursorXY(screen tcell.Screen, textbuffer textbuffer.TextBuffer, options util.Options) (x int, y int) {
	scr_w, scr_h := screen.Size()
	scr_w--
	scr_h--
	x, y = 0, 0

	cursor_i := textbuffer.GetIndex()
	
	// if wordwrap is on, then conversion is simple
	if options.WordWrap {
		x = cursor_i % scr_w
		
	}

	// otherwise, we need to determine the horizontal scroll accordingly
	return x, y
}

// Draw Title Bar
func DrawTitleBar(screen tcell.Screen, appname string, filename string, fileModified bool) {
	scr_w, scr_h := screen.Size()
	barStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	scr_w--
	scr_h--
	if fileModified {
		// Add modified symbol
		filename = "*" + filename
	}
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
	menuRow := 2
	drawMenuButton(screen,  0, 4, menuRow, barStyle, "File", 0)
	drawMenuButton(screen,  6, 4, menuRow, barStyle, "Edit", 0)
	drawMenuButton(screen, 12, 6, menuRow, barStyle, "Format", 1)
	drawMenuButton(screen, 20, 4, menuRow, barStyle, "View", 0)
	drawMenuButton(screen, 26, 4, menuRow, barStyle, "Help", 0)

	// Draw Divider
	drawHorizontalLine(screen, 0, scr_w, menuRow+1, barStyle)
}

func drawMenuButton(screen tcell.Screen, x1, width, y int, style tcell.Style, text string, hotkeyIndex int) {
	drawText(screen, x1, y, x1+width, y, style, text)
	drawText(screen, x1+hotkeyIndex, y, x1+hotkeyIndex, y, style.Underline(true), string(text[hotkeyIndex]))
}

// Draw Text Box
func DrawTextBox(screen tcell.Screen, cursor_i int, options util.Options, text string) {
	scr_w, scr_h := screen.Size()
	scr_w--
	scr_h--
	textboxStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	cursor_x, cursor_y := 0, 0

	startRow := 4
	endRow := scr_h - 2

	rows := make([]string, endRow-startRow+1)

	if options.WordWrap {
		// convert text into a series of lines short enough to fit within scr_w
		curLine := 0
		for _, str := range(strings.Split(text[:cursor_i], "\n")) {
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
		cursor_y = curLine
		
		for _, str := range(strings.Split(text[cursor_i:], "\n")) {
			if (startRow + curLine) > endRow {
				break
			}

			// Split string
			if len(str) <= scr_w {
				rows[curLine] += str
			} else {
				fullString := str
				for len(fullString) > scr_w {
					if (startRow + curLine) >= endRow {
						break
					}
					rows[curLine] += fullString[:scr_w]
					curLine++
					fullString = fullString[scr_w:]
				}
				// Add the final line
				if len(fullString) <= scr_w {
					rows[curLine] += fullString
				} else {
					rows[curLine] += fullString[:scr_w]
				}
			}
			curLine++
		}
	} else {
		// convert text into a series of truncated lines
		//TODO: scroll right if cursor is viewing stuff beyond
		panic("not implemented")
		curLine := 0
		for _, str := range(strings.Split(text[cursor_i:], "\n")) {
			if (startRow + curLine) > endRow {
				break
			}
			
			// Truncate string
			if len(str) <= scr_w {
				rows[curLine] = str
			} else {
				rows[curLine] = str[:scr_w]
			}
			curLine++
		}
		curLine--
		cursor_x = len(rows[curLine])
		cursor_y = curLine
		
		for _, str := range(strings.Split(text[:cursor_i], "\n")) {
			if (startRow + curLine) > endRow {
				break
			}
			
			// Truncate string
			if len(str) <= scr_w {
				rows[curLine] += str
			} else {
				rows[curLine] += str[:scr_w]
			}
			curLine++

		}
	}

	for i, str := range(rows) {
		if startRow + i > endRow {
			break
		}
		drawText(screen, 0, startRow + i, scr_w, startRow + i, textboxStyle, str)
	}

	// Cursor
	screen.ShowCursor(cursor_x, cursor_y + startRow)
	DrawStatusBar(screen, cursor_x, cursor_y, options)
}

// Draw the Status Bar
// (Cursor Position)           | 100% | (Line End Mode) | (Encoding)
func DrawStatusBar(screen tcell.Screen, cursor_x, cursor_y int, options util.Options) {
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
	line_end_mode_text := options.LineEndModeString()
	otherText := fmt.Sprintf("| 100%% | %v   | %v   ", line_end_mode_text, options.Encoding)
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
