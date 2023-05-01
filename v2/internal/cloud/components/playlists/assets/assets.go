package assets

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/nuclei/v2/internal/cloud/components/ui/styledlist"
)

// initAssets initializes the assets list
func (model *Model) initAssets(items []styledlist.Item) {
	keyMap := &styledlist.DelegateKeyMap{
		Entries: []key.Binding{
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "select"),
			),
			key.NewBinding(
				key.WithKeys("c"),
				key.WithHelp("c", "create"),
			),
		},
		UpdateFunc: []func(msg tea.Msg, m *list.Model) tea.Cmd{
			model.updateFuncAssetsEnterKey,
			model.updateFuncAssetsCreateKey,
		},
	}
	listConfig := styledlist.Config{
		Title:        "Assets to add to playlist",
		TitleStyle:   titleStyle,
		InitialItems: items,
		KeyMap:       keyMap,
	}
	model.assets = styledlist.New(listConfig)
}

// updateFuncAssetsEnterKey is the update function for enter key
func (model *Model) updateFuncAssetsEnterKey(msg tea.Msg, m *list.Model) tea.Cmd {
	var title string
	if i, ok := m.SelectedItem().(styledlist.Item); ok {
		title = i.Title()
	} else {
		return nil
	}
	index := m.Index()
	m.RemoveItem(index)
	model.assetsChoices = append(model.assetsChoices, model.assets.GetItemFromTitle(title))
	var cmd tea.Cmd
	return tea.Batch(m.NewStatusMessage(statusMessageStyle("Selected "+title)), cmd)
}

func (model *Model) updateFuncAssetsCreateKey(msg tea.Msg, m *list.Model) tea.Cmd {
	model.state = statePlaylist
	err := model.addAssetPlaylist(model.formValues[0], styledlist.IDsFromStyledList(model.assetsChoices))
	if err != nil {
		gologger.Error().Msgf("Could not add playlist: %s\n", err)
	}
	newList, err := model.getAssetsPlaylists()
	if err != nil {
		gologger.Error().Msgf("Could not get playlists: %s\n", err)
	}
	cmd := model.playlists.SetItems(newList, m)
	return tea.Batch(cmd, tea.ClearScreen)
}
