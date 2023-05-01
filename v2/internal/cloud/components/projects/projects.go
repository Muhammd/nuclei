package projects

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pkg/errors"
	"github.com/projectdiscovery/gologger"
	api "github.com/projectdiscovery/nuclei-cloud-api-go"
	"github.com/projectdiscovery/nuclei/v2/internal/cloud/components/projects/add"
	"github.com/projectdiscovery/nuclei/v2/internal/cloud/components/utils"
	"github.com/projectdiscovery/nuclei/v2/internal/cloud/config"
)

// Model is a model for the Projects API TUI.
type Model struct {
	client           *api.ClientWithResponses
	list             list.Model
	choice           string
	quitting         bool
	state            state
	addModel         *add.Model
	projectsNameToID map[string]int64

	config *config.Config
}

type state int

const (
	stateNormal state = iota
	stateAdd
)

// New returns a new model for the project API TUI
func New(client *api.ClientWithResponses, config *config.Config) (*Model, error) {
	const defaultWidth = 20

	model := &Model{client: client, config: config}
	items, err := model.fetchLatestData()
	if err != nil {
		return nil, errors.Wrap(err, "could not fetch latest data")
	}
	l := list.New(items, utils.ListItemDelegate{}, defaultWidth, listHeight)
	l.Title = "Project to use"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{a}
	}
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{a}
	}
	if len(items) == 0 {
		gologger.Info().Msg("No projects found. Creating a new project...")
		model.state = stateAdd
	} else {
		model.state = stateNormal
	}
	model.list = l

	addModel := add.InitialModel()
	model.addModel = addModel
	return model, nil
}

// fetchLatestData fetches the latest data from the API
func (m *Model) fetchLatestData() ([]list.Item, error) {
	response, err := m.client.GetProjectsWithResponse(context.Background(), m.config.InternalIDs.WorkspaceID)
	if err != nil {
		return nil, errors.Wrap(err, "could not get projects")
	}
	items := []list.Item{}
	projectNamesToID := map[string]int64{}

	for _, item := range *response.JSON200 {
		item := item
		if item.Name == "" {
			continue
		}
		projectNamesToID[item.Name] = item.Id
		items = append(items, utils.ListItem(item.Name))
	}
	m.projectsNameToID = projectNamesToID
	return items, nil
}

const listHeight = 14

var (
	titleStyle = lipgloss.NewStyle().MarginLeft(2).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#0F5257")).
			Padding(0, 1)
	paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle       = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle   = lipgloss.NewStyle().Margin(1, 0, 2, 4)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
	selectionCheckmark = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
)

func (m *Model) Init() tea.Cmd {
	return nil
}

var (
	// Horizonal and vertical padding for the list
	appStyle = lipgloss.NewStyle().Padding(1, 2)
)

var a = key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "Add new"))

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case stateNormal:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			h, v := appStyle.GetFrameSize()
			m.list.SetSize(msg.Width-h, msg.Height-v)
			return m, nil

		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "a":
				m.state = stateAdd
				return m, nil
			case "ctrl+c":
				m.quitting = true
				return m, tea.Quit
			case "enter":
				i, ok := m.list.SelectedItem().(utils.ListItem)
				if ok {
					m.choice = string(i)
				}
				return m, tea.Quit
			}
		}
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, tea.Batch(cmd)
	case stateAdd:
		var cmds []tea.Cmd = make([]tea.Cmd, len(m.addModel.Inputs))

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				if m.addModel.Focused == len(m.addModel.Inputs)-1 {
					// last input
					project, err := m.addNewProject()
					if err != nil {
						gologger.Error().Msgf("Could not add new project: %s\n", err)
					}
					statusCmd := m.list.NewStatusMessage(statusMessageStyle("Added " + project))
					m.state = stateNormal
					m.list.Title = ""

					gologger.Info().Msgf("Fetching latest data\n")
					items, err := m.fetchLatestData()
					if err != nil {
						gologger.Error().Msgf("Could not fetch latest data: %s\n", err)
						return m, tea.Quit
					}

					insCmd := m.list.SetItems(items)
					var cmd tea.Cmd
					m.list, cmd = m.list.Update(msg)
					return m, tea.Batch(statusCmd, cmd, insCmd, tea.ClearScreen)
				}
				m.addModel.NextInput()
			case tea.KeyCtrlC, tea.KeyEsc:
				return m, tea.Quit
			case tea.KeyShiftTab, tea.KeyCtrlP:
				m.addModel.PrevInput()
			case tea.KeyTab, tea.KeyCtrlN:
				m.addModel.NextInput()
			}
			for i := range m.addModel.Inputs {
				m.addModel.Inputs[i].Blur()
			}
			m.addModel.Inputs[m.addModel.Focused].Focus()

		// We handle errors just like any other message
		case add.ErrMsg:
			m.addModel.Err = msg
			return m, nil
		}

		for i := range m.addModel.Inputs {
			m.addModel.Inputs[i], cmds[i] = m.addModel.Inputs[i].Update(msg)
		}
		return m, tea.Batch(cmds...)
	}
	return m, nil
}

func (m *Model) addNewProject() (string, error) {
	name := m.addModel.Inputs[add.Name].Value()

	_, err := m.client.PostWorkspacesWorkspaceIdProjectsWithResponse(context.Background(), m.config.InternalIDs.WorkspaceID, api.PostWorkspacesWorkspaceIdProjectsJSONRequestBody{
		Name: name,
	})
	if err != nil {
		return "", err
	}
	gologger.Info().Msgf("Adding new project %s\n", name)
	return name, nil
}

func (m *Model) View() string {
	switch m.state {
	case stateNormal:
		if m.choice != "" {
			return quitTextStyle.Render(fmt.Sprintf("%s Using project %s", selectionCheckmark, m.choice))
		}
		if m.quitting {
			return quitTextStyle.Render("Exiting projects selection")
		}
		return "\n" + m.list.View()
	case stateAdd:
		return m.addModel.View()
	}
	return ""
}

func (m *Model) Run() (string, int64, error) {
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		return "", -1, err
	}
	id, ok := m.projectsNameToID[m.choice]
	if !ok {
		return "", -1, fmt.Errorf("invalid project name")
	}
	return m.choice, id, nil
}
