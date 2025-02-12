package accessibility

import (
	"os"

	"github.com/cli/go-gh/v2/pkg/config"
)

// ACCESSIBILITY_ENV is the name of the environment variable that can be used to enable
// accessibility features. If the value is empty, "0", or "false", the accessibility
// features are disabled. Any other value enables the accessibility features. Note that
// this environment variable supercedes the configuration file's accessible setting.
const ACCESSIBILITY_ENV = "ACCESSIBLE"

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
