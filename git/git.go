package git

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

var excludeFromDiff = []string{
	// Go module artifacts
	"go.sum",
	"go.work.sum",
	"coverage.out",

	// Generic build outputs
	"**/bin/",
	"**/build/",
	"**/dist/",
	"**/out/",
	"**/release/",
	"**/tmp/",
	"**/coverage/",

	// JavaScript / TypeScript
	"**/node_modules/",
	"**/.next/",
	"**/.nuxt/",
	"**/.turbo/",

	// Python
	"**/__pycache__/",
	"**/*.pyc",
	"**/.pytest_cache/",
	"**/.mypy_cache/",
	"**/.ruff_cache/",

	// Java / Kotlin / JVM
	"**/.gradle/",
	"**/*.class",
	"**/target/",

	// .NET
	"**/obj/",

	// Apple / Xcode
	"**/DerivedData/",
	"**/*.xcuserstate",

	// IDE / OS metadata
	"**/.DS_Store",
	"**/.idea/",
	"**/.vscode/",
	"**/*.iml",

	// Compiled binaries and archives
	"**/*.exe",

	// C, C++, Rust, Go
	"**/*.obj",
	"**/*.pdb",
	"**/*.lib",
	"**/*.gcda",
	"**/*.gcno",
	"**/*.gch",
	"**/*.pch",
	"**/*.rlib",
	"**/*.dSYM/",
	"**/*.profraw",
	"**/*.profdata",
	"**/*.prof",
	"**/*.test",
	"**/*.dll",
	"**/*.so",
	"**/*.dylib",
	"**/*.a",
	"**/*.o",
	"**/*.jar",
	"**/*.war",
	"**/*.zip",
	"**/*.tar",
	"**/*.tgz",
	"**/*.gz",

	// Logs and editor backups
	"**/*.log",
	"**/*.tmp",
	"**/*.swp",
	"**/*.swo",
}

// Command wraps git operations used by ReviewBot.
// It encapsulates default diff settings, exclusion patterns, and amend mode.
type Command struct {
	diffUnified  int
	excludedList []string
	isAmend      bool
}

// run executes a git command and returns trimmed stdout or an error with stderr.
func run(args ...string) (string, error) {
	if len(args) == 0 {
		return "", errors.New("no git command provided")
	}

	cmd := exec.Command("git", args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}

	return strings.TrimSpace(stdout.String()), nil
}

// excludedFiles returns a list of file paths prefixed with ":(exclude,top)",
// representing files to be excluded based on the command's excluded list.
func (cmd *Command) excludedFiles() []string {
	excludedFiles := make([]string, 0, len(cmd.excludedList))
	for _, file := range cmd.excludedList {
		excludedFiles = append(excludedFiles, ":(exclude,top)"+file)
	}
	return excludedFiles
}

// diffNameArgs builds the git diff command (name-only) with exclusions applied.
func (cmd *Command) diffNameArgs() []string {
	args := []string{
		"diff",
		"--name-only",
	}

	if cmd.isAmend {
		args = append(args, "HEAD^", "HEAD")
	} else {
		args = append(args, "--staged")
	}

	excludedFiles := cmd.excludedFiles()
	args = append(args, excludedFiles...)

	return args
}

// diffFilesArgs builds git diff arguments with ReviewBot defaults.
func (cmd *Command) diffFilesArgs() []string {
	args := []string{
		"diff",
		"--ignore-all-space",
		"--diff-algorithm=minimal",
		"--unified=" + strconv.Itoa(cmd.diffUnified),
	}

	if cmd.isAmend {
		args = append(args, "HEAD^", "HEAD")
	} else {
		args = append(args, "--staged")
	}

	excludedFiles := cmd.excludedFiles()
	args = append(args, excludedFiles...)

	return args
}

// commitArgs builds git commit arguments (always signoff, optional amend).
func (cmd *Command) commitArgs(msg string) []string {
	args := []string{
		"commit",
		"--signoff",
		fmt.Sprintf("--message=%s", msg),
	}

	if cmd.isAmend {
		args = append(args, "--amend")
	}

	return args
}

// Commit runs git commit and returns stdout. It errors if nothing is staged.
func (cmd *Command) Commit(msg string) (string, error) {
	output, err := run(cmd.commitArgs(msg)...)
	if err != nil {
		return "", err
	}
	if output == "" {
		return "", errors.New("please add your staged changes using git add <files...>")
	}

	return output, nil
}

// DiffFiles returns the staged diff (or HEAD^ vs HEAD when amending).
func (cmd *Command) DiffFiles() (string, error) {
	names, err := run(cmd.diffNameArgs()...)
	if err != nil {
		return "", err
	}
	if names == "" {
		return "", errors.New("please add your staged changes using git add <files...>")
	}

	return run(cmd.diffFilesArgs()...)
}

// Add stages the provided paths using git add. Paths must be non-empty.
func (cmd *Command) Add(paths ...string) (string, error) {
	if len(paths) == 0 {
		return "", errors.New("no paths provided to git add")
	}

	args := append([]string{"add"}, paths...)
	return run(args...)
}

type Config struct {
	DiffUnified  int
	ExcludedList []string
	IsAmend      bool
}

// New creates a new Command instance with the provided options.
// It applies the given options to configure the Command and returns it.
func (cfg *Config) New() *Command {
	command := &Command{
		diffUnified:  cfg.DiffUnified,
		excludedList: append(excludeFromDiff, cfg.ExcludedList...),
		isAmend:      cfg.IsAmend,
	}

	return command
}
