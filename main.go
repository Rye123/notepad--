package main

import (
	"log"
	"os"
	"github.com/Rye123/notepad--/tui"
	"github.com/gdamore/tcell/v2"
)

const AppName = "Notepad--"

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
		tui.DrawTitleBar(screen, AppName, filename)
		tui.DrawMenuBar(screen)
		tui.DrawStatusBar(screen, cursor_x, cursor_y, "CRLF", "UTF-8")
		
		screen.Show()
		
	}

}
