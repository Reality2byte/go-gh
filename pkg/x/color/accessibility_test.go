package color

import (
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/go-gh/v2/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestIsAccessibleColorsEnabled(t *testing.T) {
	tests := []struct {
		name        string
		envVarValue string
		cfgStr      string
		wantOut     bool
	}{
		{
			name:        "When the accessibility configuration and env var are both unset, it should return false",
			envVarValue: "",
			cfgStr:      "",
			wantOut:     false,
		},
		{
			name:        "When the accessibility configuration is unset but the env var is set to something truthy (not '0' or 'false'), it should return true",
			envVarValue: "1",
			cfgStr:      "",
			wantOut:     true,
		},
		{
			name:        "When the accessibility configuration is unset and the env var returns '0', it should return false",
			envVarValue: "0",
			cfgStr:      "",
			wantOut:     false,
		},
		{
			name:        "When the accessibility configuration is unset and the env var returns 'false', it should return false",
			envVarValue: "false",
			cfgStr:      "",
			wantOut:     false,
		},
		{
			name:        "When the accessibility configuration is set to enabled and the env var is unset, it should return true",
			envVarValue: "",
			cfgStr:      accessibilityEnabledConfig(),
			wantOut:     true,
		},
		{
			name:        "When the accessibility configuration is set to disabled and the env var is unset, it should return false",
			envVarValue: "",
			cfgStr:      accessibilityDisabledConfig(),
			wantOut:     false,
		},
		{
			name:        "When the accessibility configuration is set to disabled and the env var is set to something truthy (not '0' or 'false'), it should return true",
			envVarValue: "true",
			cfgStr:      accessibilityDisabledConfig(),
			wantOut:     true,
		},
		{
			name:        "When the accessibility configuration is set to enabled and the env var is set to '0', it should return false",
			envVarValue: "0",
			cfgStr:      accessibilityEnabledConfig(),
			wantOut:     false,
		},
		{
			name:        "When the accessibility configuration is set to enabled and the env var is set to 'false', it should return false",
			envVarValue: "false",
			cfgStr:      accessibilityEnabledConfig(),
			wantOut:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("GH_ACCESSIBLE_COLORS", tt.envVarValue)
			testutils.StubConfig(t, tt.cfgStr)
			assert.Equal(t, tt.wantOut, IsAccessibleColorsEnabled())
		})
	}
}

func accessibilityEnabledConfig() string {
	return heredoc.Docf(`
		%s: enabled
	`, AccessibleColorsSetting)
}

func accessibilityDisabledConfig() string {
	return heredoc.Docf(`
		%s: disabled
	`, AccessibleColorsSetting)
}
