// Package main demonstrates how to integrate Phoenix TUI with Cobra CLI
//
// This example shows the hybrid CLI+TUI pattern:
// - When flags are provided → CLI mode (scriptable, automation-friendly)
// - When no flags → TUI mode (interactive, user-friendly)
//
// This is the recommended pattern for production CLI tools.
package main

import (
	"fmt"
	"os"

	"github.com/phoenix-tui/phoenix/components"
	"github.com/phoenix-tui/phoenix/style"
	"github.com/phoenix-tui/phoenix/tea"
	"github.com/spf13/cobra"
)

// CLI flags
var (
	name    string
	email   string
	message string
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "mytool",
	Short: "A demo tool showing Cobra + Phoenix integration",
	Long: `MyTool demonstrates the hybrid CLI+TUI pattern:

• CLI Mode: Use flags for scripting and automation
  $ mytool --name "John" --email "john@example.com" --message "Hello!"

• TUI Mode: Run without flags for interactive experience
  $ mytool

This pattern gives you the best of both worlds:
- Scripts use CLI flags (reproducible, automation-friendly)
- Humans use TUI (intuitive, guided, beautiful)`,
	Run: func(cmd *cobra.Command, args []string) {
		// If any flags are provided → CLI mode
		if cmd.Flags().NFlag() > 0 {
			runCLIMode()
			return
		}

		// No flags → TUI mode
		runTUIMode()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Define CLI flags
	rootCmd.Flags().StringVarP(&name, "name", "n", "", "Your name")
	rootCmd.Flags().StringVarP(&email, "email", "e", "", "Your email")
	rootCmd.Flags().StringVarP(&message, "message", "m", "", "Your message")
}

// runCLIMode processes the command using flags (scriptable mode)
func runCLIMode() {
	fmt.Println("CLI Mode 🖥️")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if name == "" {
		fmt.Fprintln(os.Stderr, "Error: --name is required")
		os.Exit(1)
	}

	if email == "" {
		fmt.Fprintln(os.Stderr, "Error: --email is required")
		os.Exit(1)
	}

	if message == "" {
		fmt.Fprintln(os.Stderr, "Error: --message is required")
		os.Exit(1)
	}

	// Process with provided flags
	fmt.Printf("✓ Name:    %s\n", name)
	fmt.Printf("✓ Email:   %s\n", email)
	fmt.Printf("✓ Message: %s\n", message)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✓ Data processed successfully!")
}

// runTUIMode launches interactive Phoenix TUI
func runTUIMode() {
	// Create initial model
	m := newFormModel()

	// Run Phoenix TUI program
	p := tea.NewProgram(m, tea.WithAltScreen())
	if finalModel, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	} else {
		// Show results after TUI exits
		if fm, ok := finalModel.(formModel); ok && fm.submitted {
			showResults(fm)
		}
	}
}

// formModel implements tea.Model for our interactive form
type formModel struct {
	// Form fields
	nameInput    components.TextInput
	emailInput   components.TextInput
	messageInput components.TextInput

	// State
	focusIndex int
	submitted  bool
	width      int
	height     int

	// Collected data
	formData struct {
		name    string
		email   string
		message string
	}
}

func newFormModel() formModel {
	// Create form inputs with Phoenix components
	nameInput := components.NewTextInput(
		components.WithPlaceholder("Enter your name..."),
		components.WithPrompt("Name: "),
		components.WithStyle(style.New().
			Foreground(style.Color("#00FF00"))),
	)
	nameInput.Focus()

	emailInput := components.NewTextInput(
		components.WithPlaceholder("your.email@example.com"),
		components.WithPrompt("Email: "),
		components.WithStyle(style.New().
			Foreground(style.Color("#00AAFF"))),
	)

	messageInput := components.NewTextInput(
		components.WithPlaceholder("Enter your message..."),
		components.WithPrompt("Message: "),
		components.WithWidth(60),
		components.WithStyle(style.New().
			Foreground(style.Color("#FFAA00"))),
	)

	return formModel{
		nameInput:    nameInput,
		emailInput:   emailInput,
		messageInput: messageInput,
		focusIndex:   0,
	}
}

func (m formModel) Init() tea.Cmd {
	return nil
}

func (m formModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "tab", "down":
			// Next field
			m.focusIndex = (m.focusIndex + 1) % 3
			return m, m.updateFocus()

		case "shift+tab", "up":
			// Previous field
			m.focusIndex = (m.focusIndex - 1 + 3) % 3
			return m, m.updateFocus()

		case "enter":
			// Submit form
			if m.isFormValid() {
				m.formData.name = m.nameInput.Value()
				m.formData.email = m.emailInput.Value()
				m.formData.message = m.messageInput.Value()
				m.submitted = true
				return m, tea.Quit
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	// Update focused input
	switch m.focusIndex {
	case 0:
		m.nameInput, cmd = m.nameInput.Update(msg)
	case 1:
		m.emailInput, cmd = m.emailInput.Update(msg)
	case 2:
		m.messageInput, cmd = m.messageInput.Update(msg)
	}

	return m, cmd
}

func (m formModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	// Title
	titleStyle := style.New().
		Bold(true).
		Foreground(style.Color("#FF00FF")).
		MarginBottom(1)

	title := titleStyle.Render("🚀 Phoenix + Cobra Demo")

	// Instructions
	helpStyle := style.New().
		Foreground(style.Color("#888888")).
		MarginTop(1)

	help := helpStyle.Render("Tab/↓: Next • Shift+Tab/↑: Previous • Enter: Submit • Esc: Quit")

	// Form fields
	var fields string
	fields += m.nameInput.View() + "\n"
	fields += m.emailInput.View() + "\n"
	fields += m.messageInput.View() + "\n"

	// Validation message
	var validation string
	if !m.isFormValid() {
		validationStyle := style.New().
			Foreground(style.Color("#FF0000")).
			MarginTop(1)
		validation = validationStyle.Render("⚠ Please fill in all fields")
	} else {
		validationStyle := style.New().
			Foreground(style.Color("#00FF00")).
			MarginTop(1)
		validation = validationStyle.Render("✓ Press Enter to submit")
	}

	// Combine all parts
	container := style.New().
		Border(style.BorderRounded).
		BorderForeground(style.Color("#00FFFF")).
		Padding(1, 2).
		Width(m.width - 4).
		Align(style.AlignCenter)

	content := fmt.Sprintf("%s\n\n%s\n%s\n%s", title, fields, validation, help)

	return container.Render(content)
}

func (m formModel) updateFocus() tea.Cmd {
	// Update focus state for all inputs
	if m.focusIndex == 0 {
		m.nameInput.Focus()
		m.emailInput.Blur()
		m.messageInput.Blur()
	} else if m.focusIndex == 1 {
		m.nameInput.Blur()
		m.emailInput.Focus()
		m.messageInput.Blur()
	} else {
		m.nameInput.Blur()
		m.emailInput.Blur()
		m.messageInput.Focus()
	}
	return nil
}

func (m formModel) isFormValid() bool {
	return m.nameInput.Value() != "" &&
		m.emailInput.Value() != "" &&
		m.messageInput.Value() != ""
}

func showResults(m formModel) {
	resultStyle := style.New().
		Bold(true).
		Foreground(style.Color("#00FF00")).
		Border(style.BorderDouble).
		Padding(1, 2)

	result := fmt.Sprintf(`✓ Form Submitted Successfully!

Name:    %s
Email:   %s
Message: %s

Thank you for using Phoenix + Cobra! 🎉`,
		m.formData.name,
		m.formData.email,
		m.formData.message)

	fmt.Println(resultStyle.Render(result))
}

func main() {
	Execute()
}
