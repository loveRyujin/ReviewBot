package form

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	noStyle      = lipgloss.NewStyle()
)

type Model struct {
	cursor int
	steps  int

	choices []string

	// git command parameters
	diffFile     textinput.Model
	maxInputSize textinput.Model
	diffUnified  textinput.Model
	excludedList textinput.Model
	amend        textinput.Model
	lang         textinput.Model

	// ai command parameters
	provider textinput.Model
	apiKey   textinput.Model
	model    textinput.Model
	baseURL  textinput.Model

	// proxy command parameters
	proxyURL   textinput.Model
	socksURL   textinput.Model
	timeout    textinput.Model
	headers    textinput.Model
	skipVerify textinput.Model
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
	case 1:
		return chooseIfUpdateGitConfig(msg, m)
	case 2:
		return updateGitConfig(msg, m)
	case 3:
		return chooseIfUpdateProxyConfig(msg, m)
	case 4:
		return updateProxyConfig(msg, m)
	}

	return m, tea.Quit
}

func (m *Model) View() string {
	var s string

	switch m.steps {
	case 0:
		s = aiConfigView(m)
	case 1:
		s = chooseIfUpdateGitConfigView(m)
	case 2:
		s = gitConfigView(m)
	case 3:
		s = chooseIfUpdateProxyConfigView(m)
	case 4:
		s = proxyConfigView(m)
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

func (m *Model) DiffFile() string {
	return m.diffFile.Value()
}

func (m *Model) MaxInputSize() string {
	return m.maxInputSize.Value()
}

func (m *Model) DiffUnified() string {
	return m.diffUnified.Value()
}

func (m *Model) ExcludedList() []string {
	return strings.Split(m.excludedList.Value(), ",")
}

func (m *Model) Amend() string {
	return m.amend.Value()
}

func (m *Model) Lang() string {
	return m.lang.Value()
}

func (m *Model) ProxyURL() string {
	return m.proxyURL.Value()
}

func (m *Model) SocksURL() string {
	return m.socksURL.Value()
}

func (m *Model) Timeout() string {
	return m.timeout.Value()
}

func (m *Model) Headers() []string {
	return strings.Split(m.headers.Value(), ",")
}

func (m *Model) SkipVerify() string {
	return m.skipVerify.Value()
}

func checkBox(label string, checked bool) string {
	if checked {
		return color.RGB(127, 255, 212).Sprintf("[x] %s", label)
	}

	return fmt.Sprintf("[ ] %s", label)
}

func initModel() Model {
	// Initialize text input fields for AI configuration
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

	// Initialize text input fields for Git configuration
	diffFile := textinput.New()
	diffFile.Placeholder = "Diff File"
	diffFile.PromptStyle = focusedStyle
	diffFile.TextStyle = focusedStyle
	diffFile.CharLimit = 200
	diffFile.Width = 20
	diffFile.Focus()

	maxInputSize := textinput.New()
	maxInputSize.Placeholder = "Max Input Size"
	maxInputSize.CharLimit = 200
	maxInputSize.Width = 20

	diffUnified := textinput.New()
	diffUnified.Placeholder = "Diff Unified"
	diffUnified.CharLimit = 200
	diffUnified.Width = 20

	excludedList := textinput.New()
	excludedList.Placeholder = "Excluded List"
	excludedList.CharLimit = 200
	excludedList.Width = 20

	amend := textinput.New()
	amend.Placeholder = "Amend"
	amend.CharLimit = 200
	amend.Width = 20

	lang := textinput.New()
	lang.Placeholder = "Language"
	lang.CharLimit = 200
	lang.Width = 20

	// Initialize text input fields for Proxy configuration
	proxyURL := textinput.New()
	proxyURL.Placeholder = "Proxy URL"
	proxyURL.PromptStyle = focusedStyle
	proxyURL.TextStyle = focusedStyle
	proxyURL.CharLimit = 200
	proxyURL.Width = 20
	proxyURL.Focus()

	socksURL := textinput.New()
	socksURL.Placeholder = "SOCKS URL"
	socksURL.CharLimit = 200
	socksURL.Width = 20

	timeout := textinput.New()
	timeout.Placeholder = "Timeout"
	timeout.CharLimit = 200
	timeout.Width = 20

	headers := textinput.New()
	headers.Placeholder = "Headers"
	headers.CharLimit = 200
	headers.Width = 20

	skipVerify := textinput.New()
	skipVerify.Placeholder = "Skip Verify"
	skipVerify.CharLimit = 200
	skipVerify.Width = 20

	m := Model{
		provider:     provider,
		apiKey:       apiKey,
		model:        model,
		baseURL:      baseURL,
		diffFile:     diffFile,
		maxInputSize: maxInputSize,
		diffUnified:  diffUnified,
		excludedList: excludedList,
		amend:        amend,
		lang:         lang,
		proxyURL:     proxyURL,
		socksURL:     socksURL,
		timeout:      timeout,
		headers:      headers,
		skipVerify:   skipVerify,
		choices:      []string{"true", "false"},
	}

	return m
}

func Run() (*Model, error) {
	m := initModel()
	_, err := tea.NewProgram(&m).Run()

	return &m, err
}
