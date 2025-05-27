// Package browser facilitates opening of URLs in a web browser.
package browser

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	cliBrowser "github.com/cli/browser"
	"github.com/cli/go-gh/v2/pkg/config"
	"github.com/cli/safeexec"
	"github.com/google/shlex"
)

// Browser represents a web browser that can be used to open up URLs.
type Browser struct {
	launcher string
	stderr   io.Writer
	stdout   io.Writer
}

// New initializes a Browser. If a launcher is not specified
// one is determined based on environment variables or from the
// configuration file.
// The order of precedence for determining a launcher is:
// - Specified launcher;
// - GH_BROWSER environment variable;
// - browser option from configuration file;
// - BROWSER environment variable.
func New(launcher string, stdout, stderr io.Writer) *Browser {
	if launcher == "" {
		launcher = resolveLauncher()
	}
	b := &Browser{
		launcher: launcher,
		stderr:   stderr,
		stdout:   stdout,
	}
	return b
}

// Browse opens the launcher and navigates to the specified URL.
func (b *Browser) Browse(url string) error {
	return b.browse(url, nil)
}

func (b *Browser) browse(url string, env []string) error {
	if err := isPossibleProtocol(url); err != nil {
		return err
	}

	if b.launcher == "" {
		return cliBrowser.OpenURL(url)
	}
	launcherArgs, err := shlex.Split(b.launcher)
	if err != nil {
		return err
	}
	launcherExe, err := safeexec.LookPath(launcherArgs[0])
	if err != nil {
		return err
	}
	args := append(launcherArgs[1:], url)
	cmd := exec.Command(launcherExe, args...)
	cmd.Stdout = b.stdout
	cmd.Stderr = b.stderr
	if env != nil {
		cmd.Env = env
	}
	return cmd.Run()
}

func resolveLauncher() string {
	if ghBrowser := os.Getenv("GH_BROWSER"); ghBrowser != "" {
		return ghBrowser
	}
	cfg, err := config.Read(nil)
	if err == nil {
		if cfgBrowser, _ := cfg.Get([]string{"browser"}); cfgBrowser != "" {
			return cfgBrowser
		}
	}
	return os.Getenv("BROWSER")
}

func isSupportedProtocol(u string) bool {
	return strings.HasPrefix(u, "http:") ||
		strings.HasPrefix(u, "https:") ||
		strings.HasPrefix(u, "vscode:") ||
		strings.HasPrefix(u, "vscode-insiders:")
}

func isPossibleProtocol(u string) error {
	if isSupportedProtocol(u) {
		return nil
	}

	// Disallow URLs using alternative `file:` protocol
	if strings.HasPrefix(u, "file:") {
		return fmt.Errorf("opening files or directories is unsupported: %s", u)
	}

	// Disallow URLs that match existing files or directories on the filesystem
	// as these could be executables or executed by the launcher browser due to
	// the file extension and/or associated application.
	if fileInfo, _ := os.Lstat(u); fileInfo != nil {
		return fmt.Errorf("opening files or directories is unsupported: %s", u)
	}

	// Disallow URLs that match executables found in the user path.
	exec, _ := safeexec.LookPath(u)
	if exec != "" {
		return fmt.Errorf("opening executables is unsupported: %s", u)
	}

	return nil
}
