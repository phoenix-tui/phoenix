package core

import (
	"github.com/phoenix-tui/phoenix/core/internal/domain/service"
)

// Unicode service instance (package-level singleton for performance).
var unicodeSvc = service.NewUnicodeService()

// StringWidth returns the visual width of a string in terminal cells.
// This correctly handles Unicode edge cases that many libraries get wrong:
//   - Emoji (🔥, 👍, etc.): width 2
//   - CJK characters (中文, 日本語, 한국어): width 2 per character
//   - Combining characters (é = e + ́): width 0 for combiner
//   - Zero-width joiners (ZWJ): width 0
//   - Control characters: width 0
//   - ASCII: width 1
//
// This function uses the latest Unicode 16.0 data and implements
// the Unicode Standard Annex #11 (East Asian Width) correctly.
//
// Performance: Optimized with tiered lookup (9-23x faster than go-runewidth):
//   - O(1) for ASCII, CJK, simple emoji (90-95% of cases)
//   - O(log n) for rare characters
//   - Grapheme clustering only for complex Unicode (ZWJ, modifiers)
//
// Example:
//
//	core.StringWidth("Hello")        // 5
//	core.StringWidth("Hello 🔥")     // 8 (5 + 1 space + 2 emoji)
//	core.StringWidth("中文")          // 4 (2 + 2)
//	core.StringWidth("Café")         // 4 (C + a + f + é)
//	core.StringWidth("👋🏻")          // 2 (emoji + skin tone modifier)
//
// Use this function when:
//   - Calculating text layout in TUI
//   - Truncating strings to fit terminal width
//   - Aligning text in columns
//   - Rendering bordered boxes
func StringWidth(s string) int {
	return unicodeSvc.StringWidth(s)
}
