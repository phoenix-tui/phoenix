// Package main demonstrates basic textarea usage.
package main

import (
	"fmt"
	"log"

	"github.com/phoenix-tui/phoenix/components/input/textarea/api"
	tea "github.com/phoenix-tui/phoenix/tea/api"
)

type model struct {
	textarea api.TextArea
	quitting bool
}

func (m model) Init() tea.Cmd {
	return m.textarea.Init()
}

func (m model) Update(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || (msg.Ctrl && msg.Rune == 'c') {
			m.quitting = true
			return m, tea.Quit()
		}
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	row, col := m.textarea.CursorPosition()
	return fmt.Sprintf(
		"TextArea Demo - Press Ctrl+C to quit\n\n%s\n\nCursor: row=%d, col=%d | Lines: %d | Chars: %d\n",
		m.textarea.View(),
		row,
		col,
		m.textarea.LineCount(),
		len(m.textarea.Value()),
	)
}

func main() {
	initialModel := model{
		textarea: api.New().
			Size(60, 10).
			Placeholder("Type something... (Emacs keybindings enabled)").
			Keybindings(api.KeybindingsEmacs),
	}

	p := tea.New(initialModel)
	if err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
