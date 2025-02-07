package accessibility

import (
	"os"

	"github.com/cli/go-gh/v2/pkg/config"
)

const ACCESSIBILITY_ENV = "ACCESSIBILITY"

// IsEnabled returns true if accessibility features are enabled via the ACCESSIBILITY
// environment variable or the configuration file.
func IsEnabled() bool {
	envVar := os.Getenv(ACCESSIBILITY_ENV)
	if envVar != "" {
		return isEnvVarEnabled(envVar)
	}

	// We are not handling errors because we don't want to fail if the config is not
	// read. Instead, we assume an empty configuration is equivalent to "disabled".
	cfg, _ := config.Read(nil)
	accessibleConfigValue, _ := cfg.Get([]string{"accessible"})

	return accessibleConfigValue == "enabled"
}

func isEnvVarEnabled(envVar string) bool {
	return envVar != "" && envVar != "0" && envVar != "false"
}
