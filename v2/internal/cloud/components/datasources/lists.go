package datasources

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/projectdiscovery/nuclei/v2/internal/cloud/components/ui/styledlist"
)

func (model *Model) initLists(items []styledlist.Item) {
	keyMap := &styledlist.DelegateKeyMap{
		Entries: []key.Binding{
			key.NewBinding(
				key.WithKeys("a"),
				key.WithHelp("a", "add new"),
			),
		},
		UpdateFunc: []func(msg tea.Msg, m *list.Model) tea.Cmd{
			func(msg tea.Msg, m *list.Model) tea.Cmd {
				model.state = stateChoose
				return tea.ClearScreen
			},
		},
	}
	listConfig := styledlist.Config{
		Title:        "Datasources list",
		TitleStyle:   titleStyle,
		InitialItems: items,
		KeyMap:       keyMap,
	}
	model.lists = styledlist.New(listConfig)
}
