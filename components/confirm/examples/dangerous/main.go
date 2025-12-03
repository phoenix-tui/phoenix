//nolint:gocritic
// Dangerous action example demonstrates safe defaults for destructive operations.
package main

import (
	"fmt"
	"os"

	"github.com/phoenix-tui/phoenix/components/confirm"
	"github.com/phoenix-tui/phoenix/tea"
)

type model struct {
	confirm *confirm.Confirm
	result  string
}

func (m model) Init() tea.Cmd {
	return m.confirm.Init()
}

func (m model) Update(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			m.result = "Canceled"
			return m, tea.Quit()
		}

	case confirm.ConfirmResultMsg:
		// User made a selection - store result before quitting
		if m.confirm.IsYes() {
			m.result = "Delete"
		} else if m.confirm.IsNo() {
			m.result = "Keep"
		} else {
			m.result = "Canceled"
		}
		return m, tea.Quit()
	}

	newConfirm, cmd := m.confirm.Update(msg)
	m.confirm = newConfirm

	return m, cmd
}

func (m model) View() string {
	header := "Dangerous Action Example\n"
	header += "========================\n\n"

	if m.result != "" {
		switch m.result {
		case "Delete":
			return header + "🗑️  Deleting all files... (not really in this example)\n"
		case "Keep":
			return header + "✅ Operation canceled. Files are safe.\n"
		default:
			return header + "❌ Canceled.\n"
		}
	}

	return header + m.confirm.View()
}

func main() {
	// For dangerous actions:
	// 1. Use descriptive title
	// 2. Add warning description
	// 3. Use custom labels (Delete instead of Yes)
	// 4. DefaultNo() for safe default
	// 5. Optional: WithCancel(true) for explicit escape
	c := confirm.New("Delete all files?").
		Description("This action cannot be undone. All data will be permanently lost.").
		Affirmative("Delete").
		Negative("Cancel").
		DefaultNo() // Safe default - user must explicitly choose Delete

	p := tea.New(model{confirm: c})
	err := p.Run()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
