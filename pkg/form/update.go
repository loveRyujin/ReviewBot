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

	m, cmd = updateAiInputs(msg, m)

	return m, cmd
}

func updateGitConfig(msg tea.Msg, m *Model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "down", "enter":
			inputs := []textinput.Model{
				m.diffFile,
				m.maxInputSize,
				m.diffUnified,
				m.excludedList,
				m.amend,
				m.lang,
			}

			s := msg.String()

			if s == "enter" && m.cursor == len(inputs)-1 {
				m.cursor = 0
				m.steps = 3
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

			m.diffFile = inputs[0]
			m.maxInputSize = inputs[1]
			m.diffUnified = inputs[2]
			m.excludedList = inputs[3]
			m.amend = inputs[4]
			m.lang = inputs[5]

			return m, nil
		}

	}

	m, cmd = updateGitInputs(msg, m)

	return m, cmd
}

func chooseIfUpdateGitConfig(msg tea.Msg, m *Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down", "k", "up":
			s := msg.String()

			if s == "k" || s == "up" {
				m.cursor--
			} else if s == "j" || s == "down" {
				m.cursor++
			}

			if m.cursor >= 2 {
				m.cursor = 0
			} else if m.cursor < 0 {
				m.cursor = 1
			}
		case "enter":
			mode := m.choices[m.cursor]
			if mode == "false" {
				return m, tea.Quit
			}

			m.steps = 2
			m.cursor = 0
		}
	}

	return m, nil
}

func updateProxyConfig(msg tea.Msg, m *Model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "down", "enter":
			inputs := []textinput.Model{
				m.proxyURL,
				m.timeout,
				m.socksURL,
				m.headers,
				m.skipVerify,
			}

			s := msg.String()

			if s == "enter" && m.cursor == len(inputs)-1 {
				return m, tea.Quit
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

			m.proxyURL = inputs[0]
			m.timeout = inputs[1]
			m.socksURL = inputs[2]
			m.headers = inputs[3]
			m.skipVerify = inputs[4]

			return m, nil
		}
	}

	m, cmd = updateProxyInputs(msg, m)

	return m, cmd
}

func chooseIfUpdateProxyConfig(msg tea.Msg, m *Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down", "k", "up":
			s := msg.String()

			if s == "k" || s == "up" {
				m.cursor--
			} else if s == "j" || s == "down" {
				m.cursor++
			}

			if m.cursor >= 2 {
				m.cursor = 0
			} else if m.cursor < 0 {
				m.cursor = 1
			}
		case "enter":
			mode := m.choices[m.cursor]
			if mode == "false" {
				return m, tea.Quit
			}

			m.steps = 4
			m.cursor = 0
		}
	}

	return m, nil
}

func updateAiInputs(msg tea.Msg, m *Model) (*Model, tea.Cmd) {
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

func updateGitInputs(msg tea.Msg, m *Model) (*Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.diffFile, cmd = m.diffFile.Update(msg)
	cmds = append(cmds, cmd)

	m.maxInputSize, cmd = m.maxInputSize.Update(msg)
	cmds = append(cmds, cmd)

	m.diffUnified, cmd = m.diffUnified.Update(msg)
	cmds = append(cmds, cmd)

	m.excludedList, cmd = m.excludedList.Update(msg)
	cmds = append(cmds, cmd)

	m.amend, cmd = m.amend.Update(msg)
	cmds = append(cmds, cmd)

	m.lang, cmd = m.lang.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func updateProxyInputs(msg tea.Msg, m *Model) (*Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.proxyURL, cmd = m.proxyURL.Update(msg)
	cmds = append(cmds, cmd)

	m.timeout, cmd = m.timeout.Update(msg)
	cmds = append(cmds, cmd)

	m.socksURL, cmd = m.socksURL.Update(msg)
	cmds = append(cmds, cmd)

	m.headers, cmd = m.headers.Update(msg)
	cmds = append(cmds, cmd)

	m.skipVerify, cmd = m.skipVerify.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
