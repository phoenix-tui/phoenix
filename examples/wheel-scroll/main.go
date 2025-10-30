package main

import (
	"fmt"
	"strings"

	"github.com/phoenix-tui/phoenix/components/viewport"
	"github.com/phoenix-tui/phoenix/tea"
)

// model represents the application state
type model struct {
	viewport        *viewport.Viewport
	scrollSpeed     int  // Lines per wheel tick (1, 3, 5, 10)
	ready           bool // Window size received
	width           int  // Terminal width
	height          int  // Terminal height
	showHelp        bool // Show help overlay
	totalScrolled   int  // Track total lines scrolled
	wheelEventCount int  // Count wheel events
}

func initialModel() model {
	// Generate 150 lines of demo content
	lines := make([]string, 150)
	for i := range lines {
		lines[i] = fmt.Sprintf("Line %3d: This is content line number %d. Try scrolling with your mouse wheel!", i+1, i+1)
	}

	vp := viewport.New(80, 24).
		MouseEnabled(true).
		SetWheelScrollLines(3). // Default: 3 lines per wheel tick
		SetLines(lines)

	return model{
		viewport:    vp,
		scrollSpeed: 3,
		showHelp:    true,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit()

		case "h", "?":
			m.showHelp = !m.showHelp
			return m, nil

		case "1":
			// Slow scroll: 1 line per wheel tick
			m.scrollSpeed = 1
			m.viewport = m.viewport.SetWheelScrollLines(1)
			return m, nil

		case "2":
			// Default scroll: 3 lines per wheel tick
			m.scrollSpeed = 3
			m.viewport = m.viewport.SetWheelScrollLines(3)
			return m, nil

		case "3":
			// Fast scroll: 5 lines per wheel tick
			m.scrollSpeed = 5
			m.viewport = m.viewport.SetWheelScrollLines(5)
			return m, nil

		case "4":
			// Very fast scroll: 10 lines per wheel tick
			m.scrollSpeed = 10
			m.viewport = m.viewport.SetWheelScrollLines(10)
			return m, nil

		case "r":
			// Reset to top
			m.viewport = m.viewport.ScrollToTop()
			m.totalScrolled = 0
			m.wheelEventCount = 0
			return m, nil
		}

	case tea.MouseMsg:
		// Track wheel events for statistics
		if msg.Button == tea.MouseButtonWheelUp || msg.Button == tea.MouseButtonWheelDown {
			m.wheelEventCount++

			oldOffset := m.viewport.ScrollOffset()

			// Update viewport (will handle wheel scrolling)
			vp, cmd := m.viewport.Update(msg)
			m.viewport = vp

			newOffset := m.viewport.ScrollOffset()
			delta := newOffset - oldOffset
			if delta < 0 {
				delta = -delta
			}
			m.totalScrolled += delta

			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

		// Reserve space for header (3 lines) and footer (2 lines)
		viewportHeight := msg.Height - 5
		if viewportHeight < 5 {
			viewportHeight = 5
		}

		m.viewport = m.viewport.SetSize(msg.Width, viewportHeight)
		return m, nil
	}

	// Pass other messages to viewport
	vp, cmd := m.viewport.Update(msg)
	m.viewport = vp
	return m, cmd
}

func (m model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	var b strings.Builder

	// Header
	b.WriteString("╔════════════════════════════════════════════════════════════════════════════╗\n")
	b.WriteString(fmt.Sprintf("║ Mouse Wheel Scrolling Demo - Speed: %d lines/tick %-24s ║\n",
		m.scrollSpeed, "[Press 'h' for help]"))
	b.WriteString("╚════════════════════════════════════════════════════════════════════════════╝\n")

	// Viewport content
	b.WriteString(m.viewport.View())
	b.WriteString("\n")

	// Footer with stats
	scrollPercent := 0
	if m.viewport.TotalLines() > m.viewport.Height() {
		maxScroll := m.viewport.TotalLines() - m.viewport.Height()
		if maxScroll > 0 {
			scrollPercent = (m.viewport.ScrollOffset() * 100) / maxScroll
		}
	}

	b.WriteString(fmt.Sprintf("Line %d/%d (%.0f%%) | Scrolled: %d lines | Events: %d | Speed: %d lines/tick",
		m.viewport.ScrollOffset()+1,
		m.viewport.TotalLines(),
		float64(scrollPercent),
		m.totalScrolled,
		m.wheelEventCount,
		m.scrollSpeed,
	))

	// Help overlay
	if m.showHelp {
		help := `
┌─────────────────────────────────────┐
│         Wheel Scroll Demo           │
├─────────────────────────────────────┤
│ Mouse Wheel Up/Down: Scroll content │
│ 1: Slow (1 line/tick)               │
│ 2: Default (3 lines/tick)           │
│ 3: Fast (5 lines/tick)              │
│ 4: Very Fast (10 lines/tick)        │
│ r: Reset to top                     │
│ h: Toggle help                      │
│ q: Quit                             │
└─────────────────────────────────────┘`

		// Center the help box
		helpLines := strings.Split(help, "\n")
		topPadding := (m.height - len(helpLines)) / 2
		leftPadding := (m.width - 40) / 2

		var overlay strings.Builder
		for i := 0; i < topPadding; i++ {
			overlay.WriteString("\n")
		}
		for _, line := range helpLines {
			if line != "" {
				overlay.WriteString(strings.Repeat(" ", leftPadding))
				overlay.WriteString(line)
			}
			overlay.WriteString("\n")
		}

		return overlay.String()
	}

	return b.String()
}

func main() {
	p := tea.New(
		initialModel(),
		tea.WithAltScreen[model](),
		tea.WithMouseAllMotion[model](), // Enable mouse support
	)

	if err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
