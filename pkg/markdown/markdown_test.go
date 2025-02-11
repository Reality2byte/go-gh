package markdown

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/cli/go-gh/v2/pkg/accessibility"
	"github.com/stretchr/testify/assert"
)

// glamour theme colors found at https://github.com/charmbracelet/glamour/tree/master/styles
const (
	glamourLightH2Color     = "27"
	glamourDarkH2Color      = "39"
	customH2Color           = "61"
	magentaEscapeCode       = "35"
	brightMagentaEscapeCode = "95"
)

func Test_Render(t *testing.T) {
	t.Setenv("GLAMOUR_STYLE", "")
	tests := []struct {
		name             string
		text             string
		theme            string
		styleEnvVar      string
		accessibleEnvVar string
		wantOut          string
	}{
		{
			name:    "light style",
			text:    "## h2",
			theme:   "light",
			wantOut: fmt.Sprintf("\x1b[0m\x1b[38;5;%s;1mh2", glamourLightH2Color),
		},
		{
			name:    "dark style",
			text:    "## h2",
			theme:   "dark",
			wantOut: fmt.Sprintf("\x1b[0m\x1b[38;5;%s;1mh2", glamourDarkH2Color),
		},
		{
			name:    "notty style",
			text:    "## h2",
			theme:   "none",
			wantOut: "## h2", // no tty should maintain the original text
		},
		{
			name:        "when the style env var is set, we override the theme whith that style",
			text:        "## h2",
			theme:       "light",
			styleEnvVar: "customStyle",
			wantOut:     fmt.Sprintf("\x1b[0m\x1b[38;5;%s;1mh2", customH2Color),
		},
		{
			name:             "when the accessible env var is set and the light theme is selected, we use the light accessible style",
			text:             "## h2",
			theme:            "light",
			accessibleEnvVar: "true",
			wantOut:          fmt.Sprintf("\x1b[0m\x1b[%s;1mh2", magentaEscapeCode),
		},
		{
			name:             "when the accessible env var is set and the dark theme is selected, we use the dark accessible style",
			text:             "## h2",
			theme:            "dark",
			accessibleEnvVar: "true",
			wantOut:          fmt.Sprintf("\x1b[0m\x1b[%s;1mh2", brightMagentaEscapeCode),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(accessibility.ACCESSIBILITY_ENV, tt.accessibleEnvVar)

			if tt.styleEnvVar != "" {
				tmpDir := t.TempDir()
				path := filepath.Join(tmpDir, fmt.Sprintf("%s.json", tt.styleEnvVar))
				err := os.WriteFile(path, []byte(customGlamourStyle(t)), 0644)
				if err != nil {
					t.Fatal(err)
				}
				t.Setenv("GLAMOUR_STYLE", path)
			}

			out, err := Render(tt.text, WithTheme(tt.theme))
			assert.NoError(t, err)
			assert.Contains(t, out, tt.wantOut)
		})
	}
}

func customGlamourStyle(t *testing.T) string {
	t.Helper()
	return fmt.Sprintf(`
{
	"heading": {
		"block_suffix": "\n",
		"color": "%s",
		"bold": true
	},
	"h2": {
		"prefix": "## "
	}
}`, customH2Color)
}
