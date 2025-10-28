// Package model provides domain models for the table component.
package model

import (
	"github.com/phoenix-tui/phoenix/components/table/internal/domain/value"
)

// Column represents a table column definition.
// It's an immutable value object that defines how a column should be displayed.
type Column struct {
	key       string                   // Unique identifier for this column
	title     string                   // Display title in header
	width     int                      // Column width in characters
	alignment value.Alignment          // Cell alignment (left/center/right)
	sortable  bool                     // Can this column be sorted?
	renderer  func(interface{}) string // Custom cell renderer (optional)
}

// NewColumn creates a new column with left alignment and no custom renderer.
func NewColumn(key, title string, width int) *Column {
	return &Column{
		key:       key,
		title:     title,
		width:     width,
		alignment: value.AlignmentLeft,
		sortable:  false,
		renderer:  nil,
	}
}

// NewColumnWithAlignment creates a new column with specified alignment.
func NewColumnWithAlignment(key, title string, width int, alignment value.Alignment) *Column {
	return &Column{
		key:       key,
		title:     title,
		width:     width,
		alignment: alignment,
		sortable:  false,
		renderer:  nil,
	}
}

// WithWidth returns a new column with the specified width.
func (c *Column) WithWidth(width int) *Column {
	return &Column{
		key:       c.key,
		title:     c.title,
		width:     width,
		alignment: c.alignment,
		sortable:  c.sortable,
		renderer:  c.renderer,
	}
}

// WithAlignment returns a new column with the specified alignment.
func (c *Column) WithAlignment(alignment value.Alignment) *Column {
	return &Column{
		key:       c.key,
		title:     c.title,
		width:     c.width,
		alignment: alignment,
		sortable:  c.sortable,
		renderer:  c.renderer,
	}
}

// WithSortable returns a new column with sortable flag set.
func (c *Column) WithSortable(sortable bool) *Column {
	return &Column{
		key:       c.key,
		title:     c.title,
		width:     c.width,
		alignment: c.alignment,
		sortable:  sortable,
		renderer:  c.renderer,
	}
}

// WithRenderer returns a new column with a custom cell renderer.
func (c *Column) WithRenderer(renderer func(interface{}) string) *Column {
	return &Column{
		key:       c.key,
		title:     c.title,
		width:     c.width,
		alignment: c.alignment,
		sortable:  c.sortable,
		renderer:  renderer,
	}
}

// Key returns the column's unique identifier.
func (c *Column) Key() string {
	return c.key
}

// Title returns the column's display title.
func (c *Column) Title() string {
	return c.title
}

// Width returns the column's width in characters.
func (c *Column) Width() int {
	return c.width
}

// Alignment returns the column's cell alignment.
func (c *Column) Alignment() value.Alignment {
	return c.alignment
}

// IsSortable returns whether this column can be sorted.
func (c *Column) IsSortable() bool {
	return c.sortable
}

// Renderer returns the custom cell renderer, or nil if none is set.
func (c *Column) Renderer() func(interface{}) string {
	return c.renderer
}
