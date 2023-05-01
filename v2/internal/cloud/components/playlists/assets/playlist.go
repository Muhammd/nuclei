package assets

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/projectdiscovery/nuclei/v2/internal/cloud/components/ui/styledlist"
)

// initPlaylists initializes the playlists list
func (model *Model) initPlaylists(items []styledlist.Item) {
	keyMap := &styledlist.DelegateKeyMap{
		Entries: []key.Binding{
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "choose"),
			),
			key.NewBinding(
				key.WithKeys("a"),
				key.WithHelp("a", "add new"),
			),
		},
		UpdateFunc: []func(msg tea.Msg, m *list.Model) tea.Cmd{
			model.updateFuncPlaylistEnterKey,
			func(msg tea.Msg, m *list.Model) tea.Cmd {
				model.state = stateAdd
				return tea.ClearScreen
			},
		},
	}
	listConfig := styledlist.Config{
		Title:        "Assets playlists to use",
		TitleStyle:   titleStyle,
		InitialItems: items,
		KeyMap:       keyMap,
	}
	model.playlists = styledlist.New(listConfig)
}

// updateFuncPlaylistEnterKey is the update function for enter key
func (model *Model) updateFuncPlaylistEnterKey(msg tea.Msg, m *list.Model) tea.Cmd {
	var title string
	if i, ok := m.SelectedItem().(styledlist.Item); ok {
		title = i.Title()
	} else {
		return nil
	}
	index := m.Index()
	m.RemoveItem(index)
	model.playlistChoices = append(model.playlistChoices, model.playlists.GetItemFromTitle(title))
	var cmd tea.Cmd
	if len(m.Items()) == 0 {
		cmd = tea.Quit
	}
	return tea.Batch(m.NewStatusMessage(statusMessageStyle("Selected "+title)), cmd)
}
