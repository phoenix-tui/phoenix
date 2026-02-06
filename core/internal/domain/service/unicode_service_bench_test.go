package service

import (
	"strings"
	"testing"
)

// Benchmark data sets
var (
	// ASCII strings
	benchASCIIShort = "Hello"
	benchASCIIMed   = "The quick brown fox jumps over the lazy dog"
	benchASCIILong  = strings.Repeat("The quick brown fox jumps over the lazy dog. ", 20) // ~900 chars

	// Emoji strings
	benchEmojiShort = "👋😀🎉"
	benchEmojiMed   = "👋😀🎉❤️🚀🌟💡🔥✨🎯"
	benchEmojiLong  = strings.Repeat("👋😀🎉❤️🚀🌟💡🔥✨🎯", 20) // ~200 emoji

	// CJK strings
	benchCJKShort = "你好世界"
	benchCJKMed   = "你好世界，这是一个测试字符串"
	benchCJKLong  = strings.Repeat("你好世界，这是一个测试字符串。", 20) // ~300 chars

	// Mixed content
	benchMixedShort = "Hello 👋 世界"
	benchMixedMed   = "Hello 👋 世界! The quick brown 🦊 jumps over 懒狗"
	benchMixedLong  = strings.Repeat("Hello 👋 世界! ASCII, emoji 🎉, and CJK 字符. ", 20) // ~900 chars

	// Complex Unicode (emoji with modifiers, ZWJ sequences)
	benchComplexEmoji = "👋🏻👨‍👩‍👧‍👦🇺🇸"
	benchCombining    = "Café résumé naïve"

	// Real-world examples (what users actually render)
	benchRealWorld1 = "📝 TODO: Implement feature #123"
	benchRealWorld2 = "✅ Tests passed! 100% coverage 🎉"
	benchRealWorld3 = "🚀 Phoenix TUI v1.0.0 - 性能优化"
)

// ========================================
// StringWidth Benchmarks
// ========================================

func BenchmarkStringWidth_ASCII_Short(b *testing.B) {
	us := NewUnicodeService()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = us.StringWidth(benchASCIIShort)
	}
}

func BenchmarkStringWidth_ASCII_Med(b *testing.B) {
	us := NewUnicodeService()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = us.StringWidth(benchASCIIMed)
	}
}

func BenchmarkStringWidth_ASCII_Long(b *testing.B) {
	us := NewUnicodeService()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = us.StringWidth(benchASCIILong)
	}
}

func BenchmarkStringWidth_Emoji_Short(b *testing.B) {
	us := NewUnicodeService()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = us.StringWidth(benchEmojiShort)
	}
}

func BenchmarkStringWidth_Emoji_Med(b *testing.B) {
	us := NewUnicodeService()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = us.StringWidth(benchEmojiMed)
	}
}

func BenchmarkStringWidth_Emoji_Long(b *testing.B) {
	us := NewUnicodeService()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = us.StringWidth(benchEmojiLong)
	}
}

func BenchmarkStringWidth_CJK_Short(b *testing.B) {
	us := NewUnicodeService()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = us.StringWidth(benchCJKShort)
	}
}

func BenchmarkStringWidth_CJK_Med(b *testing.B) {
	us := NewUnicodeService()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = us.StringWidth(benchCJKMed)
	}
}

func BenchmarkStringWidth_CJK_Long(b *testing.B) {
	us := NewUnicodeService()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = us.StringWidth(benchCJKLong)
	}
}

func BenchmarkStringWidth_Mixed_Short(b *testing.B) {
	us := NewUnicodeService()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = us.StringWidth(benchMixedShort)
	}
}

func BenchmarkStringWidth_Mixed_Med(b *testing.B) {
	us := NewUnicodeService()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = us.StringWidth(benchMixedMed)
	}
}

func BenchmarkStringWidth_Mixed_Long(b *testing.B) {
	us := NewUnicodeService()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = us.StringWidth(benchMixedLong)
	}
}

func BenchmarkStringWidth_ComplexEmoji(b *testing.B) {
	us := NewUnicodeService()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = us.StringWidth(benchComplexEmoji)
	}
}

func BenchmarkStringWidth_Combining(b *testing.B) {
	us := NewUnicodeService()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = us.StringWidth(benchCombining)
	}
}

func BenchmarkStringWidth_RealWorld1(b *testing.B) {
	us := NewUnicodeService()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = us.StringWidth(benchRealWorld1)
	}
}

func BenchmarkStringWidth_RealWorld2(b *testing.B) {
	us := NewUnicodeService()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = us.StringWidth(benchRealWorld2)
	}
}

func BenchmarkStringWidth_RealWorld3(b *testing.B) {
	us := NewUnicodeService()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = us.StringWidth(benchRealWorld3)
	}
}

// ========================================
// ClusterWidth Benchmarks
// ========================================

func BenchmarkClusterWidth_ASCII(b *testing.B) {
	us := NewUnicodeService()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = us.ClusterWidth("A")
	}
}

func BenchmarkClusterWidth_Emoji(b *testing.B) {
	us := NewUnicodeService()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = us.ClusterWidth("👋")
	}
}

func BenchmarkClusterWidth_CJK(b *testing.B) {
	us := NewUnicodeService()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = us.ClusterWidth("中")
	}
}

func BenchmarkClusterWidth_ComplexEmoji(b *testing.B) {
	us := NewUnicodeService()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = us.ClusterWidth("👋🏻")
	}
}

// ========================================
// Memory Allocation Benchmarks
// ========================================

func BenchmarkStringWidth_Memory_ASCII(b *testing.B) {
	us := NewUnicodeService()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = us.StringWidth(benchASCIIMed)
	}
}

func BenchmarkStringWidth_Memory_Emoji(b *testing.B) {
	us := NewUnicodeService()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = us.StringWidth(benchEmojiMed)
	}
}

func BenchmarkStringWidth_Memory_Mixed(b *testing.B) {
	us := NewUnicodeService()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = us.StringWidth(benchMixedMed)
	}
}

// ========================================
// Comparison: Service Creation Overhead
// ========================================

func BenchmarkNewUnicodeService(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewUnicodeService()
	}
}

// ========================================
// Target Verification (1000+ chars < 10μs)
// ========================================

// BenchmarkStringWidth_1000Chars verifies target: < 10μs for 1000 char string
func BenchmarkStringWidth_1000Chars_ASCII(b *testing.B) {
	us := NewUnicodeService()
	longString := strings.Repeat("a", 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = us.StringWidth(longString)
	}
}

// BenchmarkStringWidth_1000Chars_Mixed verifies realistic 1000 char case
func BenchmarkStringWidth_1000Chars_Mixed(b *testing.B) {
	us := NewUnicodeService()
	// Mix of ASCII, emoji, CJK (~1000 visual chars)
	longString := strings.Repeat("Hello 👋 世界! ", 50) // ~750 chars
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = us.StringWidth(longString)
	}
}
