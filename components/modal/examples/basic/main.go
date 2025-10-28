// Package main demonstrates a basic modal dialog.
//
// This example shows:
//   - Creating a simple modal with content.
//   - Showing and hiding the modal.
//   - Closing with Esc key.
//
// Run: go run main.go.
package main

import (
	"fmt"
	"os"

	"github.com/phoenix-tui/phoenix/components/modal"
	tea "github.com/phoenix-tui/phoenix/tea"
)

// Model represents the application state.
type Model struct {
	modal *modal.Modal
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Check for quit keys when modal is not visible.
		if !m.modal.IsVisible() && (msg.String() == "q" || msg.String() == "ctrl+c") {
			return m, tea.Quit()
		}

		// Show modal on spacebar.
		if msg.String() == " " && !m.modal.IsVisible() {
			m.modal = m.modal.Show()
			return m, nil
		}

	case tea.WindowSizeMsg:
		// Pass window size to modal.
		updatedModal, cmd := m.modal.Update(msg)
		m.modal = updatedModal
		return m, cmd
	}

	// Forward all messages to modal when visible.
	if m.modal.IsVisible() {
		updatedModal, cmd := m.modal.Update(msg)
		m.modal = updatedModal
		return m, cmd
	}

	return m, nil
}

// View renders the UI.
func (m Model) View() string {
	if m.modal.IsVisible() {
		return m.modal.View()
	}

	return `Basic Modal Example

Press SPACE to show modal
Press Q or Ctrl+C to quit`
}

func main() {
	// Create modal.
	m := modal.New("This is a simple modal dialog.\n\nPress Esc to close.").
		Size(40, 10).
		DimBackground(true)

	// Create program.
	p := tea.New(Model{modal: m}, tea.WithAltScreen[Model]())

	// Run.
	if err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
