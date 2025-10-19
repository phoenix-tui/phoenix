package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/phoenix-tui/phoenix/components/input/api"
	tea "github.com/phoenix-tui/phoenix/tea/api"
)

// validatedModel demonstrates input with various validators.
type validatedModel struct {
	emailInput input.Input
	phoneInput input.Input
	urlInput   input.Input
	focused    int // 0=email, 1=phone, 2=url
}

func (m validatedModel) Init() tea.Cmd {
	return nil
}

func (m validatedModel) Update(msg tea.Msg) (validatedModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit()

		case "enter":
			// Validate all and display results
			if m.emailInput.IsValid() && m.phoneInput.IsValid() && m.urlInput.IsValid() {
				fmt.Printf("\n✓ All inputs valid!\n")
				fmt.Printf("Email: %s\n", m.emailInput.Value())
				fmt.Printf("Phone: %s\n", m.phoneInput.Value())
				fmt.Printf("URL: %s\n", m.urlInput.Value())
			} else {
				fmt.Printf("\n✗ Some inputs are invalid\n")
			}
			return m, tea.Quit()

		case "tab", "down":
			// Move to next field
			m.focused = (m.focused + 1) % 3
			m = m.updateFocus()
			return m, nil

		case "shift+tab", "up":
			// Move to previous field
			m.focused = (m.focused + 2) % 3
			m = m.updateFocus()
			return m, nil
		}

	case tea.WindowSizeMsg:
		width := msg.Width - 20
		m.emailInput = m.emailInput.Width(width)
		m.phoneInput = m.phoneInput.Width(width)
		m.urlInput = m.urlInput.Width(width)
		return m, nil
	}

	// Forward to focused input
	var cmd tea.Cmd
	var updated input.Input

	switch m.focused {
	case 0:
		updated, cmd = m.emailInput.Update(msg)
		m.emailInput = updated
	case 1:
		updated, cmd = m.phoneInput.Update(msg)
		m.phoneInput = updated
	case 2:
		updated, cmd = m.urlInput.Update(msg)
		m.urlInput = updated
	}

	return m, cmd
}

func (m validatedModel) updateFocus() validatedModel {
	m.emailInput = m.emailInput.Focused(m.focused == 0)
	m.phoneInput = m.phoneInput.Focused(m.focused == 1)
	m.urlInput = m.urlInput.Focused(m.focused == 2)
	return m
}

func (m validatedModel) View() string {
	var b strings.Builder

	b.WriteString("Validated Input Example\n\n")

	// Email field
	b.WriteString("Email: ")
	b.WriteString(m.emailInput.View())
	if m.emailInput.Value() != "" {
		if m.emailInput.IsValid() {
			b.WriteString(" ✓")
		} else {
			b.WriteString(" ✗ (must be valid email)")
		}
	}
	b.WriteString("\n\n")

	// Phone field
	b.WriteString("Phone: ")
	b.WriteString(m.phoneInput.View())
	if m.phoneInput.Value() != "" {
		if m.phoneInput.IsValid() {
			b.WriteString(" ✓")
		} else {
			b.WriteString(" ✗ (format: XXX-XXX-XXXX)")
		}
	}
	b.WriteString("\n\n")

	// URL field
	b.WriteString("URL:   ")
	b.WriteString(m.urlInput.View())
	if m.urlInput.Value() != "" {
		if m.urlInput.IsValid() {
			b.WriteString(" ✓")
		} else {
			b.WriteString(" ✗ (must start with http:// or https://)")
		}
	}
	b.WriteString("\n\n")

	b.WriteString("Tab: Next field | Enter: Submit | Esc: Quit")

	return b.String()
}

func main() {
	// Email validator
	emailValidator := func(s string) error {
		if s == "" {
			return fmt.Errorf("email required")
		}
		// Simple email regex
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(s) {
			return fmt.Errorf("invalid email format")
		}
		return nil
	}

	// Phone validator
	phoneValidator := func(s string) error {
		if s == "" {
			return fmt.Errorf("phone required")
		}
		// Format: XXX-XXX-XXXX
		phoneRegex := regexp.MustCompile(`^\d{3}-\d{3}-\d{4}$`)
		if !phoneRegex.MatchString(s) {
			return fmt.Errorf("invalid phone format")
		}
		return nil
	}

	// URL validator
	urlValidator := func(s string) error {
		if s == "" {
			return fmt.Errorf("url required")
		}
		if !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "https://") {
			return fmt.Errorf("url must start with http:// or https://")
		}
		return nil
	}

	// Create inputs
	model := validatedModel{
		emailInput: input.New(40).
			Placeholder("user@example.com").
			Validator(emailValidator).
			Focused(true),
		phoneInput: input.New(40).
			Placeholder("555-123-4567").
			Validator(phoneValidator).
			Focused(false),
		urlInput: input.New(40).
			Placeholder("https://example.com").
			Validator(urlValidator).
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
