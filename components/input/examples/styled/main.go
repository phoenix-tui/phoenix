package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/phoenix-tui/phoenix/components/input/api"
	tea "github.com/phoenix-tui/phoenix/tea/api"
)

// styledModel demonstrates styled inputs.
// Note: Full styling will be available when phoenix/style is integrated.
// This example shows the API structure.
type styledModel struct {
	input1  input.Input
	input2  input.Input
	input3  input.Input
	focused int
}

func (m styledModel) Init() tea.Cmd {
	return nil
}

func (m styledModel) Update(msg tea.Msg) (styledModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit()

		case "tab", "down":
			m.focused = (m.focused + 1) % 3
			m = m.updateFocus()
			return m, nil

		case "shift+tab", "up":
			m.focused = (m.focused + 2) % 3
			m = m.updateFocus()
			return m, nil
		}

	case tea.WindowSizeMsg:
		width := msg.Width - 10
		m.input1 = m.input1.Width(width)
		m.input2 = m.input2.Width(width)
		m.input3 = m.input3.Width(width)
		return m, nil
	}

	// Forward to focused input
	var cmd tea.Cmd
	var updated input.Input

	switch m.focused {
	case 0:
		updated, cmd = m.input1.Update(msg)
		m.input1 = updated
	case 1:
		updated, cmd = m.input2.Update(msg)
		m.input2 = updated
	case 2:
		updated, cmd = m.input3.Update(msg)
		m.input3 = updated
	}

	return m, cmd
}

func (m styledModel) updateFocus() styledModel {
	m.input1 = m.input1.Focused(m.focused == 0)
	m.input2 = m.input2.Focused(m.focused == 1)
	m.input3 = m.input3.Focused(m.focused == 2)
	return m
}

func (m styledModel) View() string {
	var b strings.Builder

	b.WriteString("Styled Input Example\n")
	b.WriteString("(Full styling coming with phoenix/style integration)\n\n")

	// Style 1: Simple with border
	b.WriteString(m.renderWithBorder("Normal Input", m.input1, m.focused == 0))
	b.WriteString("\n")

	// Style 2: Emphasized
	b.WriteString(m.renderWithBorder("Emphasized Input", m.input2, m.focused == 1))
	b.WriteString("\n")

	// Style 3: Compact
	b.WriteString(m.renderCompact("Compact", m.input3, m.focused == 2))
	b.WriteString("\n")

	b.WriteString("\nTab: Next field | Esc: Quit")

	return b.String()
}

func (m styledModel) renderWithBorder(label string, input input.Input, isFocused bool) string {
	var b strings.Builder

	// Label
	b.WriteString(label)
	b.WriteString("\n")

	// Top border
	if isFocused {
		b.WriteString("╔═══════════════════════════════════════╗\n")
	} else {
		b.WriteString("┌───────────────────────────────────────┐\n")
	}

	// Content
	if isFocused {
		b.WriteString("║ ")
	} else {
		b.WriteString("│ ")
	}
	b.WriteString(input.View())
	if isFocused {
		b.WriteString(" ║")
	} else {
		b.WriteString(" │")
	}
	b.WriteString("\n")

	// Bottom border
	if isFocused {
		b.WriteString("╚═══════════════════════════════════════╝")
	} else {
		b.WriteString("└───────────────────────────────────────┘")
	}

	return b.String()
}

func (m styledModel) renderCompact(label string, input input.Input, isFocused bool) string {
	prefix := "  "
	if isFocused {
		prefix = "▶ "
	}
	return fmt.Sprintf("%s%s: %s", prefix, label, input.View())
}

func main() {
	// Create styled inputs
	model := styledModel{
		input1: input.New(40).
			Placeholder("Type here...").
			Focused(true),
		input2: input.New(40).
			Placeholder("Important field...").
			Focused(false),
		input3: input.New(40).
			Placeholder("Quick input...").
			Focused(false),
		focused: 0,
	}

	// Run program
	p := tea.New(model)
	if err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
