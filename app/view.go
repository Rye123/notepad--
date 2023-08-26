// Contains the main application view and controller logic.
package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/Rye123/notepad--/tui"
	"github.com/Rye123/notepad--/util"
)

type UI struct {
	appstate *util.AppState
	elements []tui.TUIElem	
}

func NewUI(appstate *util.AppState, elements []tui.TUIElem) *UI {
	return &UI{appstate, elements}
}

func (ui *UI) Display() {
	defer ui.Quit()
	screen := ui.appstate.Screen
	screen.SetCursorStyle(tcell.CursorStyleBlinkingBar)

renderLoop:
	for {
		// Process Events
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			quit := ui.handleKeyEvent(ev)
			if quit {
				break renderLoop
			}
		}

		// Draw Screen
		screen.Clear()

		for _, elem := range(ui.elements) {
			elem.Draw()
		}

		screen.Show()
	}
}

// Save. Will only save if this is a modified, pre-existing file. If it doesn't exist beforehand (i.e. filename == ""), SaveAs is called.
func (ui *UI) Save() {
	if !ui.appstate.FileModified {
		return
	}

	// Check if file existed before
	if ui.appstate.Filename == "" {
		ui.SaveAs()
		return
	}

	ui.appstate.Save()
	
}

// Save, with a prompt. Will return immediately if the file never existed before and no content has been written
func (ui *UI) SaveAs() {
	filename := ui.appstate.Filename
	if filename == "" {
		filename = util.GetTemporaryTitle(ui.appstate.TextBuffer.String())
		if filename == "" {
			return
		}
	}

	ui.appstate.Filename = filename
	//TODO: Prompt, then call save from appstate.
	panic("SaveAs not implemented.")
}

func (ui *UI) Quit() {
	if ui.appstate.FileModified {
		ui.SaveAs()
	}
}

// Handles a key event. Returns true if UI is to quit after returning from this function.
func (ui *UI) handleKeyEvent(keyEvent *tcell.EventKey) bool {
	mod, key := keyEvent.Modifiers(), keyEvent.Key()

	// CONTROL KEYS
	switch key {
	case tcell.KeyCtrlS: // Ctrl-S: Save, Ctrl-Alt-S: Save As
		if mod & tcell.ModAlt != 0 {
			ui.SaveAs()
		} else {
			ui.Save()
		}
		return false
	case tcell.KeyCtrlW: // Ctrl-W: Close
		return true
	case tcell.KeyEscape: // ESC: Refocus on textbox
		for _, elem := range(ui.elements) {
			if _, ok := elem.(*tui.Textbox); ok {
				elem.Focus()
			} else {
				elem.Unfocus()
			}
		}
		
		return false
	}

	// If Alt is pressed along with a key, control handed to menubar.
	if mod & tcell.ModAlt != 0 {
		for _, elem := range(ui.elements) {
			if _, ok := elem.(*tui.MenuBar); ok {
				elem.Focus()
			} else {
				elem.Unfocus()
			}
		}
	}

	// Hand over event to all the elements
	for _, elem := range(ui.elements) {
		elem.HandleKey(keyEvent)
	}

	return false
}
