package markdown

import (
	"testing"

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
