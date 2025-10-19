package service

import (
	"testing"

	"github.com/phoenix-tui/phoenix/render/domain/model"
	"github.com/phoenix-tui/phoenix/render/domain/value"
	"github.com/stretchr/testify/assert"
)

func TestNewDiffService(t *testing.T) {
	svc := NewDiffService()
	assert.NotNil(t, svc)
}

func TestDiffService_Diff_NoDifference(t *testing.T) {
	svc := NewDiffService()

	oldBuf := model.NewBuffer(10, 5)
	newBuf := model.NewBuffer(10, 5)

	ops := svc.Diff(oldBuf, newBuf)
	assert.Empty(t, ops)
}

func TestDiffService_Diff_SingleCellChange(t *testing.T) {
	svc := NewDiffService()

	oldBuf := model.NewBuffer(10, 5)
	newBuf := model.NewBuffer(10, 5)

	// Change one cell
	pos := value.NewPosition(5, 2)
	cell := model.NewCell('A', value.NewStyleWithFg(value.ColorRed))
	newBuf.Set(pos, cell)

	ops := svc.Diff(oldBuf, newBuf)

	assert.Len(t, ops, 1)
	assert.Equal(t, OpTypeSet, ops[0].Type)
	assert.True(t, pos.Equals(ops[0].Position))
	assert.True(t, cell.Equals(ops[0].Cell))
}

func TestDiffService_Diff_MultipleCellChanges(t *testing.T) {
	svc := NewDiffService()

	oldBuf := model.NewBuffer(10, 5)
	newBuf := model.NewBuffer(10, 5)

	// Change multiple cells
	style := value.NewStyleWithFg(value.ColorRed)
	positions := []value.Position{
		value.NewPosition(0, 0),
		value.NewPosition(5, 2),
		value.NewPosition(9, 4),
	}

	for _, pos := range positions {
		newBuf.Set(pos, model.NewCell('X', style))
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

	oldBuf := model.NewBuffer(20, 5)
	newBuf := model.NewBuffer(20, 5)

	style := value.NewStyleWithFg(value.ColorRed)
	newBuf.SetString(value.NewPosition(0, 0), "Hello", style)

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

	oldBuf := model.NewBuffer(10, 5)
	newBuf := model.NewBuffer(10, 5)

	pos := value.NewPosition(5, 2)
	style1 := value.NewStyleWithFg(value.ColorRed)
	style2 := value.NewStyleWithFg(value.ColorBlue)

	// Same char, different style
	oldBuf.Set(pos, model.NewCell('A', style1))
	newBuf.Set(pos, model.NewCell('A', style2))

	ops := svc.Diff(oldBuf, newBuf)

	assert.Len(t, ops, 1)
	assert.Equal(t, OpTypeSet, ops[0].Type)
	assert.True(t, style2.Equals(ops[0].Cell.Style()))
}

func TestDiffService_Diff_DifferentDimensions(t *testing.T) {
	svc := NewDiffService()

	oldBuf := model.NewBuffer(10, 5)
	newBuf := model.NewBuffer(20, 10)

	// Fill new buffer
	style := value.NewStyleWithFg(value.ColorRed)
	newBuf.SetString(value.NewPosition(0, 0), "Test", style)

	ops := svc.Diff(oldBuf, newBuf)

	// Should return full buffer diff
	assert.NotEmpty(t, ops)
}

func TestDiffService_Diff_NilBuffers(t *testing.T) {
	svc := NewDiffService()

	tests := []struct {
		name   string
		oldBuf *model.Buffer
		newBuf *model.Buffer
	}{
		{"both nil", nil, nil},
		{"old nil", nil, model.NewBuffer(10, 5)},
		{"new nil", model.NewBuffer(10, 5), nil},
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

	style1 := value.NewStyleWithFg(value.ColorRed)
	style2 := value.NewStyleWithFg(value.ColorBlue)

	oldLine := []model.Cell{
		model.NewCell('A', style1),
		model.NewCell('B', style1),
		model.NewCell('C', style1),
	}

	newLine := []model.Cell{
		model.NewCell('A', style1), // Same
		model.NewCell('X', style2), // Changed
		model.NewCell('C', style1), // Same
	}

	ops := svc.DiffLine(oldLine, newLine, 0)

	assert.Len(t, ops, 1)
	assert.Equal(t, OpTypeSet, ops[0].Type)
	assert.Equal(t, 1, ops[0].Position.X())
	assert.Equal(t, 'X', ops[0].Cell.Char())
}

func TestDiffService_DiffLine_DifferentLength(t *testing.T) {
	svc := NewDiffService()

	oldLine := []model.Cell{
		model.NewCell('A', value.NewStyle()),
	}
	newLine := []model.Cell{
		model.NewCell('A', value.NewStyle()),
		model.NewCell('B', value.NewStyle()),
	}

	ops := svc.DiffLine(oldLine, newLine, 0)
	assert.Nil(t, ops)
}

func TestDiffService_CountChanges(t *testing.T) {
	svc := NewDiffService()

	oldBuf := model.NewBuffer(10, 5)
	newBuf := model.NewBuffer(10, 5)

	// No changes
	count := svc.CountChanges(oldBuf, newBuf)
	assert.Equal(t, 0, count)

	// Add changes
	style := value.NewStyleWithFg(value.ColorRed)
	newBuf.SetString(value.NewPosition(0, 0), "Test", style)

	count = svc.CountChanges(oldBuf, newBuf)
	assert.Equal(t, 4, count) // "Test" = 4 changes
}

func TestDiffService_HasChanges(t *testing.T) {
	svc := NewDiffService()

	oldBuf := model.NewBuffer(10, 5)
	newBuf := model.NewBuffer(10, 5)

	// No changes
	assert.False(t, svc.HasChanges(oldBuf, newBuf))

	// Add change
	newBuf.Set(value.NewPosition(0, 0), model.NewCell('A', value.NewStyle()))
	assert.True(t, svc.HasChanges(oldBuf, newBuf))
}

func TestDiffService_HasChanges_DifferentDimensions(t *testing.T) {
	svc := NewDiffService()

	oldBuf := model.NewBuffer(10, 5)
	newBuf := model.NewBuffer(20, 10)

	assert.True(t, svc.HasChanges(oldBuf, newBuf))
}

func TestDiffService_HasChanges_NilBuffers(t *testing.T) {
	svc := NewDiffService()

	assert.True(t, svc.HasChanges(nil, nil))
	assert.True(t, svc.HasChanges(nil, model.NewBuffer(10, 5)))
	assert.True(t, svc.HasChanges(model.NewBuffer(10, 5), nil))
}

func TestDiffService_Optimize(t *testing.T) {
	svc := NewDiffService()

	ops := []DiffOp{
		{Type: OpTypeSet, Position: value.NewPosition(0, 0), Cell: model.NewCell('A', value.NewStyle())},
		{Type: OpTypeSet, Position: value.NewPosition(1, 0), Cell: model.NewCell('B', value.NewStyle())},
		{Type: OpTypeSet, Position: value.NewPosition(2, 0), Cell: model.NewCell('C', value.NewStyle())},
	}

	optimized := svc.Optimize(ops)

	// Currently pass-through, but should not panic
	assert.NotNil(t, optimized)
	assert.Equal(t, len(ops), len(optimized))
}

// Benchmarks
func BenchmarkDiffService_Diff_NoChanges(b *testing.B) {
	svc := NewDiffService()
	oldBuf := model.NewBuffer(80, 24)
	newBuf := model.NewBuffer(80, 24)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = svc.Diff(oldBuf, newBuf)
	}
}

func BenchmarkDiffService_Diff_SmallChanges(b *testing.B) {
	svc := NewDiffService()
	oldBuf := model.NewBuffer(80, 24)
	newBuf := model.NewBuffer(80, 24)

	// Change 10 cells
	style := value.NewStyleWithFg(value.ColorRed)
	for i := 0; i < 10; i++ {
		newBuf.Set(value.NewPosition(i, 0), model.NewCell('X', style))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = svc.Diff(oldBuf, newBuf)
	}
}

func BenchmarkDiffService_Diff_FullScreen(b *testing.B) {
	svc := NewDiffService()
	oldBuf := model.NewBuffer(80, 24)
	newBuf := model.NewBuffer(80, 24)

	// Fill entire buffer
	style := value.NewStyleWithFg(value.ColorRed)
	newBuf.Fill('X', style)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = svc.Diff(oldBuf, newBuf)
	}
}

func BenchmarkDiffService_HasChanges(b *testing.B) {
	svc := NewDiffService()
	oldBuf := model.NewBuffer(80, 24)
	newBuf := model.NewBuffer(80, 24)
	newBuf.Set(value.NewPosition(40, 12), model.NewCell('X', value.NewStyle()))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = svc.HasChanges(oldBuf, newBuf)
	}
}
