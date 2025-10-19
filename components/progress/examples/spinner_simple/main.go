package main

import (
	"fmt"
	"os"

	progress "github.com/phoenix-tui/phoenix/components/progress/api"
	tea "github.com/phoenix-tui/phoenix/tea/api"
)

// Simple spinner example
// Demonstrates animated spinner with tea.Program
type model struct {
	spinner progress.Spinner
	count   int
}

func initialModel() model {
	return model{
		spinner: progress.NewSpinner("dots").Label("Loading..."),
		count:   0,
	}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Init()
}

func (m model) Update(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Quit on 'q' or Ctrl+C
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit()
		}

	case tea.TickMsg:
		// Update spinner and increment counter
		updated, cmd := m.spinner.Update(msg)
		m.spinner = updated // Already *progress.Spinner type
		m.count++

		// Auto-quit after 50 ticks (5 seconds at 10 FPS)
		if m.count >= 50 {
			return m, tea.Quit()
		}

		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf("\n  %s\n\n  Press 'q' to quit\n", m.spinner.View())
}

func main() {
	p := tea.New(initialModel())
	if err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
