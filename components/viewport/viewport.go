// Package viewport provides a scrollable viewport component for TUI applications.
package viewport

import (
	"strings"

	"github.com/phoenix-tui/phoenix/components/viewport/internal/domain/model"
	"github.com/phoenix-tui/phoenix/components/viewport/internal/infrastructure"
	"github.com/phoenix-tui/phoenix/tea"
)

// Viewport is the public API for the Viewport component.
// It provides a fluent interface for configuration and implements tea.Model.
type Viewport struct {
	domain         *model.Viewport
	mouseEnabled   bool
	linesPerScroll int // Lines to scroll per wheel tick (default: 3)
	// Drag scrolling state
	isDragging   bool
	dragStartY   int
	scrollStartY int
}

// New creates a new Viewport with the given dimensions.
func New(width, height int) *Viewport {
	return &Viewport{
		domain:         model.NewViewport(width, height),
		mouseEnabled:   false,
		linesPerScroll: 3, // Default: 3 lines per wheel tick
	}
}

// NewWithContent creates a new Viewport with initial content.
// Content can be a single string with newlines or multiple lines.
func NewWithContent(content string, width, height int) *Viewport {
	lines := strings.Split(content, "\n")
	return &Viewport{
		domain:         model.NewViewportWithContent(lines, width, height),
		mouseEnabled:   false,
		linesPerScroll: 3, // Default: 3 lines per wheel tick
	}
}

// NewWithLines creates a new Viewport with initial content as separate lines.
func NewWithLines(lines []string, width, height int) *Viewport {
	return &Viewport{
		domain:         model.NewViewportWithContent(lines, width, height),
		mouseEnabled:   false,
		linesPerScroll: 3, // Default: 3 lines per wheel tick
	}
}

// FollowMode enables or disables follow mode (tail -f style auto-scrolling).
// When enabled, the viewport automatically scrolls to the bottom when content changes.
func (v *Viewport) FollowMode(enabled bool) *Viewport {
	return v.withDomain(v.domain.WithFollowMode(enabled))
}

// WrapLines enables or disables line wrapping.
// When enabled, lines wider than the viewport are wrapped to multiple lines.
func (v *Viewport) WrapLines(enabled bool) *Viewport {
	return v.withDomain(v.domain.WithWrapLines(enabled))
}

// MouseEnabled enables or disables mouse wheel scrolling and drag scrolling.
func (v *Viewport) MouseEnabled(enabled bool) *Viewport {
	return &Viewport{
		domain:         v.domain,
		mouseEnabled:   enabled,
		linesPerScroll: v.linesPerScroll,
		isDragging:     v.isDragging,
		dragStartY:     v.dragStartY,
		scrollStartY:   v.scrollStartY,
	}
}

// SetWheelScrollLines sets the number of lines to scroll per mouse wheel tick.
// Default is 3 lines. Minimum is 1 line.
// Returns a new Viewport with the updated configuration (immutable).
func (v *Viewport) SetWheelScrollLines(lines int) *Viewport {
	if lines < 1 {
		lines = 1 // Minimum 1 line
	}
	return &Viewport{
		domain:         v.domain,
		mouseEnabled:   v.mouseEnabled,
		linesPerScroll: lines,
		isDragging:     v.isDragging,
		dragStartY:     v.dragStartY,
		scrollStartY:   v.scrollStartY,
	}
}

// SetContent replaces the viewport content with the given string.
// Content is split by newlines into individual lines.
func (v *Viewport) SetContent(content string) *Viewport {
	lines := strings.Split(content, "\n")
	return v.withDomain(v.domain.WithContent(lines))
}

// SetLines replaces the viewport content with the given lines.
func (v *Viewport) SetLines(lines []string) *Viewport {
	return v.withDomain(v.domain.WithContent(lines))
}

// AppendLine appends a single line to the viewport content.
// Useful for streaming content (like log viewers, command output accumulation).
// The line is added to the end of existing content.
func (v *Viewport) AppendLine(line string) *Viewport {
	currentContent := v.domain.Content()
	currentContent = append(currentContent, line)
	return v.SetLines(currentContent)
}

// AppendLines appends multiple lines to the viewport content.
// Useful for batch updates (e.g., adding multiple log entries at once).
// More efficient than multiple AppendLine() calls.
func (v *Viewport) AppendLines(lines []string) *Viewport {
	currentContent := v.domain.Content()
	currentContent = append(currentContent, lines...)
	return v.SetLines(currentContent)
}

// ScrollToBottom scrolls the viewport to the bottom (last line).
// This is a one-time scroll action, not continuous like FollowMode.
// Useful for: user presses End key, jump to bottom on command, etc.
func (v *Viewport) ScrollToBottom() *Viewport {
	return v.withDomain(v.domain.ScrollToBottom())
}

// ScrollToTop scrolls the viewport to the top (first line).
// This is a one-time scroll action.
// Useful for: user presses Home key, reset to top on command, etc.
func (v *Viewport) ScrollToTop() *Viewport {
	return v.withDomain(v.domain.ScrollToTop())
}

// SetYOffset sets the vertical scroll offset to a specific line number.
// This provides precise control over the viewport position.
// Offset is clamped to valid range [0, maxOffset].
// Useful for: jumping to specific line, restoring scroll position, etc.
func (v *Viewport) SetYOffset(offset int) *Viewport {
	return v.withDomain(v.domain.WithScrollOffset(offset))
}

// SetSize updates the viewport dimensions.
func (v *Viewport) SetSize(width, height int) *Viewport {
	return v.withDomain(v.domain.WithSize(width, height))
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
		return v.withDomain(v.domain.ScrollUp(1))
	}

	if infrastructure.IsDownKey(msg) {
		return v.withDomain(v.domain.ScrollDown(1))
	}

	if infrastructure.IsPageUpKey(msg) {
		return v.withDomain(v.domain.PageUp())
	}

	if infrastructure.IsPageDownKey(msg) {
		return v.withDomain(v.domain.PageDown())
	}

	if infrastructure.IsHomeKey(msg) {
		return v.withDomain(v.domain.ScrollToTop())
	}

	if infrastructure.IsEndKey(msg) {
		return v.withDomain(v.domain.ScrollToBottom())
	}

	if infrastructure.IsHalfPageUpKey(msg) {
		halfPage := v.domain.VisibleHeight() / 2
		return v.withDomain(v.domain.ScrollUp(halfPage))
	}

	if infrastructure.IsHalfPageDownKey(msg) {
		halfPage := v.domain.VisibleHeight() / 2
		return v.withDomain(v.domain.ScrollDown(halfPage))
	}

	return v
}

// handleMouseMsg processes mouse input for scrolling (wheel and drag).
func (v *Viewport) handleMouseMsg(msg tea.MouseMsg) *Viewport {
	// Mouse wheel scrolling
	switch msg.Button {
	case tea.MouseButtonWheelUp:
		return v.withDomain(v.domain.ScrollUp(v.linesPerScroll))

	case tea.MouseButtonWheelDown:
		return v.withDomain(v.domain.ScrollDown(v.linesPerScroll))
	}

	// Drag scrolling
	switch msg.Action {
	case tea.MouseActionPress:
		if msg.Button == tea.MouseButtonLeft {
			// Start drag: record starting Y position and current scroll offset
			return v.withDragState(true, msg.Y, v.domain.ScrollOffset())
		}

	case tea.MouseActionRelease:
		if v.isDragging {
			// End drag: clear drag state
			return v.withDragState(false, 0, 0)
		}

	case tea.MouseActionMotion:
		if v.isDragging {
			// Calculate scroll delta from drag motion
			deltaY := msg.Y - v.dragStartY

			// Calculate new scroll offset
			// When dragging down (+Y), content should scroll up (negative scroll)
			// When dragging up (-Y), content should scroll down (positive scroll)
			newScrollOffset := v.scrollStartY - deltaY

			// Apply scroll with bounds checking (domain handles clamping)
			return v.withDomain(v.domain.WithScrollOffset(newScrollOffset))
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

// withDomain returns a new Viewport with updated domain model.
// This is a helper to maintain immutability and avoid repetitive field copying.
func (v *Viewport) withDomain(domain *model.Viewport) *Viewport {
	return &Viewport{
		domain:         domain,
		mouseEnabled:   v.mouseEnabled,
		linesPerScroll: v.linesPerScroll,
		isDragging:     v.isDragging,
		dragStartY:     v.dragStartY,
		scrollStartY:   v.scrollStartY,
	}
}

// withDragState returns a new Viewport with updated drag state.
func (v *Viewport) withDragState(isDragging bool, dragStartY, scrollStartY int) *Viewport {
	return &Viewport{
		domain:         v.domain,
		mouseEnabled:   v.mouseEnabled,
		linesPerScroll: v.linesPerScroll,
		isDragging:     isDragging,
		dragStartY:     dragStartY,
		scrollStartY:   scrollStartY,
	}
}
