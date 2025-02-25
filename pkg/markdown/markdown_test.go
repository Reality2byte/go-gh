package markdown

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/cli/go-gh/v2/pkg/accessibility"
	ansi "github.com/leaanthony/go-ansi-parser"
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

func Test_Render_Codeblocks(t *testing.T) {
	t.Setenv("GLAMOUR_STYLE", "")

	text := heredoc.Docf(`
		%[1]s%[1]s%[1]sgo
		package main

		import (
			"fmt"
		)

		func main() {
			fmt.Println("Hello, world!")
		}
		%[1]s%[1]s%[1]s
	`, "`")

	tests := []struct {
		name       string
		text       string
		theme      string
		accessible bool
	}{
		{
			name:  "when the light theme is selected, the codeblock renders using 8-bit colors",
			text:  text,
			theme: "light",
		},
		{
			name:  "when the dark theme is selected, the codeblock renders using 8-bit colors",
			text:  text,
			theme: "dark",
		},
		{
			name:       "when the accessible env var is set and the light theme is selected, the codeblock renders using 4-bit colors",
			text:       text,
			theme:      "light",
			accessible: true,
		},
		{
			name:       "when the accessible env var is set and the dark theme is selected, the codeblock renders using 4-bit colors",
			text:       text,
			theme:      "dark",
			accessible: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.accessible {
				t.Setenv(accessibility.ACCESSIBILITY_ENV, "true")
			}

			out, err := Render(tt.text, WithTheme(tt.theme))
			require.NoError(t, err)

			styledText, err := ansi.Parse(out)
			require.NoError(t, err)

			for _, st := range styledText {
				if tt.accessible {
					require.Equalf(t, st.ColourMode, ansi.Default, "Inaccessible color found in '%s' at %d", st, st.Offset)
					require.Falsef(t, st.Faint(), "Inaccessible style found in '%s' at %d", st, st.Offset)
				}
			}
		})
	}
}

// Test_Render verifies that the proper ANSI color codes are applied to the rendered
// markdown by examining the ANSI escape sequences in the output for the correct color
// match. For more information on ANSI color codes, see
// https://en.wikipedia.org/wiki/ANSI_escape_code#8-bit
func Test_Render(t *testing.T) {
	codeBlock := heredoc.Docf(`
		%[1]s%[1]s%[1]sgo
		fmt.Println("Hello, world!")
		%[1]s%[1]s%[1]s
	`, "`")

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
		{
			name:    "when the light theme is selected, the codeblock renders using 8-bit colors",
			text:    codeBlock,
			theme:   "light",
			wantOut: "\x1b[0m\x1b[38;5;235mfmt\x1b[0m\x1b[38;5;210m.\x1b[0m\x1b[38;5;35mPrintln\x1b[0m\x1b[38;5;210m(\x1b[0m\x1b[38;5;95m\"Hello, world!\"\x1b[0m\x1b[38;5;210m)\x1b[0m",
		},
		{
			name:    "when the dark theme is selected, the codeblock renders using 8-bit colors",
			text:    codeBlock,
			theme:   "dark",
			wantOut: "\x1b[0m\x1b[38;5;251mfmt\x1b[0m\x1b[38;5;187m.\x1b[0m\x1b[38;5;42mPrintln\x1b[0m\x1b[38;5;187m(\x1b[0m\x1b[38;5;173m\"Hello, world!\"\x1b[0m\x1b[38;5;187m)\x1b[0m",
		},
		{
			name:             "when the accessible env var is set and the light theme is selected, the codeblock renders using 4-bit colors",
			text:             codeBlock,
			theme:            "light",
			accessibleEnvVar: "true",
			wantOut:          "\x1b[0m\x1b[30mfmt\x1b[0m\x1b[33m.\x1b[0m\x1b[36mPrintln\x1b[0m\x1b[33m(\x1b[0m\x1b[90m\"Hello, world!\"\x1b[0m\x1b[33m)\x1b[0m",
		},
		{
			name:             "when the accessible env var is set and the dark theme is selected, the codeblock renders using 4-bit colors",
			text:             codeBlock,
			theme:            "dark",
			accessibleEnvVar: "true",
			wantOut:          "\x1b[0m\x1b[37mfmt\x1b[0m\x1b[37m.\x1b[0m\x1b[36mPrintln\x1b[0m\x1b[37m(\x1b[0m\x1b[33m\"Hello, world!\"\x1b[0m\x1b[37m)\x1b[0m",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Unregister cached chroma style causing codeblock tests to fail based on previous theme
			delete(styles.Registry, "charm")
			t.Setenv(accessibility.ACCESSIBILITY_ENV, tt.accessibleEnvVar)

			if tt.styleEnvVar == "" {
				t.Setenv("GLAMOUR_STYLE", "")
			} else {
				path := filepath.Join(t.TempDir(), fmt.Sprintf("%s.json", tt.styleEnvVar))
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
