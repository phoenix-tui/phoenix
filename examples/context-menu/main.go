// Package main demonstrates Phoenix context menu positioning with smart edge detection.
// This example shows how to use phoenix/mouse CalculateMenuPosition to display
// context menus that stay fully visible at screen edges.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/phoenix-tui/phoenix/mouse"
	"github.com/phoenix-tui/phoenix/style"
	"github.com/phoenix-tui/phoenix/tea"
)

// menuItem represents a single item in the context menu.
type menuItem struct {
	label  string // Display text
	action string // Action identifier
}

// contextMenu represents a popup context menu with smart positioning.
type contextMenu struct {
	items    []menuItem // Menu items
	x        int        // Menu position X (adjusted for screen bounds)
	y        int        // Menu position Y (adjusted for screen bounds)
	width    int        // Menu width in cells
	height   int        // Menu height in cells (number of items)
	selected int        // Currently selected item index
	visible  bool       // Menu visibility
}

// model represents the application state following the Elm Architecture.
type model struct {
	mouse         *mouse.Mouse   // Mouse handler
	menu          contextMenu    // Context menu
	lastClick     mouse.Position // Last click position
	statusMessage string         // Status message
	width         int            // Terminal width
	height        int            // Terminal height
	ready         bool           // UI ready flag (after WindowSizeMsg)
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
		case "q", "ctrl+c", "esc":
			// Close menu if open, otherwise quit
			if m.menu.visible {
				m.menu.visible = false
				m.statusMessage = "Menu closed"
				return m, nil
			}
			// Disable mouse tracking before exit
			_ = m.mouse.Disable()
			return m, tea.Quit()

		case "up":
			// Navigate menu up
			if m.menu.visible && m.menu.selected > 0 {
				m.menu.selected--
			}

		case "down":
			// Navigate menu down
			if m.menu.visible && m.menu.selected < len(m.menu.items)-1 {
				m.menu.selected++
			}

		case "enter":
			// Activate selected menu item
			if m.menu.visible {
				selectedItem := m.menu.items[m.menu.selected]
				m.statusMessage = fmt.Sprintf("Activated: %s", selectedItem.action)
				m.menu.visible = false
				return m, nil
			}
		}

	case tea.MouseMsg:
		// Create mouse position
		pos := mouse.NewPosition(msg.X, msg.Y)

		// Handle mouse events based on action
		switch msg.Action {
		case tea.MouseActionPress:
			// Left-click outside menu: close menu
			if msg.Button == tea.MouseButtonLeft && m.menu.visible {
				if !m.isClickInMenu(pos) {
					m.menu.visible = false
					m.statusMessage = "Menu closed (clicked outside)"
				} else {
					// Click inside menu: activate item
					itemIndex := msg.Y - m.menu.y
					if itemIndex >= 0 && itemIndex < len(m.menu.items) {
						selectedItem := m.menu.items[itemIndex]
						m.statusMessage = fmt.Sprintf("Activated: %s", selectedItem.action)
						m.menu.visible = false
					}
				}
			}

			// Right-click: show context menu
			if msg.Button == tea.MouseButtonRight {
				m = m.showContextMenu(pos)
			}
		}

	case tea.WindowSizeMsg:
		// Handle terminal resize
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
	}

	return m, nil
}

// showContextMenu displays the context menu at the cursor position.
// Uses smart positioning to keep menu fully visible within screen bounds.
func (m model) showContextMenu(cursorPos mouse.Position) model {
	// Create menu items
	m.menu.items = []menuItem{
		{label: "Copy", action: "copy"},
		{label: "Paste", action: "paste"},
		{label: "Cut", action: "cut"},
		{label: "---", action: "separator"},
		{label: "Properties", action: "properties"},
		{label: "Delete", action: "delete"},
		{label: "Refresh", action: "refresh"},
	}

	// Calculate menu dimensions
	m.menu.width = 20 // Fixed width for this example
	m.menu.height = len(m.menu.items)

	// Calculate optimal position using MenuPositioner
	safePos := m.mouse.CalculateMenuPosition(
		cursorPos,
		m.menu.width,
		m.menu.height,
		m.width,
		m.height,
	)

	// Update menu position and visibility
	m.menu.x = safePos.X()
	m.menu.y = safePos.Y()
	m.menu.selected = 0
	m.menu.visible = true

	// Update status
	if safePos.Equals(cursorPos) {
		m.statusMessage = fmt.Sprintf("Menu at cursor (%d,%d)", cursorPos.X(), cursorPos.Y())
	} else {
		m.statusMessage = fmt.Sprintf("Menu adjusted: (%d,%d) -> (%d,%d)",
			cursorPos.X(), cursorPos.Y(), safePos.X(), safePos.Y())
	}

	m.lastClick = cursorPos

	return m
}

// isClickInMenu checks if a click position is inside the menu bounds.
func (m model) isClickInMenu(pos mouse.Position) bool {
	if !m.menu.visible {
		return false
	}

	x := pos.X()
	y := pos.Y()

	return x >= m.menu.x && x < m.menu.x+m.menu.width &&
		y >= m.menu.y && y < m.menu.y+m.menu.height
}

// View renders the current state to a string.
func (m model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	var b strings.Builder

	// Create styles
	titleStyle := style.New().
		Foreground(style.RGB(0, 255, 255)). // Cyan
		Bold(true)

	instructionStyle := style.New().
		Foreground(style.RGB(255, 255, 0)) // Yellow

	statusStyle := style.New().
		Foreground(style.RGB(0, 255, 0)) // Green

	menuItemStyle := style.New().
		Foreground(style.RGB(255, 255, 255)). // White
		Background(style.RGB(0, 0, 255))      // Blue

	menuItemSelectedStyle := style.New().
		Foreground(style.RGB(0, 0, 0)).     // Black
		Background(style.RGB(0, 255, 255)). // Cyan
		Bold(true)

	// Render title
	b.WriteString(style.Render(titleStyle, "Phoenix Context Menu Demo"))
	b.WriteString("\n\n")

	// Render instructions
	b.WriteString(style.Render(instructionStyle, "Instructions:"))
	b.WriteString("\n")
	b.WriteString("  • Right-click anywhere to open context menu\n")
	b.WriteString("  • Menu will adjust position at screen edges\n")
	b.WriteString("  • Click menu items or use ↑/↓ + Enter\n")
	b.WriteString("  • Click outside menu to close\n")
	b.WriteString("  • Press 'q' or ESC to quit\n")
	b.WriteString("\n")

	// Render status
	b.WriteString(style.Render(statusStyle, fmt.Sprintf("Status: %s", m.statusMessage)))
	b.WriteString("\n")

	// Render terminal info
	b.WriteString(fmt.Sprintf("Terminal size: %dx%d\n", m.width, m.height))

	if m.lastClick.X() >= 0 {
		b.WriteString(fmt.Sprintf("Last click: (%d,%d)\n", m.lastClick.X(), m.lastClick.Y()))
	}

	// Build screen buffer
	screen := make([][]rune, m.height)
	for i := range screen {
		screen[i] = make([]rune, m.width)
		for j := range screen[i] {
			screen[i][j] = ' '
		}
	}

	// Render instructions to buffer
	lines := strings.Split(b.String(), "\n")
	for i, line := range lines {
		if i >= m.height {
			break
		}
		runes := []rune(line)
		for j, r := range runes {
			if j >= m.width {
				break
			}
			screen[i][j] = r
		}
	}

	// Render context menu if visible (overlay on top)
	if m.menu.visible {
		// Draw menu border and items
		for i := 0; i < m.menu.height; i++ {
			menuY := m.menu.y + i
			if menuY < 0 || menuY >= m.height {
				continue
			}

			item := m.menu.items[i]
			isSelected := i == m.menu.selected

			// Render menu item
			var itemText string
			if isSelected {
				itemText = style.Render(menuItemSelectedStyle, fmt.Sprintf(" %-*s ", m.menu.width-2, item.label))
			} else {
				if item.action == "separator" {
					itemText = style.Render(menuItemStyle, strings.Repeat("─", m.menu.width))
				} else {
					itemText = style.Render(menuItemStyle, fmt.Sprintf(" %-*s ", m.menu.width-2, item.label))
				}
			}

			// Write to buffer
			runes := []rune(itemText)
			for j := 0; j < m.menu.width && j < len(runes); j++ {
				menuX := m.menu.x + j
				if menuX >= 0 && menuX < m.width {
					screen[menuY][menuX] = runes[j]
				}
			}
		}
	}

	// Convert buffer to string
	b.Reset()
	for _, line := range screen {
		b.WriteString(string(line))
		b.WriteString("\n")
	}

	return b.String()
}

func main() {
	// Create initial model
	initialModel := model{
		mouse:         mouse.New(),
		statusMessage: "Right-click anywhere to open context menu",
		lastClick:     mouse.NewPosition(-1, -1), // Invalid position initially
		width:         80,
		height:        24,
		ready:         false,
		menu: contextMenu{
			visible: false,
			items:   []menuItem{},
		},
	}

	// Create and run the program
	p := tea.New(
		initialModel,
		tea.WithAltScreen[model](),      // Use alternate screen buffer
		tea.WithMouseAllMotion[model](), // Enable mouse motion tracking
	)

	if err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nThank you for trying Phoenix context menu! 🔥")
}
