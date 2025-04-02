package cmd

import (
	"errors"

	"github.com/loveRyujin/ReviewBot/util"
	"github.com/spf13/viper"
)

type ServerOptions struct {
	GitOptions *GitOptions `mapstructure:"git"`
}

func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		GitOptions: NewGitOptions(),
	}
}

func (s *ServerOptions) Validate() error {
	if err := s.GitOptions.Validate(); err != nil {
		return err
	}

	return nil
}

// ApplyCfg applies user-provided configuration from the command line.
// It ensures the "git" command is available in the system PATH and updates
// configuration settings such as the number of unified lines for diffs
// and the list of excluded items.
func (s *ServerOptions) ApplyCfg() error {
	if !util.IsCommandAvailable("git") {
		return errors.New("git command not found in your system PATH, Please install git")
	}

	if diffUnifiedLines != 3 {
		viper.Set("git.diff_unified", diffUnifiedLines)
	}

	if len(excludedList) > 0 {
		viper.Set("git.exclude_list", excludedList)
	}

	return nil
}

type GitOptions struct {
	DiffUnified  int      `mapstructure:"diff_unified"`
	ExcludedList []string `mapstructure:"exclude_list"`
	Amend        bool     `mapstructure:"amend"`
}

func NewGitOptions() *GitOptions {
	return &GitOptions{
		DiffUnified:  3,
		ExcludedList: []string{},
		Amend:        false,
	}
}

func (g *GitOptions) Validate() error {
	if g.DiffUnified < 0 {
		return errors.New("diff_unified must be a non-negative integer")
	}

	return nil
}
