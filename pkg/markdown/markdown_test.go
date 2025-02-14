package markdown

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cli/go-gh/v2/pkg/accessibility"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// glamour theme colors found at https://github.com/charmbracelet/glamour/tree/master/styles
const (
	glamourLightH2_8bitColorSeq = "\x1b[38;5;27;"
	glamourDarkH2_8bitColorSeq  = "\x1b[38;5;39;"
	customH2_8bitColorSeq       = "\x1b[38;5;61;"
	magenta_4bitColorSeq        = "\x1b[35;"
	brightMagenta_4bitColorSeq  = "\x1b[95;"
)

// Test_Render verifies that the proper ANSI color codes are applied to the rendered
// markdown by examining the ANSI escape sequences in the output for the correct color
// match. For more information on ANSI color codes, see
// https://en.wikipedia.org/wiki/ANSI_escape_code#8-bit
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
			name:    "when the light theme is selected, the h2 renders using the 8-bit blue 27 provided by glamour",
			text:    "## h2",
			theme:   "light",
			wantOut: fmt.Sprintf("%s1mh2", glamourLightH2_8bitColorSeq),
		},
		{
			name:    "when the dark theme is selected, the h2 renders using the 8-bit blue 39 provided by glamour",
			text:    "## h2",
			theme:   "dark",
			wantOut: fmt.Sprintf("%s1mh2", glamourDarkH2_8bitColorSeq),
		},
		{
			name:    "when no theme is selected, the h2 renders in plain text without ansi coloring",
			text:    "## h2",
			theme:   "none",
			wantOut: "## h2",
		},
		{
			name:        "when the style env var is set, we override the theme with that style",
			text:        "## h2",
			theme:       "light",
			styleEnvVar: "customStyle",
			wantOut:     fmt.Sprintf("%s1mh2", customH2_8bitColorSeq),
		},
		{
			name:             "when the accessible env var is set and the light theme is selected, the h2 renders using the 4-bit magenta provided by the light accessible style",
			text:             "## h2",
			theme:            "light",
			accessibleEnvVar: "true",
			wantOut:          fmt.Sprintf("%s1mh2", magenta_4bitColorSeq),
		},
		{
			name:             "when the accessible env var is set and the dark theme is selected, the h2 renders using the 4-bit bright magenta provided by the dark accessible style",
			text:             "## h2",
			theme:            "dark",
			accessibleEnvVar: "true",
			wantOut:          fmt.Sprintf("%s1mh2", brightMagenta_4bitColorSeq),
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
			require.NoError(t, err)
			assert.Contains(t, out, tt.wantOut)
		})
	}
}

func customGlamourStyle(t *testing.T) string {
	t.Helper()
	colorCode := strings.Split(customH2_8bitColorSeq, ";")[2]
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
}`, colorCode)
}
