package main

import (
	"fmt"
	"os"

	"github.com/phoenix-tui/phoenix/components/multiselect"
	"github.com/phoenix-tui/phoenix/tea"
)

// Model wraps the multiselect component
type Model struct {
	multi *multiselect.MultiSelect[string]
	done  bool
}

func main() {
	frameworks := []string{
		"Phoenix",
		"Charm",
		"Bubbletea",
		"Ink",
		"Textual",
		"Blessed",
		"FTXUI",
		"Imtui",
	}

	m := Model{
		multi: multiselect.NewStrings("Select your favorite TUI frameworks:", frameworks).
			WithFilterable(true).
			WithHeight(6),
		done: false,
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func (m Model) Init() tea.Cmd {
	return m.multi.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.done {
		return m, tea.Quit()
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit()
		}

	case multiselect.ConfirmSelectionMsg[string]:
		m.done = true
		fmt.Println("\nYou selected:")
		for _, val := range msg.Values {
			fmt.Printf("  - %s\n", val)
		}
		return m, tea.Quit()
	}

	// Update the multiselect component
	newMulti, cmd := m.multi.Update(msg)
	m.multi = newMulti
	return m, cmd
}

func (m Model) View() string {
	if m.done {
		return "" // Don't show anything after selection
	}
	return m.multi.View()
}
