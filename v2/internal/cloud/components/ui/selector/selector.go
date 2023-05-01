package selector

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	prompt      string
	choices     []string
	cursor      int
	choice      string
	OnChoice    func(string)
	promptStyle lipgloss.Style
}

type Config struct {
	Prompt      string
	Choices     []string
	PromptStyle lipgloss.Style
}

func New(cfg Config) *Model {
	return &Model{prompt: cfg.Prompt, choices: cfg.Choices, promptStyle: cfg.PromptStyle}
}

func (m *Model) Choice() string {
	return m.choice
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
			m.choice = m.choices[m.cursor]
			m.OnChoice(m.choice)
			return m, nil

		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.choices) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.choices) - 1
			}
		}
	}

	return m, nil
}

func (m *Model) View() string {
	s := strings.Builder{}
	s.WriteString(m.promptStyle.Render(m.prompt))
	s.WriteString("\n\n")

	for i := 0; i < len(m.choices); i++ {
		if m.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(m.choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n(press q to quit)\n")

	return s.String()
}
