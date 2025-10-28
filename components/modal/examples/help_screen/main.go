// Package main demonstrates a help screen modal.
//
// This example shows:
//   - Modal with multi-line content.
//   - Larger modal size for documentation.
//   - No buttons (Esc to close)
//   - Typical help screen use case.
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

		// Show help modal on '?' or 'h'.
		if !m.modal.IsVisible() && (msg.String() == "?" || msg.String() == "h") {
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

	return `Help Screen Example

Press ? or H to show help
Press Q or Ctrl+C to quit`
}

func main() {
	// Help text content.
	helpText := `Keyboard Shortcuts:

Navigation:
  Ctrl+N     New File
  Ctrl+O     Open File
  Ctrl+S     Save File
  Ctrl+W     Close File

Editing:
  Ctrl+C     Copy
  Ctrl+X     Cut
  Ctrl+V     Paste
  Ctrl+Z     Undo
  Ctrl+Y     Redo

Search:
  Ctrl+F     Find
  Ctrl+R     Replace
  Ctrl+G     Go to Line

Press Esc to close this help screen.`

	// Create help modal.
	m := modal.NewWithTitle("Help", helpText).
		Size(60, 20).
		DimBackground(true)

	// Create program.
	p := tea.New(Model{modal: m}, tea.WithAltScreen[Model]())

	// Run.
	if err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
