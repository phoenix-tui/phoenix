// Package service provides domain services for rendering optimization.
package service

import (
	model2 "github.com/phoenix-tui/phoenix/render/internal/domain/model"
	"github.com/phoenix-tui/phoenix/render/internal/domain/value"
)

// OpType represents the type of diff operation.
type OpType int

const (
	// OpTypeSet sets a cell at a position.
	OpTypeSet OpType = iota
	// OpTypeClear clears a region (optimization for clearing multiple cells).
	OpTypeClear
	// OpTypeMoveCursor moves cursor without writing (optimization).
	OpTypeMoveCursor
)

// DiffOp represents a single diff operation.
type DiffOp struct {
	Type     OpType
	Position value.Position
	Cell     model2.Cell
}

// DiffService computes differences between buffers for minimal rendering.
type DiffService struct{}

// NewDiffService creates a new diff service.
func NewDiffService() *DiffService {
	return &DiffService{}
}

// Diff computes minimal set of operations to transform oldBuf to newBuf.
// This is the core performance optimization - only render changed cells.
func (s *DiffService) Diff(oldBuf, newBuf *model2.Buffer) []DiffOp {
	if oldBuf == nil || newBuf == nil {
		return nil
	}

	// Buffers must have same dimensions.
	if oldBuf.Width() != newBuf.Width() || oldBuf.Height() != newBuf.Height() {
		return s.diffFullBuffer(newBuf)
	}

	var ops []DiffOp

	// Compare cell by cell.
	for y := 0; y < newBuf.Height(); y++ {
		for x := 0; x < newBuf.Width(); x++ {
			pos := value.NewPosition(x, y)
			oldCell := oldBuf.Get(pos)
			newCell := newBuf.Get(pos)

			// Only emit operation if cells differ.
			if !oldCell.Equals(newCell) {
				ops = append(ops, DiffOp{
					Type:     OpTypeSet,
					Position: pos,
					Cell:     newCell,
				})
			}
		}
	}

	return ops
}

// diffFullBuffer creates ops to render entire buffer (for initial render or resize).
func (s *DiffService) diffFullBuffer(buf *model2.Buffer) []DiffOp {
	var ops []DiffOp

	for y := 0; y < buf.Height(); y++ {
		for x := 0; x < buf.Width(); x++ {
			pos := value.NewPosition(x, y)
			cell := buf.Get(pos)

			// Skip empty cells to optimize.
			if !cell.IsEmpty() {
				ops = append(ops, DiffOp{
					Type:     OpTypeSet,
					Position: pos,
					Cell:     cell,
				})
			}
		}
	}

	return ops
}

// Optimize merges adjacent operations for better performance.
// For example, adjacent cells with same style can be batched.
func (s *DiffService) Optimize(ops []DiffOp) []DiffOp {
	if len(ops) == 0 {
		return ops
	}

	optimized := make([]DiffOp, 0, len(ops))

	// For now, pass through (optimization will be added incrementally).
	// Future optimizations:.
	// 1. Merge adjacent cells with same style
	// 2. Detect clear regions (multiple empty cells)
	// 3. Optimize cursor movements (relative vs absolute)
	optimized = append(optimized, ops...)

	return optimized
}

// DiffLine compares two lines and returns operations.
// This is useful for line-based rendering optimization.
func (s *DiffService) DiffLine(oldLine, newLine []model2.Cell, y int) []DiffOp {
	if len(oldLine) != len(newLine) {
		return nil
	}

	var ops []DiffOp

	for x := 0; x < len(newLine); x++ {
		if !oldLine[x].Equals(newLine[x]) {
			ops = append(ops, DiffOp{
				Type:     OpTypeSet,
				Position: value.NewPosition(x, y),
				Cell:     newLine[x],
			})
		}
	}

	return ops
}

// CountChanges returns the number of changed cells.
// Useful for metrics and debugging.
func (s *DiffService) CountChanges(oldBuf, newBuf *model2.Buffer) int {
	ops := s.Diff(oldBuf, newBuf)
	return len(ops)
}

// HasChanges returns true if any cells differ.
// Faster than computing full diff when you only need to know if there are changes.
func (s *DiffService) HasChanges(oldBuf, newBuf *model2.Buffer) bool {
	if oldBuf == nil || newBuf == nil {
		return true
	}

	if oldBuf.Width() != newBuf.Width() || oldBuf.Height() != newBuf.Height() {
		return true
	}

	for y := 0; y < newBuf.Height(); y++ {
		for x := 0; x < newBuf.Width(); x++ {
			pos := value.NewPosition(x, y)
			if !oldBuf.Get(pos).Equals(newBuf.Get(pos)) {
				return true
			}
		}
	}

	return false
}
