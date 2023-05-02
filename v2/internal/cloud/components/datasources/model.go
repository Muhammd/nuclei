package datasources

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	api "github.com/projectdiscovery/nuclei-cloud-api-go"
	"github.com/projectdiscovery/nuclei/v2/internal/cloud/components/ui/form"
	"github.com/projectdiscovery/nuclei/v2/internal/cloud/components/ui/selector"
	"github.com/projectdiscovery/nuclei/v2/internal/cloud/components/ui/styledlist"
	"github.com/projectdiscovery/nuclei/v2/internal/cloud/config"
)

type state int

const (
	stateList state = iota
	stateChoose
	stateAdd
)

type Model struct {
	state state

	choice *selector.Model

	form               *form.Model
	formValues         []string
	cloudlistProviders []string

	lists *styledlist.Model

	client *api.ClientWithResponses
	config *config.Config

	size tea.WindowSizeMsg

	formFields          []datasourcesFormConfig
	cloudlistFormFields map[string]datasourcesFormConfig
}

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#3A435E")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render

	inputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF06B7"))
)

func New(client *api.ClientWithResponses, config *config.Config) (*Model, error) {
	model := &Model{client: client, config: config}

	items, err := model.getDatasources()
	if err != nil {
		return nil, err
	}
	form, err := model.getDatasourcesFormConfig()
	if err != nil {
		return nil, err
	}

	model.formFields = form
	if len(items) == 0 {
		model.state = stateChoose
	} else {
		model.state = stateList
	}

	// FIXME: Add support for cloudlist provider form
	// generation. Allow user to choose subtype of datasource
	// and then generate form based on that.

	model.initChoice()
	model.initLists(items)
	return model, nil
}

func (m *Model) Run() error {
	program := tea.NewProgram(m, tea.WithAltScreen())
	_, err := program.Run()
	return err
}

func (m *Model) Init() tea.Cmd {
	switch m.state {
	case stateAdd:
		return m.form.Init()
	case stateChoose:
		return m.choice.Init()
	default:
		return m.lists.Init()
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle window size changes
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.size = msg
		m.lists.SetSize(m.size.Width, m.size.Height)
		//	m.playlists.SetSize(m.size.Width, m.size.Height)
	}

	switch m.state {
	case stateAdd:
		formModel, cmd := m.form.Update(msg)
		m.form = formModel.(*form.Model)
		return m, cmd
	case stateChoose:
		chooseModel, cmd := m.choice.Update(msg)
		m.choice = chooseModel.(*selector.Model)
		return m, cmd
	default:
		listsModel, cmd := m.lists.Update(msg)
		m.lists = listsModel.(*styledlist.Model)
		return m, cmd
	}
}

func (m *Model) View() string {
	switch m.state {
	case stateAdd:
		return m.form.View()
	case stateChoose:
		return m.choice.View()
	default:
		return m.lists.View()
	}
}
