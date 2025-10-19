// Package main demonstrates a simple counter application using phoenix/tea.
//
// This example shows:
//   - Basic Model implementation
//   - Keyboard input handling
//   - State updates via Update()
//   - View rendering
//
// Controls:
//   - '+' or '=' : Increment counter
//   - '-' or '_' : Decrement counter
//   - 'q' or Ctrl+C : Quit
package main

import (
	"fmt"
	"os"

	"github.com/phoenix-tui/phoenix/tea/api"
)

// CounterModel represents the application state.
type CounterModel struct {
	count int
}

// Init initializes the model. No commands needed for this simple example.
func (m CounterModel) Init() api.Cmd {
	return nil
}

// Update handles incoming messages and updates the model.
func (m CounterModel) Update(msg api.Msg) (CounterModel, api.Cmd) {
	switch msg := msg.(type) {
	case api.KeyMsg:
		switch msg.String() {
		case "+", "=":
			// Increment
			m.count++
			return m, nil

		case "-", "_":
			// Decrement
			m.count--
			return m, nil

		case "q", "ctrl+c":
			// Quit
			return m, api.Quit()
		}
	}

	return m, nil
}

// View renders the current state as a string.
func (m CounterModel) View() string {
	return fmt.Sprintf(`
╔════════════════════════════╗
║       Counter Demo         ║
╠════════════════════════════╣
║                            ║
║    Count: %-15d ║
║                            ║
╠════════════════════════════╣
║  Controls:                 ║
║    +/= : Increment         ║
║    -/_ : Decrement         ║
║    q   : Quit              ║
╚════════════════════════════╝
`, m.count)
}

func main() {
	// Create initial model
	initialModel := CounterModel{count: 0}

	// Create program with alt screen (takes over terminal)
	p := api.New(initialModel, api.WithAltScreen[CounterModel]())

	// Run the program
	if err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
