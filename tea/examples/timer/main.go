// Package main demonstrates a countdown timer application using phoenix/tea.
//
// This example shows:
//   - Asynchronous commands (Tick)
//   - Time-based updates
//   - State machine (stopped/running/paused)
//   - Command chaining
//
// Controls:
//   - 'space' : Start/Pause timer
//   - 'r' : Reset timer
//   - '+' : Add 10 seconds
//   - '-' : Subtract 10 seconds
//   - 'q' : Quit
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/phoenix-tui/phoenix/tea/api"
	"github.com/phoenix-tui/phoenix/tea/domain/service"
)

// TimerState represents the timer's current state.
type TimerState int

const (
	StateStopped TimerState = iota
	StateRunning
	StatePaused
)

// TimerModel represents the application state.
type TimerModel struct {
	state     TimerState
	remaining time.Duration
	initial   time.Duration
}

// Init initializes the model.
func (m TimerModel) Init() api.Cmd {
	return nil
}

// Update handles incoming messages and updates the model.
func (m TimerModel) Update(msg api.Msg) (TimerModel, api.Cmd) {
	switch msg := msg.(type) {
	case api.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, api.Quit()

		case " ":
			// Start/Pause toggle
			if m.state == StateRunning {
				m.state = StatePaused
				return m, nil
			}

			if m.remaining > 0 {
				m.state = StateRunning
				// Start ticking
				return m, service.Tick(1 * time.Second)
			}
			return m, nil

		case "r":
			// Reset timer
			m.state = StateStopped
			m.remaining = m.initial
			return m, nil

		case "+", "=":
			// Add 10 seconds
			m.remaining += 10 * time.Second
			if m.remaining > 99*time.Minute+59*time.Second {
				m.remaining = 99*time.Minute + 59*time.Second
			}
			m.initial = m.remaining
			return m, nil

		case "-", "_":
			// Subtract 10 seconds
			m.remaining -= 10 * time.Second
			if m.remaining < 0 {
				m.remaining = 0
			}
			m.initial = m.remaining
			return m, nil
		}

	case service.TickMsg:
		// Timer tick
		if m.state == StateRunning {
			m.remaining -= 1 * time.Second

			if m.remaining <= 0 {
				// Timer finished!
				m.remaining = 0
				m.state = StateStopped
				return m, nil
			}

			// Continue ticking
			return m, service.Tick(1 * time.Second)
		}
	}

	return m, nil
}

// View renders the current state as a string.
func (m TimerModel) View() string {
	minutes := int(m.remaining.Minutes())
	seconds := int(m.remaining.Seconds()) % 60

	stateStr := "STOPPED"
	stateIcon := "⏹"
	switch m.state {
	case StateRunning:
		stateStr = "RUNNING"
		stateIcon = "▶"
	case StatePaused:
		stateStr = "PAUSED"
		stateIcon = "⏸"
	}

	// Create progress bar (20 chars wide)
	progress := ""
	if m.initial > 0 {
		percent := float64(m.remaining) / float64(m.initial)
		filled := int(percent * 20)
		for i := 0; i < 20; i++ {
			if i < filled {
				progress += "█"
			} else {
				progress += "░"
			}
		}
	} else {
		progress = "░░░░░░░░░░░░░░░░░░░░"
	}

	return fmt.Sprintf(`
╔════════════════════════════════╗
║       Countdown Timer          ║
╠════════════════════════════════╣
║                                ║
║      %s   %s              ║
║                                ║
║        %02d:%02d                  ║
║                                ║
║  [%s]  ║
║                                ║
╠════════════════════════════════╣
║  Controls:                     ║
║    Space : Start/Pause         ║
║    r     : Reset               ║
║    +/-   : Add/Sub 10 seconds  ║
║    q     : Quit                ║
╚════════════════════════════════╝
`, stateIcon, stateStr, minutes, seconds, progress)
}

func main() {
	// Create initial model with 60 seconds
	initialModel := TimerModel{
		state:     StateStopped,
		remaining: 60 * time.Second,
		initial:   60 * time.Second,
	}

	// Create program
	p := api.New(initialModel, api.WithAltScreen[TimerModel]())

	// Run the program
	if err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
