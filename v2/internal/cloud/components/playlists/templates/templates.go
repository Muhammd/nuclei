package templates

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/nuclei/v2/internal/cloud/components/ui/styledlist"
)

// initTemplates initializes the templates list
func (model *Model) initTemplates(items []styledlist.Item) {
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
			model.updateFuncTemplatesEnterKey,
			model.updateFuncTemplatesCreateKey,
		},
	}
	listConfig := styledlist.Config{
		Title:        "Templates to add to playlist",
		TitleStyle:   titleStyle,
		InitialItems: items,
		KeyMap:       keyMap,
	}
	model.templates = styledlist.New(listConfig)
}

// updateFuncTemplatesEnterKey is the update function for enter key
func (model *Model) updateFuncTemplatesEnterKey(msg tea.Msg, m *list.Model) tea.Cmd {
	var title string
	if i, ok := m.SelectedItem().(styledlist.Item); ok {
		title = i.Title()
	} else {
		return nil
	}
	index := m.Index()
	m.RemoveItem(index)
	model.templatesChoices = append(model.templatesChoices, model.templates.GetItemFromTitle(title))
	var cmd tea.Cmd
	return tea.Batch(m.NewStatusMessage(statusMessageStyle("Selected "+title)), cmd)
}

func (model *Model) updateFuncTemplatesCreateKey(msg tea.Msg, m *list.Model) tea.Cmd {
	model.state = statePlaylist
	err := model.addTemplatePlaylist(model.formValues[0], styledlist.IDsFromStyledList(model.templatesChoices))
	if err != nil {
		gologger.Error().Msgf("Could not add playlist: %s\n", err)
	}
	newList, err := model.getTemplatePlaylists()
	if err != nil {
		gologger.Error().Msgf("Could not get playlists: %s\n", err)
	}
	cmd := model.playlists.SetItems(newList, m)
	return tea.Batch(cmd, tea.ClearScreen)
}
