// Provides necessary models for the application.
package util

import (
	"github.com/gdamore/tcell/v2"
)

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
	BarStyle tcell.Style
	TextboxStyle tcell.Style
	ButtonStyle tcell.Style
	ButtonActiveStyle tcell.Style
	Options Options
}
