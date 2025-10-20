// Package modal provides a universal modal/overlay component for Phoenix TUI Framework.
//
// The Modal component displays dialog boxes with support for:
//   - Overlay rendering (centered or custom positioning)
//   - Focus trap (modal captures all input when visible)
//   - Keyboard dismiss (Esc to close)
//   - Custom content (any string content)
//   - Button support (optional action buttons)
//   - Background dimming (improves UX)
//
// This is a UNIVERSAL component - it works for any application (confirmation dialogs,.
// help screens, settings panels, alerts, etc.). It does NOT include application-specific.
// features.
package modal

import (
	"fmt"
	"strings"

	"github.com/phoenix-tui/phoenix/components/modal/domain/model"
	"github.com/phoenix-tui/phoenix/components/modal/domain/service"
	"github.com/phoenix-tui/phoenix/components/modal/domain/value"
	"github.com/phoenix-tui/phoenix/components/modal/infrastructure"
	tea "github.com/phoenix-tui/phoenix/tea/api"
)

// Button defines a modal action button.
type Button struct {
	Label  string // Button text (e.g., "Yes", "No", "OK")
	Key    string // Keyboard shortcut (e.g., "y", "n")
	Action string // Action identifier (e.g., "confirm", "cancel")
}

// Modal is the public API for the modal component.
// It implements tea.Model for integration with Phoenix Tea event loop.
type Modal struct {
	domain         *model.Modal
	layoutService  *service.LayoutService
	keyBindings    infrastructure.KeyBindings
	terminalWidth  int // Terminal size for rendering
	terminalHeight int
}

// New creates a new modal with the given content.
// Default: centered, 40x10, not visible, no dimming.
func New(content string) *Modal {
	return &Modal{
		domain:         model.NewModal(content),
		layoutService:  service.NewLayoutService(),
		keyBindings:    infrastructure.DefaultKeyBindings(),
		terminalWidth:  80, // Default terminal size
		terminalHeight: 24,
	}
}

// NewWithTitle creates a new modal with title and content.
func NewWithTitle(title, content string) *Modal {
	return &Modal{
		domain:         model.NewModalWithTitle(title, content),
		layoutService:  service.NewLayoutService(),
		keyBindings:    infrastructure.DefaultKeyBindings(),
		terminalWidth:  80,
		terminalHeight: 24,
	}
}

// Size returns a new modal with the specified size.
func (m *Modal) Size(width, height int) *Modal {
	return &Modal{
		domain:         m.domain.WithSize(width, height),
		layoutService:  m.layoutService,
		keyBindings:    m.keyBindings,
		terminalWidth:  m.terminalWidth,
		terminalHeight: m.terminalHeight,
	}
}

// Position returns a new modal with custom positioning.
func (m *Modal) Position(x, y int) *Modal {
	return &Modal{
		domain:         m.domain.WithPosition(value.NewPositionCustom(x, y)),
		layoutService:  m.layoutService,
		keyBindings:    m.keyBindings,
		terminalWidth:  m.terminalWidth,
		terminalHeight: m.terminalHeight,
	}
}

// Centered returns a new modal with centered positioning (default).
func (m *Modal) Centered() *Modal {
	return &Modal{
		domain:         m.domain.WithPosition(value.NewPositionCenter()),
		layoutService:  m.layoutService,
		keyBindings:    m.keyBindings,
		terminalWidth:  m.terminalWidth,
		terminalHeight: m.terminalHeight,
	}
}

// Buttons returns a new modal with the specified buttons.
func (m *Modal) Buttons(buttons []Button) *Modal {
	domainButtons := make([]*model.Button, len(buttons))
	for i, btn := range buttons {
		domainButtons[i] = model.NewButton(btn.Label, btn.Key, btn.Action)
	}

	return &Modal{
		domain:         m.domain.WithButtons(domainButtons),
		layoutService:  m.layoutService,
		keyBindings:    m.keyBindings,
		terminalWidth:  m.terminalWidth,
		terminalHeight: m.terminalHeight,
	}
}

// DimBackground returns a new modal with background dimming enabled/disabled.
func (m *Modal) DimBackground(dim bool) *Modal {
	return &Modal{
		domain:         m.domain.WithDimBackground(dim),
		layoutService:  m.layoutService,
		keyBindings:    m.keyBindings,
		terminalWidth:  m.terminalWidth,
		terminalHeight: m.terminalHeight,
	}
}

// Show returns a new modal that is visible.
func (m *Modal) Show() *Modal {
	return &Modal{
		domain:         m.domain.WithVisible(true),
		layoutService:  m.layoutService,
		keyBindings:    m.keyBindings,
		terminalWidth:  m.terminalWidth,
		terminalHeight: m.terminalHeight,
	}
}

// Hide returns a new modal that is hidden.
func (m *Modal) Hide() *Modal {
	return &Modal{
		domain:         m.domain.WithVisible(false),
		layoutService:  m.layoutService,
		keyBindings:    m.keyBindings,
		terminalWidth:  m.terminalWidth,
		terminalHeight: m.terminalHeight,
	}
}

// KeyBindings returns a new modal with custom key bindings.
func (m *Modal) KeyBindings(kb infrastructure.KeyBindings) *Modal {
	return &Modal{
		domain:         m.domain,
		layoutService:  m.layoutService,
		keyBindings:    kb,
		terminalWidth:  m.terminalWidth,
		terminalHeight: m.terminalHeight,
	}
}

// Init implements tea.Model.
func (m *Modal) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m *Modal) Update(msg tea.Msg) (*Modal, tea.Cmd) {
	// If modal is not visible, don't process any input.
	if !m.domain.IsVisible() {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Update terminal size for rendering.
		return &Modal{
			domain:         m.domain,
			layoutService:  m.layoutService,
			keyBindings:    m.keyBindings,
			terminalWidth:  msg.Width,
			terminalHeight: msg.Height,
		}, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}

	return m, nil
}

// handleKeyPress processes keyboard input.
func (m *Modal) handleKeyPress(msg tea.KeyMsg) (*Modal, tea.Cmd) {
	kb := m.keyBindings

	// Close modal (Esc by default)
	if kb.IsClose(msg) {
		return m.Hide(), nil
	}

	// Button navigation (Tab, Arrow keys)
	if kb.IsNextButton(msg) {
		return &Modal{
			domain:         m.domain.FocusNextButton(),
			layoutService:  m.layoutService,
			keyBindings:    m.keyBindings,
			terminalWidth:  m.terminalWidth,
			terminalHeight: m.terminalHeight,
		}, nil
	}

	if kb.IsPreviousButton(msg) {
		return &Modal{
			domain:         m.domain.FocusPreviousButton(),
			layoutService:  m.layoutService,
			keyBindings:    m.keyBindings,
			terminalWidth:  m.terminalWidth,
			terminalHeight: m.terminalHeight,
		}, nil
	}

	// Activate focused button (Enter by default)
	if kb.IsActivateButton(msg) {
		focusedBtn := m.domain.FocusedButton()
		if focusedBtn != nil {
			// Send ButtonPressedMsg.
			return m, func() tea.Msg {
				return ButtonPressedMsg{Action: focusedBtn.Action()}
			}
		}
	}

	// Check button shortcuts (e.g., "y" for Yes, "n" for No)
	keyStr := msg.String()
	for _, btn := range m.domain.Buttons() {
		if strings.EqualFold(keyStr, btn.Key()) {
			return m, func() tea.Msg {
				return ButtonPressedMsg{Action: btn.Action()}
			}
		}
	}

	return m, nil
}

// View implements tea.Model.
func (m *Modal) View() string {
	// If not visible, return empty string.
	if !m.domain.IsVisible() {
		return ""
	}

	var b strings.Builder

	// Calculate modal position.
	modalWidth := m.domain.Size().Width()
	modalHeight := m.domain.Size().Height()
	x, y := m.layoutService.CalculatePosition(
		m.domain.Position(),
		m.terminalWidth,
		m.terminalHeight,
		modalWidth,
		modalHeight,
	)

	// Render dimmed background (if enabled)
	if m.domain.DimBackground() {
		b.WriteString(m.renderDimmedBackground(x, y, modalWidth, modalHeight))
	}

	// Render modal box.
	b.WriteString(m.renderModalBox(x, y, modalWidth, modalHeight))

	return b.String()
}

// renderDimmedBackground renders the dimmed background overlay.
func (m *Modal) renderDimmedBackground(modalX, modalY, modalWidth, modalHeight int) string {
	var b strings.Builder

	for row := 0; row < m.terminalHeight; row++ {
		for col := 0; col < m.terminalWidth; col++ {
			// Check if this cell is inside the modal area.
			isInsideModal := col >= modalX && col < modalX+modalWidth &&
				row >= modalY && row < modalY+modalHeight

			if isInsideModal {
				b.WriteString(" ") // Modal will render here
			} else {
				b.WriteString("░") // Dim character for background
			}
		}
		b.WriteString("\n")
	}

	return b.String()
}

// renderModalBox renders the modal box at the specified position.
func (m *Modal) renderModalBox(x, y, width, height int) string {
	var b strings.Builder

	// Position cursor at modal location.
	b.WriteString(fmt.Sprintf("\x1b[%d;%dH", y+1, x+1)) // ANSI cursor positioning (1-indexed)

	// Top border.
	b.WriteString("┌" + strings.Repeat("─", width-2) + "┐\n")

	// Title row (if title exists)
	if m.domain.Title() != "" {
		b.WriteString(fmt.Sprintf("\x1b[%d;%dH", y+2, x+1))
		b.WriteString("│ " + m.padOrTruncate(m.domain.Title(), width-4) + " │\n")

		// Title separator.
		b.WriteString(fmt.Sprintf("\x1b[%d;%dH", y+3, x+1))
		b.WriteString("├" + strings.Repeat("─", width-2) + "┤\n")
	}

	// Content area.
	contentStartY := y + 2
	if m.domain.Title() != "" {
		contentStartY = y + 4
	}

	contentLines := strings.Split(m.domain.Content(), "\n")
	contentHeight := height - 4 // Reserve space for borders and buttons
	if m.domain.Title() != "" {
		contentHeight -= 2 // Reserve space for title and separator
	}

	for i := 0; i < contentHeight && i < len(contentLines); i++ {
		b.WriteString(fmt.Sprintf("\x1b[%d;%dH", contentStartY+i, x+1))
		b.WriteString("│ " + m.padOrTruncate(contentLines[i], width-4) + " │\n")
	}

	// Buttons (if any)
	if len(m.domain.Buttons()) > 0 {
		buttonY := y + height - 3
		b.WriteString(fmt.Sprintf("\x1b[%d;%dH", buttonY, x+1))
		b.WriteString("├" + strings.Repeat("─", width-2) + "┤\n")

		buttonY++
		b.WriteString(fmt.Sprintf("\x1b[%d;%dH", buttonY, x+1))
		b.WriteString("│ " + m.renderButtons(width-4) + " │\n")
	}

	// Bottom border.
	b.WriteString(fmt.Sprintf("\x1b[%d;%dH", y+height-1, x+1))
	b.WriteString("└" + strings.Repeat("─", width-2) + "┘\n")

	return b.String()
}

// renderButtons renders the button row.
func (m *Modal) renderButtons(availableWidth int) string {
	buttons := m.domain.Buttons()
	if len(buttons) == 0 {
		return strings.Repeat(" ", availableWidth)
	}

	var parts []string
	focusedIndex := m.domain.FocusedButtonIndex()

	for i, btn := range buttons {
		btnText := fmt.Sprintf("[ %s ]", btn.Label())
		if i == focusedIndex {
			btnText = fmt.Sprintf("< %s >", btn.Label()) // Focused indicator
		}
		parts = append(parts, btnText)
	}

	// Join buttons with spacing.
	result := strings.Join(parts, "  ")

	// Pad or truncate to fit available width.
	return m.padOrTruncate(result, availableWidth)
}

// padOrTruncate pads or truncates a string to the specified width.
func (m *Modal) padOrTruncate(text string, width int) string {
	if len(text) > width {
		if width > 3 {
			return text[:width-3] + "..."
		}
		return text[:width]
	}
	return text + strings.Repeat(" ", width-len(text))
}

// IsVisible returns true if the modal is currently visible.
func (m *Modal) IsVisible() bool {
	return m.domain.IsVisible()
}

// FocusedButton returns the action of the currently focused button.
// Returns empty string if no button is focused.
func (m *Modal) FocusedButton() string {
	btn := m.domain.FocusedButton()
	if btn == nil {
		return ""
	}
	return btn.Action()
}

// ButtonPressedMsg is sent when a button is activated.
type ButtonPressedMsg struct {
	Action string // Button action identifier
}
