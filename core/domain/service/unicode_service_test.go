package service

import (
	"testing"

	"github.com/phoenix-tui/phoenix/core/domain/value"
)

func TestNewUnicodeService(t *testing.T) {
	us := NewUnicodeService()
	if us == nil {
		t.Fatal("NewUnicodeService() returned nil")
	}
}

// TestStringWidth_ASCII tests width calculation for ASCII strings
func TestStringWidth_ASCII(t *testing.T) {
	us := NewUnicodeService()

	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"empty string", "", 0},
		{"single char", "a", 1},
		{"word", "Hello", 5},
		{"sentence", "Hello World", 11},
		{"numbers", "12345", 5},
		{"symbols", "!@#$%", 5},
		{"spaces", "a b c", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := us.StringWidth(tt.input)
			if got != tt.want {
				t.Errorf("StringWidth(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

// TestStringWidth_Emoji tests width calculation for emoji
func TestStringWidth_Emoji(t *testing.T) {
	us := NewUnicodeService()

	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"simple emoji", "😀", 2},
		{"waving hand", "👋", 2},
		{"emoji with modifier", "👋🏻", 2}, // Wave + skin tone = 1 cluster, 2 columns
		{"heart", "❤️", 2},
		{"family emoji", "👨‍👩‍👧‍👦", 2}, // ZWJ sequence = 1 cluster, 2 columns
		{"flag", "🇺🇸", 2},              // Regional indicators form 1 emoji, 2 columns
		{"multiple emoji", "😀😃😄", 6},
		{"emoji in text", "Hello 😀 World", 14}, // 5 + 1 + 2 + 1 + 5 = 14
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := us.StringWidth(tt.input)
			if got != tt.want {
				t.Errorf("StringWidth(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

// TestStringWidth_CJK tests width calculation for CJK characters
func TestStringWidth_CJK(t *testing.T) {
	us := NewUnicodeService()

	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"single Chinese", "中", 2},
		{"Chinese word", "你好", 4},
		{"Chinese sentence", "你好世界", 8},
		{"Japanese Hiragana", "こんにちは", 10},
		{"Japanese Katakana", "カタカナ", 8},
		{"Japanese Kanji", "日本語", 6},
		{"Korean Hangul", "한글", 4},
		{"mixed CJK", "中文日本語한글", 14}, // 4 + 6 + 4 = 14
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := us.StringWidth(tt.input)
			if got != tt.want {
				t.Errorf("StringWidth(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

// TestStringWidth_Combining tests width calculation for combining characters
func TestStringWidth_Combining(t *testing.T) {
	us := NewUnicodeService()

	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"e with acute", "é", 1},            // e + combining acute = 1 cluster, 1 column
		{"Cafe", "Café", 4},                 // C + a + f + é = 4 columns
		{"a with umlaut", "ä", 1},           // a + umlaut = 1 cluster, 1 column
		{"n with tilde", "ñ", 1},            // n + tilde = 1 cluster, 1 column
		{"German word", "Müller", 6},        // M + ü + l + l + e + r
		{"Spanish word", "español", 7},      // e + s + p + a + ñ + o + l
		{"diacritic marks", "à è ì ò ù", 9}, // 5 chars + 4 spaces = 9
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := us.StringWidth(tt.input)
			if got != tt.want {
				t.Errorf("StringWidth(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

// TestStringWidth_ZeroWidth tests width calculation for zero-width characters
func TestStringWidth_ZeroWidth(t *testing.T) {
	us := NewUnicodeService()

	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"zero-width space", "a\u200Bb", 2}, // a + ZWS + b = 2 columns
		{"zero-width joiner", "👨‍👩", 2},     // Man + ZWJ + Woman = 1 cluster, 2 columns
		{"BOM", "\uFEFFHello", 5},           // BOM + Hello = 5 columns
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := us.StringWidth(tt.input)
			if got != tt.want {
				t.Errorf("StringWidth(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

// TestStringWidth_Mixed tests width calculation for mixed content
func TestStringWidth_Mixed(t *testing.T) {
	us := NewUnicodeService()

	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"ASCII + Emoji", "Hello 😀", 8},    // 5 + 1 + 2 = 8
		{"ASCII + CJK", "Hello 世界", 10},    // 5 + 1 + 2 + 2 = 10
		{"Emoji + CJK", "😀中文", 6},          // 2 + 2 + 2 = 6
		{"Complex mix", "Hi 👋 世界!", 11},    // 2 + 1 + 2 + 1 + 2 + 2 + 1 = 11
		{"Cafe with emoji", "Café ☕", 7},   // 4 + 1 + 2 = 7
		{"URL", "https://example.com", 19}, // All ASCII
		{"Email", "user@example.com", 16},  // All ASCII
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := us.StringWidth(tt.input)
			if got != tt.want {
				t.Errorf("StringWidth(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

// TestGraphemeClusters_ASCII tests grapheme cluster splitting for ASCII
func TestGraphemeClusters_ASCII(t *testing.T) {
	us := NewUnicodeService()

	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"empty string", "", []string{}},
		{"single char", "a", []string{"a"}},
		{"word", "Hello", []string{"H", "e", "l", "l", "o"}},
		{"spaces", "a b", []string{"a", " ", "b"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := us.GraphemeClusters(tt.input)
			if !slicesEqual(got, tt.want) {
				t.Errorf("GraphemeClusters(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// TestGraphemeClusters_Emoji tests grapheme cluster splitting for emoji
func TestGraphemeClusters_Emoji(t *testing.T) {
	us := NewUnicodeService()

	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"simple emoji", "😀", []string{"😀"}},
		{"emoji with modifier", "👋🏻", []string{"👋🏻"}},
		{"multiple emoji", "😀😃", []string{"😀", "😃"}},
		{"family emoji", "👨‍👩‍👧‍👦", []string{"👨‍👩‍👧‍👦"}},
		{"flag", "🇺🇸", []string{"🇺🇸"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := us.GraphemeClusters(tt.input)
			if !slicesEqual(got, tt.want) {
				t.Errorf("GraphemeClusters(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// TestGraphemeClusters_CJK tests grapheme cluster splitting for CJK
func TestGraphemeClusters_CJK(t *testing.T) {
	us := NewUnicodeService()

	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"Chinese", "你好", []string{"你", "好"}},
		{"Japanese", "こんにちは", []string{"こ", "ん", "に", "ち", "は"}},
		{"Korean", "한글", []string{"한", "글"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := us.GraphemeClusters(tt.input)
			if !slicesEqual(got, tt.want) {
				t.Errorf("GraphemeClusters(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// TestGraphemeClusters_Combining tests grapheme cluster splitting for combining chars
func TestGraphemeClusters_Combining(t *testing.T) {
	us := NewUnicodeService()

	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"e with acute", "é", []string{"é"}},
		{"Cafe", "Café", []string{"C", "a", "f", "é"}},
		{"German word", "Müller", []string{"M", "ü", "l", "l", "e", "r"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := us.GraphemeClusters(tt.input)
			if !slicesEqual(got, tt.want) {
				t.Errorf("GraphemeClusters(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// TestClusterWidth_ASCII tests width calculation for ASCII clusters
func TestClusterWidth_ASCII(t *testing.T) {
	us := NewUnicodeService()

	tests := []struct {
		name    string
		cluster string
		want    int
	}{
		{"empty", "", 0},
		{"letter", "a", 1},
		{"digit", "5", 1},
		{"space", " ", 1},
		{"symbol", "$", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := us.ClusterWidth(tt.cluster)
			if got != tt.want {
				t.Errorf("ClusterWidth(%q) = %d, want %d", tt.cluster, got, tt.want)
			}
		})
	}
}

// TestClusterWidth_Emoji tests width calculation for emoji clusters
func TestClusterWidth_Emoji(t *testing.T) {
	us := NewUnicodeService()

	tests := []struct {
		name    string
		cluster string
		want    int
	}{
		{"smile", "😀", 2},
		{"wave", "👋", 2},
		{"wave + modifier", "👋🏻", 2},
		{"heart", "❤️", 2},
		{"family", "👨‍👩‍👧‍👦", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := us.ClusterWidth(tt.cluster)
			if got != tt.want {
				t.Errorf("ClusterWidth(%q) = %d, want %d", tt.cluster, got, tt.want)
			}
		})
	}
}

// TestClusterWidth_CJK tests width calculation for CJK clusters
func TestClusterWidth_CJK(t *testing.T) {
	us := NewUnicodeService()

	tests := []struct {
		name    string
		cluster string
		want    int
	}{
		{"Chinese", "中", 2},
		{"Japanese Hiragana", "あ", 2},
		{"Japanese Katakana", "ア", 2},
		{"Korean Hangul", "한", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := us.ClusterWidth(tt.cluster)
			if got != tt.want {
				t.Errorf("ClusterWidth(%q) = %d, want %d", tt.cluster, got, tt.want)
			}
		})
	}
}

// TestClusterWidth_Combining tests width calculation for combining characters
func TestClusterWidth_Combining(t *testing.T) {
	us := NewUnicodeService()

	tests := []struct {
		name    string
		cluster string
		want    int
	}{
		{"e with acute", "é", 1},
		{"a with umlaut", "ä", 1},
		{"n with tilde", "ñ", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := us.ClusterWidth(tt.cluster)
			if got != tt.want {
				t.Errorf("ClusterWidth(%q) = %d, want %d", tt.cluster, got, tt.want)
			}
		})
	}
}

// TestClusterWidth_ZeroWidth tests width calculation for zero-width characters
func TestClusterWidth_ZeroWidth(t *testing.T) {
	us := NewUnicodeService()

	tests := []struct {
		name    string
		cluster string
		want    int
	}{
		{"zero-width space", "\u200B", 0},
		{"BOM", "\uFEFF", 0},
		{"combining acute", "\u0301", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := us.ClusterWidth(tt.cluster)
			if got != tt.want {
				t.Errorf("ClusterWidth(%q) = %d, want %d", tt.cluster, got, tt.want)
			}
		})
	}
}

// TestClusterWidth_ControlCharacters tests width calculation for control characters
func TestClusterWidth_ControlCharacters(t *testing.T) {
	us := NewUnicodeService()

	tests := []struct {
		name    string
		cluster string
		want    int
	}{
		{"newline", "\n", 0},
		{"tab", "\t", 0},
		{"carriage return", "\r", 0},
		{"null", "\x00", 0},
		{"bell", "\x07", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := us.ClusterWidth(tt.cluster)
			if got != tt.want {
				t.Errorf("ClusterWidth(%q) = %d, want %d", tt.cluster, got, tt.want)
			}
		})
	}
}

// TestIsZeroWidth tests zero-width character detection
func TestIsZeroWidth(t *testing.T) {
	tests := []struct {
		name string
		r    rune
		want bool
	}{
		{"regular ASCII", 'a', false},
		{"space", ' ', false},
		{"zero-width space", '\u200B', true},
		{"BOM", '\uFEFF', true},
		{"combining acute", '\u0301', true},     // Mn
		{"combining tilde", '\u0303', true},     // Mn
		{"combining enclosing", '\u20DD', true}, // Me
		{"combining spacing", '\u0903', true},   // Mc
		{"format character", '\u200E', true},    // Cf - Left-to-right mark
		{"soft hyphen", '\u00AD', true},         // Cf
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isZeroWidth(tt.r)
			if got != tt.want {
				t.Errorf("isZeroWidth(%U) = %v, want %v", tt.r, got, tt.want)
			}
		})
	}
}

// TestIsEmoji tests emoji detection
func TestIsEmoji(t *testing.T) {
	tests := []struct {
		name string
		r    rune
		want bool
	}{
		{"regular ASCII", 'a', false},
		{"smile emoji", '😀', true}, // 0x1F600
		{"wave emoji", '👋', true},  // 0x1F44B
		{"heart", '❤', true},       // 0x2764
		{"sun", '☀', true},         // 0x2600
		{"Chinese", '中', false},
		{"space", ' ', false},
		{"regional indicator", '\U0001F1FA', true}, // 0x1F1FA
		{"mahjong tile", '\U0001F000', true},       // 0x1F000
		{"rocket", '\U0001F680', true},             // 0x1F680
		{"pizza", '\U0001F355', true},              // 0x1F355
		{"game die", '\U0001F3B2', true},           // 0x1F3B2
		{"yawn face", '\U0001F971', true},          // 0x1F971 (Supplemental)
		{"yo-yo", '\U0001FA80', true},              // 0x1FA80 (Extended Pictographic)
		{"scissors", '\u2702', true},               // 0x2702 (Dingbats)
		{"checkmark", '\u2714', true},              // 0x2714 (Dingbats)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isEmoji(tt.r)
			if got != tt.want {
				t.Errorf("isEmoji(%U) = %v, want %v", tt.r, got, tt.want)
			}
		})
	}
}

// TestIsCJK tests CJK character detection
func TestIsCJK(t *testing.T) {
	tests := []struct {
		name string
		r    rune
		want bool
	}{
		{"regular ASCII", 'a', false},
		{"Chinese Unified", '中', true},   // 0x4E00-0x9FFF
		{"Japanese Hiragana", 'あ', true}, // 0x3040-0x30FF
		{"Japanese Katakana", 'ア', true}, // 0x3040-0x30FF
		{"Korean Hangul", '한', true},     // 0xAC00-0xD7AF
		{"emoji", '😀', false},
		{"space", ' ', false},
		{"CJK Extension A", '\u3400', true},     // 0x3400-0x4DBF
		{"CJK Extension B", '\U00020000', true}, // 0x20000-0x2EBEF
		{"Fullwidth A", '\uFF21', true},         // 0xFF00-0xFFEF
		{"Hiragana A", '\u3042', true},          // 0x3040-0x309F
		{"Katakana A", '\u30A2', true},          // 0x30A0-0x30FF
		{"Hangul Syllable", '\uAC01', true},     // 0xAC00-0xD7AF
		{"CJK Unified start", '\u4E00', true},   // Start of 0x4E00-0x9FFF
		{"CJK Unified end", '\u9FFF', true},     // End of 0x4E00-0x9FFF
		{"CJK Ext A start", '\u3400', true},     // Start of 0x3400-0x4DBF
		{"CJK Ext A end", '\u4DBF', true},       // End of 0x3400-0x4DBF
		{"Hiragana start", '\u3040', true},      // Start of 0x3040-0x309F
		{"Hiragana end", '\u309F', true},        // End of 0x3040-0x309F
		{"Katakana start", '\u30A0', true},      // Start of 0x30A0-0x30FF
		{"Katakana end", '\u30FF', true},        // End of 0x30A0-0x30FF
		{"Hangul start", '\uAC00', true},        // Start of 0xAC00-0xD7AF
		{"Hangul end", '\uD7AF', true},          // End of 0xAC00-0xD7AF
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isCJK(tt.r)
			if got != tt.want {
				t.Errorf("isCJK(%U) = %v, want %v", tt.r, got, tt.want)
			}
		})
	}
}

// Helper function to compare string slices
func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestStringWidthWithConfig_EastAsianAmbiguous(t *testing.T) {
	us := NewUnicodeService()

	// Test characters: ± (U+00B1), ½ (U+00BD), ° (U+00B0), × (U+00D7)
	testCases := []struct {
		name           string
		input          string
		narrowExpected int
		wideExpected   int
	}{
		{
			name:           "Plus-minus sign",
			input:          "±",
			narrowExpected: 1,
			wideExpected:   2,
		},
		{
			name:           "Half fraction",
			input:          "½",
			narrowExpected: 1,
			wideExpected:   2,
		},
		{
			name:           "Degree sign",
			input:          "°",
			narrowExpected: 1,
			wideExpected:   2,
		},
		{
			name:           "Multiplication sign",
			input:          "×",
			narrowExpected: 1,
			wideExpected:   2,
		},
		{
			name:           "Multiple ambiguous chars",
			input:          "±½°×",
			narrowExpected: 4,
			wideExpected:   8,
		},
		{
			name:           "Mixed with ASCII",
			input:          "Test ± value",
			narrowExpected: 12, // Test(4) + space(1) + ±(1) + space(1) + value(5)
			wideExpected:   13, // Test(4) + space(1) + ±(2) + space(1) + value(5)
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test narrow (default English locale)
			narrowConfig := value.NewUnicodeConfig()
			narrowWidth := us.StringWidthWithConfig(tc.input, narrowConfig)
			if narrowWidth != tc.narrowExpected {
				t.Errorf("Narrow: got %d, want %d", narrowWidth, tc.narrowExpected)
			}

			// Test wide (East Asian locale)
			wideConfig := value.NewUnicodeConfig().WithEastAsianWide()
			wideWidth := us.StringWidthWithConfig(tc.input, wideConfig)
			if wideWidth != tc.wideExpected {
				t.Errorf("Wide: got %d, want %d", wideWidth, tc.wideExpected)
			}
		})
	}
}

func TestStringWidthWithConfig_NonAmbiguous(t *testing.T) {
	us := NewUnicodeService()

	// Test that non-ambiguous characters are unaffected by config
	testCases := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "ASCII",
			input:    "Hello",
			expected: 5,
		},
		{
			name:     "CJK",
			input:    "你好",
			expected: 4,
		},
		{
			name:     "Emoji",
			input:    "👋",
			expected: 2,
		},
		{
			name:     "Emoji with modifier",
			input:    "👋🏻",
			expected: 2,
		},
	}

	narrowConfig := value.NewUnicodeConfig()
	wideConfig := value.NewUnicodeConfig().WithEastAsianWide()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			narrowWidth := us.StringWidthWithConfig(tc.input, narrowConfig)
			wideWidth := us.StringWidthWithConfig(tc.input, wideConfig)

			if narrowWidth != tc.expected {
				t.Errorf("Narrow config: got %d, want %d", narrowWidth, tc.expected)
			}

			if wideWidth != tc.expected {
				t.Errorf("Wide config: got %d, want %d", wideWidth, tc.expected)
			}

			// Narrow and wide should be same for non-ambiguous
			if narrowWidth != wideWidth {
				t.Errorf("Non-ambiguous characters should have same width in both configs: narrow=%d, wide=%d", narrowWidth, wideWidth)
			}
		})
	}
}
