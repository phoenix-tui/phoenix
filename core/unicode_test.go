package core_test

import (
	"testing"

	"github.com/phoenix-tui/phoenix/core"
)

func TestStringWidth(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		// ASCII
		{"empty", "", 0},
		{"ascii", "Hello", 5},
		{"ascii with space", "Hello World", 11},

		// Emoji
		{"single emoji", "🔥", 2},
		{"emoji with text", "Hello 🔥", 8},
		{"multiple emoji", "🔥👍", 4},
		{"emoji with modifier", "👋🏻", 2}, // Emoji + skin tone

		// CJK
		{"chinese", "中文", 4},
		{"japanese", "日本語", 6},
		{"korean", "한국어", 6},

		// Combining characters
		{"combining acute", "Café", 4}, // é = e + combining acute

		// Mixed
		{"mixed", "Hello 中文 🔥", 13}, // 5 + 1 + 4 + 1 + 2
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := core.StringWidth(tt.input)
			if got != tt.want {
				t.Errorf("StringWidth(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

// Benchmark public API function
func BenchmarkStringWidth(b *testing.B) {
	tests := []struct {
		name  string
		input string
	}{
		{"ascii", "Hello World"},
		{"emoji", "Hello 🔥 World"},
		{"cjk", "中文测试"},
		{"mixed", "Hello 中文 🔥"},
		{"complex", "👋🏻 Hello 世界"},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				core.StringWidth(tt.input)
			}
		})
	}
}
