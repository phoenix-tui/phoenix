// Package main demonstrates a todo list application using phoenix/tea.
//
// This example shows:
//   - Complex state management (list of items)
//   - Multiple input modes (viewing vs editing)
//   - List navigation
//   - Item addition/removal
//
// Controls:
//   - 'a' : Add new todo
//   - 'd' : Delete selected todo
//   - 'j' or ↓ : Move down
//   - 'k' or ↑ : Move up
//   - 'space' : Toggle completion
//   - 'q' : Quit
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/phoenix-tui/phoenix/tea"
)

// TodoItem represents a single todo item.
type TodoItem struct {
	text      string
	completed bool
}

// TodoModel represents the application state.
type TodoModel struct {
	items   []TodoItem
	cursor  int
	addMode bool
	newText string
}

// Init initializes the model with some sample todos.
func (m TodoModel) Init() tea.Cmd {
	return nil
}

// Update handles incoming messages and updates the model.
//
//nolint:gocognit,gocyclo,cyclop // Example todo app logic is naturally complex for demonstration
func (m TodoModel) Update(msg tea.Msg) (TodoModel, tea.Cmd) {
	// Single type switch is clear for examples (simple pattern)
	//nolint:gocritic // singleCaseSwitch: Keep for example clarity
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.addMode {
			// Add mode - entering new todo text
			return m.handleAddMode(msg)
		}

		// Normal mode - list navigation
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit()

		case "a":
			// Enter add mode
			m.addMode = true
			m.newText = ""
			return m, nil

		case "d":
			// Delete current item
			if len(m.items) > 0 {
				m.items = append(m.items[:m.cursor], m.items[m.cursor+1:]...)
				if m.cursor >= len(m.items) && m.cursor > 0 {
					m.cursor--
				}
			}
			return m, nil

		case "j", "down":
			// Move down
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
			return m, nil

		case "k", "up":
			// Move up
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil

		case " ":
			// Toggle completion
			if len(m.items) > 0 {
				m.items[m.cursor].completed = !m.items[m.cursor].completed
			}
			return m, nil
		}
	}

	return m, nil
}

// handleAddMode handles keyboard input while adding a new todo.
func (m TodoModel) handleAddMode(msg tea.KeyMsg) (TodoModel, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Add the todo
		if m.newText != "" {
			m.items = append(m.items, TodoItem{text: m.newText, completed: false})
			m.cursor = len(m.items) - 1
		}
		m.addMode = false
		m.newText = ""
		return m, nil

	case "esc", "ctrl+c":
		// Cancel add
		m.addMode = false
		m.newText = ""
		return m, nil

	case "backspace":
		// Delete character
		if m.newText != "" {
			m.newText = m.newText[:len(m.newText)-1]
		}
		return m, nil

	default:
		// Add character (if it's a rune)
		if msg.Type == tea.KeyRune {
			m.newText += string(msg.Rune)
		}
		return m, nil
	}
}

// View renders the current state as a string.
func (m TodoModel) View() string {
	var b strings.Builder

	b.WriteString("╔═══════════════════════════════════════╗\n")
	b.WriteString("║           Todo List Demo              ║\n")
	b.WriteString("╠═══════════════════════════════════════╣\n")
	//nolint:nestif // View rendering logic for add/normal modes
	if m.addMode {
		// Add mode view
		b.WriteString("║ Adding new todo:                      ║\n")
		b.WriteString(fmt.Sprintf("║ > %-35s ║\n", m.newText+"_"))
		b.WriteString("║                                       ║\n")
		b.WriteString("║ Enter: Save  |  Esc: Cancel          ║\n")
	} else {
		// Normal mode view
		if len(m.items) == 0 {
			b.WriteString("║                                       ║\n")
			b.WriteString("║  No todos yet! Press 'a' to add one  ║\n")
			b.WriteString("║                                       ║\n")
		} else {
			// Show todos (up to 5)
			start := 0
			if m.cursor > 4 {
				start = m.cursor - 4
			}
			end := start + 5
			if end > len(m.items) {
				end = len(m.items)
			}

			for i := start; i < end; i++ {
				item := m.items[i]
				cursor := " "
				if i == m.cursor {
					cursor = ">"
				}

				checkbox := "☐"
				if item.completed {
					checkbox = "☑"
				}

				// Truncate text if too long
				text := item.text
				if len(text) > 28 {
					text = text[:25] + "..."
				}

				b.WriteString(fmt.Sprintf("║ %s %s %-28s ║\n", cursor, checkbox, text))
			}

			// Fill remaining lines
			for i := end - start; i < 5; i++ {
				b.WriteString("║                                       ║\n")
			}
		}

		b.WriteString("╠═══════════════════════════════════════╣\n")
		b.WriteString("║ a:Add  d:Delete  Space:Toggle  q:Quit ║\n")
		b.WriteString(fmt.Sprintf("║ j/k or ↑/↓ to navigate  (%d items)    ║\n", len(m.items)))
	}

	b.WriteString("╚═══════════════════════════════════════╝")

	return b.String()
}

func main() {
	// Create initial model with sample todos
	initialModel := TodoModel{
		items: []TodoItem{
			{text: "Try phoenix/tea framework", completed: false},
			{text: "Build awesome TUI apps", completed: false},
			{text: "Enjoy fast performance", completed: false},
		},
		cursor: 0,
	}

	// Create program
	p := tea.New(initialModel, tea.WithAltScreen[TodoModel]())

	// Run the program
	if err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
