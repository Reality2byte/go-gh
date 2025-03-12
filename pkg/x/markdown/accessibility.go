package markdown

import (
	"strconv"

	"github.com/charmbracelet/glamour/ansi"
	"github.com/charmbracelet/glamour/styles"
)

// glamourStyleColor represents color codes used to customize glamour style elements.
type glamourStyleColor int

// Do not change the order of the following glamour color constants,
// which matches 4-bit colors with their respective color codes.
const (
	black glamourStyleColor = iota
	red
	green
	yellow
	blue
	magenta
	cyan
	white
	brightBlack
	brightRed
	brightGreen
	brightYellow
	brightBlue
	brightMagenta
	brightCyan
	brightWhite
)

func (a glamourStyleColor) Code() *string {
	s := strconv.Itoa(int(a))
	return &s
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
	cfg := styles.DarkStyleConfig

	// Text color
	cfg.Document.StylePrimitive.Color = white.Code()

	// Link colors
	cfg.Link.Color = brightCyan.Code()

	// Heading colors
	cfg.Heading.StylePrimitive.Color = brightMagenta.Code()
	cfg.H1.StylePrimitive.Color = brightWhite.Code()
	cfg.H1.StylePrimitive.BackgroundColor = brightBlue.Code()

	// Code colors
	cfg.Code.BackgroundColor = brightWhite.Code()
	cfg.Code.Color = red.Code()

	// Image colors
	cfg.Image.Color = brightMagenta.Code()

	// Horizontal rule colors
	cfg.HorizontalRule.Color = white.Code()

	return cfg
}

func accessibleLightStyleConfig() ansi.StyleConfig {
	cfg := styles.LightStyleConfig

	// Text color
	cfg.Document.StylePrimitive.Color = black.Code()

	// Link colors
	cfg.Link.Color = brightBlue.Code()

	// Heading colors
	cfg.Heading.StylePrimitive.Color = magenta.Code()
	cfg.H1.StylePrimitive.Color = brightWhite.Code()
	cfg.H1.StylePrimitive.BackgroundColor = blue.Code()

	// Code colors
	cfg.Code.BackgroundColor = brightWhite.Code()
	cfg.Code.Color = red.Code()

	// Image colors
	cfg.Image.Color = magenta.Code()

	// Horizontal rule colors
	cfg.HorizontalRule.Color = white.Code()

	return cfg
}
