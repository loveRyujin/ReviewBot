package integration_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/loveRyujin/ReviewBot/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGitIntegration tests git package with a real git repository
func TestGitIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create test repo
	tmpDir := t.TempDir()
	require.NoError(t, os.Chdir(tmpDir))

	// Initialize git
	setupGitRepo(t)

	// Create git config
	cfg := &git.Config{
		DiffUnified: 3,
	}
	cmd := cfg.New()

	// Test adding files
	testFile := "test.txt"
	require.NoError(t, os.WriteFile(testFile, []byte("test content"), 0o644))

	_, err := cmd.Add(testFile)
	require.NoError(t, err)

	// Test getting diff
	diff, err := cmd.DiffFiles()
	require.NoError(t, err)
	assert.Contains(t, diff, "test.txt")
	assert.Contains(t, diff, "test content")

	// Test committing
	commitOutput, err := cmd.Commit("test: add test file")
	require.NoError(t, err)
	assert.NotEmpty(t, commitOutput)
}

// TestGitIntegrationWithMultipleFiles tests handling multiple files
func TestGitIntegrationWithMultipleFiles(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	require.NoError(t, os.Chdir(tmpDir))
	setupGitRepo(t)

	cfg := &git.Config{DiffUnified: 5}
	cmd := cfg.New()

	// Add multiple files
	files := map[string]string{
		"file1.go":  "package main\n\nfunc main() {}",
		"file2.go":  "package main\n\nfunc helper() {}",
		"README.md": "# Test Project",
	}

	for name, content := range files {
		require.NoError(t, os.WriteFile(name, []byte(content), 0o644))
	}

	_, err := cmd.Add(".")
	require.NoError(t, err)

	// Get diff
	diff, err := cmd.DiffFiles()
	require.NoError(t, err)

	// Verify all files in diff
	for filename := range files {
		assert.Contains(t, diff, filename)
	}
}

// setupGitRepo initializes a git repository for testing
func setupGitRepo(t *testing.T) {
	t.Helper()

	commands := [][]string{
		{"init"},
		{"config", "user.name", "Test User"},
		{"config", "user.email", "test@example.com"},
		{"config", "commit.gpgsign", "false"},
	}

	for _, args := range commands {
		cmd := exec.Command("git", args...)
		if err := cmd.Run(); err != nil {
			t.Fatalf("git %v failed: %v", args, err)
		}
	}
}
