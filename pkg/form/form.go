package form

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	noStyle      = lipgloss.NewStyle()
)

type Model struct {
	cursor int
	steps  int

	// ai command parameters
	provider textinput.Model
	apiKey   textinput.Model
	model    textinput.Model
	baseURL  textinput.Model
}

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "esc" || k == "ctrl+c" {
			return m, tea.Quit
		}
	}

	switch m.steps {
	case 0:
		return updateAiConfig(msg, m)
	}

	return m, tea.Quit
}

func (m *Model) View() string {
	var s string

	switch m.steps {
	case 0:
		s = aiConfigView(m)
	}

	return fmt.Sprint(s)
}

func (m *Model) Provider() string {
	return m.provider.Value()
}

func (m *Model) ApiKey() string {
	return m.apiKey.Value()
}

func (m *Model) ModelName() string {
	return m.model.Value()
}

func (m *Model) BaseURL() string {
	return m.baseURL.Value()
}

func initModel() Model {
	provider := textinput.New()
	provider.Placeholder = "Provider"
	provider.PromptStyle = focusedStyle
	provider.TextStyle = focusedStyle
	provider.CharLimit = 200
	provider.Width = 20
	provider.Focus()

	apiKey := textinput.New()
	apiKey.Placeholder = "API Key"
	apiKey.EchoMode = textinput.EchoPassword
	apiKey.EchoCharacter = '*'
	apiKey.CharLimit = 200
	apiKey.Width = 20

	model := textinput.New()
	model.Placeholder = "Model"
	model.CharLimit = 200
	model.Width = 20

	baseURL := textinput.New()
	baseURL.Placeholder = "Base URL"
	baseURL.CharLimit = 200
	baseURL.Width = 20

	m := Model{
		provider: provider,
		apiKey:   apiKey,
		model:    model,
		baseURL:  baseURL,
	}

	return m
}

func Run() (*Model, error) {
	m := initModel()
	_, err := tea.NewProgram(&m).Run()

	return &m, err
}
