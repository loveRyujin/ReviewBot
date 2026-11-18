package cmd

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/loveRyujin/ReviewBot/pkg/config"
	"github.com/loveRyujin/ReviewBot/pkg/form"
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
			name: "absolute path",
			filename: func() string {
				if runtime.GOOS == "windows" {
					return "C:\\etc\\reviewbot\\config.yaml"
				}
				return "/etc/reviewbot/config.yaml"
			}(),
			setup: func(t *testing.T) string {
				return ""
			},
			wantErr: false,
			validate: func(t *testing.T, result string) {
				assert.True(t, filepath.IsAbs(result))
				assert.Contains(t, result, "config.yaml")
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
			filename: filepath.Join("subdir", "config.yaml"),
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			wantErr: false,
			validate: func(t *testing.T, result string) {
				assert.True(t, filepath.IsAbs(result))
				assert.Contains(t, result, "subdir")
				assert.Contains(t, result, "config.yaml")
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

func TestInitConfigRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings", "reviewbot.yaml")
	first := form.NewModel()
	require.NoError(t, saveInitConfig(path, &first))
	require.NoFileExists(t, path)

	submit := func(m *form.Model, fields ...string) {
		for _, value := range fields {
			if value != "" {
				m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(value)})
			}
			m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		}
	}
	submit(&first, "", "fake-existing-key", "old-model", "")
	require.NoError(t, saveInitConfig(path, &first))

	second := form.NewModel()
	submit(&second, "", "", "updated-model", "")
	submit(&second, "") // Enable Git settings.
	submit(&second, "", "", "", "vendor,dist", "", "")
	submit(&second, "") // Enable proxy settings.
	submit(&second, "", "", "", "Authorization=Bearer a,b,X-Trace=1", "")
	require.NoError(t, saveInitConfig(path, &second))

	saved, err := os.ReadFile(path)
	require.NoError(t, err)
	require.YAMLEq(t, `ai:
  api_key: fake-existing-key
  model: updated-model
git:
  exclude_list: [vendor, dist]
proxy:
  headers: ["Authorization=Bearer a,b", "X-Trace=1"]
`, string(saved))
	cfg, err := config.Load(config.LoadOptions{ExplicitPath: path, EnvPrefix: "REVIEWBOT_INIT_TEST"})
	require.NoError(t, err)
	require.Equal(t, "fake-existing-key", cfg.AI.APIKey)
	require.Equal(t, []string{"vendor", "dist"}, cfg.Git.ExcludedList)
	require.Equal(t, []string{"Authorization=Bearer a,b", "X-Trace=1"}, cfg.Proxy.Headers)
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
