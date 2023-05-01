package templates

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	api "github.com/projectdiscovery/nuclei-cloud-api-go"
	"github.com/projectdiscovery/nuclei/v2/internal/cloud/components/ui/form"
	"github.com/projectdiscovery/nuclei/v2/internal/cloud/components/ui/styledlist"
	"github.com/projectdiscovery/nuclei/v2/internal/cloud/config"
)

type state int

const (
	statePlaylist state = iota
	stateAdd
	stateTemplates
)

type Model struct {
	state state

	form       *form.Model
	formValues []string

	playlists       *styledlist.Model
	playlistChoices []styledlist.Item

	templates        *styledlist.Model
	templatesChoices []styledlist.Item

	client *api.ClientWithResponses
	config *config.Config

	size tea.WindowSizeMsg
}

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render

	inputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF06B7"))
)

func New(client *api.ClientWithResponses, config *config.Config) (*Model, error) {
	model := &Model{client: client, config: config}

	items, err := model.getTemplatePlaylists()
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		model.state = stateAdd
	} else {
		model.state = statePlaylist
	}

	// FIXME: Do lazy loading on initialization of
	// templates to reduce the initial load time.
	templates, err := model.getTemplates()
	if err != nil {
		return nil, err
	}
	model.initTemplates(templates)
	model.initPlaylists(items)
	model.initAddPlaylist()
	return model, nil
}

func (m *Model) Run() ([]styledlist.Item, error) {
	program := tea.NewProgram(m, tea.WithAltScreen())
	_, err := program.Run()
	return m.playlistChoices, err
}

func (m *Model) Init() tea.Cmd {
	switch m.state {
	case stateAdd:
		return m.form.Init()
	case stateTemplates:
		return m.templates.Init()
	default:
		return m.playlists.Init()
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle window size changes
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.size = msg
		m.templates.SetSize(m.size.Width, m.size.Height)
		m.playlists.SetSize(m.size.Width, m.size.Height)
	}

	switch m.state {
	case stateAdd:
		formModel, cmd := m.form.Update(msg)
		m.form = formModel.(*form.Model)
		return m, cmd
	case stateTemplates:
		templatesModel, cmd := m.templates.Update(msg)
		m.templates = templatesModel.(*styledlist.Model)
		return m, cmd
	default:
		playlistModel, cmd := m.playlists.Update(msg)
		m.playlists = playlistModel.(*styledlist.Model)
		return m, cmd
	}
}

func (m *Model) View() string {
	switch m.state {
	case stateAdd:
		return m.form.View()
	case stateTemplates:
		return m.templates.View()
	default:
		return m.playlists.View()
	}
}
