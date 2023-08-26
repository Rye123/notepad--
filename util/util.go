// Provides necessary models for the application.
package util

import (
	"os"
	"strings"
	"log"
	"errors"
	"github.com/gdamore/tcell/v2"
	"github.com/Rye123/notepad--/textbuffer"
)

const APP_NAME = "Notepad--"

// Options for Notepad--
type Options struct {
	LineEndMode string
	Encoding string
	WordWrap bool
}

func (opt *Options) LineEndModeString() string {
	switch opt.LineEndMode {
	case "LF":
		return "Unix (LF)"
	case "CRLF":
		return "Windows (CRLF)"
	default:
		return "Windows (CRLF)"
	}
}

// Returns the line-end character.
func (opt *Options) LE() string {
	switch opt.LineEndMode {
	case "LF":
		return "\n"
	case "CRLF":
		return "\r\n"
	default:
		return "\r\n"
	}
}

type AppState struct {
	AppName string
	Filename string
	FileModified bool
	Screen tcell.Screen
	TextBuffer textbuffer.TextBuffer
	BarStyle tcell.Style
	TextboxStyle tcell.Style
	ButtonStyle tcell.Style
	ButtonActiveStyle tcell.Style
	Options Options
}

func InitialiseAppState(screen tcell.Screen, filename string, options Options) *AppState {
	defaultStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	screen.SetStyle(defaultStyle)
	screen.SetCursorStyle(tcell.CursorStyleDefault)

	textbuffer := textbuffer.NewGapBuffer()

	// Read file, if given
	if len(filename) > 0 {
		data, err := os.ReadFile(filename)
		if err != nil {
			// Report if it's an error unrelated to the file not existing
			if !errors.Is(err, os.ErrNotExist) {
				log.Fatalf("%+v", err)
			}
		} else {
			initialText := string(data)
			// Load into buffer
			if len(initialText) > 0 {
				textbuffer.Append(initialText)
			}
		}
	}

	appstate := AppState{
		AppName: APP_NAME,
		Filename: filename,
		FileModified: false,
		Screen: screen,
		TextBuffer: textbuffer,
		BarStyle: defaultStyle,
		TextboxStyle: defaultStyle,
		ButtonStyle: defaultStyle,
		ButtonActiveStyle: defaultStyle.Reverse(true),
		Options: options,
	}

	return &appstate
}

// Used to generate a temporary title if no file was used
func GetTemporaryTitle(content string) string {
	title, _, _ := strings.Cut(content, "\n")
	title = strings.TrimSpace(title)
	if len(title) > 45 {
		title = title[:45]
	}
	return title
}
