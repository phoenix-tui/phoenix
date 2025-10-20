// Package input provides a single-line text input component with cursor control and validation.
package input

import (
	"github.com/phoenix-tui/phoenix/components/input/domain/model"
	"github.com/phoenix-tui/phoenix/components/input/domain/service"
	"github.com/phoenix-tui/phoenix/components/input/infrastructure"
	tea "github.com/phoenix-tui/phoenix/tea/api"
)

// Input is the public API for the TextInput component.
// It wraps the domain model and provides a fluent interface for configuration.
// Implements tea.Model for use in Elm Architecture applications.
// Uses value semantics for immutable updates.
type Input struct {
	domain      model.TextInput // VALUE, not pointer!
	keyBindings KeyBindingHandler
}

// KeyBindingHandler handles key messages and returns updated input.
type KeyBindingHandler interface {
	Handle(model.TextInput, tea.KeyMsg) model.TextInput // VALUE parameters!
}

// New creates a new Input with the specified visible width.
// Returns pointer for initialization, but store as value in Model.
func New(width int) *Input {
	return &Input{
		domain:      *model.New(width), // Dereference!
		keyBindings: infrastructure.NewDefaultKeyBindings(),
	}
}

// Placeholder sets the placeholder text shown when input is empty and unfocused.
// Returns new Input for method chaining (value semantics).
// IMPORTANT: Must reassign: input = input.Placeholder("text").
func (i Input) Placeholder(text string) Input {
	i.domain = i.domain.WithPlaceholder(text)
	return i
}

// Validator sets the validation function for input content.
// Returns new Input for method chaining (value semantics).
// IMPORTANT: Must reassign: input = input.Validator(fn).
func (i Input) Validator(fn func(string) error) Input {
	i.domain = i.domain.WithValidator(fn)
	return i
}

// Content sets the initial content.
// Returns new Input for method chaining (value semantics).
// IMPORTANT: Must reassign: input = input.Content("text").
func (i Input) Content(text string) Input {
	i.domain = i.domain.WithContent(text)
	return i
}

// Focused sets the focus state.
// Returns new Input for method chaining (value semantics).
// IMPORTANT: Must reassign: input = input.Focused(true).
func (i Input) Focused(focused bool) Input {
	i.domain = i.domain.WithFocus(focused)
	return i
}

// ShowCursor sets whether the cursor should be rendered.
// When false, applications can use the terminal's native cursor instead.
// This is useful for shells that prefer the terminal's native blinking cursor.
// Returns new Input for method chaining (value semantics).
// IMPORTANT: Must reassign: input = input.ShowCursor(false).
func (i Input) ShowCursor(show bool) Input {
	i.domain = i.domain.WithShowCursor(show)
	return i
}

// Width sets the visible width.
// Returns new Input for method chaining (value semantics).
// IMPORTANT: Must reassign: input = input.Width(80).
func (i Input) Width(width int) Input {
	i.domain = i.domain.WithWidth(width)
	return i
}

// KeyBindings sets a custom key binding handler.
// Returns new Input for method chaining (value semantics).
// IMPORTANT: Must reassign: input = input.KeyBindings(handler).
func (i Input) KeyBindings(handler KeyBindingHandler) Input {
	i.keyBindings = handler
	return i
}

// CursorPosition returns the current cursor position (grapheme offset).
// This is a KEY DIFFERENTIATOR - enables apps to customize cursor rendering.
func (i Input) CursorPosition() int {
	return i.domain.CursorPosition()
}

// ContentParts splits content around cursor: (before, at, after).
// This is a KEY DIFFERENTIATOR - enables apps to customize cursor rendering.
func (i Input) ContentParts() (before, at, after string) {
	return i.domain.ContentParts()
}

// SetContent sets both content and cursor position atomically.
// This is a KEY DIFFERENTIATOR - prevents race conditions.
// Returns new Input for method chaining (value semantics).
// IMPORTANT: Must reassign: input = input.SetContent("text", 0).
func (i Input) SetContent(content string, cursorPos int) Input {
	i.domain = i.domain.SetContent(content, cursorPos)
	return i
}

// Value returns the current input value.
func (i Input) Value() string {
	return i.domain.Content()
}

// IsValid returns true if the content passes validation.
func (i Input) IsValid() bool {
	return i.domain.IsValid()
}

// IsFocused returns true if the input is focused.
func (i Input) IsFocused() bool {
	return i.domain.Focused()
}

// Init implements tea.Model.
func (i Input) Init() tea.Cmd {
	return nil
}

// Update implements the Update method of Elm Architecture.
// IMPORTANT: Must reassign: input = input.Update(msg).
func (i Input) Update(msg tea.Msg) (Input, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Only process keys if focused.
		if !i.domain.Focused() {
			return i, nil
		}

		// Handle key via bindings.
		result := i.keyBindings.Handle(i.domain, msg)
		i.domain = result
		return i, nil

	default:
		return i, nil
	}
}

// View implements tea.Model.
func (i Input) View() string {
	// If empty and not focused, show placeholder.
	if i.domain.Content() == "" && !i.domain.Focused() {
		if i.domain.Placeholder() != "" {
			return i.renderPlaceholder()
		}
		return ""
	}

	// Render content with cursor.
	return i.renderContent()
}

// renderPlaceholder renders the placeholder text.
func (i Input) renderPlaceholder() string {
	// Simple gray placeholder (apps can customize via styling)
	// For now, just return the placeholder text.
	// TODO: Apply styling when style system is integrated.
	return i.domain.Placeholder()
}

// renderContent renders the input content with cursor.
//
//nolint:gocognit // rendering logic requires multiple conditions
func (i Input) renderContent() string {
	before, at, after := i.domain.ContentParts()

	// Calculate visible portion based on width.
	// For now, simple truncation (TODO: proper scrolling)
	content := before + at + after

	if len(content) <= i.domain.Width() {
		// Content fits, render cursor.
		if i.domain.Focused() {
			return before + i.renderCursor(at) + after
		}
		return content
	}

	// Content too long, needs scrolling.
	// Simple approach: keep cursor in view.
	cursorPos := i.domain.CursorPosition()
	width := i.domain.Width()

	// Calculate scroll offset to keep cursor visible.
	scrollOffset := i.domain.ScrollOffset()

	// Adjust scroll if cursor is out of view.
	if cursorPos < scrollOffset {
		scrollOffset = cursorPos
	}
	if cursorPos >= scrollOffset+width {
		scrollOffset = cursorPos - width + 1
	}

	// Extract visible portion (simplified - works for ASCII, needs grapheme handling)
	visibleStart := scrollOffset
	visibleEnd := scrollOffset + width
	if visibleEnd > len(content) {
		visibleEnd = len(content)
	}

	visibleContent := content[visibleStart:visibleEnd]

	// Render cursor if in visible range.
	//nolint:nestif // rendering logic requires nested conditions for cursor positioning
	if i.domain.Focused() && cursorPos >= scrollOffset && cursorPos < scrollOffset+width {
		cursorInVisible := cursorPos - scrollOffset
		if cursorInVisible <= len(visibleContent) {
			visibleBefore := visibleContent[:cursorInVisible]
			visibleAt := ""
			if cursorInVisible < len(visibleContent) {
				visibleAt = string(visibleContent[cursorInVisible])
			}
			visibleAfter := ""
			if cursorInVisible+1 < len(visibleContent) {
				visibleAfter = visibleContent[cursorInVisible+1:]
			}
			return visibleBefore + i.renderCursor(visibleAt) + visibleAfter
		}
	}

	return visibleContent
}

// renderCursor renders the cursor at the current position.
// Returns empty string if cursor rendering is disabled (terminal cursor used instead).
func (i Input) renderCursor(at string) string {
	// If cursor rendering is disabled, return the character without highlighting.
	// This allows the terminal's native cursor to be visible instead.
	if !i.domain.ShowCursor() {
		return at
	}

	// Render Phoenix's cursor using reverse video.
	if at == "" {
		return "█" // Block cursor at end
	}
	return "\033[7m" + at + "\033[0m" // Reverse video
}

// ValidationFunc is the type for validation functions.
// Re-exported for convenience.
type ValidationFunc = service.ValidationFunc

// Common validators (re-exported for convenience).
var (
	NotEmpty  = service.NotEmpty
	MinLength = service.MinLength
	MaxLength = service.MaxLength
	Range     = service.Range
	Chain     = service.Chain
)

// Validation errors (re-exported for convenience).
var (
	ErrEmpty         = service.ErrEmpty
	ErrTooShort      = service.ErrTooShort
	ErrTooLong       = service.ErrTooLong
	ErrInvalidFormat = service.ErrInvalidFormat
)

// CustomKeyBindings creates custom key bindings for input handling.
func CustomKeyBindings(handlers ...infrastructure.KeyHandler) KeyBindingHandler {
	return infrastructure.NewCustomKeyBindings(handlers...)
}
