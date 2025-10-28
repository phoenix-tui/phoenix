package comparison_test

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"github.com/phoenix-tui/phoenix/core"
)

// ========================================
// Comparison: Phoenix vs Lipgloss vs go-runewidth
// ========================================

// Test data for comparisons
var (
	compASCII = "The quick brown fox jumps over the lazy dog"
	compEmoji = "👋😀🎉❤️🚀"
	compCJK   = "你好世界，这是测试"
	compMixed = "Hello 👋 世界! Test 🎉"
	compLong  = strings.Repeat("Hello 👋 世界 ", 50) // ~750 chars
	compBuggy = "📝 Test"                          // Lipgloss #562 bug case
)

// ========================================
// ASCII Benchmarks (should be similar)
// ========================================

func BenchmarkComparison_ASCII_Phoenix(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = core.StringWidth(compASCII)
	}
}

func BenchmarkComparison_ASCII_Lipgloss(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = lipgloss.Width(compASCII)
	}
}

func BenchmarkComparison_ASCII_Runewidth(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = runewidth.StringWidth(compASCII)
	}
}

// ========================================
// Emoji Benchmarks (Phoenix should be correct)
// ========================================

func BenchmarkComparison_Emoji_Phoenix(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = core.StringWidth(compEmoji)
	}
}

func BenchmarkComparison_Emoji_Lipgloss(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = lipgloss.Width(compEmoji)
	}
}

func BenchmarkComparison_Emoji_Runewidth(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = runewidth.StringWidth(compEmoji)
	}
}

// ========================================
// CJK Benchmarks (should be similar)
// ========================================

func BenchmarkComparison_CJK_Phoenix(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = core.StringWidth(compCJK)
	}
}

func BenchmarkComparison_CJK_Lipgloss(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = lipgloss.Width(compCJK)
	}
}

func BenchmarkComparison_CJK_Runewidth(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = runewidth.StringWidth(compCJK)
	}
}

// ========================================
// Mixed Content Benchmarks (realistic use case)
// ========================================

func BenchmarkComparison_Mixed_Phoenix(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = core.StringWidth(compMixed)
	}
}

func BenchmarkComparison_Mixed_Lipgloss(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = lipgloss.Width(compMixed)
	}
}

func BenchmarkComparison_Mixed_Runewidth(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = runewidth.StringWidth(compMixed)
	}
}

// ========================================
// Long String Benchmarks (performance test)
// ========================================

func BenchmarkComparison_Long_Phoenix(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = core.StringWidth(compLong)
	}
}

func BenchmarkComparison_Long_Lipgloss(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = lipgloss.Width(compLong)
	}
}

func BenchmarkComparison_Long_Runewidth(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = runewidth.StringWidth(compLong)
	}
}

// ========================================
// Lipgloss #562 Bug Case (Phoenix fixes this!)
// ========================================

func BenchmarkComparison_BugCase_Phoenix(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = core.StringWidth(compBuggy)
	}
}

func BenchmarkComparison_BugCase_Lipgloss(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = lipgloss.Width(compBuggy)
	}
}

func BenchmarkComparison_BugCase_Runewidth(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = runewidth.StringWidth(compBuggy)
	}
}

// ========================================
// Memory Comparison (allocations)
// ========================================

func BenchmarkComparison_Memory_Phoenix(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = core.StringWidth(compMixed)
	}
}

func BenchmarkComparison_Memory_Lipgloss(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = lipgloss.Width(compMixed)
	}
}

func BenchmarkComparison_Memory_Runewidth(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = runewidth.StringWidth(compMixed)
	}
}
