//nolint:gocritic
// Basic example demonstrates simple Yes/No confirmation.
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
			m.result = "Yes"
		} else if m.confirm.IsNo() {
			m.result = "No"
		} else {
			m.result = "Canceled"
		}
		return m, tea.Quit()
	}

	// Delegate to confirm component
	newConfirm, cmd := m.confirm.Update(msg)
	m.confirm = newConfirm

	return m, cmd
}

func (m model) View() string {
	if m.result != "" {
		return fmt.Sprintf("\nYou selected: %s\n", m.result)
	}
	return m.confirm.View()
}

func main() {
	c := confirm.New("Do you want to continue?")

	p := tea.New(model{confirm: c})
	err := p.Run()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
