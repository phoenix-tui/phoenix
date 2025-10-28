package service

import (
	"testing"

	model2 "github.com/phoenix-tui/phoenix/render/internal/domain/model"
	value2 "github.com/phoenix-tui/phoenix/render/internal/domain/value"
	"github.com/stretchr/testify/assert"
)

func TestNewDiffService(t *testing.T) {
	svc := NewDiffService()
	assert.NotNil(t, svc)
}

func TestDiffService_Diff_NoDifference(t *testing.T) {
	svc := NewDiffService()

	oldBuf := model2.NewBuffer(10, 5)
	newBuf := model2.NewBuffer(10, 5)

	ops := svc.Diff(oldBuf, newBuf)
	assert.Empty(t, ops)
}

func TestDiffService_Diff_SingleCellChange(t *testing.T) {
	svc := NewDiffService()

	oldBuf := model2.NewBuffer(10, 5)
	newBuf := model2.NewBuffer(10, 5)

	// Change one cell
	pos := value2.NewPosition(5, 2)
	cell := model2.NewCell('A', value2.NewStyleWithFg(value2.ColorRed))
	newBuf.Set(pos, cell)

	ops := svc.Diff(oldBuf, newBuf)

	assert.Len(t, ops, 1)
	assert.Equal(t, OpTypeSet, ops[0].Type)
	assert.True(t, pos.Equals(ops[0].Position))
	assert.True(t, cell.Equals(ops[0].Cell))
}

func TestDiffService_Diff_MultipleCellChanges(t *testing.T) {
	svc := NewDiffService()

	oldBuf := model2.NewBuffer(10, 5)
	newBuf := model2.NewBuffer(10, 5)

	// Change multiple cells
	style := value2.NewStyleWithFg(value2.ColorRed)
	positions := []value2.Position{
		value2.NewPosition(0, 0),
		value2.NewPosition(5, 2),
		value2.NewPosition(9, 4),
	}

	for _, pos := range positions {
		newBuf.Set(pos, model2.NewCell('X', style))
	}

	ops := svc.Diff(oldBuf, newBuf)

	assert.Len(t, ops, 3)
	for _, op := range ops {
		assert.Equal(t, OpTypeSet, op.Type)
		assert.Equal(t, 'X', op.Cell.Char())
	}
}

func TestDiffService_Diff_StringChange(t *testing.T) {
	svc := NewDiffService()

	oldBuf := model2.NewBuffer(20, 5)
	newBuf := model2.NewBuffer(20, 5)

	style := value2.NewStyleWithFg(value2.ColorRed)
	newBuf.SetString(value2.NewPosition(0, 0), "Hello", style)

	ops := svc.Diff(oldBuf, newBuf)

	assert.Len(t, ops, 5) // "Hello" = 5 chars
	assert.Equal(t, 'H', ops[0].Cell.Char())
	assert.Equal(t, 'e', ops[1].Cell.Char())
	assert.Equal(t, 'l', ops[2].Cell.Char())
	assert.Equal(t, 'l', ops[3].Cell.Char())
	assert.Equal(t, 'o', ops[4].Cell.Char())
}

func TestDiffService_Diff_StyleOnlyChange(t *testing.T) {
	svc := NewDiffService()

	oldBuf := model2.NewBuffer(10, 5)
	newBuf := model2.NewBuffer(10, 5)

	pos := value2.NewPosition(5, 2)
	style1 := value2.NewStyleWithFg(value2.ColorRed)
	style2 := value2.NewStyleWithFg(value2.ColorBlue)

	// Same char, different style
	oldBuf.Set(pos, model2.NewCell('A', style1))
	newBuf.Set(pos, model2.NewCell('A', style2))

	ops := svc.Diff(oldBuf, newBuf)

	assert.Len(t, ops, 1)
	assert.Equal(t, OpTypeSet, ops[0].Type)
	assert.True(t, style2.Equals(ops[0].Cell.Style()))
}

func TestDiffService_Diff_DifferentDimensions(t *testing.T) {
	svc := NewDiffService()

	oldBuf := model2.NewBuffer(10, 5)
	newBuf := model2.NewBuffer(20, 10)

	// Fill new buffer
	style := value2.NewStyleWithFg(value2.ColorRed)
	newBuf.SetString(value2.NewPosition(0, 0), "Test", style)

	ops := svc.Diff(oldBuf, newBuf)

	// Should return full buffer diff
	assert.NotEmpty(t, ops)
}

func TestDiffService_Diff_NilBuffers(t *testing.T) {
	svc := NewDiffService()

	tests := []struct {
		name   string
		oldBuf *model2.Buffer
		newBuf *model2.Buffer
	}{
		{"both nil", nil, nil},
		{"old nil", nil, model2.NewBuffer(10, 5)},
		{"new nil", model2.NewBuffer(10, 5), nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ops := svc.Diff(tt.oldBuf, tt.newBuf)
			assert.Nil(t, ops)
		})
	}
}

func TestDiffService_DiffLine(t *testing.T) {
	svc := NewDiffService()

	style1 := value2.NewStyleWithFg(value2.ColorRed)
	style2 := value2.NewStyleWithFg(value2.ColorBlue)

	oldLine := []model2.Cell{
		model2.NewCell('A', style1),
		model2.NewCell('B', style1),
		model2.NewCell('C', style1),
	}

	newLine := []model2.Cell{
		model2.NewCell('A', style1), // Same
		model2.NewCell('X', style2), // Changed
		model2.NewCell('C', style1), // Same
	}

	ops := svc.DiffLine(oldLine, newLine, 0)

	assert.Len(t, ops, 1)
	assert.Equal(t, OpTypeSet, ops[0].Type)
	assert.Equal(t, 1, ops[0].Position.X())
	assert.Equal(t, 'X', ops[0].Cell.Char())
}

func TestDiffService_DiffLine_DifferentLength(t *testing.T) {
	svc := NewDiffService()

	oldLine := []model2.Cell{
		model2.NewCell('A', value2.NewStyle()),
	}
	newLine := []model2.Cell{
		model2.NewCell('A', value2.NewStyle()),
		model2.NewCell('B', value2.NewStyle()),
	}

	ops := svc.DiffLine(oldLine, newLine, 0)
	assert.Nil(t, ops)
}

func TestDiffService_CountChanges(t *testing.T) {
	svc := NewDiffService()

	oldBuf := model2.NewBuffer(10, 5)
	newBuf := model2.NewBuffer(10, 5)

	// No changes
	count := svc.CountChanges(oldBuf, newBuf)
	assert.Equal(t, 0, count)

	// Add changes
	style := value2.NewStyleWithFg(value2.ColorRed)
	newBuf.SetString(value2.NewPosition(0, 0), "Test", style)

	count = svc.CountChanges(oldBuf, newBuf)
	assert.Equal(t, 4, count) // "Test" = 4 changes
}

func TestDiffService_HasChanges(t *testing.T) {
	svc := NewDiffService()

	oldBuf := model2.NewBuffer(10, 5)
	newBuf := model2.NewBuffer(10, 5)

	// No changes
	assert.False(t, svc.HasChanges(oldBuf, newBuf))

	// Add change
	newBuf.Set(value2.NewPosition(0, 0), model2.NewCell('A', value2.NewStyle()))
	assert.True(t, svc.HasChanges(oldBuf, newBuf))
}

func TestDiffService_HasChanges_DifferentDimensions(t *testing.T) {
	svc := NewDiffService()

	oldBuf := model2.NewBuffer(10, 5)
	newBuf := model2.NewBuffer(20, 10)

	assert.True(t, svc.HasChanges(oldBuf, newBuf))
}

func TestDiffService_HasChanges_NilBuffers(t *testing.T) {
	svc := NewDiffService()

	assert.True(t, svc.HasChanges(nil, nil))
	assert.True(t, svc.HasChanges(nil, model2.NewBuffer(10, 5)))
	assert.True(t, svc.HasChanges(model2.NewBuffer(10, 5), nil))
}

func TestDiffService_Optimize(t *testing.T) {
	svc := NewDiffService()

	ops := []DiffOp{
		{Type: OpTypeSet, Position: value2.NewPosition(0, 0), Cell: model2.NewCell('A', value2.NewStyle())},
		{Type: OpTypeSet, Position: value2.NewPosition(1, 0), Cell: model2.NewCell('B', value2.NewStyle())},
		{Type: OpTypeSet, Position: value2.NewPosition(2, 0), Cell: model2.NewCell('C', value2.NewStyle())},
	}

	optimized := svc.Optimize(ops)

	// Currently pass-through, but should not panic
	assert.NotNil(t, optimized)
	assert.Equal(t, len(ops), len(optimized))
}

// Benchmarks
func BenchmarkDiffService_Diff_NoChanges(b *testing.B) {
	svc := NewDiffService()
	oldBuf := model2.NewBuffer(80, 24)
	newBuf := model2.NewBuffer(80, 24)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = svc.Diff(oldBuf, newBuf)
	}
}

func BenchmarkDiffService_Diff_SmallChanges(b *testing.B) {
	svc := NewDiffService()
	oldBuf := model2.NewBuffer(80, 24)
	newBuf := model2.NewBuffer(80, 24)

	// Change 10 cells
	style := value2.NewStyleWithFg(value2.ColorRed)
	for i := 0; i < 10; i++ {
		newBuf.Set(value2.NewPosition(i, 0), model2.NewCell('X', style))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = svc.Diff(oldBuf, newBuf)
	}
}

func BenchmarkDiffService_Diff_FullScreen(b *testing.B) {
	svc := NewDiffService()
	oldBuf := model2.NewBuffer(80, 24)
	newBuf := model2.NewBuffer(80, 24)

	// Fill entire buffer
	style := value2.NewStyleWithFg(value2.ColorRed)
	newBuf.Fill('X', style)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = svc.Diff(oldBuf, newBuf)
	}
}

func BenchmarkDiffService_HasChanges(b *testing.B) {
	svc := NewDiffService()
	oldBuf := model2.NewBuffer(80, 24)
	newBuf := model2.NewBuffer(80, 24)
	newBuf.Set(value2.NewPosition(40, 12), model2.NewCell('X', value2.NewStyle()))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = svc.HasChanges(oldBuf, newBuf)
	}
}
