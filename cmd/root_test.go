package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureDefaultConfigFileCreatesFile(t *testing.T) {
	temp := t.TempDir()

	setConfigEnv(t, temp)

	// create config file at resolved default location
	if err := ensureDefaultConfigFile(); err != nil {
		t.Fatalf("ensureDefaultConfigFile error: %v", err)
	}

	expectedDir := filepath.Join(temp, "reviewbot")
	expectedFile := filepath.Join(expectedDir, defaultConfigFile)

	// file should be created on first call
	if _, err := os.Stat(expectedFile); err != nil {
		t.Fatalf("config file not created: %v", err)
	}

	// modify content and ensure function does not overwrite existing file
	if err := os.WriteFile(expectedFile, []byte("custom"), 0o600); err != nil {
		t.Fatalf("write custom content: %v", err)
	}

	// second call must be a no-op
	if err := ensureDefaultConfigFile(); err != nil {
		t.Fatalf("ensureDefaultConfigFile second call error: %v", err)
	}

	// confirm custom content is preserved
	updated, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("read config after second call: %v", err)
	}
	if string(updated) != "custom" {
		t.Fatalf("config file should not be overwritten, got %q", string(updated))
	}
}

func setConfigEnv(t *testing.T, dir string) {
	t.Helper()

	env := map[string]string{
		"XDG_CONFIG_HOME": dir,
		"HOME":            dir,
		"APPDATA":         dir,
	}

	for key, val := range env {
		old := os.Getenv(key)
		if err := os.Setenv(key, val); err != nil {
			t.Fatalf("setenv %s: %v", key, err)
		}
		t.Cleanup(func() {
			if old == "" {
				_ = os.Unsetenv(key)
			} else {
				_ = os.Setenv(key, old)
			}
		})
	}
}
