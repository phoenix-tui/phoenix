// Package main demonstrates Phoenix mouse hover detection and highlighting.
// This example shows how to use phoenix/mouse for interactive button hover states.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/phoenix-tui/phoenix/mouse"
	"github.com/phoenix-tui/phoenix/style"
	"github.com/phoenix-tui/phoenix/tea"
)

// button represents a clickable UI element with hover detection.
type button struct {
	id   string            // Unique identifier
	text string            // Display text
	x    int               // Column position (0-based)
	y    int               // Row position (0-based)
	area mouse.BoundingBox // Hover detection area
}

// model represents the application state following the Elm Architecture.
type model struct {
	mouse         *mouse.Mouse // Mouse handler
	buttons       []button     // UI buttons
	hoveredID     string       // Currently hovered button ID
	lastClickedID string       // Last clicked button ID
	width         int          // Terminal width
	height        int          // Terminal height
	ready         bool         // UI ready flag (after WindowSizeMsg)
}

// Init initializes the model (called once at startup).
func (m model) Init() tea.Cmd {
	// Enable mouse tracking
	if err := m.mouse.Enable(); err != nil {
		// In production, you might want to log this error
		_ = err
	}
	return nil
}

// Update handles incoming messages and updates the model.
func (m model) Update(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle keyboard input
		switch msg.String() {
		case "q", "ctrl+c":
			// Cleanup: disable mouse tracking before quit
			_ = m.mouse.Disable()
			return m, tea.Quit()
		case "r":
			// Reset state
			m.hoveredID = ""
			m.lastClickedID = ""
		}

	case tea.MouseMsg:
		// Handle mouse events
		m = m.handleMouseEvent(msg)

	case tea.WindowSizeMsg:
		// Handle terminal resize
		if !msg.IsValid() {
			return m, nil
		}

		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

		// Recalculate button positions and areas based on new size
		m = m.layoutButtons()
	}

	return m, nil
}

// handleMouseEvent processes mouse events for hover detection and clicks.
func (m model) handleMouseEvent(msg tea.MouseMsg) model {
	// Create position from mouse coordinates
	pos := mouse.NewPosition(msg.X, msg.Y)

	// Build component areas for hover detection
	areas := make([]mouse.ComponentArea, len(m.buttons))
	for i, btn := range m.buttons {
		areas[i] = mouse.ComponentArea{
			ID:   btn.id,
			Area: btn.area,
		}
	}

	// Process hover state
	eventType := m.mouse.ProcessHover(pos, areas)

	switch eventType {
	case mouse.EventHoverEnter:
		// Mouse entered a component
		m.hoveredID = m.mouse.CurrentHoverComponent()

	case mouse.EventHoverLeave:
		// Mouse left all components
		m.hoveredID = ""

	case mouse.EventHoverMove:
		// Mouse moved within a component (ID unchanged)
		// No action needed, but you could track position if needed

	case mouse.EventMotion:
		// Mouse moved outside all components
		// No action needed
	}

	// Handle click (button release over a component)
	if msg.Action == tea.MouseActionRelease && msg.Button == tea.MouseButtonLeft {
		if m.hoveredID != "" {
			m.lastClickedID = m.hoveredID
		}
	}

	return m
}

// layoutButtons calculates button positions based on terminal size.
func (m model) layoutButtons() model {
	// Layout: 3 buttons in row 1, 2 buttons in row 2, 1 button in row 3
	// Center horizontally with spacing

	if m.width < 40 || m.height < 15 {
		// Terminal too small, keep existing layout
		return m
	}

	buttonWidth := 14    // Fixed button width
	buttonHeight := 3    // Fixed button height
	spacing := 2         // Horizontal spacing between buttons
	verticalSpacing := 1 // Vertical spacing between rows

	// Row 1: 3 buttons centered
	row1Y := 3
	row1Buttons := 3
	totalWidth1 := row1Buttons*buttonWidth + (row1Buttons-1)*spacing
	startX1 := (m.width - totalWidth1) / 2

	m.buttons[0] = button{
		id:   "button1",
		text: "Button 1",
		x:    startX1,
		y:    row1Y,
		area: mouse.NewBoundingBox(startX1, row1Y, buttonWidth, buttonHeight),
	}

	m.buttons[1] = button{
		id:   "button2",
		text: "Button 2",
		x:    startX1 + buttonWidth + spacing,
		y:    row1Y,
		area: mouse.NewBoundingBox(startX1+buttonWidth+spacing, row1Y, buttonWidth, buttonHeight),
	}

	m.buttons[2] = button{
		id:   "button3",
		text: "Button 3",
		x:    startX1 + 2*(buttonWidth+spacing),
		y:    row1Y,
		area: mouse.NewBoundingBox(startX1+2*(buttonWidth+spacing), row1Y, buttonWidth, buttonHeight),
	}

	// Row 2: 2 buttons centered
	row2Y := row1Y + buttonHeight + verticalSpacing
	row2Buttons := 2
	totalWidth2 := row2Buttons*buttonWidth + (row2Buttons-1)*spacing
	startX2 := (m.width - totalWidth2) / 2

	m.buttons[3] = button{
		id:   "button4",
		text: "Button 4",
		x:    startX2,
		y:    row2Y,
		area: mouse.NewBoundingBox(startX2, row2Y, buttonWidth, buttonHeight),
	}

	m.buttons[4] = button{
		id:   "button5",
		text: "Button 5",
		x:    startX2 + buttonWidth + spacing,
		y:    row2Y,
		area: mouse.NewBoundingBox(startX2+buttonWidth+spacing, row2Y, buttonWidth, buttonHeight),
	}

	// Row 3: 1 button centered
	row3Y := row2Y + buttonHeight + verticalSpacing
	startX3 := (m.width - buttonWidth) / 2

	m.buttons[5] = button{
		id:   "button6",
		text: "Button 6",
		x:    startX3,
		y:    row3Y,
		area: mouse.NewBoundingBox(startX3, row3Y, buttonWidth, buttonHeight),
	}

	return m
}

// View renders the current state of the model.
func (m model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	var b strings.Builder

	// Title
	title := "Phoenix TUI - Hover Highlight Demo"
	titleStyle := style.New().
		Foreground(style.RGB(255, 255, 255)).
		Bold(true)
	b.WriteString(style.Render(titleStyle, title))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("═", len(title)))
	b.WriteString("\n\n")

	// Render buttons
	// We'll build a 2D grid and render it
	grid := make([][]rune, m.height)
	for i := range grid {
		grid[i] = make([]rune, m.width)
		for j := range grid[i] {
			grid[i][j] = ' ' // Initialize with spaces
		}
	}

	// Draw each button
	for _, btn := range m.buttons {
		m.drawButton(&grid, btn, btn.id == m.hoveredID)
	}

	// Convert grid to string
	for y := 0; y < m.height-6; y++ { // Leave room for status
		if y >= len(grid) {
			break
		}
		b.WriteString(string(grid[y]))
		b.WriteString("\n")
	}

	// Status section (at bottom)
	b.WriteString("\n")

	statusStyle := style.New().Foreground(style.RGB(150, 150, 150))

	// Hovered status
	hoveredText := "Hovered: "
	if m.hoveredID != "" {
		hoveredText += m.hoveredID
	} else {
		hoveredText += "(none)"
	}
	b.WriteString(style.Render(statusStyle, hoveredText))
	b.WriteString("\n")

	// Last clicked status
	clickedText := "Last clicked: "
	if m.lastClickedID != "" {
		clickStyle := style.New().
			Foreground(style.RGB(0, 255, 0)).
			Bold(true)
		clickedText += style.Render(clickStyle, m.lastClickedID)
	} else {
		clickedText += "(none)"
	}
	b.WriteString(style.Render(statusStyle, clickedText))
	b.WriteString("\n\n")

	// Instructions
	instructStyle := style.New().
		Foreground(style.RGB(100, 150, 200)).
		Italic(true)
	b.WriteString(style.Render(instructStyle, "Move mouse to hover buttons • Click to select • Press 'r' to reset • Press 'q' to quit"))

	return b.String()
}

// drawButton renders a button into the grid with optional hover highlighting.
func (m model) drawButton(grid *[][]rune, btn button, isHovered bool) {
	// Button dimensions from bounding box
	x := btn.area.X()
	y := btn.area.Y()
	w := btn.area.Width()
	h := btn.area.Height()

	// Validate bounds
	if y < 0 || y+h > len(*grid) {
		return
	}

	// Choose border style based on hover state
	var topLeft, topRight, bottomLeft, bottomRight, horizontal, vertical rune
	if isHovered {
		// Hovered: double-line border (more prominent)
		topLeft, topRight = '╔', '╗'
		bottomLeft, bottomRight = '╚', '╝'
		horizontal, vertical = '═', '║'
	} else {
		// Normal: single-line border
		topLeft, topRight = '╭', '╮'
		bottomLeft, bottomRight = '╰', '╯'
		horizontal, vertical = '─', '│'
	}

	// Draw button border and background
	for dy := 0; dy < h; dy++ {
		row := y + dy
		if row >= len(*grid) {
			break
		}

		for dx := 0; dx < w; dx++ {
			col := x + dx
			if col >= len((*grid)[row]) {
				break
			}

			// Draw border
			if dy == 0 {
				// Top border
				if dx == 0 {
					(*grid)[row][col] = topLeft
				} else if dx == w-1 {
					(*grid)[row][col] = topRight
				} else {
					(*grid)[row][col] = horizontal
				}
			} else if dy == h-1 {
				// Bottom border
				if dx == 0 {
					(*grid)[row][col] = bottomLeft
				} else if dx == w-1 {
					(*grid)[row][col] = bottomRight
				} else {
					(*grid)[row][col] = horizontal
				}
			} else {
				// Middle rows
				if dx == 0 || dx == w-1 {
					(*grid)[row][col] = vertical
				} else if dy == h/2 {
					// Center row: write button text
					textStartCol := x + (w-len(btn.text))/2
					if col >= textStartCol && col < textStartCol+len(btn.text) {
						(*grid)[row][col] = rune(btn.text[col-textStartCol])
					} else {
						(*grid)[row][col] = ' '
					}
				} else {
					(*grid)[row][col] = ' '
				}
			}
		}
	}

	// Note: For simplicity, this basic renderer doesn't apply ANSI codes to individual cells.
	// In a real TUI, you'd use phoenix/render or apply styles per-cell.
	// This example focuses on demonstrating hover detection logic.
}

// initialModel creates the initial model with default state.
func initialModel() model {
	// Create 6 buttons (positions will be calculated on first WindowSizeMsg)
	buttons := make([]button, 6)
	for i := 0; i < 6; i++ {
		buttons[i] = button{
			id:   fmt.Sprintf("button%d", i+1),
			text: fmt.Sprintf("Button %d", i+1),
		}
	}

	return model{
		mouse:         mouse.New(),
		buttons:       buttons,
		hoveredID:     "",
		lastClickedID: "",
		width:         80, // Default, will be updated by WindowSizeMsg
		height:        24, // Default, will be updated by WindowSizeMsg
		ready:         false,
	}
}

func main() {
	// Create program with mouse support
	p := tea.New(
		initialModel(),
		tea.WithAltScreen[model](),      // Use alternate screen buffer
		tea.WithMouseAllMotion[model](), // Enable mouse motion tracking
	)

	// Run the TUI
	if err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nThank you for trying Phoenix hover detection! 🔥")
}
