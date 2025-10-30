// Package main demonstrates drag scrolling in the Viewport component.
// This example shows how to:
//   - Enable mouse support for drag scrolling
//   - Handle mouse events (press, motion, release)
//   - Scroll content by clicking and dragging
//   - Visual feedback with scroll indicators
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/phoenix-tui/phoenix/components/viewport"
	"github.com/phoenix-tui/phoenix/tea"
)

// Model represents the application state.
type Model struct {
	viewport *viewport.Viewport
	ready    bool
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update processes messages and updates the model.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit()
		case "r":
			// Reset to top
			m.viewport = m.viewport.ScrollToTop()
		case "e":
			// Jump to end
			m.viewport = m.viewport.ScrollToBottom()
		}

	case tea.WindowSizeMsg:
		if !m.ready {
			// First resize - initialize viewport
			m.viewport = createViewport(msg.Width, msg.Height-4)
			m.ready = true
		} else {
			// Subsequent resizes
			m.viewport = m.viewport.SetSize(msg.Width, msg.Height-4)
		}
	}

	// Update viewport (handles mouse events)
	if m.ready {
		m.viewport, cmd = m.viewport.Update(msg)
	}

	return m, cmd
}

// View renders the UI.
func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	// Render viewport content
	content := m.viewport.View()

	// Calculate scroll indicator
	scrollPercent := 0
	if m.viewport.TotalLines() > m.viewport.Height() {
		scrollPercent = (m.viewport.ScrollOffset() * 100) / (m.viewport.TotalLines() - m.viewport.Height())
	}

	// Status bar
	status := fmt.Sprintf(
		"Drag Scroll Demo | Lines: %d/%d | Offset: %d | Scroll: %d%% | Press 'q' to quit, 'r' to reset, 'e' to end",
		m.viewport.Height(),
		m.viewport.TotalLines(),
		m.viewport.ScrollOffset(),
		scrollPercent,
	)

	// Instructions
	instructions := "Click and drag with LEFT mouse button to scroll | Use mouse wheel or arrow keys"

	// Scroll indicators
	scrollIndicator := ""
	if m.viewport.CanScrollUp() {
		scrollIndicator += "↑ "
	}
	if m.viewport.CanScrollDown() {
		scrollIndicator += "↓"
	}
	if scrollIndicator == "" {
		scrollIndicator = "—"
	}

	return fmt.Sprintf("%s\n\n%s\n\n%s [%s]", instructions, content, status, scrollIndicator)
}

// createViewport creates a viewport with sample content.
func createViewport(width, height int) *viewport.Viewport {
	// Generate sample content (100 lines)
	lines := make([]string, 100)
	for i := range lines {
		lines[i] = fmt.Sprintf("Line %3d: This is sample content for the drag scroll demo. Try clicking and dragging!", i+1)
	}

	return viewport.NewWithLines(lines, width, height).
		MouseEnabled(true) // Enable mouse wheel AND drag scrolling
}

func main() {
	// Create initial model
	m := Model{
		ready: false,
	}

	// Create program with mouse support
	p := tea.New(
		m,
		tea.WithAltScreen[Model](),
		tea.WithMouseAllMotion[Model](), // Required for drag scrolling
	)

	// Run the program
	if err := p.Run(); err != nil {
		log.Fatal(err)
	}

	// Clean exit message
	fmt.Println("Thanks for trying the drag scroll demo!")
	os.Exit(0)
}
