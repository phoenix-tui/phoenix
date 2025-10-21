// Package main demonstrates using Phoenix via the umbrella module.
// This example shows how to use the convenience API for common Phoenix operations.
package main

import (
	"fmt"
	"os"

	"github.com/phoenix-tui/phoenix"
	tea "github.com/phoenix-tui/phoenix/tea/api"
)

// Simple counter model implementing the tea.Model interface
type model struct {
	count int
}

// Init initializes the model (called once at startup)
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles incoming messages and updates the model
func (m model) Update(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, phoenix.Quit()
		case "+", "=":
			m.count++
		case "-":
			if m.count > 0 {
				m.count--
			}
		}
	}
	return m, nil
}

// View renders the current state of the model
func (m model) View() string {
	// For this simple demo, we'll use plain text
	// The full Style API is available through phoenix/style package
	return fmt.Sprintf(
		"Phoenix TUI - Umbrella Module Demo\n\n"+
			"Counter: %d\n\n"+
			"Press +/- to change, q to quit\n",
		m.count,
	)
}

func main() {
	fmt.Println("Phoenix TUI Framework - Umbrella Module Example")
	fmt.Println("=================================================\n")

	// Demonstrate terminal detection
	term := phoenix.AutoDetectTerminal()
	fmt.Printf("Terminal: %dx%d cells\n", term.Size().Width, term.Size().Height)
	fmt.Printf("Color depth: %v\n", term.Capabilities().ColorDepth())
	fmt.Printf("Mouse support: %v\n\n", term.Capabilities().SupportsMouse())

	// Demonstrate that Phoenix components are accessible
	_ = phoenix.NewStyle() // Style API available
	fmt.Println("Phoenix TUI libraries initialized successfully!")
	fmt.Println()

	// Demonstrate clipboard (optional - may not work in all environments)
	testClipboard()

	// Wait for user to press Enter
	fmt.Println("\nPress Enter to start the TUI application...")
	fmt.Scanln()

	// Run the TUI application
	p := phoenix.NewProgram(
		model{count: 0},
		phoenix.WithAltScreen[model](),      // Use alternate screen
		phoenix.WithMouseAllMotion[model](), // Enable mouse support
	)

	if err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nThank you for using Phoenix! 🔥")
}

// testClipboard demonstrates clipboard operations (may not work in all environments)
func testClipboard() {
	testText := "Phoenix TUI Framework - Umbrella Module Test"

	// Try to write to clipboard
	if err := phoenix.WriteClipboard(testText); err != nil {
		fmt.Printf("Clipboard write not available: %v\n", err)
		return
	}

	// Try to read from clipboard
	text, err := phoenix.ReadClipboard()
	if err != nil {
		fmt.Printf("Clipboard read not available: %v\n", err)
		return
	}

	if text == testText {
		fmt.Println("✓ Clipboard test successful!")
	} else {
		fmt.Printf("⚠ Clipboard content mismatch (expected %q, got %q)\n", testText, text)
	}
}
