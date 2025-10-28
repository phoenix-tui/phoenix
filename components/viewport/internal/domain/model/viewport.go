// Package model contains domain models for viewport component.
package model

import (
	"strings"

	"github.com/phoenix-tui/phoenix/components/viewport/internal/domain/service"
	value2 "github.com/phoenix-tui/phoenix/components/viewport/internal/domain/value"
	"github.com/rivo/uniseg"
)

// Viewport is the aggregate root for scrollable content.
// It manages content display, scrolling, and viewport size.
// All operations are immutable - they return new Viewport instances.
type Viewport struct {
	content      []string
	size         *value2.ViewportSize
	scrollOffset *value2.ScrollOffset
	followMode   bool
	wrapLines    bool
	scrollSvc    *service.ScrollService
}

// NewViewport creates a new Viewport with the given dimensions.
func NewViewport(width, height int) *Viewport {
	return &Viewport{
		content:      []string{},
		size:         value2.NewViewportSize(width, height),
		scrollOffset: value2.NewScrollOffset(0),
		followMode:   false,
		wrapLines:    false,
		scrollSvc:    service.NewScrollService(),
	}
}

// NewViewportWithContent creates a new Viewport with initial content.
// Content is split by newlines into individual lines.
func NewViewportWithContent(content []string, width, height int) *Viewport {
	v := NewViewport(width, height)
	return v.WithContent(content)
}

// WithContent returns a new Viewport with the given content.
// If follow mode is enabled, the viewport scrolls to the bottom.
func (v *Viewport) WithContent(content []string) *Viewport {
	newV := v.clone()
	newV.content = make([]string, len(content))
	copy(newV.content, content)

	// If follow mode is enabled, scroll to bottom.
	if newV.followMode {
		offset := newV.scrollSvc.FollowModeOffset(len(newV.content), newV.size.Height())
		newV.scrollOffset = value2.NewScrollOffset(offset)
	} else {
		// Clamp existing offset to new content bounds.
		maxOffset := newV.scrollSvc.MaxScrollOffset(len(newV.content), newV.size.Height())
		newV.scrollOffset = newV.scrollOffset.Clamp(maxOffset)
	}

	return newV
}

// WithSize returns a new Viewport with the given dimensions.
func (v *Viewport) WithSize(width, height int) *Viewport {
	newV := v.clone()
	newV.size = value2.NewViewportSize(width, height)

	// Clamp scroll offset to new size bounds.
	maxOffset := newV.scrollSvc.MaxScrollOffset(len(newV.content), height)
	newV.scrollOffset = newV.scrollOffset.Clamp(maxOffset)

	return newV
}

// WithFollowMode returns a new Viewport with follow mode enabled/disabled.
// When enabled, the viewport automatically scrolls to the bottom when content changes.
func (v *Viewport) WithFollowMode(enabled bool) *Viewport {
	newV := v.clone()
	newV.followMode = enabled

	// If enabling follow mode, scroll to bottom immediately.
	if enabled {
		offset := newV.scrollSvc.FollowModeOffset(len(newV.content), newV.size.Height())
		newV.scrollOffset = value2.NewScrollOffset(offset)
	}

	return newV
}

// WithWrapLines returns a new Viewport with line wrapping enabled/disabled.
func (v *Viewport) WithWrapLines(enabled bool) *Viewport {
	newV := v.clone()
	newV.wrapLines = enabled
	return newV
}

// ScrollUp returns a new Viewport scrolled up by the given number of lines.
// If follow mode is enabled, it is automatically disabled.
func (v *Viewport) ScrollUp(lines int) *Viewport {
	newV := v.clone()
	newOffset := v.scrollSvc.ScrollUp(v.scrollOffset.Offset(), lines)
	newV.scrollOffset = value2.NewScrollOffset(newOffset)
	newV.followMode = false // Disable follow mode on manual scroll
	return newV
}

// ScrollDown returns a new Viewport scrolled down by the given number of lines.
func (v *Viewport) ScrollDown(lines int) *Viewport {
	newV := v.clone()
	maxOffset := v.scrollSvc.MaxScrollOffset(len(v.content), v.size.Height())
	newOffset := v.scrollSvc.ScrollDown(v.scrollOffset.Offset(), lines, maxOffset)
	newV.scrollOffset = value2.NewScrollOffset(newOffset)

	// Re-enable follow mode if we've scrolled to the bottom.
	if newV.scrollSvc.IsAtBottom(newOffset, len(v.content), v.size.Height()) {
		newV.followMode = true
	}

	return newV
}

// ScrollToTop returns a new Viewport scrolled to the top.
// Follow mode is disabled.
func (v *Viewport) ScrollToTop() *Viewport {
	newV := v.clone()
	newV.scrollOffset = value2.NewScrollOffset(0)
	newV.followMode = false
	return newV
}

// ScrollToBottom returns a new Viewport scrolled to the bottom.
// Follow mode is enabled.
func (v *Viewport) ScrollToBottom() *Viewport {
	newV := v.clone()
	offset := v.scrollSvc.FollowModeOffset(len(v.content), v.size.Height())
	newV.scrollOffset = value2.NewScrollOffset(offset)
	newV.followMode = true
	return newV
}

// WithScrollOffset returns a new Viewport with the scroll offset set to the given value.
// The offset is clamped to valid range [0, maxOffset].
// Follow mode is disabled when manually setting the offset.
// This provides low-level control over viewport position.
func (v *Viewport) WithScrollOffset(offset int) *Viewport {
	newV := v.clone()
	maxOffset := v.scrollSvc.MaxScrollOffset(len(v.content), v.size.Height())
	// Clamp to valid range.
	if offset < 0 {
		offset = 0
	}
	if offset > maxOffset {
		offset = maxOffset
	}
	newV.scrollOffset = value2.NewScrollOffset(offset)
	newV.followMode = false // Disable follow mode on manual scroll
	return newV
}

// PageUp returns a new Viewport scrolled up by the viewport height.
func (v *Viewport) PageUp() *Viewport {
	return v.ScrollUp(v.size.Height())
}

// PageDown returns a new Viewport scrolled down by the viewport height.
func (v *Viewport) PageDown() *Viewport {
	return v.ScrollDown(v.size.Height())
}

// VisibleLines returns the currently visible lines in the viewport.
// Lines are truncated or wrapped based on the wrapLines setting.
func (v *Viewport) VisibleLines() []string {
	visible := v.scrollSvc.VisibleLines(v.content, v.scrollOffset.Offset(), v.size.Height())

	if v.wrapLines {
		return v.wrapVisibleLines(visible)
	}

	return v.truncateVisibleLines(visible)
}

// ScrollOffset returns the current scroll offset.
func (v *Viewport) ScrollOffset() int {
	return v.scrollOffset.Offset()
}

// CanScrollUp returns true if the viewport can scroll up.
func (v *Viewport) CanScrollUp() bool {
	return v.scrollSvc.CanScrollUp(v.scrollOffset.Offset())
}

// CanScrollDown returns true if the viewport can scroll down.
func (v *Viewport) CanScrollDown() bool {
	return v.scrollSvc.CanScrollDown(v.scrollOffset.Offset(), len(v.content), v.size.Height())
}

// IsAtTop returns true if the viewport is at the top.
func (v *Viewport) IsAtTop() bool {
	return v.scrollSvc.IsAtTop(v.scrollOffset.Offset())
}

// IsAtBottom returns true if the viewport is at the bottom.
func (v *Viewport) IsAtBottom() bool {
	return v.scrollSvc.IsAtBottom(v.scrollOffset.Offset(), len(v.content), v.size.Height())
}

// TotalLines returns the total number of content lines.
func (v *Viewport) TotalLines() int {
	return len(v.content)
}

// VisibleHeight returns the viewport height.
func (v *Viewport) VisibleHeight() int {
	return v.size.Height()
}

// FollowMode returns true if follow mode is enabled.
func (v *Viewport) FollowMode() bool {
	return v.followMode
}

// WrapLines returns true if line wrapping is enabled.
func (v *Viewport) WrapLines() bool {
	return v.wrapLines
}

// Content returns a defensive copy of the viewport content.
// This allows API layer to implement operations like AppendLine().
// Returns a new slice to maintain encapsulation (domain remains in control).
func (v *Viewport) Content() []string {
	contentCopy := make([]string, len(v.content))
	copy(contentCopy, v.content)
	return contentCopy
}

// clone creates a shallow copy of the viewport for immutability.
func (v *Viewport) clone() *Viewport {
	return &Viewport{
		content:      v.content,
		size:         v.size,
		scrollOffset: v.scrollOffset,
		followMode:   v.followMode,
		wrapLines:    v.wrapLines,
		scrollSvc:    v.scrollSvc,
	}
}

// truncateVisibleLines truncates lines that exceed the viewport width.
func (v *Viewport) truncateVisibleLines(lines []string) []string {
	if v.size.Width() <= 0 {
		return []string{}
	}

	result := make([]string, len(lines))
	for i, line := range lines {
		result[i] = v.truncateLine(line, v.size.Width())
	}
	return result
}

// truncateLine truncates a single line to fit within the given width.
// Uses Unicode-aware width calculation.
func (v *Viewport) truncateLine(line string, width int) string {
	if width <= 0 {
		return ""
	}

	currentWidth := 0
	var result strings.Builder

	graphemes := uniseg.NewGraphemes(line)
	for graphemes.Next() {
		g := graphemes.Str()
		gWidth := uniseg.StringWidth(g)

		if currentWidth+gWidth > width {
			break
		}

		result.WriteString(g)
		currentWidth += gWidth
	}

	return result.String()
}

// wrapVisibleLines wraps lines that exceed the viewport width.
// This is a simple implementation - each line that's too wide is split.
func (v *Viewport) wrapVisibleLines(lines []string) []string {
	if v.size.Width() <= 0 {
		return []string{}
	}

	var result []string
	for _, line := range lines {
		wrapped := v.wrapLine(line, v.size.Width())
		result = append(result, wrapped...)
	}
	return result
}

// wrapLine wraps a single line to fit within the given width.
// Returns multiple lines if wrapping is needed.
func (v *Viewport) wrapLine(line string, width int) []string {
	if width <= 0 {
		return []string{}
	}

	lineWidth := uniseg.StringWidth(line)
	if lineWidth <= width {
		return []string{line}
	}

	var result []string
	currentWidth := 0
	var currentLine strings.Builder

	graphemes := uniseg.NewGraphemes(line)
	for graphemes.Next() {
		g := graphemes.Str()
		gWidth := uniseg.StringWidth(g)

		if currentWidth+gWidth > width && currentWidth > 0 {
			// Start a new line.
			result = append(result, currentLine.String())
			currentLine.Reset()
			currentWidth = 0
		}

		currentLine.WriteString(g)
		currentWidth += gWidth
	}

	// Add the last line.
	if currentLine.Len() > 0 {
		result = append(result, currentLine.String())
	}

	return result
}
