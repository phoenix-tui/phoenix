package service

import (
	"testing"

	"github.com/phoenix-tui/phoenix/components/input/textarea/domain/model"
)

func TestNewEditingService(t *testing.T) {
	svc := NewEditingService()
	if svc == nil {
		t.Error("NewEditingService() should not return nil")
	}
}

func TestEditingService_InsertChar(t *testing.T) {
	svc := NewEditingService()

	tests := []struct {
		name     string
		text     string
		row      int
		col      int
		char     rune
		wantText string
		wantRow  int
		wantCol  int
		readOnly bool
	}{
		{
			name:     "insert at start",
			text:     "hello",
			row:      0,
			col:      0,
			char:     'X',
			wantText: "Xhello",
			wantRow:  0,
			wantCol:  1,
		},
		{
			name:     "insert in middle",
			text:     "hello",
			row:      0,
			col:      2,
			char:     'X',
			wantText: "heXllo",
			wantRow:  0,
			wantCol:  3,
		},
		{
			name:     "insert at end",
			text:     "hello",
			row:      0,
			col:      5,
			char:     'X',
			wantText: "helloX",
			wantRow:  0,
			wantCol:  6,
		},
		{
			name:     "insert emoji",
			text:     "hello",
			row:      0,
			col:      2,
			char:     '👋',
			wantText: "he👋llo",
			wantRow:  0,
			wantCol:  3,
		},
		{
			name:     "insert in multiline",
			text:     "line1\nline2",
			row:      1,
			col:      3,
			char:     'X',
			wantText: "line1\nlinXe2",
			wantRow:  1,
			wantCol:  4,
		},
		{
			name:     "read-only blocks insert",
			text:     "hello",
			row:      0,
			col:      2,
			char:     'X',
			wantText: "hello",
			wantRow:  2,
			wantCol:  0,
			readOnly: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := model.NewTextArea().
				WithBuffer(model.NewBufferFromString(tt.text)).
				WithCursor(model.NewCursor(tt.row, tt.col)).
				WithReadOnly(tt.readOnly)

			result := svc.InsertChar(ta, tt.char)

			if result.Value() != tt.wantText {
				t.Errorf("InsertChar() text = %q, want %q", result.Value(), tt.wantText)
			}

			row, col := result.CursorPosition()
			if !tt.readOnly && (row != tt.wantRow || col != tt.wantCol) {
				t.Errorf("InsertChar() cursor = (%d, %d), want (%d, %d)", row, col, tt.wantRow, tt.wantCol)
			}
		})
	}
}

func TestEditingService_DeleteCharBackward(t *testing.T) {
	svc := NewEditingService()

	tests := []struct {
		name     string
		text     string
		row      int
		col      int
		wantText string
		wantRow  int
		wantCol  int
		readOnly bool
	}{
		{
			name:     "delete middle char",
			text:     "hello",
			row:      0,
			col:      3,
			wantText: "helo",
			wantRow:  0,
			wantCol:  2,
		},
		{
			name:     "delete at start does nothing",
			text:     "hello",
			row:      0,
			col:      0,
			wantText: "hello",
			wantRow:  0,
			wantCol:  0,
		},
		{
			name:     "delete at end",
			text:     "hello",
			row:      0,
			col:      5,
			wantText: "hell",
			wantRow:  0,
			wantCol:  4,
		},
		{
			name:     "delete joins lines",
			text:     "hello\nworld",
			row:      1,
			col:      0,
			wantText: "helloworld",
			wantRow:  0,
			wantCol:  5,
		},
		{
			name:     "delete emoji",
			text:     "he👋lo",
			row:      0,
			col:      3,
			wantText: "helo",
			wantRow:  0,
			wantCol:  2,
		},
		{
			name:     "read-only blocks delete",
			text:     "hello",
			row:      0,
			col:      3,
			wantText: "hello",
			wantRow:  3,
			wantCol:  0,
			readOnly: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := model.NewTextArea().
				WithBuffer(model.NewBufferFromString(tt.text)).
				WithCursor(model.NewCursor(tt.row, tt.col)).
				WithReadOnly(tt.readOnly)

			result := svc.DeleteCharBackward(ta)

			if result.Value() != tt.wantText {
				t.Errorf("DeleteCharBackward() text = %q, want %q", result.Value(), tt.wantText)
			}

			row, col := result.CursorPosition()
			if !tt.readOnly && (row != tt.wantRow || col != tt.wantCol) {
				t.Errorf("DeleteCharBackward() cursor = (%d, %d), want (%d, %d)", row, col, tt.wantRow, tt.wantCol)
			}
		})
	}
}

func TestEditingService_DeleteCharForward(t *testing.T) {
	svc := NewEditingService()

	tests := []struct {
		name     string
		text     string
		row      int
		col      int
		wantText string
		wantRow  int
		wantCol  int
		readOnly bool
	}{
		{
			name:     "delete char at cursor",
			text:     "hello",
			row:      0,
			col:      2,
			wantText: "helo",
			wantRow:  0,
			wantCol:  2,
		},
		{
			name:     "delete at end joins lines",
			text:     "hello\nworld",
			row:      0,
			col:      5,
			wantText: "helloworld",
			wantRow:  0,
			wantCol:  5,
		},
		{
			name:     "delete at end of buffer does nothing",
			text:     "hello",
			row:      0,
			col:      5,
			wantText: "hello",
			wantRow:  0,
			wantCol:  5,
		},
		{
			name:     "delete emoji",
			text:     "he👋lo",
			row:      0,
			col:      2,
			wantText: "helo",
			wantRow:  0,
			wantCol:  2,
		},
		{
			name:     "read-only blocks delete",
			text:     "hello",
			row:      0,
			col:      2,
			wantText: "hello",
			wantRow:  2,
			wantCol:  0,
			readOnly: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := model.NewTextArea().
				WithBuffer(model.NewBufferFromString(tt.text)).
				WithCursor(model.NewCursor(tt.row, tt.col)).
				WithReadOnly(tt.readOnly)

			result := svc.DeleteCharForward(ta)

			if result.Value() != tt.wantText {
				t.Errorf("DeleteCharForward() text = %q, want %q", result.Value(), tt.wantText)
			}

			row, col := result.CursorPosition()
			if !tt.readOnly && (row != tt.wantRow || col != tt.wantCol) {
				t.Errorf("DeleteCharForward() cursor = (%d, %d), want (%d, %d)", row, col, tt.wantRow, tt.wantCol)
			}
		})
	}
}

func TestEditingService_InsertNewline(t *testing.T) {
	svc := NewEditingService()

	tests := []struct {
		name     string
		text     string
		row      int
		col      int
		wantText string
		wantRow  int
		wantCol  int
		readOnly bool
	}{
		{
			name:     "insert newline at start",
			text:     "hello",
			row:      0,
			col:      0,
			wantText: "\nhello",
			wantRow:  1,
			wantCol:  0,
		},
		{
			name:     "insert newline in middle",
			text:     "hello",
			row:      0,
			col:      2,
			wantText: "he\nllo",
			wantRow:  1,
			wantCol:  0,
		},
		{
			name:     "insert newline at end",
			text:     "hello",
			row:      0,
			col:      5,
			wantText: "hello\n",
			wantRow:  1,
			wantCol:  0,
		},
		{
			name:     "read-only blocks newline",
			text:     "hello",
			row:      0,
			col:      2,
			wantText: "hello",
			wantRow:  2,
			wantCol:  0,
			readOnly: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := model.NewTextArea().
				WithBuffer(model.NewBufferFromString(tt.text)).
				WithCursor(model.NewCursor(tt.row, tt.col)).
				WithReadOnly(tt.readOnly)

			result := svc.InsertNewline(ta)

			if result.Value() != tt.wantText {
				t.Errorf("InsertNewline() text = %q, want %q", result.Value(), tt.wantText)
			}

			row, col := result.CursorPosition()
			if !tt.readOnly && (row != tt.wantRow || col != tt.wantCol) {
				t.Errorf("InsertNewline() cursor = (%d, %d), want (%d, %d)", row, col, tt.wantRow, tt.wantCol)
			}
		})
	}
}

func TestEditingService_KillLine(t *testing.T) {
	svc := NewEditingService()

	tests := []struct {
		name       string
		text       string
		row        int
		col        int
		wantText   string
		wantKilled string
		readOnly   bool
	}{
		{
			name:       "kill from middle to end",
			text:       "hello world",
			row:        0,
			col:        6,
			wantText:   "hello ",
			wantKilled: "world",
		},
		{
			name:       "kill from start",
			text:       "hello world",
			row:        0,
			col:        0,
			wantText:   "",
			wantKilled: "hello world",
		},
		{
			name:       "kill at end joins next line",
			text:       "hello\nworld",
			row:        0,
			col:        5,
			wantText:   "helloworld",
			wantKilled: "\n",
		},
		{
			name:       "kill at end of buffer does nothing",
			text:       "hello",
			row:        0,
			col:        5,
			wantText:   "hello",
			wantKilled: "",
		},
		{
			name:       "read-only blocks kill",
			text:       "hello world",
			row:        0,
			col:        6,
			wantText:   "hello world",
			wantKilled: "",
			readOnly:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := model.NewTextArea().
				WithBuffer(model.NewBufferFromString(tt.text)).
				WithCursor(model.NewCursor(tt.row, tt.col)).
				WithReadOnly(tt.readOnly)

			result := svc.KillLine(ta)

			if result.Value() != tt.wantText {
				t.Errorf("KillLine() text = %q, want %q", result.Value(), tt.wantText)
			}

			if !tt.readOnly {
				killed := result.GetKillRing().Yank()
				if killed != tt.wantKilled {
					t.Errorf("KillLine() killed = %q, want %q", killed, tt.wantKilled)
				}
			}
		})
	}
}

func TestEditingService_KillWord(t *testing.T) {
	svc := NewEditingService()

	tests := []struct {
		name       string
		text       string
		row        int
		col        int
		wantText   string
		wantKilled string
		readOnly   bool
	}{
		{
			name:       "kill word from start",
			text:       "hello world",
			row:        0,
			col:        0,
			wantText:   " world",
			wantKilled: "hello",
		},
		{
			name:       "kill word from middle",
			text:       "hello world test",
			row:        0,
			col:        6,
			wantText:   "hello  test",
			wantKilled: "world",
		},
		{
			name:       "kill at word boundary does nothing",
			text:       "hello world",
			row:        0,
			col:        5,
			wantText:   "hello world",
			wantKilled: "",
		},
		{
			name:       "read-only blocks kill",
			text:       "hello world",
			row:        0,
			col:        0,
			wantText:   "hello world",
			wantKilled: "",
			readOnly:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := model.NewTextArea().
				WithBuffer(model.NewBufferFromString(tt.text)).
				WithCursor(model.NewCursor(tt.row, tt.col)).
				WithReadOnly(tt.readOnly)

			result := svc.KillWord(ta)

			if result.Value() != tt.wantText {
				t.Errorf("KillWord() text = %q, want %q", result.Value(), tt.wantText)
			}

			if !tt.readOnly && tt.wantKilled != "" {
				killed := result.GetKillRing().Yank()
				if killed != tt.wantKilled {
					t.Errorf("KillWord() killed = %q, want %q", killed, tt.wantKilled)
				}
			}
		})
	}
}

func TestEditingService_Yank(t *testing.T) {
	svc := NewEditingService()

	tests := []struct {
		name     string
		text     string
		row      int
		col      int
		killText string
		wantText string
		wantRow  int
		wantCol  int
		readOnly bool
	}{
		{
			name:     "yank at start",
			text:     "hello",
			row:      0,
			col:      0,
			killText: "world",
			wantText: "worldhello",
			wantRow:  0,
			wantCol:  5,
		},
		{
			name:     "yank in middle",
			text:     "hello",
			row:      0,
			col:      2,
			killText: "XXX",
			wantText: "heXXXllo",
			wantRow:  0,
			wantCol:  5,
		},
		{
			name:     "yank at end",
			text:     "hello",
			row:      0,
			col:      5,
			killText: " world",
			wantText: "hello world",
			wantRow:  0,
			wantCol:  11,
		},
		{
			name:     "yank multiline",
			text:     "hello",
			row:      0,
			col:      5,
			killText: "\nworld",
			wantText: "hello\nworld",
			wantRow:  1,
			wantCol:  5,
		},
		{
			name:     "yank empty ring does nothing",
			text:     "hello",
			row:      0,
			col:      2,
			killText: "",
			wantText: "hello",
			wantRow:  0,
			wantCol:  2,
		},
		{
			name:     "read-only blocks yank",
			text:     "hello",
			row:      0,
			col:      2,
			killText: "world",
			wantText: "hello",
			wantRow:  2,
			wantCol:  0,
			readOnly: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			killRing := model.NewKillRing(10)
			if tt.killText != "" {
				killRing = killRing.Kill(tt.killText)
			}

			ta := model.NewTextArea().
				WithBuffer(model.NewBufferFromString(tt.text)).
				WithCursor(model.NewCursor(tt.row, tt.col)).
				WithKillRing(killRing).
				WithReadOnly(tt.readOnly)

			result := svc.Yank(ta)

			if result.Value() != tt.wantText {
				t.Errorf("Yank() text = %q, want %q", result.Value(), tt.wantText)
			}

			row, col := result.CursorPosition()
			if !tt.readOnly && (row != tt.wantRow || col != tt.wantCol) {
				t.Errorf("Yank() cursor = (%d, %d), want (%d, %d)", row, col, tt.wantRow, tt.wantCol)
			}
		})
	}
}

func TestEditingService_ReadOnlyMode(t *testing.T) {
	svc := NewEditingService()
	ta := model.NewTextArea().
		WithBuffer(model.NewBufferFromString("hello world")).
		WithCursor(model.NewCursor(0, 5)).
		WithReadOnly(true)

	originalText := ta.Value()

	// Try all editing operations - none should modify text.
	tests := []struct {
		name string
		op   func(*model.TextArea) *model.TextArea
	}{
		{"InsertChar", func(ta *model.TextArea) *model.TextArea { return svc.InsertChar(ta, 'X') }},
		{"DeleteCharBackward", func(ta *model.TextArea) *model.TextArea { return svc.DeleteCharBackward(ta) }},
		{"DeleteCharForward", func(ta *model.TextArea) *model.TextArea { return svc.DeleteCharForward(ta) }},
		{"InsertNewline", func(ta *model.TextArea) *model.TextArea { return svc.InsertNewline(ta) }},
		{"KillLine", func(ta *model.TextArea) *model.TextArea { return svc.KillLine(ta) }},
		{"KillWord", func(ta *model.TextArea) *model.TextArea { return svc.KillWord(ta) }},
		{"Yank", func(ta *model.TextArea) *model.TextArea { return svc.Yank(ta) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.op(ta)
			if result.Value() != originalText {
				t.Errorf("%s modified read-only text: %q, want %q", tt.name, result.Value(), originalText)
			}
		})
	}
}

func TestEditingService_Immutability(t *testing.T) {
	svc := NewEditingService()
	original := model.NewTextArea().
		WithBuffer(model.NewBufferFromString("hello world")).
		WithCursor(model.NewCursor(0, 5))

	originalText := original.Value()
	originalRow, originalCol := original.CursorPosition()

	// Apply all editing operations.
	_ = svc.InsertChar(original, 'X')
	_ = svc.DeleteCharBackward(original)
	_ = svc.DeleteCharForward(original)
	_ = svc.InsertNewline(original)
	_ = svc.KillLine(original)
	_ = svc.KillWord(original)
	_ = svc.Yank(original)

	// Original should remain unchanged.
	if original.Value() != originalText {
		t.Errorf("Original text changed: %q, want %q", original.Value(), originalText)
	}

	row, col := original.CursorPosition()
	if row != originalRow || col != originalCol {
		t.Errorf("Original cursor changed: (%d, %d), want (%d, %d)", row, col, originalRow, originalCol)
	}
}

func TestEditingService_ComplexScenarios(t *testing.T) {
	svc := NewEditingService()

	t.Run("type and delete sequence", func(t *testing.T) {
		ta := model.NewTextArea()

		// Type "hello".
		for _, ch := range "hello" {
			ta = svc.InsertChar(ta, ch)
		}
		if ta.Value() != "hello" {
			t.Errorf("After typing: %q, want %q", ta.Value(), "hello")
		}

		// Delete 2 chars.
		ta = svc.DeleteCharBackward(ta)
		ta = svc.DeleteCharBackward(ta)
		if ta.Value() != "hel" {
			t.Errorf("After deleting: %q, want %q", ta.Value(), "hel")
		}
	})

	t.Run("kill and yank sequence", func(t *testing.T) {
		ta := model.NewTextArea().
			WithBuffer(model.NewBufferFromString("hello world")).
			WithCursor(model.NewCursor(0, 6))

		// Kill "world".
		ta = svc.KillLine(ta)
		if ta.Value() != "hello " {
			t.Errorf("After kill: %q, want %q", ta.Value(), "hello ")
		}

		// Move cursor and yank.
		ta = ta.WithCursor(model.NewCursor(0, 0))
		ta = svc.Yank(ta)
		if ta.Value() != "worldhello " {
			t.Errorf("After yank: %q, want %q", ta.Value(), "worldhello ")
		}
	})

	t.Run("multiline editing", func(t *testing.T) {
		ta := model.NewTextArea().
			WithBuffer(model.NewBufferFromString("line1")).
			WithCursor(model.NewCursor(0, 5))

		// Insert newline.
		ta = svc.InsertNewline(ta)
		// Type "line2".
		for _, ch := range "line2" {
			ta = svc.InsertChar(ta, ch)
		}

		if ta.Value() != "line1\nline2" {
			t.Errorf("Multiline: %q, want %q", ta.Value(), "line1\nline2")
		}

		row, col := ta.CursorPosition()
		if row != 1 || col != 5 {
			t.Errorf("Cursor: (%d, %d), want (1, 5)", row, col)
		}
	})
}

func TestEditingService_EdgeCases(t *testing.T) {
	svc := NewEditingService()

	t.Run("insert char at various positions", func(t *testing.T) {
		positions := []struct {
			row, col int
		}{
			{0, 0}, // start
			{0, 2}, // middle
			{0, 5}, // end
		}

		for _, pos := range positions {
			ta := model.NewTextArea().
				WithBuffer(model.NewBufferFromString("hello")).
				WithCursor(model.NewCursor(pos.row, pos.col))

			result := svc.InsertChar(ta, 'X')
			if result.IsEmpty() {
				t.Errorf("InsertChar at (%d, %d) resulted in empty buffer", pos.row, pos.col)
			}
		}
	})

	t.Run("delete operations at boundaries", func(t *testing.T) {
		// Delete backward at start.
		ta1 := model.NewTextArea().
			WithBuffer(model.NewBufferFromString("hello")).
			WithCursor(model.NewCursor(0, 0))
		result1 := svc.DeleteCharBackward(ta1)
		if result1.Value() != "hello" {
			t.Error("DeleteBackward at start should not change text")
		}

		// Delete forward at end.
		ta2 := model.NewTextArea().
			WithBuffer(model.NewBufferFromString("hello")).
			WithCursor(model.NewCursor(0, 5))
		result2 := svc.DeleteCharForward(ta2)
		if result2.Value() != "hello" {
			t.Error("DeleteForward at end should not change text")
		}
	})
}
