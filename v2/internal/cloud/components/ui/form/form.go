package form

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	continueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#767676"))
)

type Model struct {
	inputs  []textinput.Model
	focused int
	err     error
	config  Config
}

type Config struct {
	Inputs     []Input
	InputStyle lipgloss.Style
	OnSubmit   func([]string)
}

type Input struct {
	PlaceHolder string
	CharLimit   int
	Width       int
	Prompt      string
}

func InitialModel(config Config) *Model {
	inputs := make([]textinput.Model, len(config.Inputs))
	for i, input := range config.Inputs {
		inputs[i] = textinput.New()
		inputs[i].Placeholder = input.PlaceHolder
		inputs[i].CharLimit = input.CharLimit
		inputs[i].Width = input.Width
		inputs[i].Prompt = ""
	}
	inputs[0].Focus()
	return &Model{inputs: inputs, focused: 0, err: nil, config: config}
}

// NextInput focuses the next input field
func (m *Model) NextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

// PrevInput focuses the previous input field
func (m *Model) PrevInput() {
	m.focused--
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) View() string {
	var builder strings.Builder
	for i, input := range m.inputs {
		origInput := m.config.Inputs[i]

		builder.WriteString("\n")
		builder.WriteString(m.config.InputStyle.Width(origInput.Width).Render(origInput.Prompt))
		builder.WriteString("\n")
		builder.WriteString(input.View())
		builder.WriteString("\n\n")
	}
	builder.WriteString(continueStyle.Render("Continue ->"))
	return builder.String()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 {
				values := make([]string, len(m.inputs))
				for i, input := range m.inputs {
					values[i] = input.Value()
				}
				m.config.OnSubmit(values)
				return m, nil
			}
			m.NextInput()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.PrevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.NextInput()
		}
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}
