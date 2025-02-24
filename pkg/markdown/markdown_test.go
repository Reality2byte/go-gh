package markdown

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc"
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

	// TODO: Include a little more context on SGR including link to https://en.wikipedia.org/wiki/ANSI_escape_code#Select_Graphic_Rendition_parameters for more info
	// sgrSequencePattern identifies ANSI escape sequences containing display attributes
	// that affect the color, emphasis, and other aspects of displaying text.
	sgrSequencePattern = `\x1b\[(.+?)m`

	// sgrAttributePattern analyzes separate display attributes within an ANSI escape sequence
	// for detecting color depth (3-bit, 4-bit, 8-bit, 24-bit, etc) or other effects.
	//
	// This is a separate regex from sgrSequencePattern as `FindAllStringSubmatch()` does not
	// handle repeating capture groups well.
	sgrAttributePattern = `;?(?P<sequence>\d+)` // TODO: change the `sequence` note; if we aren't actually using the name group, remove it
)

func Test_Render_Codeblocks(t *testing.T) {
	t.Setenv("GLAMOUR_STYLE", "")

	sequencesRegex := regexp.MustCompile(sgrSequencePattern)
	attributesRegex := regexp.MustCompile(sgrAttributePattern)
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
			sequences := sequencesRegex.FindAllStringSubmatch(out, -1)
			require.NotEmpty(t, sequences, "Failed to find expected SGR sequences in rendered output")

			// TODO: Review use of a module like https://github.com/leaanthony/go-ansi-parser/blob/main/ansi.go#L264 to do the sequence and attribute parsing such that we just have to iterate over the results
			for _, sequence := range sequences {
				attributes := attributesRegex.FindAllStringSubmatch(sequence[1], -1)
				require.NotEmpty(t, attributes, "Failed to extract SGR attributes for testing")

				// Analysis loop handles index incrementing due to unique display attribute situations like color depth
				for i := 0; i < len(attributes); {
					// TODO: Use constants for 38,48 and maybe remove the int conversion
					attribute, err := strconv.Atoi(attributes[i][1])
					require.NoError(t, err, "Failed to convert SGR attribute for testing")

					switch attribute {
					case 38, 48:
						// Display attributes for setting 8-bit and 24-bit foreground and background colors
						colorDepth, err := strconv.Atoi(attributes[i+1][1])
						require.NoError(t, err, "Failed to convert SGR color depth attribute for testing")

						switch colorDepth {
						case 2:
							// 24-bit color display attribute form
							// - ESC[38;2;⟨r⟩;⟨g⟩;⟨b⟩m for foreground colors
							// - ESC[48;2;⟨r⟩;⟨g⟩;⟨b⟩m for background colors
							require.False(t, tt.accessible, "24-bit color is not accessible, customizable")

							color24bitRed, err := strconv.Atoi(attributes[i+2][1])
							require.NoError(t, err, "Failed to convert 24-bit red color value for testing")
							require.True(t, color24bitRed >= 0 && color24bitRed <= 255, "24-bit red color value out of 0-255 range")

							color24bitGreen, err := strconv.Atoi(attributes[i+3][1])
							require.NoError(t, err, "Failed to convert 24-bit green color value for testing")
							require.True(t, color24bitGreen >= 0 && color24bitGreen <= 255, "24-bit green color value out of 0-255 range")

							color24bitBlue, err := strconv.Atoi(attributes[i+4][1])
							require.NoError(t, err, "Failed to convert 24-bit blue color value for testing")
							require.True(t, color24bitBlue >= 0 && color24bitBlue <= 255, "24-bit blue color value out of 0-255 range")

							i += 5
						case 5:
							// 8-bit color display attributes form:
							// - ESC[38;5;⟨n⟩m for foreground colors
							// - ESC[48;5;⟨n⟩m for background colors
							require.False(t, tt.accessible, "8-bit color is not accessible, customizable")

							color8bit, err := strconv.Atoi(attributes[i+2][1])
							require.NoError(t, err, "Failed to convert 8-bit color value for testing")
							require.True(t, color8bit >= 0 && color8bit <= 255, "8-bit color value out of 0-255 range")

							i += 3
						default:
							require.Fail(t, "Unexpected color depth in attribute")
						}
					default:
						// Increment index as this attribute does not affect accessibility currently
						i += 1
					}
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
