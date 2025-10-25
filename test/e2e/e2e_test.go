package e2e_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestE2EGitBasicWorkflow tests a basic git workflow with reviewbot
func TestE2EGitBasicWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	// Setup test repository
	tmpDir := t.TempDir()
	require.NoError(t, os.Chdir(tmpDir))

	// Initialize git repo
	runGitCommand(t, "init")
	runGitCommand(t, "config", "user.name", "Test User")
	runGitCommand(t, "config", "user.email", "test@example.com")
	runGitCommand(t, "config", "commit.gpgsign", "false")

	// Create initial commit
	testFile := filepath.Join(tmpDir, "test.txt")
	require.NoError(t, os.WriteFile(testFile, []byte("initial content"), 0o644))
	runGitCommand(t, "add", "test.txt")
	runGitCommand(t, "commit", "-m", "Initial commit")

	// Make changes
	require.NoError(t, os.WriteFile(testFile, []byte("modified content"), 0o644))
	runGitCommand(t, "add", "test.txt")

	// Get staged diff
	diff := runGitCommand(t, "diff", "--cached")
	assert.Contains(t, diff, "test.txt")
	assert.Contains(t, diff, "modified content")
}

// TestE2EGitDiffOperations tests git diff operations
func TestE2EGitDiffOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	tmpDir := t.TempDir()
	require.NoError(t, os.Chdir(tmpDir))

	// Setup
	runGitCommand(t, "init")
	runGitCommand(t, "config", "user.name", "Test User")
	runGitCommand(t, "config", "user.email", "test@example.com")
	runGitCommand(t, "config", "commit.gpgsign", "false")

	// Create multiple files
	for i, content := range []string{"file1 content", "file2 content", "file3 content"} {
		filename := filepath.Join(tmpDir, "file"+string(rune('1'+i))+".txt")
		require.NoError(t, os.WriteFile(filename, []byte(content), 0o644))
	}

	runGitCommand(t, "add", ".")
	runGitCommand(t, "commit", "-m", "Add multiple files")

	// Modify files
	file1 := filepath.Join(tmpDir, "file1.txt")
	require.NoError(t, os.WriteFile(file1, []byte("file1 modified"), 0o644))
	runGitCommand(t, "add", "file1.txt")

	// Check diff
	diff := runGitCommand(t, "diff", "--cached")
	assert.Contains(t, diff, "file1.txt")
	assert.Contains(t, diff, "file1 modified")
}

// TestE2EMultipleCommits tests handling multiple commits
func TestE2EMultipleCommits(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	tmpDir := t.TempDir()
	require.NoError(t, os.Chdir(tmpDir))

	runGitCommand(t, "init")
	runGitCommand(t, "config", "user.name", "Test User")
	runGitCommand(t, "config", "user.email", "test@example.com")
	runGitCommand(t, "config", "commit.gpgsign", "false")

	// Create several commits
	commits := []struct {
		file    string
		content string
		message string
	}{
		{"feature1.txt", "feature 1 implementation", "feat: add feature 1"},
		{"feature2.txt", "feature 2 implementation", "feat: add feature 2"},
		{"bugfix.txt", "fix critical bug", "fix: resolve critical bug"},
	}

	for _, commit := range commits {
		filename := filepath.Join(tmpDir, commit.file)
		require.NoError(t, os.WriteFile(filename, []byte(commit.content), 0o644))
		runGitCommand(t, "add", commit.file)
		runGitCommand(t, "commit", "-m", commit.message)
	}

	// Verify commit history
	log := runGitCommand(t, "log", "--oneline")
	assert.Contains(t, log, "feat: add feature 1")
	assert.Contains(t, log, "feat: add feature 2")
	assert.Contains(t, log, "fix: resolve critical bug")
}

// runGitCommand executes a git command and returns output
func runGitCommand(t *testing.T, args ...string) string {
	t.Helper()

	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s failed: %v\nOutput: %s", strings.Join(args, " "), err, string(output))
	}

	return string(output)
}
