// Package main provides a basic example of using the Select component.
package main

import (
	"fmt"
	"os"

	selectcomponent "github.com/phoenix-tui/phoenix/components/select"
	"github.com/phoenix-tui/phoenix/tea"
)

// model is the application model containing the select component.
type model struct {
	selector *selectcomponent.Select[string]
	selected string
	quitting bool
}

func initialModel() model {
	frameworks := []string{
		"Phoenix TUI",
		"Charm/Bubbletea",
		"Ink (React for CLI)",
		"Textual (Python)",
		"Ratatui (Rust)",
		"FTXUI (C++)",
		"blessed (Node.js)",
		"urwid (Python)",
	}

	sel := selectcomponent.NewString("Which TUI framework do you prefer?", frameworks).
		WithHeight(10).
		WithFilterable(true)

	return model{
		selector: sel,
	}
}

func (m model) Init() tea.Cmd {
	return m.selector.Init()
}

func (m model) Update(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			m.quitting = true
			return m, tea.Quit()
		}

	case selectcomponent.ConfirmSelectionMsg[string]:
		if msg.OK {
			m.selected = msg.Value
			m.quitting = true
			return m, tea.Quit()
		}
	}

	// Update the selector
	newSelector, cmd := m.selector.Update(msg)
	m.selector = newSelector
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		if m.selected != "" {
			return fmt.Sprintf("\nYou selected: %s\n", m.selected)
		}
		return "\nCanceled.\n"
	}

	return m.selector.View() + "\n\nPress Enter to select, Esc to clear filter, Ctrl+C to quit"
}

func main() {
	p := tea.New(initialModel())
	if err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
