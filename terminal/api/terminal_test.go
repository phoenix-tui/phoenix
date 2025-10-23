package api

import (
	"strings"
	"testing"
)

// ============================================================================
// Platform Type Tests
// ============================================================================

func TestPlatform_String(t *testing.T) {
	testCases := []struct {
		platform Platform
		expected string
	}{
		{PlatformWindowsConsole, "Windows Console (Win32 API)"},
		{PlatformWindowsANSI, "Windows ANSI (Git Bash)"},
		{PlatformUnix, "Unix (ANSI)"},
		{PlatformUnknown, "Unknown"},
		{Platform(999), "Unknown"}, // Invalid value
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			result := tc.platform.String()
			if result != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, result)
			}
		})
	}
}

// ============================================================================
// CursorStyle Type Tests
// ============================================================================

func TestCursorStyle_String(t *testing.T) {
	testCases := []struct {
		style    CursorStyle
		expected string
	}{
		{CursorBlock, "Block"},
		{CursorUnderline, "Underline"},
		{CursorBar, "Bar"},
		{CursorStyle(999), "Unknown"}, // Invalid value
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			result := tc.style.String()
			if result != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, result)
			}
		})
	}
}

// NOTE: Terminal interface implementation tests are in infrastructure package
// to avoid import cycles. This test file focuses on API types only.

// ============================================================================
// String Formatting Tests (Documentation)
// ============================================================================

func TestPlatform_StringsAreDescriptive(t *testing.T) {
	platforms := []Platform{
		PlatformWindowsConsole,
		PlatformWindowsANSI,
		PlatformUnix,
		PlatformUnknown,
	}

	for _, p := range platforms {
		str := p.String()
		if str == "" {
			t.Errorf("Platform %d has empty String()", p)
		}
		if !strings.Contains(str, "Unknown") && p != PlatformUnknown {
			// Valid platforms should have descriptive names
			if len(str) < 4 {
				t.Errorf("Platform %d has too short String: %s", p, str)
			}
		}
	}
}

func TestCursorStyle_StringsAreDescriptive(t *testing.T) {
	styles := []CursorStyle{
		CursorBlock,
		CursorUnderline,
		CursorBar,
	}

	for _, s := range styles {
		str := s.String()
		if str == "" {
			t.Errorf("CursorStyle %d has empty String()", s)
		}
		if str == "Unknown" {
			t.Errorf("Valid CursorStyle %d has 'Unknown' String()", s)
		}
	}
}

// ============================================================================
// Constants Validation Tests
// ============================================================================

func TestPlatformConstants_Distinct(t *testing.T) {
	platforms := []Platform{
		PlatformUnknown,
		PlatformWindowsConsole,
		PlatformWindowsANSI,
		PlatformUnix,
	}

	// Verify all constants are distinct
	seen := make(map[Platform]bool)
	for _, p := range platforms {
		if seen[p] {
			t.Errorf("duplicate Platform constant: %d", p)
		}
		seen[p] = true
	}
}

func TestCursorStyleConstants_Distinct(t *testing.T) {
	styles := []CursorStyle{
		CursorBlock,
		CursorUnderline,
		CursorBar,
	}

	// Verify all constants are distinct
	seen := make(map[CursorStyle]bool)
	for _, s := range styles {
		if seen[s] {
			t.Errorf("duplicate CursorStyle constant: %d", s)
		}
		seen[s] = true
	}
}
