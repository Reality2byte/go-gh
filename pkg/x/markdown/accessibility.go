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

func (a glamourStyleColor) code() *string {
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
	cfg.Document.StylePrimitive.Color = white.code()

	// Link colors
	cfg.Link.Color = brightCyan.code()

	// Heading colors
	cfg.Heading.StylePrimitive.Color = brightMagenta.code()
	cfg.H1.StylePrimitive.Color = brightWhite.code()
	cfg.H1.StylePrimitive.BackgroundColor = brightBlue.code()

	// Code colors
	cfg.Code.BackgroundColor = brightWhite.code()
	cfg.Code.Color = red.code()

	// Image colors
	cfg.Image.Color = brightMagenta.code()

	// Horizontal rule colors
	cfg.HorizontalRule.Color = white.code()

	return cfg
}

func accessibleLightStyleConfig() ansi.StyleConfig {
	cfg := styles.LightStyleConfig

	// Text color
	cfg.Document.StylePrimitive.Color = black.code()

	// Link colors
	cfg.Link.Color = brightBlue.code()

	// Heading colors
	cfg.Heading.StylePrimitive.Color = magenta.code()
	cfg.H1.StylePrimitive.Color = brightWhite.code()
	cfg.H1.StylePrimitive.BackgroundColor = blue.code()

	// Code colors
	cfg.Code.BackgroundColor = brightWhite.code()
	cfg.Code.Color = red.code()

	// Image colors
	cfg.Image.Color = magenta.code()

	// Horizontal rule colors
	cfg.HorizontalRule.Color = white.code()

	return cfg
}
