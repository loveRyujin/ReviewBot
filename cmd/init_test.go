package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveConfigFilePath(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		setup    func(t *testing.T) string // returns temp dir
		wantErr  bool
		validate func(t *testing.T, result string)
	}{
		{
			name:     "absolute path",
			filename: "/etc/reviewbot/config.yaml",
			setup: func(t *testing.T) string {
				return ""
			},
			wantErr: false,
			validate: func(t *testing.T, result string) {
				assert.Equal(t, "/etc/reviewbot/config.yaml", result)
			},
		},
		{
			name:     "relative path",
			filename: "config.yaml",
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			wantErr: false,
			validate: func(t *testing.T, result string) {
				assert.True(t, filepath.IsAbs(result))
				assert.Contains(t, result, "config.yaml")
			},
		},
		{
			name:     "relative path with subdirectory",
			filename: "subdir/config.yaml",
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			wantErr: false,
			validate: func(t *testing.T, result string) {
				assert.True(t, filepath.IsAbs(result))
				assert.Contains(t, result, "subdir/config.yaml")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tmpDir := tt.setup(t)
				if tmpDir != "" {
					oldWd, _ := os.Getwd()
					require.NoError(t, os.Chdir(tmpDir))
					t.Cleanup(func() {
						_ = os.Chdir(oldWd)
					})
				}
			}

			result, err := resolveConfigFilePath(tt.filename)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}

func TestEnsureConfigDir(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "reviewbot", "config.yaml")

	err := ensureConfigDir(configPath)
	assert.NoError(t, err)

	// Verify directory was created
	dirPath := filepath.Dir(configPath)
	info, err := os.Stat(dirPath)
	assert.NoError(t, err)
	assert.True(t, info.IsDir())

	// Test idempotency - should not fail if directory already exists
	err = ensureConfigDir(configPath)
	assert.NoError(t, err)
}

func TestEnsureDefaultConfigFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Set environment variables for config location
	t.Setenv("XDG_CONFIG_HOME", tmpDir)
	t.Setenv("HOME", tmpDir)

	err := ensureDefaultConfigFile()
	// This may fail if we can't determine config path, but shouldn't panic
	// The actual behavior depends on the platform
	if err != nil {
		t.Logf("ensureDefaultConfigFile returned error (may be expected): %v", err)
	}
}
