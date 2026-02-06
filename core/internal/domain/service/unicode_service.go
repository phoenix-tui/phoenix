package service

import (
	"unicode"

	"github.com/phoenix-tui/phoenix/core/internal/domain/value"
	"github.com/unilibs/uniwidth"
)

// UnicodeService provides Unicode text analysis for correct width calculation.
// This is a domain service because width calculation is core business logic
// needed by Cell value object and ALL Phoenix libraries.
//
// This service fixes Charm's lipgloss#562 bug by correctly calculating
// visual width of grapheme clusters including emoji, CJK, and combining chars.
type UnicodeService struct{}

// NewUnicodeService creates a new Unicode service instance.
func NewUnicodeService() *UnicodeService {
	return &UnicodeService{}
}

// StringWidth calculates the visual width of a string in terminal columns.
// Correctly handles:
//   - ASCII: 1 column each
//   - Emoji: 2 columns (including modifiers, ZWJ sequences)
//   - CJK characters: 2 columns
//   - Zero-width characters: 0 columns
//   - Combining characters: 0 columns
//
// Powered by uniwidth which provides tiered O(1) lookup for all Unicode categories
// including ZWJ sequences, emoji modifiers, and variation selectors.
//
// Example:
//
//	StringWidth("Hello")        // 5
//	StringWidth("👋")            // 2
//	StringWidth("👋🏻")           // 2 (emoji + modifier = 1 cluster, 2 columns)
//	StringWidth("こんにちは")      // 10 (5 CJK chars * 2 columns)
//	StringWidth("Café")         // 4 (C + a + f + é)
func (us *UnicodeService) StringWidth(s string) int {
	if s == "" {
		return 0
	}
	return uniwidth.StringWidth(s)
}

// ClusterWidth calculates the visual width of a single grapheme cluster.
// Returns:
//   - 0 for zero-width/combining characters
//   - 1 for ASCII and most characters
//   - 2 for emoji, CJK characters
//
// A grapheme cluster is a user-perceived character that may consist of multiple runes:
//   - Simple: "a" (1 rune) → width 1
//   - Emoji: "👋" (1 rune) → width 2
//   - Emoji + modifier: "👋🏻" (2 runes) → width 2
//   - ZWJ sequence: "👨‍👩‍👧‍👦" (7 runes) → width 2
//   - Combining: "é" (2 runes: e + combining acute) → width 1
//
// Example:
//
//	ClusterWidth("a")      // 1
//	ClusterWidth("👋")     // 2
//	ClusterWidth("👋🏻")    // 2 (emoji with modifier)
//	ClusterWidth("👨‍👩‍👧‍👦") // 2 (ZWJ sequence)
//	ClusterWidth("中")     // 2 (CJK)
//	ClusterWidth("é")      // 1 (e + combining acute)
//	ClusterWidth("\u0301") // 0 (combining acute accent alone)
func (us *UnicodeService) ClusterWidth(cluster string) int {
	if cluster == "" {
		return 0
	}
	return uniwidth.StringWidth(cluster)
}

// isZeroWidth checks if a rune is zero-width.
// Zero-width characters include:
//   - Combining marks (Mn, Me, Mc)
//   - Format characters (Cf)
//   - Non-spacing marks
func isZeroWidth(r rune) bool {
	// Combining marks
	if unicode.In(r, unicode.Mn, unicode.Me, unicode.Mc) {
		return true
	}

	// Format characters
	if unicode.In(r, unicode.Cf) {
		return true
	}

	// Zero-width space
	if r == '\u200B' || r == '\uFEFF' {
		return true
	}

	return false
}

// isEmoji checks if a rune is an emoji.
// Emoji detection based on Unicode emoji properties and PR #563 ranges.
// Includes:
//   - Emoticons (U+1F600 - U+1F64F)
//   - Misc Symbols and Pictographs (U+1F300 - U+1F5FF)
//   - Transport and Map Symbols (U+1F680 - U+1F6FF)
//   - Miscellaneous Technical (U+2300 - U+23FF) - includes ⏰ clock!
//   - Miscellaneous Symbols (U+2600 - U+26FF)
//   - Dingbats (U+2700 - U+27BF)
//   - Regional Indicator Symbols (U+1F1E6 - U+1F1FF) for flags
//   - Extended Pictographic
func isEmoji(r rune) bool {
	// PR #563 emoji ranges (tested and correct)
	if (r >= 0x1F600 && r <= 0x1F64F) || // Emoticons
		(r >= 0x1F300 && r <= 0x1F5FF) || // Misc Symbols and Pictographs
		(r >= 0x1F680 && r <= 0x1F6FF) || // Transport and Map Symbols
		(r >= 0x1F700 && r <= 0x1F77F) || // Alchemical Symbols
		(r >= 0x2300 && r <= 0x23FF) || // Miscellaneous Technical (⏰ clocks!)
		(r >= 0x2600 && r <= 0x26FF) || // Miscellaneous Symbols
		(r >= 0x2700 && r <= 0x27BF) { // Dingbats
		return true
	}

	// Regional Indicator Symbols (flags like 🇺🇸)
	if r >= 0x1F1E6 && r <= 0x1F1FF {
		return true
	}

	// Additional emoji blocks
	if (r >= 0x1F000 && r <= 0x1F02F) || // Mahjong, Domino tiles
		(r >= 0x1F900 && r <= 0x1F9FF) || // Supplemental Symbols and Pictographs
		(r >= 0x1FA00 && r <= 0x1FAFF) { // Extended Pictographic
		return true
	}

	return false
}

// isCJK checks if a rune is a CJK (Chinese, Japanese, Korean) character.
// CJK characters are wide (2 columns) in terminals.
// Ranges based on Unicode standard:
//   - CJK Unified Ideographs (U+4E00 - U+9FFF)
//   - CJK Extension A (U+3400 - U+4DBF)
//   - CJK Extension B+ (U+20000 - U+2EBEF)
//   - Hangul Syllables (U+AC00 - U+D7AF)
//   - Hiragana/Katakana (U+3040 - U+30FF)
func isCJK(r rune) bool {
	// CJK Unified Ideographs
	if r >= 0x4E00 && r <= 0x9FFF {
		return true
	}

	// CJK Extension A
	if r >= 0x3400 && r <= 0x4DBF {
		return true
	}

	// CJK Extensions B, C, D, E, F, G
	if r >= 0x20000 && r <= 0x2EBEF {
		return true
	}

	// Hangul Syllables
	if r >= 0xAC00 && r <= 0xD7AF {
		return true
	}

	// Hiragana and Katakana
	if r >= 0x3040 && r <= 0x30FF {
		return true
	}

	// Halfwidth and Fullwidth Forms (fullwidth only)
	if r >= 0xFF00 && r <= 0xFFEF {
		return true
	}

	return false
}

// StringWidthWithConfig calculates the visual width of a string with custom Unicode configuration.
// This allows locale-specific width calculation, particularly for East Asian Ambiguous characters.
//
// East Asian Ambiguous characters (±, ½, °, ×, §, etc.) have different widths in different locales:
// - Narrow (width 1): Default for English and neutral locales
// - Wide (width 2): For East Asian locales (Japanese, Chinese, Korean)
//
// Example:
//
//	// English locale (default)
//	config := value.NewUnicodeConfig()
//	width := us.StringWidthWithConfig("±", config)  // 1
//
//	// Japanese locale
//	config := value.NewUnicodeConfig().WithEastAsianWide()
//	width := us.StringWidthWithConfig("±", config)  // 2
//
// For most use cases, use StringWidth() which uses neutral locale defaults.
// Use StringWidthWithConfig() when you need locale-specific rendering.
func (us *UnicodeService) StringWidthWithConfig(s string, config value.UnicodeConfig) int {
	if s == "" {
		return 0
	}

	// Base width (handles emoji/ZWJ/modifiers correctly via grapheme awareness)
	width := uniwidth.StringWidth(s)

	// Adjust for East Asian Ambiguous characters in wide mode.
	// We compute per-rune delta because StringWidthWithOptions lacks
	// grapheme awareness for emoji modifiers/ZWJ sequences.
	if config.IsEastAsianWide() {
		for _, r := range s {
			wideW := uniwidth.RuneWidthWithOptions(r,
				uniwidth.WithEastAsianAmbiguous(config.EastAsianAmbiguous()))
			narrowW := uniwidth.RuneWidth(r)
			if wideW > narrowW {
				width += wideW - narrowW
			}
		}
	}

	return width
}

// ClusterWidthWithConfig calculates the width of a grapheme cluster with custom configuration.
// This is the locale-aware version of ClusterWidth().
//
// Example:
//
//	// English locale
//	config := value.NewUnicodeConfig()
//	width := us.ClusterWidthWithConfig("±", config)  // 1
//
//	// Japanese locale
//	config := value.NewUnicodeConfig().WithEastAsianWide()
//	width := us.ClusterWidthWithConfig("±", config)  // 2
func (us *UnicodeService) ClusterWidthWithConfig(cluster string, config value.UnicodeConfig) int {
	if cluster == "" {
		return 0
	}

	// Base width (handles emoji/ZWJ/modifiers correctly via grapheme awareness)
	width := uniwidth.StringWidth(cluster)

	// Adjust for East Asian Ambiguous characters in wide mode.
	// We compute per-rune delta because StringWidthWithOptions lacks
	// grapheme awareness for emoji modifiers/ZWJ sequences.
	if config.IsEastAsianWide() {
		for _, r := range cluster {
			wideW := uniwidth.RuneWidthWithOptions(r,
				uniwidth.WithEastAsianAmbiguous(config.EastAsianAmbiguous()))
			narrowW := uniwidth.RuneWidth(r)
			if wideW > narrowW {
				width += wideW - narrowW
			}
		}
	}

	return width
}
