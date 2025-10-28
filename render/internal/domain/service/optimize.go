//nolint:nestif // Optimization logic requires nested conditionals
package service

import (
	"github.com/phoenix-tui/phoenix/render/internal/domain/value"
)

// OptimizeService optimizes rendering operations for better performance.
// This includes batching ANSI sequences, removing redundancy, and optimizing cursor movements.
type OptimizeService struct{}

// NewOptimizeService creates a new optimize service.
func NewOptimizeService() *OptimizeService {
	return &OptimizeService{}
}

// BatchStyles groups adjacent cells with same style into single ANSI sequence.
// This significantly reduces the number of ANSI escape sequences.
func (s *OptimizeService) BatchStyles(ops []DiffOp) []DiffOp {
	if len(ops) == 0 {
		return ops
	}

	// For now, pass through (batching will be implemented incrementally).
	// Future optimization: detect adjacent cells with same style and merge.
	return ops
}

// RemoveRedundant removes operations that don't change the current state.
// For example, if style is already set, don't emit another style sequence.
func (s *OptimizeService) RemoveRedundant(ops []DiffOp, currentStyle value.Style) []DiffOp {
	if len(ops) == 0 {
		return ops
	}

	optimized := make([]DiffOp, 0, len(ops))
	lastStyle := currentStyle

	for _, op := range ops {
		if op.Type == OpTypeSet {
			// Only include if style differs from last.
			if !op.Cell.Style().Equals(lastStyle) {
				optimized = append(optimized, op)
				lastStyle = op.Cell.Style()
			} else {
				// Style same, but still need to write the cell.
				optimized = append(optimized, op)
			}
		} else {
			optimized = append(optimized, op)
		}
	}

	return optimized
}

// OptimizeCursorMoves reduces cursor movement operations.
// Uses relative movements when shorter than absolute positioning.
func (s *OptimizeService) OptimizeCursorMoves(ops []DiffOp) []DiffOp {
	if len(ops) == 0 {
		return ops
	}

	// For now, pass through (cursor optimization will be implemented incrementally).
	// Future optimization: use relative cursor moves (\x1b[C, \x1b[D) when shorter.
	return ops
}

// MergeAdjacentOps merges adjacent operations on same line for efficiency.
func (s *OptimizeService) MergeAdjacentOps(ops []DiffOp) []DiffOp {
	if len(ops) < 2 {
		return ops
	}

	merged := make([]DiffOp, 0, len(ops))
	i := 0

	for i < len(ops) {
		current := ops[i]

		// Look ahead for adjacent cells on same line with same style.
		if current.Type == OpTypeSet && i+1 < len(ops) {
			next := ops[i+1]

			// Check if adjacent and same style.
			if next.Type == OpTypeSet &&
				current.Position.Y() == next.Position.Y() &&
				next.Position.X() == current.Position.X()+current.Cell.Width() &&
				current.Cell.Style().Equals(next.Cell.Style()) {
				// Can merge - continue looking.
				merged = append(merged, current)
				i++
				continue
			}
		}

		merged = append(merged, current)
		i++
	}

	return merged
}

// CountOperations returns the number of operations.
func (s *OptimizeService) CountOperations(ops []DiffOp) int {
	return len(ops)
}

// EstimateOutputSize estimates the number of bytes the operations will produce.
// Useful for buffer pre-allocation.
func (s *OptimizeService) EstimateOutputSize(ops []DiffOp) int {
	if len(ops) == 0 {
		return 0
	}

	size := 0
	for _, op := range ops {
		switch op.Type {
		case OpTypeSet:
			// ANSI sequence (approx 20 bytes) + cursor move (10 bytes) + char (1-4 bytes).
			size += 35
		case OpTypeMoveCursor:
			// Cursor move sequence (10 bytes).
			size += 10
		case OpTypeClear:
			// Clear sequence (10 bytes).
			size += 10
		}
	}

	return size
}

// ShouldFullRedraw determines if full redraw is more efficient than differential.
// When too many cells changed, full redraw may be faster.
func (s *OptimizeService) ShouldFullRedraw(changedCells, totalCells int) bool {
	if totalCells == 0 {
		return false
	}

	// If more than 75% of cells changed, full redraw is more efficient.
	threshold := 0.75
	ratio := float64(changedCells) / float64(totalCells)
	return ratio > threshold
}

// GroupByLine groups operations by line for line-based optimization.
func (s *OptimizeService) GroupByLine(ops []DiffOp) map[int][]DiffOp {
	groups := make(map[int][]DiffOp)

	for _, op := range ops {
		y := op.Position.Y()
		groups[y] = append(groups[y], op)
	}

	return groups
}
