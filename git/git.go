package git

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
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

type Command struct {
	diffUnified  int
	excludedList []string
	isAmend      bool
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

// diffName constructs and returns a git command to list file names that have been changed.
// It supports both staged changes and amendments, and allows excluding specific files.
func (cmd *Command) diffName() *exec.Cmd {
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

	return exec.Command(
		"git",
		args...,
	)
}

// diffFiles constructs and returns an *exec.Cmd to execute a Git diff command.
// It supports generating diffs for staged changes or between the last commit
// and its parent (in case of an amend). The method also applies options such as
// ignoring whitespace changes, using a minimal diff algorithm, and excluding
// specified files.
func (cmd *Command) diffFiles() *exec.Cmd {
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

	return exec.Command(
		"git",
		args...,
	)
}

// commit constructs a git commit command with the provided commit message.
// It includes the --signoff flag and optionally the --amend flag if cmd.isAmend is true.
func (cmd *Command) commit(msg string) *exec.Cmd {
	args := []string{
		"commit",
		"--signoff",
		fmt.Sprintf("--message=%s", msg),
	}

	if cmd.isAmend {
		args = append(args, "--amend")
	}

	return exec.Command(
		"git",
		args...,
	)
}

// Commit executes a git commit with the provided message and returns the output.
// If there are no staged changes, it returns an error prompting the user to stage files.
func (cmd *Command) Commit(msg string) (string, error) {
	output, err := cmd.commit(msg).Output()
	if err != nil {
		return "", err
	}
	if len(output) == 0 {
		return "", errors.New("please add your staged changes using git add <files...>")
	}

	return string(output), nil
}

func (cmd *Command) DiffFiles() (string, error) {
	output, err := cmd.diffName().Output()
	if err != nil {
		return "", err
	}
	if len(output) == 0 {
		return "", errors.New("please add your staged changes using git add <files...>")
	}

	output, err = cmd.diffFiles().Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
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
