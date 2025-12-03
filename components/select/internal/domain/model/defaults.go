package model

import (
	"fmt"
	"strings"

	"github.com/phoenix-tui/phoenix/components/select/internal/domain/value"
)

// defaultFilter provides basic substring matching filter.
func defaultFilter[T any](opt *value.Option[T], query string) bool {
	if query == "" {
		return true
	}
	return strings.Contains(
		strings.ToLower(opt.Label()),
		strings.ToLower(query),
	)
}

// defaultRender provides basic option rendering with cursor indicator.
func defaultRender[T any](opt *value.Option[T], _ int, focused bool) string {
	cursor := "  "
	if focused {
		cursor = "> "
	}

	label := opt.Label()
	if opt.Disabled() {
		label = fmt.Sprintf("%s (disabled)", label)
	}

	result := cursor + label

	if opt.Description() != "" {
		result = fmt.Sprintf("%s - %s", result, opt.Description())
	}

	return result
}
