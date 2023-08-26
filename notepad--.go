package main

import (
	"log"
	"os"
	"github.com/Rye123/notepad--/app"
	"github.com/Rye123/notepad--/tui"
	"github.com/Rye123/notepad--/util"
	"github.com/gdamore/tcell/v2"
)

func cleanup(screen tcell.Screen) {
	// Catch any panic
	maybePanic := recover()
	screen.SetCursorStyle(tcell.CursorStyleDefault)
	screen.Fini()

	if maybePanic != nil {
		log.Fatalf("%+v\n", maybePanic)
	}
}

func main() {
	// Get command line arguments
	commandArgs := os.Args
	
	// Initialise variables
	filename := ""
	if len(commandArgs) > 1 {
		filename = commandArgs[1]
	}
	
	// Initialise screen
	screen, err := tcell.NewScreen()
	defer cleanup(screen)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	// Initialise app state
	options := util.Options{
		LineEndMode: "CRLF",
		Encoding: "UTF-8",
		WordWrap: true,
	}
	appstate := util.InitialiseAppState(screen, filename, options)

	// Setup Screen
	screen.Clear()
	titlebar := tui.NewTitleBar(appstate)
	menubar := tui.NewMenuBar(appstate)
	textbox := tui.NewTextbox(appstate)
	statusbar := tui.NewStatusBar(appstate, textbox)

	elems := []tui.TUIElem{
		titlebar,
		menubar,
		textbox,
		statusbar,
	}
		
	// Event Loop
	ui := app.NewUI(appstate, elems)
	ui.Display()
}
