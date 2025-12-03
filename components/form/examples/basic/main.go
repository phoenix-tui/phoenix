// Package main demonstrates basic form usage with simple text fields.
package main

import (
	"fmt"
	"os"

	"github.com/phoenix-tui/phoenix/components/form"
	"github.com/phoenix-tui/phoenix/components/form/internal/domain/value"
	tea "github.com/phoenix-tui/phoenix/tea"
)

// simpleInput is a minimal input model for demonstration.
type simpleInput struct {
	label string
	val   string
	focus bool
}

func newSimpleInput(label, initial string) *simpleInput {
	return &simpleInput{
		label: label,
		val:   initial,
		focus: false,
	}
}

func (s *simpleInput) Value() string {
	return s.val
}

func (s *simpleInput) Init() tea.Cmd {
	return nil
}

func (s *simpleInput) Update(msg tea.Msg) (*simpleInput, tea.Cmd) {
	if !s.focus {
		return s, nil
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.Type {
		case tea.KeyRune:
			s.val += string(keyMsg.Rune)
		case tea.KeyBackspace:
			if s.val != "" {
				s.val = s.val[:len(s.val)-1]
			}
		}
	}

	return s, nil
}

func (s *simpleInput) View() string {
	cursor := " "
	if s.focus {
		cursor = "█"
	}
	return "[" + s.val + cursor + "]"
}

func (s *simpleInput) SetFocus(focus bool) *simpleInput {
	s.focus = focus
	return s
}

// model is the application model.
type model struct {
	formComponent *form.Form
	nameInput     *simpleInput
	emailInput    *simpleInput
	submitted     bool
	result        string
}

func initialModel() model {
	nameInput := newSimpleInput("Name", "")
	emailInput := newSimpleInput("Email", "")

	// Create form with validation
	f := form.New("User Registration").
		Field("name", "Name", nameInput,
			value.Required(),
			value.MinLength(2)).
		Field("email", "Email", emailInput,
			value.Required(),
			value.Email())

	return model{
		formComponent: f,
		nameInput:     nameInput,
		emailInput:    emailInput,
		submitted:     false,
	}
}

func (m model) Init() tea.Cmd {
	return m.formComponent.Init()
}

func (m model) Update(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case form.SubmitMsg:
		// Form submitted successfully (all validations passed)
		m.submitted = true
		m.result = fmt.Sprintf("Submitted!\nName: %s\nEmail: %s",
			m.nameInput.Value(),
			m.emailInput.Value())
		return m, tea.Quit()

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit()
		}
	}

	// Update form
	newForm, cmd := m.formComponent.Update(msg)
	m.formComponent = newForm

	return m, cmd
}

func (m model) View() string {
	if m.submitted {
		return m.result + "\n"
	}

	return m.formComponent.View() + "\n\n(Ctrl+C to quit)\n"
}

func main() {
	p := tea.New(initialModel())
	if err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
