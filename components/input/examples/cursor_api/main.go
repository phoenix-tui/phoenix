// Package main demonstrates the cursor control API.
// This example shows Phoenix TextInput's public cursor API - a key differentiator.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/phoenix-tui/phoenix/components/input/api"
	tea "github.com/phoenix-tui/phoenix/tea/api"
)

// cursorAPIModel demonstrates the public cursor API.
// This is a KEY DIFFERENTIATOR of Phoenix TextInput.
type cursorAPIModel struct {
	input   input.Input
	info    string
	lastKey string
}

func (m cursorAPIModel) Init() tea.Cmd {
	return nil
}

func (m cursorAPIModel) Update(msg tea.Msg) (cursorAPIModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.lastKey = msg.String()

		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit()

		case "ctrl+s":
			// Demonstrate SetContent with cursor position.
			m.input = m.input.SetContent("Saved content!", 6)
			m.info = "Content set atomically with cursor at position 6"
			return m, nil

		case "ctrl+i":
			// Show cursor info.
			pos := m.input.CursorPosition()
			before, at, after := m.input.ContentParts()
			m.info = fmt.Sprintf(
				"Cursor at position %d | Before: %q | At: %q | After: %q",
				pos, before, at, after,
			)
			return m, nil

		case "ctrl+j":
			// Jump to middle.
			content := m.input.Value()
			middle := len(content) / 2
			m.input = m.input.SetContent(content, middle)
			m.info = fmt.Sprintf("Jumped to middle (position %d)", middle)
			return m, nil

		case "ctrl+p":
			// Programmatic manipulation using cursor API.
			before, at, after := m.input.ContentParts()
			// Insert "[CURSOR]" marker.
			newContent := before + "[" + at + "]" + after
			cursorPos := m.input.CursorPosition()
			m.input = m.input.SetContent(newContent, cursorPos)
			m.info = "Inserted cursor marker using ContentParts()"
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.input = m.input.Width(msg.Width - 4)
		return m, nil
	}

	// Forward to input.
	updated, cmd := m.input.Update(msg)
	m.input = updated

	// Update info.
	if m.lastKey != "" && m.lastKey != "ctrl+i" {
		pos := m.input.CursorPosition()
		m.info = fmt.Sprintf("Cursor at position %d (press Ctrl-I for details)", pos)
	}

	return m, cmd
}

func (m cursorAPIModel) View() string {
	var b strings.Builder

	b.WriteString("Cursor API Example - KEY DIFFERENTIATOR!\n\n")

	b.WriteString("Input: ")
	b.WriteString(m.input.View())
	b.WriteString("\n\n")

	b.WriteString("═══ Public Cursor API ═══\n")
	b.WriteString(m.info)
	b.WriteString("\n\n")

	b.WriteString("Commands:\n")
	b.WriteString("  Ctrl-I: Show cursor info (position, parts)\n")
	b.WriteString("  Ctrl-J: Jump to middle of content\n")
	b.WriteString("  Ctrl-S: Set content atomically\n")
	b.WriteString("  Ctrl-P: Insert cursor marker (demo ContentParts)\n")
	b.WriteString("  Esc: Quit\n\n")

	b.WriteString("Why This Matters:\n")
	b.WriteString("━━━━━━━━━━━━━━━━━━\n")
	b.WriteString("• CursorPosition() - Get exact cursor offset\n")
	b.WriteString("• ContentParts() - Split around cursor for custom rendering\n")
	b.WriteString("• SetContent() - Atomic content + cursor update (race-free)\n")
	b.WriteString("\n")
	b.WriteString("This API enables:\n")
	b.WriteString("  ✓ Custom cursor rendering (gosh will use for shell prompt)\n")
	b.WriteString("  ✓ Syntax highlighting with cursor position\n")
	b.WriteString("  ✓ Autocomplete at exact cursor location\n")
	b.WriteString("  ✓ Multi-line content splitting\n")
	b.WriteString("  ✓ History navigation without race conditions\n")

	return b.String()
}

func main() {
	// Create input with initial content.
	model := cursorAPIModel{
		input: input.New(60).
			Content("Try moving the cursor and pressing Ctrl-I").
			Focused(true),
		info: "Press Ctrl-I to see cursor information",
	}

	// Move cursor to middle initially.
	content := model.input.Value()
	middle := len(content) / 2
	model.input = model.input.SetContent(content, middle)

	// Run program.
	p := tea.New(model)
	if err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
