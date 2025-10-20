// Package main demonstrates shell prompt textarea integration.
// This example shows cursor control and movement validation for shell-like behavior.
package main

import (
	"fmt"
	"log"

	"github.com/phoenix-tui/phoenix/components/input/textarea/api"
	tea "github.com/phoenix-tui/phoenix/tea/api"
)

// ShellModel demonstrates shell-like boundary protection using Phoenix TextArea.
type ShellModel struct {
	textarea    api.TextArea
	promptLen   int
	borderHits  int
	lastMessage string
}

// NewShellModel creates a new shell model with prompt protection.
func NewShellModel() ShellModel {
	prompt := "> "
	promptLen := len(prompt)

	ta := api.New().
		SetValue(prompt).
		SetCursorPosition(0, promptLen).
		Size(80, 1).
		ShowCursor(true).
		// Boundary protection - prevent cursor from moving before prompt.
		OnMovement(func(_, to api.CursorPos) bool {
			// Block movement before prompt.
			if to.Row == 0 && to.Col < promptLen {
				return false
			}
			return true
		}).
		// Observer - track when cursor moves.
		OnCursorMoved(func(from, to api.CursorPos) {
			fmt.Printf("Cursor moved from (%d,%d) to (%d,%d)\n", from.Row, from.Col, to.Row, to.Col)
		}).
		// Feedback - notify user when hitting boundary.
		OnBoundaryHit(func(attemptedPos api.CursorPos, reason string) {
			fmt.Printf("Boundary hit! Attempted to move to (%d,%d): %s\n",
				attemptedPos.Row, attemptedPos.Col, reason)
		})

	return ShellModel{
		textarea:    ta,
		promptLen:   promptLen,
		borderHits:  0,
		lastMessage: "Shell prompt demo - try moving cursor before '> ' prompt!",
	}
}

// Init initializes the shell model.
func (m ShellModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and returns updated model.
func (m ShellModel) Update(msg tea.Msg) (ShellModel, tea.Cmd) {
	//nolint:gocritic // switch preferred for extensibility in examples
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Check for quit.
		if msg.Type == tea.KeyCtrlC || (msg.Type == tea.KeyRune && msg.Rune == 'q') {
			return m, tea.Quit()
		}

		// Handle Enter - execute command.
		if msg.Type == tea.KeyEnter {
			command := m.getCommand()
			m.lastMessage = fmt.Sprintf("Executing: %s", command)

			// Reset prompt.
			prompt := "> "
			m.textarea = m.textarea.SetValue(prompt).SetCursorPosition(0, len(prompt))
			return m, nil
		}

		// Update textarea.
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View renders the shell interface.
func (m ShellModel) View() string {
	return fmt.Sprintf(
		"%s\n\n"+
			"Last message: %s\n"+
			"Command: %s\n\n"+
			"Keys:\n"+
			"  ← → Home End  : Navigate (try moving before prompt!)\n"+
			"  Enter         : Execute command\n"+
			"  q / Ctrl+C    : Quit\n",
		m.textarea.View(),
		m.lastMessage,
		m.getCommand(),
	)
}

// getCommand extracts command (text after prompt).
func (m ShellModel) getCommand() string {
	value := m.textarea.Value()
	if len(value) <= m.promptLen {
		return ""
	}
	return value[m.promptLen:]
}

func main() {
	// Create program.
	p := tea.New(NewShellModel())

	// Run.
	if err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
