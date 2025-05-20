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
	if !isPossibleProtocol(url) {
		return fmt.Errorf("cannot browse due to unsupported protocol: %s", url)
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

func (b *Browser) IsURL(u string) bool {
	return isPossibleProtocol(u)
}

func isSupportedProtocol(u string) bool {
	return strings.HasPrefix(u, "http:") ||
		strings.HasPrefix(u, "https:") ||
		strings.HasPrefix(u, "vscode:") ||
		strings.HasPrefix(u, "vscode-insiders:")
}

func isPossibleProtocol(u string) bool {
	if isSupportedProtocol(u) {
		return true
	}

	// Disallow URLs that match existing files or directorys on the filesystem
	if fileInfo, _ := os.Lstat(u); fileInfo != nil {
		return false
	}

	// Disallow URLs using alternative `file:` protocol not located previously
	if strings.HasPrefix(u, "file:") {
		return false
	}

	return true
}
