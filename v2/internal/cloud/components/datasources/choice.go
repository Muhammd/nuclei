package datasources

import (
	"github.com/projectdiscovery/nuclei/v2/internal/cloud/components/ui/selector"
)

func (m *Model) initChoice() {
	cfg := selector.Config{
		Prompt: selector.Prompt{
			Prompt: "Select a datasource type to create",
			Choices: []string{
				"Github",
				"S3",
				"Cloudlist",
			},
		},
		SubPrompt: map[string]selector.Prompt{
			"Cloudlist": {
				Prompt:  "Select a cloudlist type to create",
				Choices: m.cloudlistProviders,
			},
		},
		PromptStyle: inputStyle,
		OnChoice: func(choice []string) {
			m.state = stateAdd
			m.initForm()
		},
	}
	m.choice = selector.New(cfg)
}
