package markdown

import (
	"strconv"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/glamour/ansi"
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

func AccessibleStyleConfig(theme string) ansi.StyleConfig {
	switch theme {
	case "light":
		return accessibleLightStyleConfig()
	case "dark":
		return accessibleDarkStyleConfig()
	default:
		return ansi.StyleConfig{}
	}
}

func accessibleDarkStyleConfig() ansi.StyleConfig {
	cfg := glamour.DarkStyleConfig

	// Text color
	cfg.Document.StylePrimitive.Color = White.ValuePtr()

	// Link colors
	cfg.Link.Color = Black.ValuePtr()
	cfg.LinkText.Color = BrightCyan.ValuePtr()

	// Heading colors
	cfg.Heading.StylePrimitive.Color = BrightMagenta.ValuePtr()
	cfg.H1.StylePrimitive.Color = BrightWhite.ValuePtr()
	cfg.H1.StylePrimitive.BackgroundColor = BrightBlue.ValuePtr()

	// Code colors
	cfg.Code.BackgroundColor = BrightWhite.ValuePtr()
	cfg.Code.Color = Red.ValuePtr()

	// Image colors
	cfg.Image.Color = BrightMagenta.ValuePtr()

	// Horizontal rule colors
	cfg.HorizontalRule.Color = White.ValuePtr()

	return cfg
}

func accessibleLightStyleConfig() ansi.StyleConfig {
	cfg := glamour.LightStyleConfig

	// Text color
	cfg.Document.StylePrimitive.Color = Black.ValuePtr()

	// Link colors
	cfg.Link.Color = Black.ValuePtr()
	cfg.LinkText.Color = BrightBlue.ValuePtr()

	// Heading colors
	cfg.Heading.StylePrimitive.Color = Magenta.ValuePtr()
	cfg.H1.StylePrimitive.Color = BrightWhite.ValuePtr()
	cfg.H1.StylePrimitive.BackgroundColor = Blue.ValuePtr()

	// Code colors
	cfg.Code.BackgroundColor = BrightWhite.ValuePtr()
	cfg.Code.Color = Red.ValuePtr()

	// Image colors
	cfg.Image.Color = Magenta.ValuePtr()

	// Horizontal rule colors
	cfg.HorizontalRule.Color = White.ValuePtr()

	return cfg
}

func strPtr(s string) *string {
	return &s
}
