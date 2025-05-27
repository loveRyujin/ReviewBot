package form

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
			s += "\n\n"
		}
	}

	return s
}
