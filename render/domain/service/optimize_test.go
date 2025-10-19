package service

import (
	"testing"

	"github.com/phoenix-tui/phoenix/render/domain/model"
	"github.com/phoenix-tui/phoenix/render/domain/value"
	"github.com/stretchr/testify/assert"
)

func TestNewOptimizeService(t *testing.T) {
	svc := NewOptimizeService()
	assert.NotNil(t, svc)
}

func TestOptimizeService_BatchStyles(t *testing.T) {
	svc := NewOptimizeService()

	style := value.NewStyleWithFg(value.ColorRed)
	ops := []DiffOp{
		{Type: OpTypeSet, Position: value.NewPosition(0, 0), Cell: model.NewCell('A', style)},
		{Type: OpTypeSet, Position: value.NewPosition(1, 0), Cell: model.NewCell('B', style)},
		{Type: OpTypeSet, Position: value.NewPosition(2, 0), Cell: model.NewCell('C', style)},
	}

	result := svc.BatchStyles(ops)
	assert.NotNil(t, result)
	// Currently pass-through
	assert.Equal(t, len(ops), len(result))
}

func TestOptimizeService_RemoveRedundant(t *testing.T) {
	svc := NewOptimizeService()

	style1 := value.NewStyleWithFg(value.ColorRed)
	style2 := value.NewStyleWithFg(value.ColorBlue)

	ops := []DiffOp{
		{Type: OpTypeSet, Position: value.NewPosition(0, 0), Cell: model.NewCell('A', style1)},
		{Type: OpTypeSet, Position: value.NewPosition(1, 0), Cell: model.NewCell('B', style1)},
		{Type: OpTypeSet, Position: value.NewPosition(2, 0), Cell: model.NewCell('C', style2)},
	}

	currentStyle := value.NewStyle()
	result := svc.RemoveRedundant(ops, currentStyle)

	assert.NotNil(t, result)
	assert.Equal(t, len(ops), len(result))
}

func TestOptimizeService_RemoveRedundant_EmptyOps(t *testing.T) {
	svc := NewOptimizeService()
	var ops []DiffOp

	result := svc.RemoveRedundant(ops, value.NewStyle())
	assert.Empty(t, result)
}

func TestOptimizeService_OptimizeCursorMoves(t *testing.T) {
	svc := NewOptimizeService()

	ops := []DiffOp{
		{Type: OpTypeMoveCursor, Position: value.NewPosition(0, 0)},
		{Type: OpTypeSet, Position: value.NewPosition(0, 0), Cell: model.NewCell('A', value.NewStyle())},
		{Type: OpTypeMoveCursor, Position: value.NewPosition(1, 0)},
		{Type: OpTypeSet, Position: value.NewPosition(1, 0), Cell: model.NewCell('B', value.NewStyle())},
	}

	result := svc.OptimizeCursorMoves(ops)
	assert.NotNil(t, result)
	// Currently pass-through
	assert.Equal(t, len(ops), len(result))
}

func TestOptimizeService_MergeAdjacentOps(t *testing.T) {
	svc := NewOptimizeService()

	style := value.NewStyleWithFg(value.ColorRed)

	tests := []struct {
		name     string
		ops      []DiffOp
		expected int // Expected number of operations after merge
	}{
		{
			"no adjacent ops",
			[]DiffOp{
				{Type: OpTypeSet, Position: value.NewPosition(0, 0), Cell: model.NewCell('A', style)},
				{Type: OpTypeSet, Position: value.NewPosition(5, 0), Cell: model.NewCell('B', style)},
			},
			2,
		},
		{
			"adjacent same style",
			[]DiffOp{
				{Type: OpTypeSet, Position: value.NewPosition(0, 0), Cell: model.NewCell('A', style)},
				{Type: OpTypeSet, Position: value.NewPosition(1, 0), Cell: model.NewCell('B', style)},
				{Type: OpTypeSet, Position: value.NewPosition(2, 0), Cell: model.NewCell('C', style)},
			},
			3, // Currently no merging implemented
		},
		{
			"empty ops",
			[]DiffOp{},
			0,
		},
		{
			"single op",
			[]DiffOp{
				{Type: OpTypeSet, Position: value.NewPosition(0, 0), Cell: model.NewCell('A', style)},
			},
			1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.MergeAdjacentOps(tt.ops)
			assert.Equal(t, tt.expected, len(result))
		})
	}
}

func TestOptimizeService_CountOperations(t *testing.T) {
	svc := NewOptimizeService()

	tests := []struct {
		name     string
		ops      []DiffOp
		expected int
	}{
		{"empty", []DiffOp{}, 0},
		{"single", []DiffOp{{Type: OpTypeSet}}, 1},
		{"multiple", []DiffOp{{Type: OpTypeSet}, {Type: OpTypeSet}, {Type: OpTypeSet}}, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := svc.CountOperations(tt.ops)
			assert.Equal(t, tt.expected, count)
		})
	}
}

func TestOptimizeService_EstimateOutputSize(t *testing.T) {
	svc := NewOptimizeService()

	tests := []struct {
		name string
		ops  []DiffOp
	}{
		{"empty", []DiffOp{}},
		{
			"set operations",
			[]DiffOp{
				{Type: OpTypeSet, Position: value.NewPosition(0, 0), Cell: model.NewCell('A', value.NewStyle())},
				{Type: OpTypeSet, Position: value.NewPosition(1, 0), Cell: model.NewCell('B', value.NewStyle())},
			},
		},
		{
			"move operations",
			[]DiffOp{
				{Type: OpTypeMoveCursor, Position: value.NewPosition(0, 0)},
				{Type: OpTypeMoveCursor, Position: value.NewPosition(10, 5)},
			},
		},
		{
			"clear operations",
			[]DiffOp{
				{Type: OpTypeClear, Position: value.NewPosition(0, 0)},
			},
		},
		{
			"mixed operations",
			[]DiffOp{
				{Type: OpTypeSet, Position: value.NewPosition(0, 0), Cell: model.NewCell('A', value.NewStyle())},
				{Type: OpTypeMoveCursor, Position: value.NewPosition(1, 0)},
				{Type: OpTypeClear, Position: value.NewPosition(2, 0)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			size := svc.EstimateOutputSize(tt.ops)
			if len(tt.ops) == 0 {
				assert.Equal(t, 0, size)
			} else {
				assert.Greater(t, size, 0)
			}
		})
	}
}

func TestOptimizeService_ShouldFullRedraw(t *testing.T) {
	svc := NewOptimizeService()

	tests := []struct {
		name         string
		changedCells int
		totalCells   int
		expected     bool
	}{
		{"no changes", 0, 100, false},
		{"small changes", 10, 100, false},
		{"half changed", 50, 100, false},
		{"most changed", 80, 100, true},
		{"all changed", 100, 100, true},
		{"zero total", 0, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.ShouldFullRedraw(tt.changedCells, tt.totalCells)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestOptimizeService_GroupByLine(t *testing.T) {
	svc := NewOptimizeService()

	ops := []DiffOp{
		{Type: OpTypeSet, Position: value.NewPosition(0, 0), Cell: model.NewCell('A', value.NewStyle())},
		{Type: OpTypeSet, Position: value.NewPosition(1, 0), Cell: model.NewCell('B', value.NewStyle())},
		{Type: OpTypeSet, Position: value.NewPosition(0, 1), Cell: model.NewCell('C', value.NewStyle())},
		{Type: OpTypeSet, Position: value.NewPosition(0, 2), Cell: model.NewCell('D', value.NewStyle())},
		{Type: OpTypeSet, Position: value.NewPosition(1, 2), Cell: model.NewCell('E', value.NewStyle())},
	}

	groups := svc.GroupByLine(ops)

	assert.Len(t, groups, 3)    // 3 lines
	assert.Len(t, groups[0], 2) // Line 0: 2 ops
	assert.Len(t, groups[1], 1) // Line 1: 1 op
	assert.Len(t, groups[2], 2) // Line 2: 2 ops
}

func TestOptimizeService_GroupByLine_EmptyOps(t *testing.T) {
	svc := NewOptimizeService()
	var ops []DiffOp

	groups := svc.GroupByLine(ops)
	assert.Empty(t, groups)
}

// Benchmarks
func BenchmarkOptimizeService_RemoveRedundant(b *testing.B) {
	svc := NewOptimizeService()
	style := value.NewStyleWithFg(value.ColorRed)

	ops := make([]DiffOp, 100)
	for i := 0; i < 100; i++ {
		ops[i] = DiffOp{
			Type:     OpTypeSet,
			Position: value.NewPosition(i, 0),
			Cell:     model.NewCell('A', style),
		}
	}

	currentStyle := value.NewStyle()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = svc.RemoveRedundant(ops, currentStyle)
	}
}

func BenchmarkOptimizeService_MergeAdjacentOps(b *testing.B) {
	svc := NewOptimizeService()
	style := value.NewStyleWithFg(value.ColorRed)

	ops := make([]DiffOp, 100)
	for i := 0; i < 100; i++ {
		ops[i] = DiffOp{
			Type:     OpTypeSet,
			Position: value.NewPosition(i, 0),
			Cell:     model.NewCell('A', style),
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = svc.MergeAdjacentOps(ops)
	}
}

func BenchmarkOptimizeService_GroupByLine(b *testing.B) {
	svc := NewOptimizeService()
	style := value.NewStyleWithFg(value.ColorRed)

	ops := make([]DiffOp, 1000)
	for i := 0; i < 1000; i++ {
		ops[i] = DiffOp{
			Type:     OpTypeSet,
			Position: value.NewPosition(i%80, i/80),
			Cell:     model.NewCell('A', style),
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = svc.GroupByLine(ops)
	}
}
