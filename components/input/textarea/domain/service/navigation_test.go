package service

import (
	"testing"

	"github.com/phoenix-tui/phoenix/components/input/textarea/domain/model"
)

func TestNewNavigationService(t *testing.T) {
	svc := NewNavigationService()
	if svc == nil {
		t.Error("NewNavigationService() should not return nil")
	}
}

func TestNavigationService_MoveLeft(t *testing.T) {
	svc := NewNavigationService()

	tests := []struct {
		name    string
		text    string
		fromRow int
		fromCol int
		wantRow int
		wantCol int
	}{
		{
			name:    "move left within line",
			text:    "hello",
			fromRow: 0,
			fromCol: 3,
			wantRow: 0,
			wantCol: 2,
		},
		{
			name:    "move left from start of line to end of previous line",
			text:    "hello\nworld",
			fromRow: 1,
			fromCol: 0,
			wantRow: 0,
			wantCol: 5, // End of "hello"
		},
		{
			name:    "cannot move left from start of buffer",
			text:    "hello",
			fromRow: 0,
			fromCol: 0,
			wantRow: 0,
			wantCol: 0,
		},
		{
			name:    "move left in unicode text",
			text:    "he👋lo",
			fromRow: 0,
			fromCol: 3,
			wantRow: 0,
			wantCol: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := model.NewTextArea().
				WithBuffer(model.NewBufferFromString(tt.text)).
				WithCursor(model.NewCursor(tt.fromRow, tt.fromCol))

			result := svc.MoveLeft(ta)
			row, col := result.CursorPosition()

			if row != tt.wantRow || col != tt.wantCol {
				t.Errorf("MoveLeft() cursor = (%d, %d), want (%d, %d)", row, col, tt.wantRow, tt.wantCol)
			}
		})
	}
}

func TestNavigationService_MoveRight(t *testing.T) {
	svc := NewNavigationService()

	tests := []struct {
		name    string
		text    string
		fromRow int
		fromCol int
		wantRow int
		wantCol int
	}{
		{
			name:    "move right within line",
			text:    "hello",
			fromRow: 0,
			fromCol: 2,
			wantRow: 0,
			wantCol: 3,
		},
		{
			name:    "move right from end of line to start of next line",
			text:    "hello\nworld",
			fromRow: 0,
			fromCol: 5, // End of "hello"
			wantRow: 1,
			wantCol: 0,
		},
		{
			name:    "cannot move right from end of buffer",
			text:    "hello",
			fromRow: 0,
			fromCol: 5,
			wantRow: 0,
			wantCol: 5,
		},
		{
			name:    "move right in unicode text",
			text:    "he👋lo",
			fromRow: 0,
			fromCol: 2,
			wantRow: 0,
			wantCol: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := model.NewTextArea().
				WithBuffer(model.NewBufferFromString(tt.text)).
				WithCursor(model.NewCursor(tt.fromRow, tt.fromCol))

			result := svc.MoveRight(ta)
			row, col := result.CursorPosition()

			if row != tt.wantRow || col != tt.wantCol {
				t.Errorf("MoveRight() cursor = (%d, %d), want (%d, %d)", row, col, tt.wantRow, tt.wantCol)
			}
		})
	}
}

func TestNavigationService_MoveUp(t *testing.T) {
	svc := NewNavigationService()

	tests := []struct {
		name    string
		text    string
		fromRow int
		fromCol int
		wantRow int
		wantCol int
	}{
		{
			name:    "move up maintaining column",
			text:    "hello\nworld\ntest",
			fromRow: 2,
			fromCol: 3,
			wantRow: 1,
			wantCol: 3,
		},
		{
			name:    "move up to shorter line clamps column",
			text:    "hi\nhello world",
			fromRow: 1,
			fromCol: 10,
			wantRow: 0,
			wantCol: 2, // End of "hi"
		},
		{
			name:    "cannot move up from first line",
			text:    "hello\nworld",
			fromRow: 0,
			fromCol: 3,
			wantRow: 0,
			wantCol: 3,
		},
		{
			name:    "move up from end of long line",
			text:    "short\nvery long line",
			fromRow: 1,
			fromCol: 10,
			wantRow: 0,
			wantCol: 5, // End of "short"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := model.NewTextArea().
				WithBuffer(model.NewBufferFromString(tt.text)).
				WithCursor(model.NewCursor(tt.fromRow, tt.fromCol))

			result := svc.MoveUp(ta)
			row, col := result.CursorPosition()

			if row != tt.wantRow || col != tt.wantCol {
				t.Errorf("MoveUp() cursor = (%d, %d), want (%d, %d)", row, col, tt.wantRow, tt.wantCol)
			}
		})
	}
}

func TestNavigationService_MoveDown(t *testing.T) {
	svc := NewNavigationService()

	tests := []struct {
		name    string
		text    string
		fromRow int
		fromCol int
		wantRow int
		wantCol int
	}{
		{
			name:    "move down maintaining column",
			text:    "hello\nworld\ntest",
			fromRow: 0,
			fromCol: 3,
			wantRow: 1,
			wantCol: 3,
		},
		{
			name:    "move down to shorter line clamps column",
			text:    "hello world\nhi",
			fromRow: 0,
			fromCol: 10,
			wantRow: 1,
			wantCol: 2, // End of "hi"
		},
		{
			name:    "cannot move down from last line",
			text:    "hello\nworld",
			fromRow: 1,
			fromCol: 3,
			wantRow: 1,
			wantCol: 3,
		},
		{
			name:    "move down from start of line",
			text:    "first\nsecond\nthird",
			fromRow: 0,
			fromCol: 0,
			wantRow: 1,
			wantCol: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := model.NewTextArea().
				WithBuffer(model.NewBufferFromString(tt.text)).
				WithCursor(model.NewCursor(tt.fromRow, tt.fromCol))

			result := svc.MoveDown(ta)
			row, col := result.CursorPosition()

			if row != tt.wantRow || col != tt.wantCol {
				t.Errorf("MoveDown() cursor = (%d, %d), want (%d, %d)", row, col, tt.wantRow, tt.wantCol)
			}
		})
	}
}

func TestNavigationService_MoveToLineStart(t *testing.T) {
	svc := NewNavigationService()

	tests := []struct {
		name    string
		text    string
		fromRow int
		fromCol int
		wantRow int
		wantCol int
	}{
		{
			name:    "move to start from middle",
			text:    "hello world",
			fromRow: 0,
			fromCol: 6,
			wantRow: 0,
			wantCol: 0,
		},
		{
			name:    "move to start from end",
			text:    "hello world",
			fromRow: 0,
			fromCol: 11,
			wantRow: 0,
			wantCol: 0,
		},
		{
			name:    "already at start",
			text:    "hello world",
			fromRow: 0,
			fromCol: 0,
			wantRow: 0,
			wantCol: 0,
		},
		{
			name:    "move to start on second line",
			text:    "first\nsecond line",
			fromRow: 1,
			fromCol: 7,
			wantRow: 1,
			wantCol: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := model.NewTextArea().
				WithBuffer(model.NewBufferFromString(tt.text)).
				WithCursor(model.NewCursor(tt.fromRow, tt.fromCol))

			result := svc.MoveToLineStart(ta)
			row, col := result.CursorPosition()

			if row != tt.wantRow || col != tt.wantCol {
				t.Errorf("MoveToLineStart() cursor = (%d, %d), want (%d, %d)", row, col, tt.wantRow, tt.wantCol)
			}
		})
	}
}

func TestNavigationService_MoveToLineEnd(t *testing.T) {
	svc := NewNavigationService()

	tests := []struct {
		name    string
		text    string
		fromRow int
		fromCol int
		wantRow int
		wantCol int
	}{
		{
			name:    "move to end from start",
			text:    "hello world",
			fromRow: 0,
			fromCol: 0,
			wantRow: 0,
			wantCol: 11,
		},
		{
			name:    "move to end from middle",
			text:    "hello world",
			fromRow: 0,
			fromCol: 6,
			wantRow: 0,
			wantCol: 11,
		},
		{
			name:    "already at end",
			text:    "hello world",
			fromRow: 0,
			fromCol: 11,
			wantRow: 0,
			wantCol: 11,
		},
		{
			name:    "move to end on second line",
			text:    "first\nsecond line",
			fromRow: 1,
			fromCol: 0,
			wantRow: 1,
			wantCol: 11,
		},
		{
			name:    "move to end of empty line",
			text:    "",
			fromRow: 0,
			fromCol: 0,
			wantRow: 0,
			wantCol: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := model.NewTextArea().
				WithBuffer(model.NewBufferFromString(tt.text)).
				WithCursor(model.NewCursor(tt.fromRow, tt.fromCol))

			result := svc.MoveToLineEnd(ta)
			row, col := result.CursorPosition()

			if row != tt.wantRow || col != tt.wantCol {
				t.Errorf("MoveToLineEnd() cursor = (%d, %d), want (%d, %d)", row, col, tt.wantRow, tt.wantCol)
			}
		})
	}
}

func TestNavigationService_MoveToBufferStart(t *testing.T) {
	svc := NewNavigationService()

	tests := []struct {
		name    string
		text    string
		fromRow int
		fromCol int
	}{
		{
			name:    "from middle of buffer",
			text:    "line1\nline2\nline3",
			fromRow: 1,
			fromCol: 3,
		},
		{
			name:    "from end of buffer",
			text:    "line1\nline2\nline3",
			fromRow: 2,
			fromCol: 5,
		},
		{
			name:    "already at start",
			text:    "line1\nline2\nline3",
			fromRow: 0,
			fromCol: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := model.NewTextArea().
				WithBuffer(model.NewBufferFromString(tt.text)).
				WithCursor(model.NewCursor(tt.fromRow, tt.fromCol))

			result := svc.MoveToBufferStart(ta)
			row, col := result.CursorPosition()

			if row != 0 || col != 0 {
				t.Errorf("MoveToBufferStart() cursor = (%d, %d), want (0, 0)", row, col)
			}
		})
	}
}

func TestNavigationService_MoveToBufferEnd(t *testing.T) {
	svc := NewNavigationService()

	tests := []struct {
		name    string
		text    string
		fromRow int
		fromCol int
		wantRow int
		wantCol int
	}{
		{
			name:    "from start of buffer",
			text:    "line1\nline2\nline3",
			fromRow: 0,
			fromCol: 0,
			wantRow: 2,
			wantCol: 5,
		},
		{
			name:    "from middle of buffer",
			text:    "line1\nline2\nline3",
			fromRow: 1,
			fromCol: 3,
			wantRow: 2,
			wantCol: 5,
		},
		{
			name:    "already at end",
			text:    "line1\nline2\nline3",
			fromRow: 2,
			fromCol: 5,
			wantRow: 2,
			wantCol: 5,
		},
		{
			name:    "single line",
			text:    "hello",
			fromRow: 0,
			fromCol: 0,
			wantRow: 0,
			wantCol: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := model.NewTextArea().
				WithBuffer(model.NewBufferFromString(tt.text)).
				WithCursor(model.NewCursor(tt.fromRow, tt.fromCol))

			result := svc.MoveToBufferEnd(ta)
			row, col := result.CursorPosition()

			if row != tt.wantRow || col != tt.wantCol {
				t.Errorf("MoveToBufferEnd() cursor = (%d, %d), want (%d, %d)", row, col, tt.wantRow, tt.wantCol)
			}
		})
	}
}

func TestNavigationService_ForwardWord(t *testing.T) {
	svc := NewNavigationService()

	tests := []struct {
		name    string
		text    string
		fromRow int
		fromCol int
		wantRow int
		wantCol int
	}{
		{
			name:    "move forward one word",
			text:    "hello world",
			fromRow: 0,
			fromCol: 0,
			wantRow: 0,
			wantCol: 6, // After "hello "
		},
		{
			name:    "move forward from middle of word",
			text:    "hello world",
			fromRow: 0,
			fromCol: 2,
			wantRow: 0,
			wantCol: 6,
		},
		{
			name:    "move forward over punctuation",
			text:    "hello, world",
			fromRow: 0,
			fromCol: 0,
			wantRow: 0,
			wantCol: 7, // After "hello, "
		},
		{
			name:    "at end of line",
			text:    "hello",
			fromRow: 0,
			fromCol: 5,
			wantRow: 0,
			wantCol: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := model.NewTextArea().
				WithBuffer(model.NewBufferFromString(tt.text)).
				WithCursor(model.NewCursor(tt.fromRow, tt.fromCol))

			result := svc.ForwardWord(ta)
			row, col := result.CursorPosition()

			if row != tt.wantRow || col != tt.wantCol {
				t.Errorf("ForwardWord() cursor = (%d, %d), want (%d, %d)", row, col, tt.wantRow, tt.wantCol)
			}
		})
	}
}

func TestNavigationService_BackwardWord(t *testing.T) {
	svc := NewNavigationService()

	tests := []struct {
		name    string
		text    string
		fromRow int
		fromCol int
		wantRow int
		wantCol int
	}{
		{
			name:    "move backward one word",
			text:    "hello world",
			fromRow: 0,
			fromCol: 11,
			wantRow: 0,
			wantCol: 6, // Start of "world"
		},
		{
			name:    "move backward from middle of word",
			text:    "hello world",
			fromRow: 0,
			fromCol: 8,
			wantRow: 0,
			wantCol: 6,
		},
		{
			name:    "move backward over punctuation",
			text:    "hello, world",
			fromRow: 0,
			fromCol: 12,
			wantRow: 0,
			wantCol: 7, // Start of "world"
		},
		{
			name:    "at start of line",
			text:    "hello world",
			fromRow: 0,
			fromCol: 0,
			wantRow: 0,
			wantCol: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := model.NewTextArea().
				WithBuffer(model.NewBufferFromString(tt.text)).
				WithCursor(model.NewCursor(tt.fromRow, tt.fromCol))

			result := svc.BackwardWord(ta)
			row, col := result.CursorPosition()

			if row != tt.wantRow || col != tt.wantCol {
				t.Errorf("BackwardWord() cursor = (%d, %d), want (%d, %d)", row, col, tt.wantRow, tt.wantCol)
			}
		})
	}
}

func TestNavigationService_Immutability(t *testing.T) {
	svc := NewNavigationService()
	original := model.NewTextArea().
		WithBuffer(model.NewBufferFromString("line1\nline2\nline3")).
		WithCursor(model.NewCursor(1, 2))

	// Apply all navigation operations
	_ = svc.MoveLeft(original)
	_ = svc.MoveRight(original)
	_ = svc.MoveUp(original)
	_ = svc.MoveDown(original)
	_ = svc.MoveToLineStart(original)
	_ = svc.MoveToLineEnd(original)
	_ = svc.MoveToBufferStart(original)
	_ = svc.MoveToBufferEnd(original)
	_ = svc.ForwardWord(original)
	_ = svc.BackwardWord(original)

	// Original should remain unchanged
	row, col := original.CursorPosition()
	if row != 1 || col != 2 {
		t.Errorf("Original cursor changed: (%d, %d), want (1, 2)", row, col)
	}
	if original.Value() != "line1\nline2\nline3" {
		t.Error("Original buffer was modified")
	}
}
