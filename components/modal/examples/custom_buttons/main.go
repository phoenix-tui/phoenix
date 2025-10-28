// Package main demonstrates custom button actions.
//
// This example shows:
//   - Modal with multiple custom buttons.
//   - Different button actions.
//   - Handling multiple action types.
//   - Real-world use case (save/discard/cancel)
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
	modal             *modal.Modal
	result            string // Result of user's choice
	hasUnsavedChanges bool
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages.
//
//nolint:gocyclo,cyclop // Example code demonstrates multiple modal interactions and state handling
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Check for quit keys when modal is not visible.
		if !m.modal.IsVisible() && (msg.String() == "q" || msg.String() == "ctrl+c") {
			// If there are unsaved changes, show confirmation modal.
			if m.hasUnsavedChanges {
				m.modal = m.modal.Show()
				return m, nil
			}
			return m, tea.Quit()
		}

		// Simulate making changes.
		if msg.String() == "e" && !m.modal.IsVisible() {
			m.hasUnsavedChanges = true
			m.result = "Document has been edited. Try quitting now (Q or Ctrl+C)."
			return m, nil
		}

	case modal.ButtonPressedMsg:
		// Handle button press.
		switch msg.Action {
		case "save":
			m.result = "Changes saved successfully!"
			m.hasUnsavedChanges = false
			m.modal = m.modal.Hide()
			return m, tea.Quit()
		case "discard":
			m.result = "Changes discarded."
			m.hasUnsavedChanges = false
			m.modal = m.modal.Hide()
			return m, tea.Quit()
		case "cancel":
			m.result = "Operation canceled. Continue editing."
			m.modal = m.modal.Hide()
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

	view := `Custom Buttons Example

Press E to simulate editing a document
Press Q or Ctrl+C to quit

Status: `

	if m.hasUnsavedChanges {
		view += "UNSAVED CHANGES"
	} else {
		view += "No changes"
	}

	view += "\n\n"

	if m.result != "" {
		view += "Last result: " + m.result
	}

	return view
}

func main() {
	// Create modal with custom buttons.
	m := modal.NewWithTitle("Unsaved Changes", "You have unsaved changes. What would you like to do?").
		Size(60, 10).
		Buttons([]modal.Button{
			{Label: "Save", Key: "s", Action: "save"},
			{Label: "Discard", Key: "d", Action: "discard"},
			{Label: "Cancel", Key: "c", Action: "cancel"},
		}).
		DimBackground(true)

	// Create program.
	p := tea.New(Model{
		modal:             m,
		hasUnsavedChanges: false,
	}, tea.WithAltScreen[Model]())

	// Run.
	if err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
