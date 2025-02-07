package accessibility

import (
	"testing"

	"github.com/cli/go-gh/v2/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestIsEnabled(t *testing.T) {
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
			name:        "When the accessibility configuration is unset but the ACCESSIBLE env var is set to something truthy (not '0' or 'false'), it should return true",
			envVarValue: "1",
			cfgStr:      "",
			wantOut:     true,
		},
		{
			name:        "When the accessibility configuration is unset and the ACCESSIBLE env var returns '0', it should return false",
			envVarValue: "0",
			cfgStr:      "",
			wantOut:     false,
		},
		{
			name:        "When the accessibility configuration is unset and the ACCESSIBLE env var returns 'false', it should return false",
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
			t.Setenv("GH_ACCESSIBILE", tt.envVarValue)
			testutils.StubConfig(t, tt.cfgStr)
			assert.Equal(t, tt.wantOut, IsEnabled())
		})
	}
}

func accessibilityEnabledConfig() string {
	return `
accessible: enabled
`
}

func accessibilityDisabledConfig() string {
	return `
accessible: disabled
`
}
