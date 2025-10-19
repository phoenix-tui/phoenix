package viewport

import (
	"strings"

	"github.com/phoenix-tui/phoenix/components/viewport/domain/model"
	"github.com/phoenix-tui/phoenix/components/viewport/infrastructure"
	tea "github.com/phoenix-tui/phoenix/tea/api"
)

// Viewport is the public API for the Viewport component.
// It provides a fluent interface for configuration and implements tea.Model.
type Viewport struct {
	domain       *model.Viewport
	mouseEnabled bool
}

// New creates a new Viewport with the given dimensions.
func New(width, height int) *Viewport {
	return &Viewport{
		domain:       model.NewViewport(width, height),
		mouseEnabled: false,
	}
}

// NewWithContent creates a new Viewport with initial content.
// Content can be a single string with newlines or multiple lines.
func NewWithContent(content string, width, height int) *Viewport {
	lines := strings.Split(content, "\n")
	return &Viewport{
		domain:       model.NewViewportWithContent(lines, width, height),
		mouseEnabled: false,
	}
}

// NewWithLines creates a new Viewport with initial content as separate lines.
func NewWithLines(lines []string, width, height int) *Viewport {
	return &Viewport{
		domain:       model.NewViewportWithContent(lines, width, height),
		mouseEnabled: false,
	}
}

// FollowMode enables or disables follow mode (tail -f style auto-scrolling).
// When enabled, the viewport automatically scrolls to the bottom when content changes.
func (v *Viewport) FollowMode(enabled bool) *Viewport {
	return &Viewport{
		domain:       v.domain.WithFollowMode(enabled),
		mouseEnabled: v.mouseEnabled,
	}
}

// WrapLines enables or disables line wrapping.
// When enabled, lines wider than the viewport are wrapped to multiple lines.
func (v *Viewport) WrapLines(enabled bool) *Viewport {
	return &Viewport{
		domain:       v.domain.WithWrapLines(enabled),
		mouseEnabled: v.mouseEnabled,
	}
}

// MouseEnabled enables or disables mouse wheel scrolling.
func (v *Viewport) MouseEnabled(enabled bool) *Viewport {
	return &Viewport{
		domain:       v.domain,
		mouseEnabled: enabled,
	}
}

// SetContent replaces the viewport content with the given string.
// Content is split by newlines into individual lines.
func (v *Viewport) SetContent(content string) *Viewport {
	lines := strings.Split(content, "\n")
	return &Viewport{
		domain:       v.domain.WithContent(lines),
		mouseEnabled: v.mouseEnabled,
	}
}

// SetLines replaces the viewport content with the given lines.
func (v *Viewport) SetLines(lines []string) *Viewport {
	return &Viewport{
		domain:       v.domain.WithContent(lines),
		mouseEnabled: v.mouseEnabled,
	}
}

// AppendLine appends a single line to the viewport content.
// Useful for streaming content (like log viewers, command output accumulation).
// The line is added to the end of existing content.
func (v *Viewport) AppendLine(line string) *Viewport {
	currentContent := v.domain.Content()
	newContent := append(currentContent, line)
	return v.SetLines(newContent)
}

// AppendLines appends multiple lines to the viewport content.
// Useful for batch updates (e.g., adding multiple log entries at once).
// More efficient than multiple AppendLine() calls.
func (v *Viewport) AppendLines(lines []string) *Viewport {
	currentContent := v.domain.Content()
	newContent := append(currentContent, lines...)
	return v.SetLines(newContent)
}

// ScrollToBottom scrolls the viewport to the bottom (last line).
// This is a one-time scroll action, not continuous like FollowMode.
// Useful for: user presses End key, jump to bottom on command, etc.
func (v *Viewport) ScrollToBottom() *Viewport {
	return &Viewport{
		domain:       v.domain.ScrollToBottom(),
		mouseEnabled: v.mouseEnabled,
	}
}

// ScrollToTop scrolls the viewport to the top (first line).
// This is a one-time scroll action.
// Useful for: user presses Home key, reset to top on command, etc.
func (v *Viewport) ScrollToTop() *Viewport {
	return &Viewport{
		domain:       v.domain.ScrollToTop(),
		mouseEnabled: v.mouseEnabled,
	}
}

// SetYOffset sets the vertical scroll offset to a specific line number.
// This provides precise control over the viewport position.
// Offset is clamped to valid range [0, maxOffset].
// Useful for: jumping to specific line, restoring scroll position, etc.
func (v *Viewport) SetYOffset(offset int) *Viewport {
	return &Viewport{
		domain:       v.domain.WithScrollOffset(offset),
		mouseEnabled: v.mouseEnabled,
	}
}

// SetSize updates the viewport dimensions.
func (v *Viewport) SetSize(width, height int) *Viewport {
	return &Viewport{
		domain:       v.domain.WithSize(width, height),
		mouseEnabled: v.mouseEnabled,
	}
}

// Init initializes the viewport (implements Init() Cmd for tea.Model constraint).
func (v *Viewport) Init() tea.Cmd {
	return nil
}

// Update processes messages and returns updated viewport (implements Update(Msg) (*Viewport, Cmd)).
// Note: Returns concrete *Viewport instead of tea.Model interface for better usability.
func (v *Viewport) Update(msg tea.Msg) (*Viewport, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return v.handleKeyMsg(msg), nil

	case tea.MouseMsg:
		if v.mouseEnabled {
			return v.handleMouseMsg(msg), nil
		}

	case tea.WindowSizeMsg:
		return v.SetSize(msg.Width, msg.Height), nil
	}

	return v, nil
}

// handleKeyMsg processes keyboard input for scrolling.
func (v *Viewport) handleKeyMsg(msg tea.KeyMsg) *Viewport {
	if infrastructure.IsUpKey(msg) {
		return &Viewport{
			domain:       v.domain.ScrollUp(1),
			mouseEnabled: v.mouseEnabled,
		}
	}

	if infrastructure.IsDownKey(msg) {
		return &Viewport{
			domain:       v.domain.ScrollDown(1),
			mouseEnabled: v.mouseEnabled,
		}
	}

	if infrastructure.IsPageUpKey(msg) {
		return &Viewport{
			domain:       v.domain.PageUp(),
			mouseEnabled: v.mouseEnabled,
		}
	}

	if infrastructure.IsPageDownKey(msg) {
		return &Viewport{
			domain:       v.domain.PageDown(),
			mouseEnabled: v.mouseEnabled,
		}
	}

	if infrastructure.IsHomeKey(msg) {
		return &Viewport{
			domain:       v.domain.ScrollToTop(),
			mouseEnabled: v.mouseEnabled,
		}
	}

	if infrastructure.IsEndKey(msg) {
		return &Viewport{
			domain:       v.domain.ScrollToBottom(),
			mouseEnabled: v.mouseEnabled,
		}
	}

	if infrastructure.IsHalfPageUpKey(msg) {
		halfPage := v.domain.VisibleHeight() / 2
		return &Viewport{
			domain:       v.domain.ScrollUp(halfPage),
			mouseEnabled: v.mouseEnabled,
		}
	}

	if infrastructure.IsHalfPageDownKey(msg) {
		halfPage := v.domain.VisibleHeight() / 2
		return &Viewport{
			domain:       v.domain.ScrollDown(halfPage),
			mouseEnabled: v.mouseEnabled,
		}
	}

	return v
}

// handleMouseMsg processes mouse input for scrolling.
func (v *Viewport) handleMouseMsg(msg tea.MouseMsg) *Viewport {
	switch msg.Button {
	case tea.MouseButtonWheelUp:
		return &Viewport{
			domain:       v.domain.ScrollUp(3), // Scroll 3 lines per wheel tick
			mouseEnabled: v.mouseEnabled,
		}

	case tea.MouseButtonWheelDown:
		return &Viewport{
			domain:       v.domain.ScrollDown(3), // Scroll 3 lines per wheel tick
			mouseEnabled: v.mouseEnabled,
		}
	}

	return v
}

// View implements tea.Model.
func (v *Viewport) View() string {
	lines := v.domain.VisibleLines()

	if len(lines) == 0 {
		return ""
	}

	return strings.Join(lines, "\n")
}

// VisibleLines returns the currently visible lines.
func (v *Viewport) VisibleLines() []string {
	return v.domain.VisibleLines()
}

// ScrollOffset returns the current scroll offset.
func (v *Viewport) ScrollOffset() int {
	return v.domain.ScrollOffset()
}

// CanScrollUp returns true if the viewport can scroll up.
func (v *Viewport) CanScrollUp() bool {
	return v.domain.CanScrollUp()
}

// CanScrollDown returns true if the viewport can scroll down.
func (v *Viewport) CanScrollDown() bool {
	return v.domain.CanScrollDown()
}

// IsAtTop returns true if the viewport is at the top.
func (v *Viewport) IsAtTop() bool {
	return v.domain.IsAtTop()
}

// IsAtBottom returns true if the viewport is at the bottom.
func (v *Viewport) IsAtBottom() bool {
	return v.domain.IsAtBottom()
}

// TotalLines returns the total number of content lines.
func (v *Viewport) TotalLines() int {
	return v.domain.TotalLines()
}

// Height returns the viewport height.
func (v *Viewport) Height() int {
	return v.domain.VisibleHeight()
}
