package form

import "fmt"

func aiConfigView(m *Model) string {
	s := "Introduce your AI configuration:\n\n"

	inputs := []string{
		m.provider.View(),
		m.apiKey.View(),
		m.model.View(),
		m.baseURL.View(),
	}

	for i := range inputs {
		s += inputs[i]
		if i < len(inputs)-1 {
			s += "\n"
		}
	}

	s += "\n"

	return s
}

func chooseIfUpdateGitConfigView(m *Model) string {
	s := "Do you want to update your Git configuration? \n\n"

	for i := range m.choices {
		s += fmt.Sprintf("%s\n", checkBox(m.choices[i], m.cursor == i))
	}

	return s
}

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

	for i := range inputs {
		s += inputs[i]
		if i < len(inputs)-1 {
			s += "\n"
		}
	}

	s += "\n"

	return s
}

func chooseIfUpdateProxyConfigView(m *Model) string {
	s := "Do you want to update your Proxy configuration? \n\n"

	for i := range m.choices {
		s += fmt.Sprintf("%s\n", checkBox(m.choices[i], m.cursor == i))
	}
	return s
}

func proxyConfigView(m *Model) string {
	s := "Introduce your Proxy configuration:\n\n"

	inputs := []string{
		m.proxyURL.View(),
		m.timeout.View(),
		m.socksURL.View(),
		m.headers.View(),
		m.skipVerify.View(),
	}

	for i := range inputs {
		s += inputs[i]
		if i < len(inputs)-1 {
			s += "\n"
		}
	}

	s += "\n"

	return s
}
