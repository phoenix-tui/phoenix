package ansi

import (
	"strings"
	"testing"
)

func TestEnableMouseTracking(t *testing.T) {
	tests := []struct {
		mode     MouseMode
		expected string
	}{
		{ModeX10, "\x1b[?1000h"},
		{ModeVT200, "\x1b[?1002h"},
		{ModeSGR, "\x1b[?1006h"},
		{ModeURxvt, "\x1b[?1015h"},
		{ModeAnyMotion, "\x1b[?1003h"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := EnableMouseTracking(tt.mode)
			if result != tt.expected {
				t.Errorf("EnableMouseTracking() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestDisableMouseTracking(t *testing.T) {
	tests := []struct {
		mode     MouseMode
		expected string
	}{
		{ModeX10, "\x1b[?1000l"},
		{ModeVT200, "\x1b[?1002l"},
		{ModeSGR, "\x1b[?1006l"},
		{ModeURxvt, "\x1b[?1015l"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := DisableMouseTracking(tt.mode)
			if result != tt.expected {
				t.Errorf("DisableMouseTracking() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestEnableFocusEvents(t *testing.T) {
	expected := "\x1b[?1004h"
	result := EnableFocusEvents()
	if result != expected {
		t.Errorf("EnableFocusEvents() = %q, want %q", result, expected)
	}
}

func TestDisableFocusEvents(t *testing.T) {
	expected := "\x1b[?1004l"
	result := DisableFocusEvents()
	if result != expected {
		t.Errorf("DisableFocusEvents() = %q, want %q", result, expected)
	}
}

func TestEnableMouseAll(t *testing.T) {
	result := EnableMouseAll()

	// Should contain both SGR and ButtonMotion enables
	if !strings.Contains(result, "\x1b[?1006h") {
		t.Error("EnableMouseAll() should contain SGR enable sequence")
	}
	if !strings.Contains(result, "\x1b[?1002h") {
		t.Error("EnableMouseAll() should contain ButtonMotion enable sequence")
	}
}

func TestDisableMouseAll(t *testing.T) {
	result := DisableMouseAll()

	// Should contain both SGR and ButtonMotion disables
	if !strings.Contains(result, "\x1b[?1006l") {
		t.Error("DisableMouseAll() should contain SGR disable sequence")
	}
	if !strings.Contains(result, "\x1b[?1002l") {
		t.Error("DisableMouseAll() should contain ButtonMotion disable sequence")
	}
}

func Test_itoa(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{10, "10"},
		{123, "123"},
		{1000, "1000"},
		{1006, "1006"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := itoa(tt.input)
			if result != tt.expected {
				t.Errorf("itoa(%d) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}
