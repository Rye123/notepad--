package main

import (
	"fmt"
	"log"
	"os"
	"errors"
	"github.com/Rye123/notepad--/tui"
	"github.com/Rye123/notepad--/util"
	"github.com/gdamore/tcell/v2"
)

const APP_NAME = "Notepad--"

func main() {
	// Get command line arguments
	commandArgs := os.Args
	
	// Initialise variables
	filename := ""
	initialText := ""
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
			initialText = string(data)
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

	appstate := util.AppState{
		AppName: APP_NAME,
		Filename: filename,
		FileModified: false,
		Screen: screen,
		BarStyle: defaultStyle,
		TextboxStyle: defaultStyle,
		ButtonStyle: defaultStyle,
		ButtonActiveStyle: defaultStyle.Reverse(true),
		Options: util.Options{
			LineEndMode: "CRLF",
			Encoding: "UTF-8",
			WordWrap: true,
		},
	}

	// Setup Screen
	screen.Clear()
	menubar := tui.NewMenuBar(2, 3, &appstate)
	textbox := tui.NewTextbox(4, initialText, &appstate)
	elements := []tui.TUIElem{
		tui.NewTitleBar(0, 1, &appstate, textbox),
		menubar,
		textbox,
		tui.NewStatusBar(textbox, &appstate),
	}
	activeElem := elements[0]
		
	// Event Loop
	save := func() {
		// TODO: Add a separate prompt function to handle both the later prompt and the quit prompt
		if !appstate.FileModified {
			return
		}
		content := textbox.Content()

		if appstate.Filename == "" {
			if len(content) == 0 {
				return
			}
			panic("not implemented")
			//filename := util.GetTemporaryTitle(textbox.Content())
			// TODO: Prompt to save (i.e. create) new file based on the temp title
		}
		
		err := os.WriteFile(
			appstate.Filename,
			[]byte(textbox.Content()),
			0660, // RW for user and group
		)
		if err != nil {
			// TODO: Save backup before crash
			panic("Error Saving File.")
		}
		appstate.FileModified = false
	}
	quit := func() {
		// Save Dialog
		if appstate.FileModified {
			//TODO: Add dialog
		}
		// Catch panics
		maybePanic := recover()
		screen.SetCursorStyle(tcell.CursorStyleDefault)
		screen.Fini()

		if maybePanic != nil {
			cursor_x, cursor_y := textbox.GetCursorXY()
			scr_w, scr_h := screen.Size()
			panic(fmt.Sprintf("%+v\nCurrent Screen Size: (%d, %d)\nCursor Data: (%d, %d), %d. Left: %d", maybePanic, scr_w, scr_h, cursor_x, cursor_y, textbox.GetCursorIndex(), textbox.GetLeft()))
		}
		os.Exit(0)
	}
	defer quit()

	screen.SetCursorStyle(tcell.CursorStyleBlinkingBar)

renderLoop:
	for {
		screen.Show()

		// Identify active element
		for _, elem := range(elements) {
			if elem.IsActive() {
				activeElem = elem
				break
			}
		}

		// Process Events
		event := screen.PollEvent()
		switch eventType := event.(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			// Note: SHIFT doesn't appear to work.
			mod, key, ch := eventType.Modifiers(), eventType.Key(), eventType.Rune()

			if mod & tcell.ModAlt != 0 {
				menubar.Focus()
				textbox.Unfocus()
			} else {
				menubar.Unfocus()
				textbox.Focus()
			}

			// Arrow Keys: Get active element and movecursorindex there
			if key == tcell.KeyLeft {
				activeElem.SetCursorIndex(activeElem.GetCursorIndex() - 1)
				break
			} else if key == tcell.KeyRight {
				activeElem.SetCursorIndex(activeElem.GetCursorIndex() + 1)
				break
			}

			// Ctrl-S: Save
			if key == tcell.KeyCtrlS {
				save()
				break
			}
			
			// Ctrl-W: Exit (if no more tabs left)
			if key == tcell.KeyCtrlW {
				break renderLoop
			}

			// (Debug: ESC to Exit)
			if key == tcell.KeyEscape {
				break renderLoop
			}

			// If no modifiers are pressed, it's normal input
			if (mod&tcell.ModCtrl == 0) && (mod&tcell.ModAlt == 0) && (mod&tcell.ModMeta == 0) {
				if textbox.IsActive() {
					if key == tcell.KeyBackspace || key == tcell.KeyBackspace2 {
						textbox.Backspace()
					} else if key == tcell.KeyDelete {
						textbox.Delete()
					} else if key == tcell.KeyEnter {
						textbox.Insert('\n')
						
					} else if key == tcell.KeyTab {
						textbox.Insert('\t')
					} else {
						textbox.Insert(ch)
					}

					if !appstate.FileModified {
						appstate.FileModified = true
					}

				} else if menubar.IsActive() {
					//TODO: Add dropdown for menu
				}
			}
		}

		// Draw Screen
		screen.Clear()
		
		//// TUI
		for _, elem := range(elements) {
			elem.Draw()
		}
		
		screen.Show()
		
	}

}
