package styledlist

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// A list delegate that adds additional key bindings and help entries.
//
// UpdateFunc and Entries need to be of the same length.
type DelegateKeyMap struct {
	Entries    []key.Binding
	UpdateFunc []func(msg tea.Msg, m *list.Model) tea.Cmd
}

func (d DelegateKeyMap) ShortHelp() []key.Binding {
	return d.Entries
}

func (d DelegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{d.Entries}
}

// newItemDelegate returns a list delegate that adds additional key bindings and help entries.
func newItemDelegate(keys *DelegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			for i, entry := range keys.Entries {
				entry := entry
				if key.Matches(msg, entry) {
					return keys.UpdateFunc[i](msg, m)
				}
			}
		}
		return nil
	}
	d.ShortHelpFunc = func() []key.Binding {
		return keys.Entries
	}
	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{keys.Entries}
	}
	return d
}

// listKeyMap is a set of key bindings for the list.
type listKeyMap struct {
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
}

// newListKeyMap returns a new listKeyMap.
func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}
