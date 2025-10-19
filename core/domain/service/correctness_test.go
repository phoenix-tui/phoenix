package service

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

// TestCorrectness_Lipgloss562Bug verifies if Lipgloss #562 is still broken
func TestCorrectness_Lipgloss562Bug(t *testing.T) {
	us := NewUnicodeService()

	testCases := []struct {
		name     string
		input    string
		expected int // Correct width
	}{
		{
			name:     "Lipgloss #562 bug case",
			input:    "📝 Test",
			expected: 7, // 📝(2) + space(1) + Test(4) = 7
		},
		{
			name:     "Simple emoji",
			input:    "👋",
			expected: 2,
		},
		{
			name:     "Emoji with text",
			input:    "Hello 👋 World",
			expected: 14, // Hello_=6, 👋=2, _World=6 → 6+2+6=14
		},
		{
			name:     "Multiple emoji",
			input:    "👋😀🎉",
			expected: 6, // 2 + 2 + 2 = 6
		},
		{
			name:     "CJK",
			input:    "你好",
			expected: 4, // 2 + 2 = 4
		},
		{
			name:     "Mixed complex",
			input:    "Hello 👋 世界!",
			expected: 14, // 5 + 1 + 2 + 1 + 2 + 2 + 1 = 14
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			phoenixWidth := us.StringWidth(tc.input)
			lipglossWidth := lipgloss.Width(tc.input)
			runewidthWidth := runewidth.StringWidth(tc.input)

			t.Logf("Input: %q", tc.input)
			t.Logf("Expected: %d", tc.expected)
			t.Logf("Phoenix:   %d (correct: %v)", phoenixWidth, phoenixWidth == tc.expected)
			t.Logf("Lipgloss:  %d (correct: %v)", lipglossWidth, lipglossWidth == tc.expected)
			t.Logf("Runewidth: %d (correct: %v)", runewidthWidth, runewidthWidth == tc.expected)

			// Verify Phoenix is correct
			if phoenixWidth != tc.expected {
				t.Errorf("Phoenix INCORRECT: got %d, want %d", phoenixWidth, tc.expected)
			}

			// Report if Lipgloss is wrong
			if lipglossWidth != tc.expected {
				t.Logf("⚠️  Lipgloss INCORRECT: got %d, want %d (bug still present)", lipglossWidth, tc.expected)
			}

			// Report if go-runewidth is wrong
			if runewidthWidth != tc.expected {
				t.Logf("⚠️  go-runewidth INCORRECT: got %d, want %d", runewidthWidth, tc.expected)
			}
		})
	}
}
