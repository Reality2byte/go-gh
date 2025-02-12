package markdown

import (
	"strconv"
)

type AnsiColor int

const (
	Black AnsiColor = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	BrightBlack
	BrightRed
	BrightGreen
	BrightYellow
	BrightBlue
	BrightMagenta
	BrightCyan
	BrightWhite
)

func (a AnsiColor) Name() string {
	return [...]string{
		"black",
		"red",
		"green",
		"yellow",
		"blue",
		"magenta",
		"cyan",
		"white",
		"brightBlack",
		"brightRed",
		"brightGreen",
		"brightYellow",
		"brightBlue",
		"brightMagenta",
		"brightCyan",
		"brightWhite",
	}[a]
}

func (a AnsiColor) Value() string {
	return strconv.Itoa(int(a))
}

func (a AnsiColor) ValuePtr() *string {
	return strPtr(a.Value())
}

func strPtr(s string) *string {
	return &s
}
