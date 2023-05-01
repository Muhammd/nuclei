package add

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	ErrMsg error
)

const (
	Name = iota
	TeamSize
	Url
)

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
)

var (
	inputStyle    = lipgloss.NewStyle().Foreground(hotPink)
	continueStyle = lipgloss.NewStyle().Foreground(darkGray)
)

type Model struct {
	Inputs  []textinput.Model
	Focused int
	Err     error
}

func InitialModel() *Model {
	var inputs []textinput.Model = make([]textinput.Model, 3)
	inputs[Name] = textinput.New()
	inputs[Name].Placeholder = "new-glorious-workspace"
	inputs[Name].Focus()
	inputs[Name].CharLimit = 20
	inputs[Name].Width = 20
	inputs[Name].Prompt = ""

	inputs[TeamSize] = textinput.New()
	inputs[TeamSize].Placeholder = "1-10"
	inputs[TeamSize].CharLimit = 20
	inputs[TeamSize].Width = 20
	inputs[TeamSize].Prompt = ""

	inputs[Url] = textinput.New()
	inputs[Url].Placeholder = "www.example.com"
	inputs[Url].CharLimit = 20
	inputs[Url].Width = 20
	inputs[Url].Prompt = ""

	return &Model{
		Inputs:  inputs,
		Focused: 0,
		Err:     nil,
	}
}

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) View() string {
	return fmt.Sprintf(
		`
 %s
 %s

 %s
 %s

 %s
 %s

 %s
`,
		inputStyle.Width(30).Render("Workspace Name"),
		m.Inputs[Name].View(),
		inputStyle.Width(20).Render("Team Size"),
		m.Inputs[TeamSize].View(),
		inputStyle.Width(20).Render("Org URL"),
		m.Inputs[Url].View(),
		continueStyle.Render("Continue ->"),
	)
}

// NextInput focuses the next input field
func (m *Model) NextInput() {
	m.Focused = (m.Focused + 1) % len(m.Inputs)
}

// PrevInput focuses the previous input field
func (m *Model) PrevInput() {
	m.Focused--
	// Wrap around
	if m.Focused < 0 {
		m.Focused = len(m.Inputs) - 1
	}
}
