package datasources

import (
	"fmt"
	"strings"

	"github.com/projectdiscovery/nuclei/v2/internal/cloud/components/ui/form"
)

func (m *Model) initForm() {
	var inputs datasourcesFormConfig
	for _, input := range m.formFields {
		if strings.EqualFold(string(input.Type), m.choice.Choice()) {
			inputs = input
			break
		}
	}
	cfg := form.Config{
		Inputs:     inputs.Fields,
		InputStyle: inputStyle,
		OnSubmit: func(inputs []string) {
			fmt.Printf("Got inputs: %v\n", inputs)
		},
	}
	m.form = form.InitialModel(cfg)
}
