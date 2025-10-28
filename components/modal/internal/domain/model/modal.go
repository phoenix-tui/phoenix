package model

import (
	value2 "github.com/phoenix-tui/phoenix/components/modal/internal/domain/value"
)

// Modal is the aggregate root for modal/overlay component.
// It represents a dialog box that can be centered or positioned at custom coordinates,.
// with optional title, content, buttons, and background dimming.
//
// All operations are immutable - methods return new instances rather than modifying in place.
type Modal struct {
	title         string           // Modal title (optional, can be empty)
	content       string           // Modal content (text or rendered sub-component)
	buttons       []*Button        // Action buttons (optional, can be empty)
	size          *value2.Size     // Modal size (width, height)
	position      *value2.Position // Position (center or custom x, y)
	focusedButton int              // Currently focused button index
	visible       bool             // Is modal visible?
	dimBackground bool             // Dim background when visible?
}

// NewModal creates a new modal with the given content.
// Default configuration: centered, 40x10, not visible, no dimming.
func NewModal(content string) *Modal {
	return &Modal{
		title:         "",
		content:       content,
		buttons:       []*Button{},
		size:          value2.NewSize(40, 10),
		position:      value2.NewPositionCenter(),
		focusedButton: 0,
		visible:       false,
		dimBackground: false,
	}
}

// NewModalWithTitle creates a new modal with title and content.
func NewModalWithTitle(title, content string) *Modal {
	return &Modal{
		title:         title,
		content:       content,
		buttons:       []*Button{},
		size:          value2.NewSize(40, 10),
		position:      value2.NewPositionCenter(),
		focusedButton: 0,
		visible:       false,
		dimBackground: false,
	}
}

// WithTitle returns a new modal with the specified title.
func (m *Modal) WithTitle(title string) *Modal {
	return &Modal{
		title:         title,
		content:       m.content,
		buttons:       m.buttons,
		size:          m.size,
		position:      m.position,
		focusedButton: m.focusedButton,
		visible:       m.visible,
		dimBackground: m.dimBackground,
	}
}

// WithContent returns a new modal with the specified content.
func (m *Modal) WithContent(content string) *Modal {
	return &Modal{
		title:         m.title,
		content:       content,
		buttons:       m.buttons,
		size:          m.size,
		position:      m.position,
		focusedButton: m.focusedButton,
		visible:       m.visible,
		dimBackground: m.dimBackground,
	}
}

// WithSize returns a new modal with the specified size.
func (m *Modal) WithSize(width, height int) *Modal {
	return &Modal{
		title:         m.title,
		content:       m.content,
		buttons:       m.buttons,
		size:          value2.NewSize(width, height),
		position:      m.position,
		focusedButton: m.focusedButton,
		visible:       m.visible,
		dimBackground: m.dimBackground,
	}
}

// WithPosition returns a new modal with the specified position.
func (m *Modal) WithPosition(position *value2.Position) *Modal {
	return &Modal{
		title:         m.title,
		content:       m.content,
		buttons:       m.buttons,
		size:          m.size,
		position:      position,
		focusedButton: m.focusedButton,
		visible:       m.visible,
		dimBackground: m.dimBackground,
	}
}

// WithButtons returns a new modal with the specified buttons.
// Focused button index is reset to 0.
func (m *Modal) WithButtons(buttons []*Button) *Modal {
	return &Modal{
		title:         m.title,
		content:       m.content,
		buttons:       buttons,
		size:          m.size,
		position:      m.position,
		focusedButton: 0, // Reset focus when buttons change
		visible:       m.visible,
		dimBackground: m.dimBackground,
	}
}

// WithDimBackground returns a new modal with dimming enabled/disabled.
func (m *Modal) WithDimBackground(dim bool) *Modal {
	return &Modal{
		title:         m.title,
		content:       m.content,
		buttons:       m.buttons,
		size:          m.size,
		position:      m.position,
		focusedButton: m.focusedButton,
		visible:       m.visible,
		dimBackground: dim,
	}
}

// WithVisible returns a new modal with visibility set.
func (m *Modal) WithVisible(visible bool) *Modal {
	return &Modal{
		title:         m.title,
		content:       m.content,
		buttons:       m.buttons,
		size:          m.size,
		position:      m.position,
		focusedButton: m.focusedButton,
		visible:       visible,
		dimBackground: m.dimBackground,
	}
}

// FocusNextButton returns a new modal with focus moved to the next button.
// Wraps around to first button if at the end.
func (m *Modal) FocusNextButton() *Modal {
	if len(m.buttons) == 0 {
		return m // No buttons to focus
	}

	nextFocus := (m.focusedButton + 1) % len(m.buttons)

	return &Modal{
		title:         m.title,
		content:       m.content,
		buttons:       m.buttons,
		size:          m.size,
		position:      m.position,
		focusedButton: nextFocus,
		visible:       m.visible,
		dimBackground: m.dimBackground,
	}
}

// FocusPreviousButton returns a new modal with focus moved to the previous button.
// Wraps around to last button if at the beginning.
func (m *Modal) FocusPreviousButton() *Modal {
	if len(m.buttons) == 0 {
		return m // No buttons to focus
	}

	prevFocus := m.focusedButton - 1
	if prevFocus < 0 {
		prevFocus = len(m.buttons) - 1
	}

	return &Modal{
		title:         m.title,
		content:       m.content,
		buttons:       m.buttons,
		size:          m.size,
		position:      m.position,
		focusedButton: prevFocus,
		visible:       m.visible,
		dimBackground: m.dimBackground,
	}
}

// FocusedButton returns the currently focused button (nil if no buttons).
func (m *Modal) FocusedButton() *Button {
	if len(m.buttons) == 0 || m.focusedButton < 0 || m.focusedButton >= len(m.buttons) {
		return nil
	}
	return m.buttons[m.focusedButton]
}

// Title returns the modal title.
func (m *Modal) Title() string {
	return m.title
}

// Content returns the modal content.
func (m *Modal) Content() string {
	return m.content
}

// Buttons returns the modal buttons.
func (m *Modal) Buttons() []*Button {
	return m.buttons
}

// Size returns the modal size.
func (m *Modal) Size() *value2.Size {
	return m.size
}

// Position returns the modal position.
func (m *Modal) Position() *value2.Position {
	return m.position
}

// IsVisible returns true if the modal is visible.
func (m *Modal) IsVisible() bool {
	return m.visible
}

// DimBackground returns true if background dimming is enabled.
func (m *Modal) DimBackground() bool {
	return m.dimBackground
}

// FocusedButtonIndex returns the index of the focused button.
func (m *Modal) FocusedButtonIndex() int {
	return m.focusedButton
}
