// Package styledlist implements a styled list for nuclei
// cloud TUI resources like assets-list, templates-list etc.
package styledlist

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Horizonal and vertical padding for the list
	appStyle = lipgloss.NewStyle().Padding(1, 2)
)

// Model is the model for the styled list
type Model struct {
	list        list.Model
	items       []Item
	titleToItem map[string]Item
	keys        *listKeyMap
}

// Config is the config for the styled list
type Config struct {
	// Title is the title for the list
	Title string
	// TitleStyle is the style for the title
	TitleStyle lipgloss.Style
	// InitialItems is the initial items for the list
	InitialItems []Item
	// KeyMap is the keymap for the styled list
	KeyMap *DelegateKeyMap
}

// Item is the item for the styled list
type Item struct {
	ID              int64
	ItemTitle       string
	ItemDescription string
}

func (i Item) Title() string       { return i.ItemTitle }
func (i Item) Description() string { return i.ItemDescription }
func (i Item) FilterValue() string { return i.ItemTitle }

// IDsFromStyledList returns the ids from the styled list
func IDsFromStyledList(items []Item) []int64 {
	ids := make([]int64, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	return ids
}

// New creates a new styled list
func New(config Config) *Model {
	initialItems := make([]list.Item, 0, len(config.InitialItems))
	for _, item := range config.InitialItems {
		initialItems = append(initialItems, item)
	}
	listDelegate := newItemDelegate(config.KeyMap)
	listKeys := newListKeyMap()
	itemList := list.New(initialItems, listDelegate, 0, 0)
	itemList.Title = config.Title
	itemList.Styles.Title = config.TitleStyle
	itemList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.toggleStatusBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}
	model := &Model{
		list:        itemList,
		keys:        listKeys,
		titleToItem: make(map[string]Item, len(config.InitialItems)),
		items:       config.InitialItems,
	}
	for _, item := range config.InitialItems {
		model.titleToItem[item.ItemTitle] = item
	}
	return model
}

func (m *Model) Items() []Item {
	return m.items
}

func (m *Model) GetItemFromTitle(title string) Item {
	return m.titleToItem[title]
}

// SetItems sets the items for the styled list
func (m *Model) SetItems(items []Item, msg tea.Msg) tea.Cmd {
	m.items = items
	var list = make([]list.Item, 0, len(items))
	for _, item := range items {
		list = append(list, item)
	}
	m.list.SetItems(list)
	for _, item := range items {
		m.titleToItem[item.ItemTitle] = item
	}
	_, cmd := m.list.Update(msg)
	return cmd
}

func (m *Model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m *Model) SetSize(width, height int) {
	h, v := appStyle.GetFrameSize()
	m.list.SetSize(width-h, height-v)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		m, cmd, handled := m.handleUpdateKeypress(msg)
		if handled {
			return m, cmd
		}
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *Model) handleUpdateKeypress(msg tea.KeyMsg) (tea.Model, tea.Cmd, bool) {
	switch {
	case key.Matches(msg, m.keys.toggleStatusBar):
		m.list.SetShowStatusBar(!m.list.ShowStatusBar())
		return m, nil, true

	case key.Matches(msg, m.keys.togglePagination):
		m.list.SetShowPagination(!m.list.ShowPagination())
		return m, nil, true

	case key.Matches(msg, m.keys.toggleHelpMenu):
		m.list.SetShowHelp(!m.list.ShowHelp())
		return m, nil, true
	}
	return m, nil, false
}

func (m *Model) View() string {
	return appStyle.Render(m.list.View())
}
