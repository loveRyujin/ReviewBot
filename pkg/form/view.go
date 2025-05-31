package form

import "fmt"

// aiConfigView generates the view for AI configuration input.
func aiConfigView(m *Model) string {
	s := "Introduce your AI configuration:\n\n"

	inputs := []string{
		m.provider.View(),
		m.apiKey.View(),
		m.model.View(),
		m.baseURL.View(),
	}

	return s + formatInputsWithNewlines(inputs)
}

// chooseIfUpdateGitConfigView generates the view for choosing whether to update Git configuration.
func chooseIfUpdateGitConfigView(m *Model) string {
	s := "Do you want to update your Git configuration? \n\n"

	return s + formatChoicesWithCheckboxes(m.choices, m.cursor)
}

// gitConfigView generates the view for Git configuration input.
func gitConfigView(m *Model) string {
	s := "Introduce your Git configuration:\n\n"

	inputs := []string{
		m.diffFile.View(),
		m.maxInputSize.View(),
		m.diffUnified.View(),
		m.excludedList.View(),
		m.amend.View(),
		m.lang.View(),
	}

	return s + formatInputsWithNewlines(inputs)
}

// chooseIfUpdateProxyConfigView generates the view for choosing whether to update Proxy configuration.
func chooseIfUpdateProxyConfigView(m *Model) string {
	s := "Do you want to update your Proxy configuration? \n\n"

	return s + formatChoicesWithCheckboxes(m.choices, m.cursor)
}

// proxyConfigView generates the view for Proxy configuration input.
func proxyConfigView(m *Model) string {
	s := "Introduce your Proxy configuration:\n\n"

	inputs := []string{
		m.proxyURL.View(),
		m.timeout.View(),
		m.socksURL.View(),
		m.headers.View(),
		m.skipVerify.View(),
	}

	return s + formatInputsWithNewlines(inputs)
}

// formatInputsWithNewlines formats a slice of strings into a single string with each input on a new line.
func formatInputsWithNewlines(inputs []string) string {
	s := ""
	for i := range inputs {
		s += inputs[i]
		if i < len(inputs)-1 {
			s += "\n"
		}
	}
	s += "\n"

	return s
}

// formatChoicesWithCheckboxes formats a slice of choices into a string with checkboxes, highlighting the current cursor position.
func formatChoicesWithCheckboxes(choices []string, cursor int) string {
	s := ""
	for i := range choices {
		s += fmt.Sprintf("%s\n", checkBox(choices[i], cursor == i))
	}
	return s
}
