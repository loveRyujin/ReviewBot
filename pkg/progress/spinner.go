package progress

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

// Spinner wraps the spinner functionality
type Spinner struct {
	spinner *spinner.Spinner
	message string
}

// NewSpinner creates a new spinner instance
func NewSpinner(message string) *Spinner {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = message + " "
	_ = s.Color("cyan")

	return &Spinner{
		spinner: s,
		message: message,
	}
}

// Start starts the spinner
func (s *Spinner) Start() {
	s.spinner.Start()
}

// Stop stops the spinner
func (s *Spinner) Stop() {
	s.spinner.Stop()
}

// UpdateMessage updates the spinner message
func (s *Spinner) UpdateMessage(message string) {
	s.message = message
	s.spinner.Prefix = message + " "
}

// Success stops spinner and shows success message
func (s *Spinner) Success(message string) {
	s.spinner.Stop()
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("%s %s\n", green("✓"), message)
}

// Error stops spinner and shows error message
func (s *Spinner) Error(message string) {
	s.spinner.Stop()
	red := color.New(color.FgRed).SprintFunc()
	fmt.Printf("%s %s\n", red("✗"), message)
}

// WithSpinner executes operation with spinner
func WithSpinner(message string, operation func() error) error {
	s := NewSpinner(message)
	s.Start()

	if err := operation(); err != nil {
		s.Error("Failed")
		return err
	}

	s.Success("Completed")

	return nil
}

// WithSpinnerAndCustomMessages executes operation with custom success/error messages
func WithSpinnerAndCustomMessages(message, successMsg, errorMsg string, operation func() error) error {
	s := NewSpinner(message)
	s.Start()

	if err := operation(); err != nil {
		s.Error(errorMsg)
		return err
	}

	s.Success(successMsg)

	return nil
}
