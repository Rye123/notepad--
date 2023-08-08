package main

import (
	"log"
	"os"
	"github.com/Rye123/notepad--/tui"
	"github.com/Rye123/notepad--/util"
	"github.com/gdamore/tcell/v2"
)

const AppName = "Notepad--"

func main() {
	// Initialise variables
	filename := "Untitled"
	cursor_x := 0
	cursor_y := 0
	options := util.Options{
		LineEndMode: "CRLF",
		Encoding: "UTF-8",
		WordWrap: true,
	}
	buffer := "The quick\r\nbrown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog. Quick brown\r\nfox jumps over the lazy dog. \n\n\n\n\n\n \n\n\n\n\n\n\nThis is such a long line that we will probably have to word wrap several times but this is necessary so that i can actually test the full word wrap features and hopefully this works perfectly and i have accounted for everything. In the future, I'll probably need to probably implement WORD wrap and not CHARACTER wrap so that the word goes to the upper or lower line where necessary, but for now hopefully this suffices."
	
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
		tui.DrawStatusBar(screen, cursor_x, cursor_y, options)

		tui.DrawTextBox(screen, cursor_x, cursor_y, options, buffer)
		
		screen.Show()
		
	}

}
