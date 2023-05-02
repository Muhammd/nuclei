package selector

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Horizonal and vertical padding for the list
	appStyle = lipgloss.NewStyle().Padding(1, 2)
)

type Model struct {
	cfg *Config

	cursor int
	choice string

	choices []string
}

type Config struct {
	Prompt Prompt

	// SubPrompts is a map of choices to subprompts. If a choice is selected
	SubPrompt map[string]Prompt

	PromptStyle lipgloss.Style
	OnChoice    func([]string)
}

type Prompt struct {
	Prompt  string
	Choices []string
}

func New(cfg Config) *Model {
	return &Model{cfg: &cfg}
}

func (m *Model) Choice() []string {
	return m.choices
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			m.choice = m.cfg.Prompt.Choices[m.cursor]
			m.choices = append(m.choices, m.choice)
			// If we have a subprompt, do not call OnChoice yet
			if newPrompt, ok := m.cfg.SubPrompt[m.choice]; ok {
				m.choice = ""
				m.cursor = 0
				m.cfg.Prompt = newPrompt
				return m, tea.ClearScreen
			}
			m.cfg.OnChoice(m.choices)
			return m, nil

		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.cfg.Prompt.Choices) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.cfg.Prompt.Choices) - 1
			}
		}
	}

	return m, nil
}

func (m *Model) View() string {
	s := strings.Builder{}
	s.WriteString(m.cfg.PromptStyle.Render(m.cfg.Prompt.Prompt))
	s.WriteString("\n\n")

	for i := 0; i < len(m.cfg.Prompt.Choices); i++ {
		if m.cursor == i {
			s.WriteString("[>] ")
		} else {
			s.WriteString("[ ] ")
		}
		s.WriteString(m.cfg.Prompt.Choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n(press q to quit)\n")

	return appStyle.Render(s.String())
}
