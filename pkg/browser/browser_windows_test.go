//go:build windows

package browser

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBrowseWindows covers unique Windows-specific test cases that go beyond TestBrowse.
func TestBrowseWindows(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		launcher string
		expected string
		setup    func(*testing.T) error
		wantErr  bool
	}{
		{
			name:    "Explicit absolute Windows file URL errors",
			url:     `C:\Windows\System32\calc.exe`,
			wantErr: true,
		},
		{
			name:    "Explicit absolute Windows directory URL errors",
			url:     `C:\Windows\System32`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				err := tt.setup(t)
				require.NoError(t, err)
			}

			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}
			b := Browser{launcher: tt.launcher, stdout: stdout, stderr: stderr}
			err := b.browse(tt.url, []string{"GH_WANT_HELPER_PROCESS=1", fmt.Sprintf("GOCOVERDIR=%s", os.Getenv("GOCOVERDIR"))})

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, stdout.String())
				assert.Equal(t, "", stderr.String())
			}
		})
	}
}
