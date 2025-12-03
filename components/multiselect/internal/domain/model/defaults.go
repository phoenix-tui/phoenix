// Package model provides the domain model for the multiselect component.
package model

import (
	"strings"

	"github.com/phoenix-tui/phoenix/components/multiselect/internal/domain/value"
)

// defaultFilter implements fuzzy filtering for options.
func defaultFilter[T any](opt *value.Option[T], query string) bool {
	if query == "" {
		return true
	}

	// Simple case-insensitive substring match
	label := strings.ToLower(opt.Label())
	queryLower := strings.ToLower(query)

	return strings.Contains(label, queryLower)
}

// defaultRender renders an option with checkbox and focus indicator.
func defaultRender[T any](opt *value.Option[T], _ int, focused, selected bool) string {
	var b strings.Builder

	// Cursor indicator
	if focused {
		b.WriteString("> ")
	} else {
		b.WriteString("  ")
	}

	// Checkbox
	if selected {
		b.WriteString("[x] ")
	} else {
		b.WriteString("[ ] ")
	}

	// Label
	b.WriteString(opt.Label())

	// Description (if present)
	if opt.Description() != "" {
		b.WriteString(" - ")
		b.WriteString(opt.Description())
	}

	// Disabled indicator
	if opt.Disabled() {
		b.WriteString(" (disabled)")
	}

	return b.String()
}
