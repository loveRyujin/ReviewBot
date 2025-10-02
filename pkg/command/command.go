package command

import (
	"os/exec"
	"strings"
)

// IsCommandAvailable reports whether the given command exists on the current system PATH.
// It trims surrounding quotes and whitespace to remain robust across Windows, macOS, and Linux.
func IsCommandAvailable(command string) bool {
	name := strings.TrimSpace(command)
	name = strings.Trim(name, "\"'")
	if name == "" {
		return false
	}

	_, err := exec.LookPath(name)
	return err == nil
}
