package main

import (
	"log"
	"os"
	"fmt"
	"github.com/gdamore/tcell/v2"
)

const AppName = "Notepad--"

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

func drawHorizontalLine(screen tcell.Screen, x1, x2, y int, style tcell.Style) {
	if x2 < x1 {
		x1, x2 = x2, x1
	}
	for col := x1; col <= x2; col ++ {
		screen.SetContent(col, y, tcell.RuneHLine, nil, style)
	}
}


// Draw the Title Bar (Line 0 and Line 1)
func drawTitleBar(screen tcell.Screen, filename string) {
	scr_w, scr_h := screen.Size()
	barStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	scr_w--
	scr_h--

	titleText := "🗒 " + filename + " - " + AppName
	drawText(screen, 0, 0, len(titleText), 0, barStyle, titleText)
	drawHorizontalLine(screen, 0, scr_w, 1, barStyle)
}

func drawMenuButton(screen tcell.Screen, x1, width, y int, style tcell.Style, text string, hotkeyIndex int) {
	drawText(screen, x1, y, x1+width, y, style, text)
	drawText(screen, x1+hotkeyIndex, y, x1+hotkeyIndex, y, style.Underline(true), string(text[hotkeyIndex]))
}

// Draw the Menu Bar (Line 2 and Line 3)
// File | Edit | Format | View | Help
func drawMenuBar(screen tcell.Screen) {
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

// Draw the Status Bar
// (Cursor Position)           | 100% | (Line End Mode) | (Encoding)
func drawStatusBar(screen tcell.Screen, cursor_x, cursor_y int, line_end_mode string, encoding string) {
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

func main() {
	// Initialise variables
	filename := "Untitled"
	cursor_x := 0
	cursor_y := 0
	
	// Initialise screen
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	

	// Set Default Text Style
	defaultStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	screen.SetStyle(defaultStyle)

	// Setup Screen
	screen.Clear()
		
	// Event Loop
	quit := func() {
		screen.Fini()
		os.Exit(0)
	}
	for {
		screen.Show()

		// Process Events
		event := screen.PollEvent()
		switch eventType := event.(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			// Note: SHIFT doesn't appear to work.
			_, key, _ := eventType.Modifiers(), eventType.Key(), eventType.Rune()

			// Ctrl-W: Exit (if no more tabs left)
			if key == tcell.KeyCtrlW {
				quit()
			}

			// (Debug: ESC to Exit)
			if key == tcell.KeyEscape {
				quit()
			}
		}

		// Draw Screen
		screen.Clear()
		//// TUI
		drawTitleBar(screen, filename)
		drawMenuBar(screen)
		drawStatusBar(screen, cursor_x, cursor_y, "CRLF", "UTF-8")
		
		screen.Show()
		
	}

}
