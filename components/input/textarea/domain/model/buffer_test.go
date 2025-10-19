package model

import (
	"testing"

	"github.com/phoenix-tui/phoenix/components/input/textarea/domain/value"
)

func TestNewBuffer(t *testing.T) {
	buf := NewBuffer()

	if buf.LineCount() != 1 {
		t.Errorf("NewBuffer() should have 1 line, got %d", buf.LineCount())
	}

	if buf.Line(0) != "" {
		t.Errorf("NewBuffer() first line should be empty, got %q", buf.Line(0))
	}

	if !buf.IsEmpty() {
		t.Error("NewBuffer() should be empty")
	}
}

func TestNewBufferFromString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantLines int
		wantLine0 string
	}{
		{
			name:      "single line",
			input:     "hello",
			wantLines: 1,
			wantLine0: "hello",
		},
		{
			name:      "two lines",
			input:     "hello\nworld",
			wantLines: 2,
			wantLine0: "hello",
		},
		{
			name:      "empty string",
			input:     "",
			wantLines: 1,
			wantLine0: "",
		},
		{
			name:      "trailing newline",
			input:     "hello\n",
			wantLines: 2,
			wantLine0: "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := NewBufferFromString(tt.input)

			if buf.LineCount() != tt.wantLines {
				t.Errorf("LineCount() = %d, want %d", buf.LineCount(), tt.wantLines)
			}

			if buf.Line(0) != tt.wantLine0 {
				t.Errorf("Line(0) = %q, want %q", buf.Line(0), tt.wantLine0)
			}
		})
	}
}

func TestBuffer_InsertChar(t *testing.T) {
	tests := []struct {
		name     string
		initial  string
		row      int
		col      int
		char     rune
		expected string
	}{
		{
			name:     "insert at start",
			initial:  "hello",
			row:      0,
			col:      0,
			char:     'X',
			expected: "Xhello",
		},
		{
			name:     "insert in middle",
			initial:  "hello",
			row:      0,
			col:      3,
			char:     'X',
			expected: "helXlo",
		},
		{
			name:     "insert at end",
			initial:  "hello",
			row:      0,
			col:      5,
			char:     'X',
			expected: "helloX",
		},
		{
			name:     "insert emoji",
			initial:  "hello",
			row:      0,
			col:      2,
			char:     '👋',
			expected: "he👋llo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := NewBufferFromString(tt.initial)
			result := buf.InsertChar(tt.row, tt.col, tt.char)

			if result.String() != tt.expected {
				t.Errorf("InsertChar() = %q, want %q", result.String(), tt.expected)
			}

			// Check immutability
			if buf.String() != tt.initial {
				t.Errorf("Original buffer was modified: %q, want %q", buf.String(), tt.initial)
			}
		})
	}
}

func TestBuffer_DeleteChar(t *testing.T) {
	tests := []struct {
		name     string
		initial  string
		row      int
		col      int
		expected string
	}{
		{
			name:     "delete first char",
			initial:  "hello",
			row:      0,
			col:      0,
			expected: "ello",
		},
		{
			name:     "delete middle char",
			initial:  "hello",
			row:      0,
			col:      2,
			expected: "helo",
		},
		{
			name:     "delete last char",
			initial:  "hello",
			row:      0,
			col:      4,
			expected: "hell",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := NewBufferFromString(tt.initial)
			result := buf.DeleteChar(tt.row, tt.col)

			if result.String() != tt.expected {
				t.Errorf("DeleteChar() = %q, want %q", result.String(), tt.expected)
			}
		})
	}
}

func TestBuffer_InsertNewline(t *testing.T) {
	buf := NewBufferFromString("hello world")
	result := buf.InsertNewline(0, 5)

	if result.LineCount() != 2 {
		t.Errorf("InsertNewline() line count = %d, want 2", result.LineCount())
	}

	if result.Line(0) != "hello" {
		t.Errorf("InsertNewline() line 0 = %q, want %q", result.Line(0), "hello")
	}

	if result.Line(1) != " world" {
		t.Errorf("InsertNewline() line 1 = %q, want %q", result.Line(1), " world")
	}
}

func TestBuffer_DeleteLine(t *testing.T) {
	buf := NewBufferFromString("line1\nline2\nline3")
	result, deleted := buf.DeleteLine(1)

	if result.LineCount() != 2 {
		t.Errorf("DeleteLine() line count = %d, want 2", result.LineCount())
	}

	if deleted != "line2" {
		t.Errorf("DeleteLine() deleted = %q, want %q", deleted, "line2")
	}

	if result.Line(0) != "line1" {
		t.Errorf("DeleteLine() line 0 = %q, want %q", result.Line(0), "line1")
	}

	if result.Line(1) != "line3" {
		t.Errorf("DeleteLine() line 1 = %q, want %q", result.Line(1), "line3")
	}
}

func TestBuffer_DeleteLine_LastLine(t *testing.T) {
	buf := NewBufferFromString("only line")
	result, deleted := buf.DeleteLine(0)

	if result.LineCount() != 1 {
		t.Errorf("DeleteLine() on last line should keep 1 empty line, got %d", result.LineCount())
	}

	if deleted != "only line" {
		t.Errorf("DeleteLine() deleted = %q, want %q", deleted, "only line")
	}

	if result.Line(0) != "" {
		t.Errorf("DeleteLine() on last line should result in empty line, got %q", result.Line(0))
	}
}

func TestBuffer_Immutability(t *testing.T) {
	buf := NewBufferFromString("hello")
	result := buf.InsertChar(0, 0, 'X')

	// Original should be unchanged
	if buf.String() != "hello" {
		t.Errorf("Original buffer was modified: %q, want %q", buf.String(), "hello")
	}

	if result.String() != "Xhello" {
		t.Errorf("Result buffer = %q, want %q", result.String(), "Xhello")
	}
}

func TestBuffer_DeleteToLineEnd(t *testing.T) {
	tests := []struct {
		name     string
		initial  string
		row      int
		col      int
		wantText string
		wantKill string
	}{
		{
			name:     "delete from start",
			initial:  "hello world",
			row:      0,
			col:      0,
			wantText: "",
			wantKill: "hello world",
		},
		{
			name:     "delete from middle",
			initial:  "hello world",
			row:      0,
			col:      6,
			wantText: "hello ",
			wantKill: "world",
		},
		{
			name:     "delete from end (empty)",
			initial:  "hello",
			row:      0,
			col:      5,
			wantText: "hello",
			wantKill: "",
		},
		{
			name:     "multiline - delete first line end",
			initial:  "hello\nworld",
			row:      0,
			col:      2,
			wantText: "he\nworld",
			wantKill: "llo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := NewBufferFromString(tt.initial)
			result, killed := buf.DeleteToLineEnd(tt.row, tt.col)

			if result.String() != tt.wantText {
				t.Errorf("DeleteToLineEnd() text = %q, want %q", result.String(), tt.wantText)
			}

			if killed != tt.wantKill {
				t.Errorf("DeleteToLineEnd() killed = %q, want %q", killed, tt.wantKill)
			}

			// Check immutability
			if buf.String() != tt.initial {
				t.Errorf("Original buffer modified")
			}
		})
	}
}

func TestBuffer_SetLine(t *testing.T) {
	tests := []struct {
		name    string
		initial string
		row     int
		newText string
		want    string
	}{
		{
			name:    "set first line",
			initial: "hello\nworld",
			row:     0,
			newText: "HELLO",
			want:    "HELLO\nworld",
		},
		{
			name:    "set second line",
			initial: "hello\nworld",
			row:     1,
			newText: "WORLD",
			want:    "hello\nWORLD",
		},
		{
			name:    "set to empty",
			initial: "hello",
			row:     0,
			newText: "",
			want:    "",
		},
		{
			name:    "set invalid row (out of bounds)",
			initial: "hello",
			row:     5,
			newText: "test",
			want:    "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := NewBufferFromString(tt.initial)
			result := buf.SetLine(tt.row, tt.newText)

			if result.String() != tt.want {
				t.Errorf("SetLine() = %q, want %q", result.String(), tt.want)
			}

			// Check immutability
			if buf.String() != tt.initial {
				t.Errorf("Original buffer modified")
			}
		})
	}
}

func TestBuffer_JoinWithNextLine(t *testing.T) {
	tests := []struct {
		name    string
		initial string
		row     int
		want    string
	}{
		{
			name:    "join two lines",
			initial: "hello\nworld",
			row:     0,
			want:    "helloworld",
		},
		{
			name:    "join middle lines",
			initial: "line1\nline2\nline3",
			row:     1,
			want:    "line1\nline2line3",
		},
		{
			name:    "join last line (no change)",
			initial: "hello\nworld",
			row:     1,
			want:    "hello\nworld",
		},
		{
			name:    "join with spaces",
			initial: "hello \n world",
			row:     0,
			want:    "hello  world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := NewBufferFromString(tt.initial)
			result := buf.JoinWithNextLine(tt.row)

			if result.String() != tt.want {
				t.Errorf("JoinWithNextLine() = %q, want %q", result.String(), tt.want)
			}

			// Check immutability
			if buf.String() != tt.initial {
				t.Errorf("Original buffer modified")
			}
		})
	}
}

func TestBuffer_TextInRange(t *testing.T) {
	tests := []struct {
		name     string
		initial  string
		startRow int
		startCol int
		endRow   int
		endCol   int
		want     string
	}{
		{
			name:     "single line range",
			initial:  "hello world",
			startRow: 0,
			startCol: 0,
			endRow:   0,
			endCol:   5,
			want:     "hello",
		},
		{
			name:     "single line middle",
			initial:  "hello world",
			startRow: 0,
			startCol: 6,
			endRow:   0,
			endCol:   11,
			want:     "world",
		},
		{
			name:     "multiline range",
			initial:  "line1\nline2\nline3",
			startRow: 0,
			startCol: 2,
			endRow:   1,
			endCol:   3,
			want:     "ne1\nlin",
		},
		{
			name:     "entire buffer",
			initial:  "hello\nworld",
			startRow: 0,
			startCol: 0,
			endRow:   1,
			endCol:   5,
			want:     "hello\nworld",
		},
		{
			name:     "empty range",
			initial:  "hello",
			startRow: 0,
			startCol: 2,
			endRow:   0,
			endCol:   2,
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := NewBufferFromString(tt.initial)
			// Create range using value package
			r := value.NewRange(
				value.NewPosition(tt.startRow, tt.startCol),
				value.NewPosition(tt.endRow, tt.endCol),
			)
			result := buf.TextInRange(r)

			if result != tt.want {
				t.Errorf("TextInRange() = %q, want %q", result, tt.want)
			}
		})
	}
}

func TestBuffer_Line_OutOfBounds(t *testing.T) {
	buf := NewBufferFromString("hello\nworld")

	tests := []struct {
		name string
		row  int
		want string
	}{
		{
			name: "valid row 0",
			row:  0,
			want: "hello",
		},
		{
			name: "valid row 1",
			row:  1,
			want: "world",
		},
		{
			name: "out of bounds positive",
			row:  10,
			want: "",
		},
		{
			name: "out of bounds negative",
			row:  -1,
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buf.Line(tt.row)
			if result != tt.want {
				t.Errorf("Line(%d) = %q, want %q", tt.row, result, tt.want)
			}
		})
	}
}

func TestBuffer_InsertChar_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		initial  string
		row      int
		col      int
		char     rune
		expected string
	}{
		{
			name:     "insert beyond line end (clamped)",
			initial:  "hi",
			row:      0,
			col:      100,
			char:     'X',
			expected: "hiX",
		},
		{
			name:     "insert on invalid row (no change)",
			initial:  "hello",
			row:      10,
			col:      0,
			char:     'X',
			expected: "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := NewBufferFromString(tt.initial)
			result := buf.InsertChar(tt.row, tt.col, tt.char)

			if result.String() != tt.expected {
				t.Errorf("InsertChar() = %q, want %q", result.String(), tt.expected)
			}
		})
	}
}

func TestBuffer_InsertNewline_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		initial   string
		row       int
		col       int
		wantLines int
		wantLine0 string
		wantLine1 string
	}{
		{
			name:      "insert at start",
			initial:   "hello",
			row:       0,
			col:       0,
			wantLines: 2,
			wantLine0: "",
			wantLine1: "hello",
		},
		{
			name:      "insert at end",
			initial:   "hello",
			row:       0,
			col:       5,
			wantLines: 2,
			wantLine0: "hello",
			wantLine1: "",
		},
		{
			name:      "insert beyond line end (clamped)",
			initial:   "hi",
			row:       0,
			col:       100,
			wantLines: 2,
			wantLine0: "hi",
			wantLine1: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := NewBufferFromString(tt.initial)
			result := buf.InsertNewline(tt.row, tt.col)

			if result.LineCount() != tt.wantLines {
				t.Errorf("InsertNewline() line count = %d, want %d", result.LineCount(), tt.wantLines)
			}

			if result.Line(0) != tt.wantLine0 {
				t.Errorf("InsertNewline() line 0 = %q, want %q", result.Line(0), tt.wantLine0)
			}

			if result.Line(1) != tt.wantLine1 {
				t.Errorf("InsertNewline() line 1 = %q, want %q", result.Line(1), tt.wantLine1)
			}
		})
	}
}

func TestBuffer_NewBufferFromString_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantLines int
		wantLine0 string
	}{
		{
			name:      "multiple newlines",
			input:     "\n\n\n",
			wantLines: 4,
			wantLine0: "",
		},
		{
			name:      "windows line endings (keeps \\r)",
			input:     "line1\r\nline2",
			wantLines: 2,
			wantLine0: "line1\r", // \r is preserved (only \n splits lines)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := NewBufferFromString(tt.input)

			if buf.LineCount() != tt.wantLines {
				t.Errorf("LineCount() = %d, want %d", buf.LineCount(), tt.wantLines)
			}

			if buf.Line(0) != tt.wantLine0 {
				t.Errorf("Line(0) = %q, want %q", buf.Line(0), tt.wantLine0)
			}
		})
	}
}

func TestBuffer_DeleteChar_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		initial  string
		row      int
		col      int
		expected string
	}{
		{
			name:     "delete beyond end (no change)",
			initial:  "hello",
			row:      0,
			col:      10,
			expected: "hello",
		},
		{
			name:     "delete on invalid row",
			initial:  "hello",
			row:      5,
			col:      0,
			expected: "hello",
		},
		{
			name:     "delete emoji",
			initial:  "he👋lo",
			row:      0,
			col:      2,
			expected: "helo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := NewBufferFromString(tt.initial)
			result := buf.DeleteChar(tt.row, tt.col)

			if result.String() != tt.expected {
				t.Errorf("DeleteChar() = %q, want %q", result.String(), tt.expected)
			}
		})
	}
}

func TestBuffer_TextInRange_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		initial  string
		startRow int
		startCol int
		endRow   int
		endCol   int
		want     string
	}{
		{
			name:     "range with emojis",
			initial:  "hello 👋 world",
			startRow: 0,
			startCol: 6,
			endRow:   0,
			endCol:   8,
			want:     "👋 ",
		},
		{
			name:     "three line range",
			initial:  "line1\nline2\nline3\nline4",
			startRow: 1,
			startCol: 2,
			endRow:   3,
			endCol:   2,
			want:     "ne2\nline3\nli",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := NewBufferFromString(tt.initial)
			r := value.NewRange(
				value.NewPosition(tt.startRow, tt.startCol),
				value.NewPosition(tt.endRow, tt.endCol),
			)
			result := buf.TextInRange(r)

			if result != tt.want {
				t.Errorf("TextInRange() = %q, want %q", result, tt.want)
			}
		})
	}
}
