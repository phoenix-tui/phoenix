package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	progress "github.com/phoenix-tui/phoenix/components/progress/api"
	tea "github.com/phoenix-tui/phoenix/tea/api"
)

// Multi-progress example
// Demonstrates multiple progress bars and a spinner
type model struct {
	spinner progress.Spinner
	bars    []progress.Bar
	speeds  []int // Progress increment per tick
	count   int
}

func initialModel() model {
	return model{
		spinner: progress.NewSpinner("dots").Label("Overall progress"),
		bars: []progress.Bar{
			progress.NewBar(40).Label("Task 1").ShowPercent(true),
			progress.NewBar(40).Label("Task 2").ShowPercent(true),
			progress.NewBar(40).Label("Task 3").ShowPercent(true),
		},
		speeds: []int{3, 2, 1}, // Different speeds
		count:  0,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Init(),
		tickCmd(),
	)
}

func (m model) Update(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit()
		}

	case tea.TickMsg:
		// Update spinner
		updated, cmd := m.spinner.Update(msg)
		m.spinner = updated // Already *progress.Spinner type
		return m, cmd

	case tickProgressMsg:
		// Update progress bars
		allComplete := true
		for i := range m.bars {
			if !m.bars[i].IsComplete() {
				m.bars[i] = m.bars[i].Increment(m.speeds[i])
				allComplete = false
			}
		}

		m.count++

		// Quit if all bars complete or after 100 ticks
		if allComplete || m.count >= 100 {
			return m, tea.Quit()
		}

		return m, tickCmd()
	}

	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString("\n  Multi-Progress Example\n")
	b.WriteString("  ======================\n\n")

	// Render spinner
	b.WriteString("  ")
	b.WriteString(m.spinner.View())
	b.WriteString("\n\n")

	// Render all progress bars
	for _, bar := range m.bars {
		b.WriteString("  ")
		b.WriteString(bar.View())
		b.WriteString("\n")
	}

	b.WriteString("\n  Press 'q' to quit\n")

	return b.String()
}

// tickProgressMsg is sent to update progress bars
type tickProgressMsg struct{}

func tickCmd() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(200 * time.Millisecond) // 5 FPS
		return tickProgressMsg{}
	}
}

func main() {
	p := tea.New(initialModel())
	if err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
