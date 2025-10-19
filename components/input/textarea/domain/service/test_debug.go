package service

import (
	"fmt"
	"github.com/phoenix-tui/phoenix/components/input/textarea/domain/model"
	"testing"
)

func TestDebugInsertNewline(t *testing.T) {
	svc := NewEditingService()

	ta := model.NewTextArea().
		WithBuffer(model.NewBufferFromString("hello")).
		WithCursor(model.NewCursor(0, 2))

	row1, col1 := ta.CursorPosition()
	fmt.Printf("Before: lines=%v, cursor=(%d,%d)\n", ta.Lines(), row1, col1)

	result := svc.InsertNewline(ta)

	row, col := result.CursorPosition()
	fmt.Printf("After: lines=%v, cursor=(%d,%d)\n", result.Lines(), row, col)
	fmt.Printf("Value: %q\n", result.Value())
}
