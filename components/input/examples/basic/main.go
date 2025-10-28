// Package main demonstrates basic input component usage.
// This example shows a simple text input with basic editing capabilities.
package main

import (
	"fmt"
	"os"

	"github.com/phoenix-tui/phoenix/components/input"
	"github.com/phoenix-tui/phoenix/tea"
)

// basicModel wraps the input for demonstration.
type basicModel struct {
	input input.Input
}

func (m basicModel) Init() tea.Cmd {
	return m.input.Init()
}

func (m basicModel) Update(msg tea.Msg) (basicModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit()

		case "enter":
			// Print the value and quit.
			fmt.Printf("\nYou entered: %q\n", m.input.Value())
			return m, tea.Quit()
		}

	case tea.WindowSizeMsg:
		// Adjust input width based on window size.
		m.input = m.input.Width(msg.Width - 4)
		return m, nil
	}

	// Forward all other messages to input.
	updated, cmd := m.input.Update(msg)
	m.input = updated
	return m, cmd
}

func (m basicModel) View() string {
	return fmt.Sprintf(
		"Basic Text Input Example\n\n"+
			"Type something: %s\n\n"+
			"Press Enter to submit, Esc to quit",
		m.input.View(),
	)
}

func main() {
	// Create input with placeholder.
	inputField := input.New(40).
		Placeholder("Enter your name...").
		Focused(true)

	// Create model.
	model := basicModel{
		input: inputField,
	}

	// Run program.
	p := tea.New(model)
	if err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
