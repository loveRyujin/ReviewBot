package pkg

import "testing"

func Test_IsCommandAvailable(t *testing.T) {
	testCases := []struct {
		name     string
		command  string
		expected bool
	}{
		{
			name:     "valid command-1",
			command:  "pwd",
			expected: true,
		},
		{
			name:     "valid command-2",
			command:  "git",
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
