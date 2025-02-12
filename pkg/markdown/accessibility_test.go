package markdown

import (
	"testing"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/glamour/ansi"
	"github.com/stretchr/testify/assert"
)

func TestAnsiColorName(t *testing.T) {
	tests := []struct {
		name string
		c    AnsiColor
		want string
	}{
		{
			name: "black",
			c:    Black,
			want: "black",
		},
		{
			name: "bright black",
			c:    BrightBlack,
			want: "brightBlack",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.c.Name())
		})
	}
}

func TestAnsiColorValue(t *testing.T) {
	tests := []struct {
		name string
		c    AnsiColor
		want string
	}{
		{
			name: "red",
			c:    Red,
			want: "1",
		},
		{
			name: "bright red",
			c:    BrightRed,
			want: "9",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.c.Value())
		})
	}
}

func TestAnsiColorValuePtr(t *testing.T) {
	tests := []struct {
		name string
		c    AnsiColor
		want *string
	}{
		{
			name: "green",
			c:    Green,
			want: strPtr("2"),
		},
		{
			name: "bright green",
			c:    BrightGreen,
			want: strPtr("10"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.c.ValuePtr())
		})
	}
}

func TestAccessibleStyleConfig(t *testing.T) {
	tests := []struct {
		name  string
		theme string
		want  ansi.StyleConfig
	}{
		{
			name:  "light",
			theme: "light",
			want:  accessibleLightStyleConfig(),
		},
		{
			name:  "dark",
			theme: "dark",
			want:  accessibleDarkStyleConfig(),
		},
		{
			name:  "fallback",
			theme: "foo",
			want:  ansi.StyleConfig{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, AccessibleStyleConfig(tt.theme))
		})
	}
}

func Test_accessibleDarkStyleConfig(t *testing.T) {
	cfg := accessibleDarkStyleConfig()
	assert.Equal(t, White.ValuePtr(), cfg.Document.StylePrimitive.Color)
	assert.Equal(t, Black.ValuePtr(), cfg.Link.Color)
	assert.Equal(t, BrightCyan.ValuePtr(), cfg.LinkText.Color)
	assert.Equal(t, BrightMagenta.ValuePtr(), cfg.Heading.StylePrimitive.Color)
	assert.Equal(t, BrightWhite.ValuePtr(), cfg.H1.StylePrimitive.Color)
	assert.Equal(t, BrightBlue.ValuePtr(), cfg.H1.StylePrimitive.BackgroundColor)
	assert.Equal(t, BrightWhite.ValuePtr(), cfg.Code.BackgroundColor)
	assert.Equal(t, Red.ValuePtr(), cfg.Code.Color)
	assert.Equal(t, BrightMagenta.ValuePtr(), cfg.Image.Color)
	assert.Equal(t, White.ValuePtr(), cfg.HorizontalRule.Color)

	// Test that we haven't changed the original style
	assert.Equal(t, glamour.LightStyleConfig.H2, cfg.H2)
}

func Test_accessibleLightStyleConfig(t *testing.T) {
	cfg := accessibleLightStyleConfig()
	assert.Equal(t, Black.ValuePtr(), cfg.Document.StylePrimitive.Color)
	assert.Equal(t, Black.ValuePtr(), cfg.Link.Color)
	assert.Equal(t, BrightBlue.ValuePtr(), cfg.LinkText.Color)
	assert.Equal(t, Magenta.ValuePtr(), cfg.Heading.StylePrimitive.Color)
	assert.Equal(t, BrightWhite.ValuePtr(), cfg.H1.StylePrimitive.Color)
	assert.Equal(t, Blue.ValuePtr(), cfg.H1.StylePrimitive.BackgroundColor)
	assert.Equal(t, BrightWhite.ValuePtr(), cfg.Code.BackgroundColor)
	assert.Equal(t, Red.ValuePtr(), cfg.Code.Color)
	assert.Equal(t, Magenta.ValuePtr(), cfg.Image.Color)
	assert.Equal(t, White.ValuePtr(), cfg.HorizontalRule.Color)

	// Test that we haven't changed the original style
	assert.Equal(t, glamour.LightStyleConfig.H2, cfg.H2)
}
