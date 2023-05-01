package datasources

import (
	"github.com/projectdiscovery/nuclei/v2/internal/cloud/components/ui/selector"
)

func (m *Model) initChoice() {
	cfg := selector.Config{
		Prompt: "Select a datasource type to create",
		Choices: []string{
			"Github",
			"S3",
		},
		PromptStyle: inputStyle,
	}
	m.choice = selector.New(cfg)
	m.choice.OnChoice = func(choice string) {
		m.state = stateAdd
		m.initForm()
	}
}
