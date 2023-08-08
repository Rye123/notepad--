package main

import (
	"log"
	"os"
	"errors"
	"github.com/Rye123/notepad--/tui"
	"github.com/Rye123/notepad--/util"
	"github.com/Rye123/notepad--/textbuffer"
	"github.com/gdamore/tcell/v2"
)

const AppName = "Notepad--"

func main() {
	// Get command line arguments
	commandArgs := os.Args
	
	// Initialise variables
	filename := "Untitled"
	options := util.Options{
		LineEndMode: "CRLF",
		Encoding: "UTF-8",
		WordWrap: true,
	}
	textbuf := textbuffer.NewGapBuffer()
	fileModified := false
	if len(commandArgs) > 1 {
		filename = commandArgs[1]
		// Read file if given
		data, err := os.ReadFile(filename)
		if err != nil {
			// Ignore error if file doesn't exist yet
			if !errors.Is(err, os.ErrNotExist) {
				log.Fatalf("%+v", err)
			}
			//TODO: add prompt to create file if doesn't exist
		} else {
			// Load data into buffer
			textbuf.Append(string(data))
		}
	}
	
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
	screen.SetCursorStyle(tcell.CursorStyleDefault)

	// Setup Screen
	screen.Clear()
		
	// Event Loop
	quit := func() {
		// Catch panics
		maybePanic := recover()
		screen.SetCursorStyle(tcell.CursorStyleDefault)
		screen.Fini()

		if maybePanic != nil {
			panic(maybePanic)
		}
		os.Exit(0)
	}
	defer quit()

	screen.SetCursorStyle(tcell.CursorStyleBlinkingBar)

renderLoop:
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

			// Arrow Keys
			if key == tcell.KeyLeft {
				textbuf.MoveIndex(textbuf.GetIndex() - 1)
			}
			if key == tcell.KeyRight {
				textbuf.MoveIndex(textbuf.GetIndex() + 1)
			}
			
			// Ctrl-W: Exit (if no more tabs left)
			if key == tcell.KeyCtrlW {
				break renderLoop
			}

			// (Debug: ESC to Exit)
			if key == tcell.KeyEscape {
				break renderLoop
			}
		}

		// Draw Screen
		screen.Clear()
		
		//// TUI
		tui.DrawTitleBar(screen, AppName, filename, fileModified)
		tui.DrawMenuBar(screen)

		tui.DrawTextBox(screen, textbuf.GetIndex(), options, textbuf.String())
		
		screen.Show()
		
	}

}
