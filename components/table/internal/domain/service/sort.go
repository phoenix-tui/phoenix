// Package service provides domain services for the table component.
package service

import (
	"fmt"
	"sort"

	"github.com/phoenix-tui/phoenix/components/table/internal/domain/model"
	"github.com/phoenix-tui/phoenix/components/table/internal/domain/value"
)

// SortService handles table sorting logic.
// This is a domain service because sorting behavior is part of the business logic.
// but doesn't belong to a single entity.
type SortService struct{}

// NewSortService creates a new sort service.
func NewSortService() *SortService {
	return &SortService{}
}

// Sort sorts the rows by the specified column and direction.
// Returns a new sorted slice (original rows are not modified).
func (s *SortService) Sort(rows []model.Row, columnKey string, direction value.SortDirection) []model.Row {
	if direction.IsNone() || len(rows) == 0 {
		return rows
	}

	// Create a copy to avoid modifying original.
	sorted := make([]model.Row, len(rows))
	copy(sorted, rows)

	// Sort using stable sort.
	sort.SliceStable(sorted, func(i, j int) bool {
		valI := sorted[i][columnKey]
		valJ := sorted[j][columnKey]

		cmp := s.Compare(valI, valJ)

		if direction.IsAscending() {
			return cmp < 0
		}
		return cmp > 0
	})

	return sorted
}

// Compare compares two cell values.
// Supports string, int, int64, float64, and bool types.
// Returns:
//
//	-1 if a < b.
//	 0 if a == b.
//	 1 if a > b.
func (s *SortService) Compare(a, b interface{}) int {
	// Handle nil values.
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}

	// Compare by type.
	switch aVal := a.(type) {
	case string:
		if bVal, ok := b.(string); ok {
			return compareStrings(aVal, bVal)
		}

	case int:
		if bVal, ok := b.(int); ok {
			return compareInts(aVal, bVal)
		}

	case int64:
		if bVal, ok := b.(int64); ok {
			return compareInt64s(aVal, bVal)
		}

	case float64:
		if bVal, ok := b.(float64); ok {
			return compareFloat64s(aVal, bVal)
		}

	case bool:
		if bVal, ok := b.(bool); ok {
			return compareBools(aVal, bVal)
		}
	}

	// Fallback: compare as strings.
	return compareStrings(fmt.Sprintf("%v", a), fmt.Sprintf("%v", b))
}

// compareStrings compares two strings lexicographically.
func compareStrings(a, b string) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// compareInts compares two integers.
func compareInts(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// compareInt64s compares two int64 values.
func compareInt64s(a, b int64) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// compareFloat64s compares two float64 values.
func compareFloat64s(a, b float64) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// compareBools compares two boolean values (false < true).
func compareBools(a, b bool) int {
	if a == b {
		return 0
	}
	if !a && b {
		return -1
	}
	return 1
}
