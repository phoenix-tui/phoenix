package service

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

// TestLipgloss562_ActualBrokenCases tests the EXACT strings from issue #562
func TestLipgloss562_ActualBrokenCases(t *testing.T) {
	us := NewUnicodeService()

	testCases := []struct {
		name     string
		input    string
		expected int // Correct width based on Unicode standard
	}{
		{
			name:     "Issue #562 case 1: Clock emoji",
			input:    "⏰ Emoji",
			expected: 8, // ⏰(2) + space(1) + Emoji(5) = 8
		},
		{
			name:     "Issue #562 case 2: Shield emoji",
			input:    "🛡️ Shield",
			expected: 9, // 🛡️(2) + space(1) + Shield(6) = 9 (but variation selector!)
		},
		{
			name:     "Issue #562 case 3: Toolbox emoji (reportedly OK)",
			input:    "🧰",
			expected: 2, // 🧰(2)
		},
		{
			name:     "Clock emoji alone",
			input:    "⏰",
			expected: 2,
		},
		{
			name:     "Shield emoji alone",
			input:    "🛡️",
			expected: 2, // Shield + variation selector = 1 grapheme cluster, width 2
		},
		{
			name:     "Shield without variation selector",
			input:    "🛡",
			expected: 2,
		},
	}

	t.Log("=" + "=" + ("="))
	t.Log("Testing EXACT cases from Lipgloss issue #562")
	t.Log("=================================")

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			phoenixWidth := us.StringWidth(tc.input)
			lipglossWidth := lipgloss.Width(tc.input)
			runewidthWidth := runewidth.StringWidth(tc.input)

			t.Logf("")
			t.Logf("Input: %q", tc.input)
			t.Logf("Expected:  %d", tc.expected)
			t.Logf("Phoenix:   %d %s", phoenixWidth, statusEmoji(phoenixWidth == tc.expected))
			t.Logf("Lipgloss:  %d %s", lipglossWidth, statusEmoji(lipglossWidth == tc.expected))
			t.Logf("Runewidth: %d %s", runewidthWidth, statusEmoji(runewidthWidth == tc.expected))

			// Check if Phoenix is correct (this is what we care about!)
			if phoenixWidth != tc.expected {
				t.Errorf("❌ Phoenix WRONG: got %d, want %d", phoenixWidth, tc.expected)
			}

			// Log Lipgloss status (informational only - we expect bugs!)
			if lipglossWidth != tc.expected {
				t.Logf("🔴 Lipgloss bug #562 CONFIRMED: got %d, want %d", lipglossWidth, tc.expected)
			} else {
				t.Logf("✅ Lipgloss is correct for this case (bug may be in other cases)")
			}

			if runewidthWidth != tc.expected {
				t.Logf("⚠️  go-runewidth wrong: got %d, want %d", runewidthWidth, tc.expected)
			}
		})
	}
}

func statusEmoji(correct bool) string {
	if correct {
		return "✅"
	}
	return "❌"
}

// TestLipgloss562_VariationSelectors tests emoji with variation selectors (U+FE0F)
func TestLipgloss562_VariationSelectors(t *testing.T) {
	us := NewUnicodeService()

	testCases := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "Text variation selector",
			input:    "☀︎", // Sun + text variant
			expected: 1,    // Should be narrow
		},
		{
			name:     "Emoji variation selector",
			input:    "☀️", // Sun + emoji variant
			expected: 2,    // Should be wide
		},
		{
			name:     "Shield with emoji variant",
			input:    "🛡️",
			expected: 2,
		},
		{
			name:     "Clock (no variation selector needed)",
			input:    "⏰",
			expected: 2,
		},
	}

	t.Log("\n=================================")
	t.Log("Testing variation selectors (U+FE0F)")
	t.Log("=================================")

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			phoenixWidth := us.StringWidth(tc.input)
			lipglossWidth := lipgloss.Width(tc.input)
			runewidthWidth := runewidth.StringWidth(tc.input)

			t.Logf("")
			t.Logf("Input: %q (% X)", tc.input, []byte(tc.input))
			t.Logf("Expected:  %d", tc.expected)
			t.Logf("Phoenix:   %d %s", phoenixWidth, statusEmoji(phoenixWidth == tc.expected))
			t.Logf("Lipgloss:  %d %s", lipglossWidth, statusEmoji(lipglossWidth == tc.expected))
			t.Logf("Runewidth: %d %s", runewidthWidth, statusEmoji(runewidthWidth == tc.expected))

			if phoenixWidth != tc.expected {
				t.Errorf("Phoenix wrong: got %d, want %d", phoenixWidth, tc.expected)
			}

			if lipglossWidth != tc.expected {
				t.Logf("🔴 Lipgloss variation selector issue")
			}
		})
	}
}
