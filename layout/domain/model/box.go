// Package model provides rich domain models for layout system.
// These models encapsulate business logic and maintain invariants.
package model

import (
	"fmt"
	"strings"

	"github.com/phoenix-tui/phoenix/layout/domain/value"
)

// Box represents a layout box following the CSS box model.
//
// Box Model Structure:
//
//	┌─────────────────────────────────┐ ← Margin (outside)
//	│  ┌─────────────────────────┐    │
//	│  │ Border                  │    │
//	│  │  ┌─────────────────┐    │    │
//	│  │  │ Padding         │    │    │
//	│  │  │  ┌───────────┐  │    │    │
//	│  │  │  │  Content  │  │    │    │
//	│  │  │  └───────────┘  │    │    │
//	│  │  └─────────────────┘    │    │
//	│  └─────────────────────────┘    │
//	└─────────────────────────────────┘
//
// Design Philosophy:
//   - Immutable aggregate root (rich domain model)
//   - Content-first design (content drives size)
//   - Fluent API for composability
//   - Size calculations respect box model layers
//   - Alignment determines positioning within parent
//
// Size Calculation Methods:
//   - ContentSize(): Inner content dimensions
//   - PaddedSize(): Content + padding
//   - BorderedSize(): Content + padding + border
//   - TotalSize(): Content + padding + border + margin
//
// Example:
//
//	box := NewBox("Hello, World!").
//		WithPadding(value.NewSpacingAll(1)).
//		WithBorder(true).
//		WithMargin(value.NewSpacingVH(1, 2))
//
//	totalSize := box.TotalSize() // Full outer size
type Box struct {
	content   string          // Content text (drives size)
	padding   value.Spacing   // Inner spacing (inside border)
	margin    value.Spacing   // Outer spacing (outside border)
	hasBorder bool            // Whether box has border
	size      value.Size      // Size constraints
	alignment value.Alignment // Alignment within parent
}

// NewBox creates a Box with the given content.
// Content cannot be empty (panics on empty string).
//
// Default values:
//   - Padding: Zero spacing
//   - Margin: Zero spacing
//   - Border: Disabled
//   - Size: Unconstrained
//   - Alignment: Top-left
//
// Example:
//
//	box := NewBox("Hello") // Simple box with defaults
func NewBox(content string) *Box {
	if content == "" {
		panic("box: content cannot be empty")
	}
	return &Box{
		content:   content,
		padding:   value.NewSpacingZero(),
		margin:    value.NewSpacingZero(),
		hasBorder: false,
		size:      value.NewSizeUnconstrained(),
		alignment: value.NewAlignmentDefault(),
	}
}

// Content returns the box content.
func (b *Box) Content() string {
	return b.content
}

// Padding returns the padding spacing.
func (b *Box) Padding() value.Spacing {
	return b.padding
}

// Margin returns the margin spacing.
func (b *Box) Margin() value.Spacing {
	return b.margin
}

// HasBorder returns true if box has a border.
func (b *Box) HasBorder() bool {
	return b.hasBorder
}

// Size returns the size constraints.
func (b *Box) Size() value.Size {
	return b.size
}

// Alignment returns the alignment within parent.
func (b *Box) Alignment() value.Alignment {
	return b.alignment
}

// WithContent returns a new Box with the given content.
// Panics if content is empty.
func (b *Box) WithContent(content string) *Box {
	if content == "" {
		panic("box: content cannot be empty")
	}
	result := *b
	result.content = content
	return &result
}

// WithPadding returns a new Box with the given padding.
// Padding is applied inside the border.
//
// Example:
//
//	box := NewBox("Text").WithPadding(value.NewSpacingAll(1))
func (b *Box) WithPadding(p value.Spacing) *Box {
	result := *b
	result.padding = p
	return &result
}

// WithMargin returns a new Box with the given margin.
// Margin is applied outside the border.
//
// Example:
//
//	box := NewBox("Text").WithMargin(value.NewSpacingVH(1, 2))
func (b *Box) WithMargin(m value.Spacing) *Box {
	result := *b
	result.margin = m
	return &result
}

// WithBorder returns a new Box with border enabled or disabled.
// Border adds 1 cell on each side when enabled.
//
// Example:
//
//	box := NewBox("Text").WithBorder(true)
func (b *Box) WithBorder(hasBorder bool) *Box {
	result := *b
	result.hasBorder = hasBorder
	return &result
}

// WithSize returns a new Box with the given size constraints.
// Size constraints are applied during layout pass.
//
// Example:
//
//	box := NewBox("Text").WithSize(value.NewSizeExact(80, 24))
func (b *Box) WithSize(s value.Size) *Box {
	result := *b
	result.size = s
	return &result
}

// WithAlignment returns a new Box with the given alignment.
// Alignment determines positioning within parent container.
//
// Example:
//
//	box := NewBox("Text").WithAlignment(value.NewAlignmentCenter())
func (b *Box) WithAlignment(a value.Alignment) *Box {
	result := *b
	result.alignment = a
	return &result
}

// ContentSize calculates the size of the content area.
// For now, this measures string length (simple approach).
// Later (Day 3), this will integrate with phoenix/core.UnicodeService
// for proper grapheme cluster width calculation.
//
// Current implementation:
//   - Width: Length of longest line
//   - Height: Number of lines
//
// Returns:
//   - Size with exact width and height
func (b *Box) ContentSize() value.Size {
	lines := strings.Split(b.content, "\n")

	// Calculate width (longest line)
	maxWidth := 0
	for _, line := range lines {
		// TODO(Day 3): Use phoenix/core.UnicodeService.Width(line)
		width := len(line)
		if width > maxWidth {
			maxWidth = width
		}
	}

	// Height is number of lines
	height := len(lines)

	return value.NewSizeExact(maxWidth, height)
}

// PaddedSize calculates size including padding.
// This is content size + padding on all sides.
//
// Example:
//
//	box := NewBox("Hi").WithPadding(value.NewSpacingAll(1))
//	size := box.PaddedSize() // Width: 2 + 2 = 4, Height: 1 + 2 = 3
//
// Returns:
//   - Size with width and height including padding
func (b *Box) PaddedSize() value.Size {
	content := b.ContentSize()

	// Add padding (both sides)
	width := content.Width() + b.padding.Horizontal()
	height := content.Height() + b.padding.Vertical()

	return value.NewSizeExact(width, height)
}

// BorderedSize calculates size including padding and border.
// Border adds 1 cell on each side when enabled.
//
// Example:
//
//	box := NewBox("Hi").
//		WithPadding(value.NewSpacingAll(1)).
//		WithBorder(true)
//	size := box.BorderedSize() // Width: 2 + 2 + 2 = 6, Height: 1 + 2 + 2 = 5
//
// Returns:
//   - Size with width and height including border
func (b *Box) BorderedSize() value.Size {
	padded := b.PaddedSize()

	width := padded.Width()
	height := padded.Height()

	// Add border (1 cell per side)
	if b.hasBorder {
		width += 2
		height += 2
	}

	return value.NewSizeExact(width, height)
}

// TotalSize calculates the total outer size of the box.
// This includes content + padding + border + margin.
// This is the full space the box occupies in its parent.
//
// Example:
//
//	box := NewBox("Hi").
//		WithPadding(value.NewSpacingAll(1)).
//		WithBorder(true).
//		WithMargin(value.NewSpacingVH(1, 2))
//	size := box.TotalSize() // Full outer size
//
// Returns:
//   - Size with full width and height including margin
func (b *Box) TotalSize() value.Size {
	bordered := b.BorderedSize()

	// Add margin (both sides)
	width := bordered.Width() + b.margin.Horizontal()
	height := bordered.Height() + b.margin.Vertical()

	return value.NewSizeExact(width, height)
}

// String returns a human-readable debug representation.
// Shows all box properties and calculated sizes.
func (b *Box) String() string {
	// Truncate content if too long
	contentPreview := b.content
	if len(contentPreview) > 30 {
		contentPreview = contentPreview[:27] + "..."
	}
	// Replace newlines for readability
	contentPreview = strings.ReplaceAll(contentPreview, "\n", "\\n")

	totalSize := b.TotalSize()

	var parts []string
	parts = append(parts, fmt.Sprintf("content=%q", contentPreview))

	if !b.padding.IsZero() {
		parts = append(parts, fmt.Sprintf("padding=%s", b.padding))
	}

	if !b.margin.IsZero() {
		parts = append(parts, fmt.Sprintf("margin=%s", b.margin))
	}

	if b.hasBorder {
		parts = append(parts, "border=true")
	}

	if !b.size.IsUnconstrained() {
		parts = append(parts, fmt.Sprintf("size=%s", b.size))
	}

	if !b.alignment.IsDefault() {
		parts = append(parts, fmt.Sprintf("align=%s", b.alignment))
	}

	parts = append(parts, fmt.Sprintf("total=%dx%d", totalSize.Width(), totalSize.Height()))

	return fmt.Sprintf("Box{%s}", strings.Join(parts, " "))
}
