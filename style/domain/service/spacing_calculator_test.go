package service

import (
	"strings"
	"testing"

	coreService "github.com/phoenix-tui/phoenix/core/domain/service"
	"github.com/phoenix-tui/phoenix/style/domain/value"
)

// --- Helper Functions ---

func newTestSpacingCalculator() SpacingCalculator {
	unicodeService := coreService.NewUnicodeService()
	return NewSpacingCalculator(unicodeService)
}

// --- CalculateTotalWidth Tests ---

func TestCalculateTotalWidth(t *testing.T) {
	calc := newTestSpacingCalculator()

	tests := []struct {
		name         string
		contentWidth int
		padding      value.Padding
		margin       value.Margin
		want         int
	}{
		{
			name:         "no padding or margin",
			contentWidth: 10,
			padding:      value.UniformPadding(0),
			margin:       value.UniformMargin(0),
			want:         10,
		},
		{
			name:         "uniform padding only",
			contentWidth: 10,
			padding:      value.UniformPadding(2),
			margin:       value.UniformMargin(0),
			want:         14, // 10 + 2*2
		},
		{
			name:         "uniform margin only",
			contentWidth: 10,
			padding:      value.UniformPadding(0),
			margin:       value.UniformMargin(3),
			want:         16, // 10 + 3*2
		},
		{
			name:         "both padding and margin",
			contentWidth: 10,
			padding:      value.UniformPadding(2),
			margin:       value.UniformMargin(3),
			want:         20, // 10 + 2*2 + 3*2
		},
		{
			name:         "asymmetric padding",
			contentWidth: 10,
			padding:      value.NewPadding(0, 2, 0, 3), // top, right, bottom, left
			margin:       value.UniformMargin(0),
			want:         15, // 10 + 2 + 3
		},
		{
			name:         "asymmetric margin",
			contentWidth: 10,
			padding:      value.UniformPadding(0),
			margin:       value.NewMargin(0, 4, 0, 5), // top, right, bottom, left
			want:         19,                          // 10 + 4 + 5
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.CalculateTotalWidth(tt.contentWidth, tt.padding, tt.margin)
			if got != tt.want {
				t.Errorf("CalculateTotalWidth() = %d, want %d", got, tt.want)
			}
		})
	}
}

// --- CalculateTotalHeight Tests ---

func TestCalculateTotalHeight(t *testing.T) {
	calc := newTestSpacingCalculator()

	tests := []struct {
		name          string
		contentHeight int
		padding       value.Padding
		margin        value.Margin
		want          int
	}{
		{
			name:          "no padding or margin",
			contentHeight: 5,
			padding:       value.UniformPadding(0),
			margin:        value.UniformMargin(0),
			want:          5,
		},
		{
			name:          "uniform padding only",
			contentHeight: 5,
			padding:       value.UniformPadding(2),
			margin:        value.UniformMargin(0),
			want:          9, // 5 + 2*2
		},
		{
			name:          "uniform margin only",
			contentHeight: 5,
			padding:       value.UniformPadding(0),
			margin:        value.UniformMargin(3),
			want:          11, // 5 + 3*2
		},
		{
			name:          "both padding and margin",
			contentHeight: 5,
			padding:       value.UniformPadding(2),
			margin:        value.UniformMargin(3),
			want:          15, // 5 + 2*2 + 3*2
		},
		{
			name:          "asymmetric padding",
			contentHeight: 5,
			padding:       value.NewPadding(2, 0, 3, 0), // top, right, bottom, left
			margin:        value.UniformMargin(0),
			want:          10, // 5 + 2 + 3
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.CalculateTotalHeight(tt.contentHeight, tt.padding, tt.margin)
			if got != tt.want {
				t.Errorf("CalculateTotalHeight() = %d, want %d", got, tt.want)
			}
		})
	}
}

// --- ApplyPadding Tests ---

func TestApplyPadding(t *testing.T) {
	calc := newTestSpacingCalculator()

	tests := []struct {
		name    string
		content string
		padding value.Padding
		want    string
	}{
		{
			name:    "no padding",
			content: "Hello",
			padding: value.UniformPadding(0),
			want:    "Hello",
		},
		{
			name:    "uniform padding",
			content: "Hello",
			padding: value.UniformPadding(1),
			want: strings.Join([]string{
				"       ", // top padding (7 spaces: 1 left + 5 content + 1 right)
				" Hello ",
				"       ", // bottom padding
			}, "\n"),
		},
		{
			name:    "left padding only",
			content: "Hello",
			padding: value.NewPadding(0, 0, 0, 2),
			want:    "  Hello",
		},
		{
			name:    "right padding only",
			content: "Hello",
			padding: value.NewPadding(0, 2, 0, 0),
			want:    "Hello  ",
		},
		{
			name:    "top padding only",
			content: "Hello",
			padding: value.NewPadding(2, 0, 0, 0),
			want: strings.Join([]string{
				"     ", // empty line (width 5)
				"     ", // empty line
				"Hello",
			}, "\n"),
		},
		{
			name:    "multi-line content",
			content: "Hello\nWorld",
			padding: value.NewPadding(1, 1, 1, 1),
			want: strings.Join([]string{
				"       ", // top padding (1+5+1)
				" Hello ",
				" World ",
				"       ", // bottom padding
			}, "\n"),
		},
		{
			name:    "empty content",
			content: "",
			padding: value.UniformPadding(2),
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.ApplyPadding(tt.content, tt.padding)
			if got != tt.want {
				t.Errorf("ApplyPadding() = %q, want %q", got, tt.want)
			}
		})
	}
}

// --- ApplyMargin Tests ---

func TestApplyMargin(t *testing.T) {
	calc := newTestSpacingCalculator()

	tests := []struct {
		name    string
		content string
		margin  value.Margin
		want    string
	}{
		{
			name:    "no margin",
			content: "Hello",
			margin:  value.UniformMargin(0),
			want:    "Hello",
		},
		{
			name:    "uniform margin",
			content: "Hello",
			margin:  value.UniformMargin(1),
			want: strings.Join([]string{
				"       ", // top margin (7 spaces: 1 left + 5 content + 1 right)
				" Hello ",
				"       ", // bottom margin
			}, "\n"),
		},
		{
			name:    "left margin only",
			content: "Hello",
			margin:  value.NewMargin(0, 0, 0, 2),
			want:    "  Hello",
		},
		{
			name:    "multi-line content",
			content: "Hello\nWorld",
			margin:  value.NewMargin(1, 1, 1, 1),
			want: strings.Join([]string{
				"       ", // top margin
				" Hello ",
				" World ",
				"       ", // bottom margin
			}, "\n"),
		},
		{
			name:    "empty content",
			content: "",
			margin:  value.UniformMargin(2),
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.ApplyMargin(tt.content, tt.margin)
			if got != tt.want {
				t.Errorf("ApplyMargin() = %q, want %q", got, tt.want)
			}
		})
	}
}

// --- ApplyBoth Tests ---

func TestApplyBoth(t *testing.T) {
	calc := newTestSpacingCalculator()

	tests := []struct {
		name    string
		content string
		padding value.Padding
		margin  value.Margin
		want    string
	}{
		{
			name:    "no padding or margin",
			content: "Hello",
			padding: value.UniformPadding(0),
			margin:  value.UniformMargin(0),
			want:    "Hello",
		},
		{
			name:    "padding and margin",
			content: "X",
			padding: value.UniformPadding(1),
			margin:  value.UniformMargin(1),
			want: strings.Join([]string{
				"     ", // top margin (5 spaces: 1 margin + 3 padded width + 1 margin)
				"     ", // top padding + margin
				"  X  ", // margin(1) + padding(1) + X + padding(1) + margin(1) = 5
				"     ", // bottom padding + margin
				"     ", // bottom margin
			}, "\n"),
		},
		{
			name:    "verify padding applied first",
			content: "A",
			padding: value.NewPadding(0, 1, 0, 1), // horizontal padding only
			margin:  value.NewMargin(0, 2, 0, 2),  // horizontal margin only
			want:    "   A   ",                    // 2 margin + 1 padding + A + 1 padding + 2 margin
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.ApplyBoth(tt.content, tt.padding, tt.margin)
			if got != tt.want {
				t.Errorf("ApplyBoth() =\n%q\nwant\n%q", got, tt.want)
			}
		})
	}
}

// --- Unicode Tests ---

func TestApplyPaddingWithUnicode(t *testing.T) {
	calc := newTestSpacingCalculator()

	tests := []struct {
		name    string
		content string
		padding value.Padding
		check   func(t *testing.T, got string)
	}{
		{
			name:    "emoji content",
			content: "👋",
			padding: value.NewPadding(0, 1, 0, 1),
			check: func(t *testing.T, got string) {
				want := " 👋 "
				if got != want {
					t.Errorf("ApplyPadding() = %q, want %q", got, want)
				}
			},
		},
		{
			name:    "CJK content",
			content: "你好",
			padding: value.NewPadding(0, 1, 0, 1),
			check: func(t *testing.T, got string) {
				want := " 你好 "
				if got != want {
					t.Errorf("ApplyPadding() = %q, want %q", got, want)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.ApplyPadding(tt.content, tt.padding)
			tt.check(t, got)
		})
	}
}

// --- Edge Case Tests ---

func TestSpacingCalculatorEdgeCases(t *testing.T) {
	calc := newTestSpacingCalculator()

	t.Run("zero width content", func(t *testing.T) {
		width := calc.CalculateTotalWidth(0, value.UniformPadding(2), value.UniformMargin(3))
		want := 10 // 0 + 2*2 + 3*2
		if width != want {
			t.Errorf("CalculateTotalWidth(0, ...) = %d, want %d", width, want)
		}
	})

	t.Run("single space content", func(t *testing.T) {
		got := calc.ApplyPadding(" ", value.UniformPadding(1))
		// Should wrap space with padding
		if !strings.Contains(got, " ") {
			t.Errorf("ApplyPadding with space should preserve space")
		}
	})

	t.Run("content with trailing newline", func(t *testing.T) {
		got := calc.ApplyPadding("Hello\n", value.NewPadding(0, 1, 0, 1))
		lines := strings.Split(got, "\n")
		if len(lines) != 2 {
			t.Errorf("ApplyPadding should preserve trailing newline, got %d lines", len(lines))
		}
	})
}
