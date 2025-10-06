package git

import (
	"os"
	"strings"
	"testing"
)

// TestCommandAddRequiresPaths verifies Add rejects empty input.
func TestCommandAddRequiresPaths(t *testing.T) {
	cmd := (&Config{DiffUnified: 3}).New()
	if _, err := cmd.Add(); err == nil {
		t.Fatalf("expected error when calling Add with no paths")
	}
}

// TestCommandAddStagesFile ensures Add stages the file.
func TestCommandAddStagesFile(t *testing.T) {
	setupRepo(t)

	if err := os.WriteFile("foo.txt", []byte("hello"), 0o644); err != nil {
		t.Fatalf("write foo.txt: %v", err)
	}

	cmd := (&Config{DiffUnified: 3}).New()
	if _, err := cmd.Add("foo.txt"); err != nil {
		t.Fatalf("cmd.Add: %v", err)
	}

	status := gitRun(t, "status", "--short")
	if !strings.Contains(status, "A  foo.txt") {
		t.Fatalf("expected foo.txt to be staged, git status: %q", status)
	}
}

// TestCommandDiffFiles verifies staged diffs are returned.
func TestCommandDiffFiles(t *testing.T) {
	setupRepo(t)

	if err := os.WriteFile("foo.txt", []byte("hello"), 0o644); err != nil {
		t.Fatalf("write foo.txt: %v", err)
	}

	cmd := (&Config{DiffUnified: 3}).New()
	if _, err := cmd.Add("foo.txt"); err != nil {
		t.Fatalf("cmd.Add: %v", err)
	}

	diff, err := cmd.DiffFiles()
	if err != nil {
		t.Fatalf("DiffFiles: %v", err)
	}
	if !strings.Contains(diff, "foo.txt") {
		t.Fatalf("diff should include foo.txt, got: %q", diff)
	}
}

// TestCommandDiffFilesNoChanges expects error when nothing staged.
func TestCommandDiffFilesNoChanges(t *testing.T) {
	setupRepo(t)

	cmd := (&Config{DiffUnified: 3}).New()
	if _, err := cmd.DiffFiles(); err == nil {
		t.Fatalf("expected error when no staged changes")
	}
}

// TestCommandCommit ensures commit succeeds and cleans working tree.
func TestCommandCommit(t *testing.T) {
	setupRepo(t)

	if err := os.WriteFile("foo.txt", []byte("hello"), 0o644); err != nil {
		t.Fatalf("write foo.txt: %v", err)
	}

	cmd := (&Config{DiffUnified: 3}).New()
	if _, err := cmd.Add("foo.txt"); err != nil {
		t.Fatalf("cmd.Add: %v", err)
	}

	output, err := cmd.Commit("feat: add foo")
	if err != nil {
		t.Fatalf("Commit: %v", err)
	}
	if output == "" {
		t.Fatalf("expected git commit output, got empty string")
	}

	status := gitRun(t, "status", "--short")
	if strings.TrimSpace(status) != "" {
		t.Fatalf("working tree should be clean after commit, status: %q", status)
	}
}

// setupRepo initializes a temporary git repository for testing.
func setupRepo(t *testing.T) {
	t.Helper()

	dir := t.TempDir()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	t.Cleanup(func() {
		_ = os.Chdir(wd)
	})

	gitRun(t, "init")
	gitRun(t, "config", "user.name", "ReviewBot")
	gitRun(t, "config", "user.email", "bot@example.com")
	gitRun(t, "config", "commit.gpgsign", "false")
}

// gitRun wraps run to fail tests with context on errors.
func gitRun(t *testing.T, args ...string) string {
	t.Helper()

	out, err := run(args...)
	if err != nil {
		t.Fatalf("git %s: %v", strings.Join(args, " "), err)
	}
	return out
}
