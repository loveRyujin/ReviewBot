package command

import (
	"runtime"
	"testing"
)

func Test_IsCommandAvailable(t *testing.T) {
	testCases := []struct {
		name     string
		command  string
		expected bool
	}{
		{
			name:     "valid command-2",
			command:  "git",
			expected: true,
		},
		{
			name:     "valid shell command",
			command:  defaultShellCommand(),
			expected: true,
		},

		{
			name:     "invalid command",
			command:  "nonexistentcommand",
			expected: false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			result := IsCommandAvailable(test.command)
			if result != test.expected {
				t.Errorf("expected %v, got %v", test.expected, result)
			}
		})
	}
}

func defaultShellCommand() string {
	if runtime.GOOS == "windows" {
		return "cmd"
	}
	return "sh"
}
