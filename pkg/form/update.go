package form

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func updateAiConfig(msg tea.Msg, m *Model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "down", "enter":
			inputs := []textinput.Model{
				m.provider,
				m.apiKey,
				m.model,
				m.baseURL,
			}

			s := msg.String()

			if s == "enter" && m.cursor == len(inputs)-1 {
				m.cursor = 0
				m.steps = 1
				return m, nil
			}

			if s == "up" {
				m.cursor--
			} else {
				m.cursor++
			}

			if m.cursor >= len(inputs) {
				m.cursor = 0
			} else if m.cursor < 0 {
				m.cursor = len(inputs) - 1
			}

			for i := range len(inputs) {
				if i == m.cursor {
					inputs[i].Focus()
					inputs[i].PromptStyle = focusedStyle
					inputs[i].TextStyle = focusedStyle
					continue
				}
				inputs[i].Blur()
				inputs[i].PromptStyle = noStyle
				inputs[i].TextStyle = noStyle
			}

			m.provider = inputs[0]
			m.apiKey = inputs[1]
			m.model = inputs[2]
			m.baseURL = inputs[3]

			return m, nil
		}

	}

	m, cmd = updateInputs(msg, m)

	return m, cmd
}

func updateInputs(msg tea.Msg, m *Model) (*Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.provider, cmd = m.provider.Update(msg)
	cmds = append(cmds, cmd)

	m.apiKey, cmd = m.apiKey.Update(msg)
	cmds = append(cmds, cmd)

	m.model, cmd = m.model.Update(msg)
	cmds = append(cmds, cmd)

	m.baseURL, cmd = m.baseURL.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
