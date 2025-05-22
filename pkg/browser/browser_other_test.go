//go:build !windows

package browser

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBrowseOthers covers unique Unix and Linux specific test cases that go beyond TestBrowse.
func TestBrowseOthers(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		launcher string
		expected string
		setup    func(*testing.T) error
		wantErr  bool
	}{
		{
			name:    "Explicit absolute `file://` URL errors",
			url:     "file:///System/Applications/Calculator.app",
			wantErr: true,
		},
		{
			name:    "Implicit absolute file URL errors",
			url:     "/bin/sh",
			wantErr: true,
		},
		{
			name:    "Implicit absolute directory URL errors",
			url:     "/System/Applications/Calculator.app",
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
