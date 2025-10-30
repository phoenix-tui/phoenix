package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/phoenix-tui/phoenix/clipboard"
	"github.com/phoenix-tui/phoenix/clipboard/internal/domain/service"
	"github.com/phoenix-tui/phoenix/clipboard/internal/domain/value"
	tea "github.com/phoenix-tui/phoenix/tea"
)

// model represents the application state
type model struct {
	clipboard     *clipboard.Clipboard
	codec         *service.RichTextCodec
	mode          string // "menu", "write-html", "write-rtf", "read-html", "read-rtf", "convert"
	input         string
	output        string
	styles        value.TextStyles
	cursorPos     int
	lastError     string
	convertSource string // "html" or "rtf"
	quitting      bool
}

func initialModel() model {
	clip, err := clipboard.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize clipboard: %v\n", err)
		os.Exit(1)
	}

	return model{
		clipboard: clip,
		codec:     service.NewRichTextCodec(),
		mode:      "menu",
		styles:    value.NewTextStyles(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.mode {
		case "menu":
			return m.updateMenu(msg)
		case "write-html", "write-rtf":
			return m.updateWrite(msg)
		case "read-html", "read-rtf", "read-plain":
			return m.updateRead(msg)
		case "convert":
			return m.updateConvert(msg)
		}
	}
	return m, nil
}

func (m model) updateMenu(msg tea.KeyMsg) (model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		m.quitting = true
		return m, tea.Quit()
	case "1":
		m.mode = "write-html"
		m.input = ""
		m.lastError = ""
		return m, nil
	case "2":
		m.mode = "write-rtf"
		m.input = ""
		m.lastError = ""
		return m, nil
	case "3":
		m.mode = "read-html"
		m.output = ""
		m.lastError = ""
		// Read HTML from clipboard
		html, err := m.clipboard.ReadHTML()
		if err != nil {
			m.lastError = fmt.Sprintf("Error: %v", err)
		} else {
			m.output = html
			// Try to decode and show styles
			text, styles, err := m.codec.DecodeHTML(html)
			if err == nil {
				m.output = fmt.Sprintf("HTML:\n%s\n\nPlain Text: %s\nStyles: %s", html, text, styles.String())
			}
		}
		return m, nil
	case "4":
		m.mode = "read-rtf"
		m.output = ""
		m.lastError = ""
		// Read RTF from clipboard
		rtf, err := m.clipboard.ReadRTF()
		if err != nil {
			m.lastError = fmt.Sprintf("Error: %v", err)
		} else {
			m.output = rtf
			// Try to decode and show styles
			text, styles, err := m.codec.DecodeRTF(rtf)
			if err == nil {
				m.output = fmt.Sprintf("RTF:\n%s\n\nPlain Text: %s\nStyles: %s", rtf, text, styles.String())
			}
		}
		return m, nil
	case "5":
		m.mode = "read-plain"
		m.output = ""
		m.lastError = ""
		// Read HTML and strip tags
		plain, err := m.clipboard.ReadHTMLAsPlainText()
		if err != nil {
			m.lastError = fmt.Sprintf("Error: %v", err)
		} else {
			m.output = fmt.Sprintf("Plain text (HTML stripped):\n%s", plain)
		}
		return m, nil
	case "6":
		m.mode = "convert"
		m.convertSource = "html"
		m.input = ""
		m.lastError = ""
		return m, nil
	case "7":
		m.mode = "convert"
		m.convertSource = "rtf"
		m.input = ""
		m.lastError = ""
		return m, nil
	}
	return m, nil
}

func (m model) updateWrite(msg tea.KeyMsg) (model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		m.quitting = true
		return m, tea.Quit()
	case "esc":
		m.mode = "menu"
		m.input = ""
		m.lastError = ""
		return m, nil
	case "enter":
		// Write to clipboard
		if m.mode == "write-html" {
			// Encode with current styles
			html, err := m.codec.EncodeHTML(m.input, m.styles)
			if err != nil {
				m.lastError = fmt.Sprintf("Error encoding HTML: %v", err)
			} else {
				err = m.clipboard.WriteHTML(html)
				if err != nil {
					m.lastError = fmt.Sprintf("Error writing to clipboard: %v", err)
				} else {
					m.lastError = "Success! HTML written to clipboard"
				}
			}
		} else if m.mode == "write-rtf" {
			// Encode with current styles
			rtf, err := m.codec.EncodeRTF(m.input, m.styles)
			if err != nil {
				m.lastError = fmt.Sprintf("Error encoding RTF: %v", err)
			} else {
				err = m.clipboard.WriteRTF(rtf)
				if err != nil {
					m.lastError = fmt.Sprintf("Error writing to clipboard: %v", err)
				} else {
					m.lastError = "Success! RTF written to clipboard"
				}
			}
		}
		return m, nil
	case "ctrl+b":
		m.styles = m.styles.WithBold(!m.styles.Bold)
		return m, nil
	case "ctrl+i":
		m.styles = m.styles.WithItalic(!m.styles.Italic)
		return m, nil
	case "ctrl+u":
		m.styles = m.styles.WithUnderline(!m.styles.Underline)
		return m, nil
	case "ctrl+r":
		// Set red color
		m.styles, _ = m.styles.WithColor("#FF0000")
		return m, nil
	case "ctrl+g":
		// Set green color
		m.styles, _ = m.styles.WithColor("#00FF00")
		return m, nil
	case "ctrl+l":
		// Clear color
		m.styles, _ = m.styles.WithColor("")
		return m, nil
	case "backspace":
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}
		return m, nil
	default:
		// Add character to input
		if len(msg.String()) == 1 {
			m.input += msg.String()
		}
		return m, nil
	}
}

func (m model) updateRead(msg tea.KeyMsg) (model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		m.quitting = true
		return m, tea.Quit()
	case "esc", "enter", "q":
		m.mode = "menu"
		m.output = ""
		m.lastError = ""
		return m, nil
	}
	return m, nil
}

func (m model) updateConvert(msg tea.KeyMsg) (model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		m.quitting = true
		return m, tea.Quit()
	case "esc":
		m.mode = "menu"
		m.input = ""
		m.output = ""
		m.lastError = ""
		return m, nil
	case "enter":
		// Convert
		if m.convertSource == "html" {
			rtf, err := m.clipboard.ConvertHTMLToRTF(m.input)
			if err != nil {
				m.lastError = fmt.Sprintf("Error: %v", err)
			} else {
				m.output = fmt.Sprintf("Converted to RTF:\n%s", rtf)
			}
		} else {
			html, err := m.clipboard.ConvertRTFToHTML(m.input)
			if err != nil {
				m.lastError = fmt.Sprintf("Error: %v", err)
			} else {
				m.output = fmt.Sprintf("Converted to HTML:\n%s", html)
			}
		}
		return m, nil
	case "backspace":
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}
		return m, nil
	default:
		// Add character to input
		if len(msg.String()) == 1 {
			m.input += msg.String()
		}
		return m, nil
	}
}

func (m model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	var s strings.Builder

	s.WriteString("╔═══════════════════════════════════════════════════════════════╗\n")
	s.WriteString("║         Phoenix Clipboard - Rich Text Example                ║\n")
	s.WriteString("╚═══════════════════════════════════════════════════════════════╝\n\n")

	switch m.mode {
	case "menu":
		s.WriteString("Choose an option:\n\n")
		s.WriteString("  1. Write HTML to clipboard\n")
		s.WriteString("  2. Write RTF to clipboard\n")
		s.WriteString("  3. Read HTML from clipboard\n")
		s.WriteString("  4. Read RTF from clipboard\n")
		s.WriteString("  5. Read as plain text (strip HTML)\n")
		s.WriteString("  6. Convert HTML to RTF\n")
		s.WriteString("  7. Convert RTF to HTML\n")
		s.WriteString("  q. Quit\n")

	case "write-html", "write-rtf":
		format := "HTML"
		if m.mode == "write-rtf" {
			format = "RTF"
		}
		s.WriteString(fmt.Sprintf("Write %s to Clipboard\n", format))
		s.WriteString("─────────────────────────────────────────────────────────────────\n\n")
		s.WriteString("Current styles:\n")
		s.WriteString(fmt.Sprintf("  Bold: %v (Ctrl+B to toggle)\n", m.styles.Bold))
		s.WriteString(fmt.Sprintf("  Italic: %v (Ctrl+I to toggle)\n", m.styles.Italic))
		s.WriteString(fmt.Sprintf("  Underline: %v (Ctrl+U to toggle)\n", m.styles.Underline))
		color := "none"
		if m.styles.Color != "" {
			color = m.styles.Color
		}
		s.WriteString(fmt.Sprintf("  Color: %s (Ctrl+R=red, Ctrl+G=green, Ctrl+L=clear)\n", color))
		s.WriteString("\nType your text:\n")
		s.WriteString(fmt.Sprintf("> %s\n", m.input))
		s.WriteString("\nPress Enter to write to clipboard, Esc to cancel\n")

		if m.lastError != "" {
			s.WriteString(fmt.Sprintf("\n%s\n", m.lastError))
		}

	case "read-html", "read-rtf", "read-plain":
		s.WriteString("Clipboard Content\n")
		s.WriteString("─────────────────────────────────────────────────────────────────\n\n")
		if m.lastError != "" {
			s.WriteString(fmt.Sprintf("%s\n", m.lastError))
		} else {
			s.WriteString(fmt.Sprintf("%s\n", m.output))
		}
		s.WriteString("\nPress any key to return to menu\n")

	case "convert":
		source := "HTML"
		target := "RTF"
		if m.convertSource == "rtf" {
			source = "RTF"
			target = "HTML"
		}
		s.WriteString(fmt.Sprintf("Convert %s to %s\n", source, target))
		s.WriteString("─────────────────────────────────────────────────────────────────\n\n")
		s.WriteString(fmt.Sprintf("Enter %s content:\n", source))
		s.WriteString(fmt.Sprintf("> %s\n", m.input))
		s.WriteString("\nPress Enter to convert, Esc to cancel\n")

		if m.output != "" {
			s.WriteString(fmt.Sprintf("\n%s\n", m.output))
		}
		if m.lastError != "" {
			s.WriteString(fmt.Sprintf("\n%s\n", m.lastError))
		}
	}

	return s.String()
}

func main() {
	p := tea.New(initialModel(), tea.WithAltScreen[model]())
	if err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
