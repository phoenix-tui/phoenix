// Package main demonstrates a confirmation dialog.
//
// This example shows:
//   - Modal with Yes/No buttons
//   - Handling button press events
//   - Button navigation (Tab, arrows)
//   - Keyboard shortcuts ('y', 'n')
//
// Run: go run main.go
package main

import (
	"fmt"
	"os"

	"github.com/phoenix-tui/phoenix/components/modal/api"
	tea "github.com/phoenix-tui/phoenix/tea/api"
)

// Model represents the application state.
type Model struct {
	modal  *modal.Modal
	result string // Result of user's choice
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Check for quit keys when modal is not visible
		if !m.modal.IsVisible() && (msg.String() == "q" || msg.String() == "ctrl+c") {
			return m, tea.Quit()
		}

		// Show modal on spacebar
		if msg.String() == " " && !m.modal.IsVisible() {
			m.result = "" // Reset result
			m.modal = m.modal.Show()
			return m, nil
		}

	case modal.ButtonPressedMsg:
		// Handle button press
		if msg.Action == "confirm" {
			m.result = "You clicked YES! File will be deleted."
			m.modal = m.modal.Hide()
			return m, nil
		} else if msg.Action == "cancel" {
			m.result = "You clicked NO. Operation cancelled."
			m.modal = m.modal.Hide()
			return m, nil
		}

	case tea.WindowSizeMsg:
		// Pass window size to modal
		updatedModal, cmd := m.modal.Update(msg)
		m.modal = updatedModal
		return m, cmd
	}

	// Forward all messages to modal when visible
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

	view := `Confirmation Dialog Example

Press SPACE to show confirmation dialog
Press Q or Ctrl+C to quit

`

	if m.result != "" {
		view += "\nLast result: " + m.result
	}

	return view
}

func main() {
	// Create confirmation modal
	m := modal.NewWithTitle("Confirm Action", "Are you sure you want to delete this file?").
		Size(50, 8).
		Buttons([]modal.Button{
			{Label: "Yes", Key: "y", Action: "confirm"},
			{Label: "No", Key: "n", Action: "cancel"},
		}).
		DimBackground(true)

	// Create program
	p := tea.New(Model{modal: m}, tea.WithAltScreen[Model]())

	// Run
	if err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
