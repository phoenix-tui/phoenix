// Package value provides value objects for the multiselect component domain.
package value

// Selection represents the selection state for a multi-choice selection.
// It's a value object that tracks which indices are selected and enforces min/max constraints.
type Selection struct {
	indices map[int]bool // Selected indices
	min     int          // Minimum required selections (0 = no minimum)
	max     int          // Maximum allowed selections (0 = unlimited)
}

// NewSelection creates a new selection state with optional min/max constraints.
func NewSelection(minCount, maxCount int) *Selection {
	if minCount < 0 {
		minCount = 0
	}
	if maxCount < 0 {
		maxCount = 0
	}
	// Ensure min <= max (if max is set)
	if maxCount > 0 && minCount > maxCount {
		minCount = maxCount
	}

	return &Selection{
		indices: make(map[int]bool),
		min:     minCount,
		max:     maxCount,
	}
}

// WithSelected returns a new Selection with the specified indices selected.
func (s *Selection) WithSelected(indices ...int) *Selection {
	newIndices := make(map[int]bool, len(s.indices)+len(indices))
	for k, v := range s.indices {
		newIndices[k] = v
	}

	for _, idx := range indices {
		if idx >= 0 {
			// Only add if max constraint is not violated
			if s.max == 0 || len(newIndices) < s.max {
				newIndices[idx] = true
			}
		}
	}

	return &Selection{
		indices: newIndices,
		min:     s.min,
		max:     s.max,
	}
}

// Toggle returns a new Selection with the specified index toggled.
func (s *Selection) Toggle(index int) *Selection {
	if index < 0 {
		return s
	}

	newIndices := make(map[int]bool, len(s.indices))
	for k, v := range s.indices {
		newIndices[k] = v
	}

	if newIndices[index] {
		// Deselecting - always allowed
		delete(newIndices, index)
	} else if s.max == 0 || len(newIndices) < s.max {
		newIndices[index] = true
	}

	return &Selection{
		indices: newIndices,
		min:     s.min,
		max:     s.max,
	}
}

// SelectAll returns a new Selection with all indices (up to maxIndex) selected.
// Respects max constraint if set.
func (s *Selection) SelectAll(maxIndex int) *Selection {
	if maxIndex < 0 {
		return s
	}

	newIndices := make(map[int]bool)

	// If max is set, only select up to max count
	limit := maxIndex + 1
	if s.max > 0 && s.max < limit {
		limit = s.max
	}

	for i := 0; i < limit; i++ {
		newIndices[i] = true
	}

	return &Selection{
		indices: newIndices,
		min:     s.min,
		max:     s.max,
	}
}

// Clear returns a new Selection with all selections cleared.
func (s *Selection) Clear() *Selection {
	return &Selection{
		indices: make(map[int]bool),
		min:     s.min,
		max:     s.max,
	}
}

// IsSelected returns true if the specified index is selected.
func (s *Selection) IsSelected(index int) bool {
	return s.indices[index]
}

// Count returns the number of selected items.
func (s *Selection) Count() int {
	return len(s.indices)
}

// Indices returns a slice of selected indices in ascending order.
func (s *Selection) Indices() []int {
	if len(s.indices) == 0 {
		return []int{}
	}

	result := make([]int, 0, len(s.indices))
	for idx := range s.indices {
		result = append(result, idx)
	}

	// Sort in ascending order
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i] > result[j] {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result
}

// CanSelect returns true if another item can be selected (max not reached).
func (s *Selection) CanSelect() bool {
	return s.max == 0 || len(s.indices) < s.max
}

// CanConfirm returns true if enough items are selected (min reached).
func (s *Selection) CanConfirm() bool {
	return len(s.indices) >= s.min
}

// Min returns the minimum required selections.
func (s *Selection) Min() int {
	return s.min
}

// Max returns the maximum allowed selections (0 = unlimited).
func (s *Selection) Max() int {
	return s.max
}
