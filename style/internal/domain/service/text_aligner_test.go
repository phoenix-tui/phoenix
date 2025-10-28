package service

import (
	"strings"
	"testing"

	"github.com/phoenix-tui/phoenix/style/internal/domain/value"
)

// --- Helper Functions ---

func newTestTextAligner() TextAligner {
	return NewTextAligner()
}

// --- AlignHorizontal Tests ---

func TestAlignHorizontalLeft(t *testing.T) {
	aligner := newTestTextAligner()

	tests := []struct {
		name  string
		text  string
		width int
		want  string
	}{
		{
			name:  "exact fit",
			text:  "Hello",
			width: 5,
			want:  "Hello",
		},
		{
			name:  "needs padding",
			text:  "Hi",
			width: 10,
			want:  "Hi        ", // 2 + 8 spaces
		},
		{
			name:  "empty text",
			text:  "",
			width: 5,
			want:  "     ",
		},
		{
			name:  "single char",
			text:  "X",
			width: 5,
			want:  "X    ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aligner.AlignHorizontal(tt.text, tt.width, value.AlignLeft)
			if got != tt.want {
				t.Errorf("AlignHorizontal(Left) = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAlignHorizontalCenter(t *testing.T) {
	aligner := newTestTextAligner()

	tests := []struct {
		name  string
		text  string
		width int
		want  string
	}{
		{
			name:  "exact fit",
			text:  "Hello",
			width: 5,
			want:  "Hello",
		},
		{
			name:  "even padding",
			text:  "Hi",
			width: 10,
			want:  "    Hi    ", // 4 + 2 + 4
		},
		{
			name:  "odd padding - left gets extra",
			text:  "Hi",
			width: 9,
			want:  "    Hi   ", // 4 + 2 + 3 (left gets extra)
		},
		{
			name:  "single char even",
			text:  "X",
			width: 6,
			want:  "   X  ", // 3 + 1 + 2 (left gets extra)
		},
		{
			name:  "single char odd",
			text:  "X",
			width: 5,
			want:  "  X  ", // 2 + 1 + 2
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aligner.AlignHorizontal(tt.text, tt.width, value.AlignCenter)
			if got != tt.want {
				t.Errorf("AlignHorizontal(Center) = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAlignHorizontalRight(t *testing.T) {
	aligner := newTestTextAligner()

	tests := []struct {
		name  string
		text  string
		width int
		want  string
	}{
		{
			name:  "exact fit",
			text:  "Hello",
			width: 5,
			want:  "Hello",
		},
		{
			name:  "needs padding",
			text:  "Hi",
			width: 10,
			want:  "        Hi", // 8 spaces + 2
		},
		{
			name:  "single char",
			text:  "X",
			width: 5,
			want:  "    X",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aligner.AlignHorizontal(tt.text, tt.width, value.AlignRight)
			if got != tt.want {
				t.Errorf("AlignHorizontal(Right) = %q, want %q", got, tt.want)
			}
		})
	}
}

// --- AlignVertical Tests ---

func TestAlignVerticalTop(t *testing.T) {
	aligner := newTestTextAligner()

	tests := []struct {
		name   string
		text   string
		height int
		want   string
	}{
		{
			name:   "exact fit",
			text:   "Line1\nLine2\nLine3",
			height: 3,
			want:   "Line1\nLine2\nLine3",
		},
		{
			name:   "needs padding",
			text:   "Line1\nLine2",
			height: 5,
			want:   "Line1\nLine2\n     \n     \n     ", // 2 lines + 3 empty
		},
		{
			name:   "single line",
			text:   "Hello",
			height: 3,
			want:   "Hello\n     \n     ", // 1 line + 2 empty
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aligner.AlignVertical(tt.text, tt.height, value.AlignTop)
			if got != tt.want {
				t.Errorf("AlignVertical(Top) =\n%q\nwant\n%q", got, tt.want)
			}
		})
	}
}

func TestAlignVerticalMiddle(t *testing.T) {
	aligner := newTestTextAligner()

	tests := []struct {
		name   string
		text   string
		height int
		want   string
	}{
		{
			name:   "exact fit",
			text:   "Line1\nLine2",
			height: 2,
			want:   "Line1\nLine2",
		},
		{
			name:   "even padding",
			text:   "Hello",
			height: 5,
			want:   "     \n     \nHello\n     \n     ", // 2 empty + 1 line + 2 empty
		},
		{
			name:   "odd padding - top gets extra",
			text:   "Hello",
			height: 4,
			want:   "     \n     \nHello\n     ", // 2 empty + 1 line + 1 empty (top gets extra)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aligner.AlignVertical(tt.text, tt.height, value.AlignMiddle)
			if got != tt.want {
				t.Errorf("AlignVertical(Middle) =\n%q\nwant\n%q", got, tt.want)
			}
		})
	}
}

func TestAlignVerticalBottom(t *testing.T) {
	aligner := newTestTextAligner()

	tests := []struct {
		name   string
		text   string
		height int
		want   string
	}{
		{
			name:   "exact fit",
			text:   "Line1\nLine2",
			height: 2,
			want:   "Line1\nLine2",
		},
		{
			name:   "needs padding",
			text:   "Line1\nLine2",
			height: 5,
			want:   "     \n     \n     \nLine1\nLine2", // 3 empty + 2 lines
		},
		{
			name:   "single line",
			text:   "Hello",
			height: 3,
			want:   "     \n     \nHello", // 2 empty + 1 line
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aligner.AlignVertical(tt.text, tt.height, value.AlignBottom)
			if got != tt.want {
				t.Errorf("AlignVertical(Bottom) =\n%q\nwant\n%q", got, tt.want)
			}
		})
	}
}

// --- AlignBoth Tests ---

func TestAlignBoth(t *testing.T) {
	aligner := newTestTextAligner()

	tests := []struct {
		name      string
		text      string
		width     int
		height    int
		alignment value.Alignment
		want      string
	}{
		{
			name:      "center middle - single line",
			text:      "Hi",
			width:     10,
			height:    3,
			alignment: value.CenterMiddle(),
			want:      "          \n    Hi    \n          ", // center + middle
		},
		{
			name:      "left top",
			text:      "Hi",
			width:     10,
			height:    3,
			alignment: value.LeftTop(),
			want:      "Hi        \n          \n          ",
		},
		{
			name:      "right bottom",
			text:      "Hi",
			width:     10,
			height:    3,
			alignment: value.RightBottom(),
			want:      "          \n          \n        Hi",
		},
		{
			name:      "multi-line center middle",
			text:      "A\nB",
			width:     5,
			height:    4,
			alignment: value.CenterMiddle(),
			want:      "     \n  A  \n  B  \n     ", // 1 empty + 2 centered lines + 1 empty
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aligner.AlignBoth(tt.text, tt.width, tt.height, tt.alignment)
			if got != tt.want {
				t.Errorf("AlignBoth() =\n%q\nwant\n%q", got, tt.want)
			}
		})
	}
}

// --- Unicode Tests ---

func TestAlignHorizontalWithUnicode(t *testing.T) {
	aligner := newTestTextAligner()

	tests := []struct {
		name  string
		text  string
		width int
		align value.HorizontalAlignment
		want  string
	}{
		{
			name:  "emoji left align",
			text:  "👋",
			width: 5,
			align: value.AlignLeft,
			want:  "👋   ", // emoji width 2 + 3 spaces
		},
		{
			name:  "emoji center align",
			text:  "👋",
			width: 6,
			align: value.AlignCenter,
			want:  "  👋  ", // 2 + emoji(2) + 2
		},
		{
			name:  "emoji right align",
			text:  "👋",
			width: 5,
			align: value.AlignRight,
			want:  "   👋", // 3 spaces + emoji(2)
		},
		{
			name:  "CJK center align",
			text:  "你好",
			width: 10,
			align: value.AlignCenter,
			want:  "   你好   ", // 3 + CJK(4) + 3
		},
		{
			name:  "mixed ASCII and emoji",
			text:  "Hi👋",
			width: 10,
			align: value.AlignCenter,
			want:  "   Hi👋   ", // 3 + Hi(2)+emoji(2) + 3
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aligner.AlignHorizontal(tt.text, tt.width, tt.align)
			if got != tt.want {
				t.Errorf("AlignHorizontal(%s) = %q, want %q", tt.text, got, tt.want)
			}
		})
	}
}

func TestAlignVerticalWithMultiLineUnicode(t *testing.T) {
	aligner := newTestTextAligner()

	// Multi-line with different widths (emoji = 2, ASCII = 1)
	text := "👋\nHello"
	got := aligner.AlignVertical(text, 4, value.AlignMiddle)

	lines := strings.Split(got, "\n")
	if len(lines) != 4 {
		t.Errorf("Expected 4 lines, got %d", len(lines))
	}

	// Empty lines should match the widest line (Hello = 5)
	expectedEmptyLine := "     "
	if lines[0] != expectedEmptyLine {
		t.Errorf("Empty line = %q, want %q", lines[0], expectedEmptyLine)
	}
}

// --- Edge Case Tests ---

func TestTextAlignerEdgeCases(t *testing.T) {
	aligner := newTestTextAligner()

	t.Run("text wider than target width", func(t *testing.T) {
		text := "VeryLongText"
		got := aligner.AlignHorizontal(text, 5, value.AlignCenter)
		// Should truncate (simple truncation for now)
		if len(got) > 12 { // Original length
			t.Errorf("Should handle text wider than width")
		}
	})

	t.Run("text taller than target height", func(t *testing.T) {
		text := "Line1\nLine2\nLine3\nLine4"
		got := aligner.AlignVertical(text, 2, value.AlignTop)
		lines := strings.Split(got, "\n")
		if len(lines) > 2 {
			t.Errorf("Should truncate to target height, got %d lines", len(lines))
		}
	})

	t.Run("zero width", func(_ *testing.T) {
		got := aligner.AlignHorizontal("Hi", 0, value.AlignCenter)
		// Should handle gracefully (may truncate or return empty)
		_ = got
	})

	t.Run("zero height", func(_ *testing.T) {
		got := aligner.AlignVertical("Hello\nWorld", 0, value.AlignMiddle)
		// Should handle gracefully (may truncate to empty)
		_ = got
	})

	t.Run("empty text horizontal", func(t *testing.T) {
		got := aligner.AlignHorizontal("", 5, value.AlignCenter)
		if len(got) != 5 {
			t.Errorf("Empty text should pad to width, got %q", got)
		}
	})

	t.Run("single newline", func(t *testing.T) {
		got := aligner.AlignVertical("\n", 3, value.AlignTop)
		lines := strings.Split(got, "\n")
		if len(lines) < 2 {
			t.Errorf("Should preserve newline structure")
		}
	})
}

// --- All Alignment Combinations Test ---

func TestAllAlignmentCombinations(t *testing.T) {
	aligner := newTestTextAligner()

	horizontalAlignments := []value.HorizontalAlignment{
		value.AlignLeft,
		value.AlignCenter,
		value.AlignRight,
	}

	verticalAlignments := []value.VerticalAlignment{
		value.AlignTop,
		value.AlignMiddle,
		value.AlignBottom,
	}

	// Test all 9 combinations
	for _, h := range horizontalAlignments {
		for _, v := range verticalAlignments {
			alignment := value.NewAlignment(h, v)
			got := aligner.AlignBoth("X", 5, 3, alignment)

			lines := strings.Split(got, "\n")
			if len(lines) != 3 {
				t.Errorf("AlignBoth(%s) should produce 3 lines, got %d",
					alignment.String(), len(lines))
			}

			// Verify each line has width 5
			for i, line := range lines {
				if len([]rune(line)) != 5 {
					t.Errorf("AlignBoth(%s) line %d has wrong width: %q",
						alignment.String(), i, line)
				}
			}
		}
	}
}
