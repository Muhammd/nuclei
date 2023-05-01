package templates

import (
	"github.com/projectdiscovery/nuclei/v2/internal/cloud/components/ui/form"
)

// initAddPlaylist initializes the add playlist input
func (model *Model) initAddPlaylist() {
	cfg := form.Config{
		Inputs: []form.Input{
			{
				PlaceHolder: "my-new-playlist",
				CharLimit:   30,
				Width:       30,
				Prompt:      "Template Playlist Name",
			},
		},
		InputStyle: inputStyle,
		OnSubmit: func(inputs []string) {
			model.formValues = inputs
			model.state = stateTemplates
		},
	}
	model.form = form.InitialModel(cfg)
}
