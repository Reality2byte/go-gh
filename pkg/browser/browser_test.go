package browser

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/cli/go-gh/v2/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GH_WANT_HELPER_PROCESS") != "1" {
		return
	}
	fmt.Fprintf(os.Stdout, "%v", os.Args[3:])
	os.Exit(0)
}

func TestBrowse(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		launcher string
		expected string
		setup    func(*testing.T) error
		wantErr  bool
	}{
		{
			name:     "Explicit `http` URL works",
			url:      "http://github.com",
			launcher: fmt.Sprintf("%q -test.run=TestHelperProcess -- explicit http", os.Args[0]),
			expected: "[explicit http http://github.com]",
		},
		{
			name:     "Explicit `https` URL works",
			url:      "https://github.com",
			launcher: fmt.Sprintf("%q -test.run=TestHelperProcess -- explicit https", os.Args[0]),
			expected: "[explicit https https://github.com]",
		},
		{
			name:     "Explicit `vscode` URL works",
			url:      "vscode:extension/GitHub.copilot",
			launcher: fmt.Sprintf("%q -test.run=TestHelperProcess -- explicit vscode", os.Args[0]),
			expected: "[explicit vscode vscode:extension/GitHub.copilot]",
		},
		{
			name:     "Explicit `vscode-insiders` URL works",
			url:      "vscode-insiders:extension/GitHub.copilot",
			launcher: fmt.Sprintf("%q -test.run=TestHelperProcess -- explicit vscode-insiders", os.Args[0]),
			expected: "[explicit vscode-insiders vscode-insiders:extension/GitHub.copilot]",
		},
		{
			name:     "Implicit `https` URL works",
			url:      "github.com",
			launcher: fmt.Sprintf("%q -test.run=TestHelperProcess -- implicit https", os.Args[0]),
			expected: "[implicit https github.com]",
		},
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
		{
			name:    "Explicit absolute Windows file URL errors",
			url:     `C:\Windows\System32\calc.exe`,
			wantErr: true,
		},
		{
			name: "Implicit relative file URL errors",
			url:  "poc.command",
			setup: func(t *testing.T) error {
				// Capture current working directory for test cleanup
				cwd, err := os.Getwd()
				if err != nil {
					return err
				}

				// Create temporary directory containing relative executable for testing
				tempDir := t.TempDir()
				err = os.Chdir(tempDir)
				if err != nil {
					return err
				}

				path := filepath.Join(tempDir, "poc.command")
				err = os.WriteFile(path, []byte("#!/bin/bash\necho hello"), 0755)
				if err != nil {
					return err
				}

				// Restore original working directory after test
				t.Cleanup(func() {
					_ = os.Chdir(cwd)
				})

				return nil
			},
			wantErr: true,
		},
		{
			name: "Implicit relative directory URL errors",
			url:  "poc.command",
			setup: func(t *testing.T) error {
				// Capture current working directory for test cleanup
				cwd, err := os.Getwd()
				if err != nil {
					return err
				}

				// Create temporary directory containing relative executable for testing
				tempDir := t.TempDir()
				err = os.Chdir(tempDir)
				if err != nil {
					return err
				}

				path := filepath.Join(tempDir, "Fake.app")
				err = os.Mkdir(path, 0755)
				if err != nil {
					return err
				}

				path = filepath.Join(path, "poc.command")
				err = os.WriteFile(path, []byte("#!/bin/bash\necho hello"), 0755)
				if err != nil {
					return err
				}

				// Restore original working directory after test
				t.Cleanup(func() {
					_ = os.Chdir(cwd)
				})

				return nil
			},
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
			err := b.browse(tt.url, []string{"GH_WANT_HELPER_PROCESS=1"})

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

func TestResolveLauncher(t *testing.T) {
	tests := []struct {
		name         string
		env          map[string]string
		config       *config.Config
		wantLauncher string
	}{
		{
			name: "GH_BROWSER set",
			env: map[string]string{
				"GH_BROWSER": "GH_BROWSER",
			},
			wantLauncher: "GH_BROWSER",
		},
		{
			name:         "config browser set",
			config:       config.ReadFromString("browser: CONFIG_BROWSER"),
			wantLauncher: "CONFIG_BROWSER",
		},
		{
			name: "BROWSER set",
			env: map[string]string{
				"BROWSER": "BROWSER",
			},
			wantLauncher: "BROWSER",
		},
		{
			name: "GH_BROWSER and config browser set",
			env: map[string]string{
				"GH_BROWSER": "GH_BROWSER",
			},
			config:       config.ReadFromString("browser: CONFIG_BROWSER"),
			wantLauncher: "GH_BROWSER",
		},
		{
			name: "config browser and BROWSER set",
			env: map[string]string{
				"BROWSER": "BROWSER",
			},
			config:       config.ReadFromString("browser: CONFIG_BROWSER"),
			wantLauncher: "CONFIG_BROWSER",
		},
		{
			name: "GH_BROWSER and BROWSER set",
			env: map[string]string{
				"BROWSER":    "BROWSER",
				"GH_BROWSER": "GH_BROWSER",
			},
			wantLauncher: "GH_BROWSER",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.env != nil {
				for k, v := range tt.env {
					t.Setenv(k, v)
				}
			}
			if tt.config != nil {
				old := config.Read
				config.Read = func(_ *config.Config) (*config.Config, error) {
					return tt.config, nil
				}
				defer func() { config.Read = old }()
			}
			launcher := resolveLauncher()
			assert.Equal(t, tt.wantLauncher, launcher)
		})
	}
}
