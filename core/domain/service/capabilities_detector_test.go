package service_test

import (
	"testing"

	"github.com/phoenix-tui/phoenix/core/domain/service"
	"github.com/phoenix-tui/phoenix/core/domain/value"
)

// MockEnvironment is a test double for EnvironmentProvider
type MockEnvironment struct {
	vars     map[string]string
	platform string
}

func NewMockEnvironment(platform string) *MockEnvironment {
	return &MockEnvironment{
		vars:     make(map[string]string),
		platform: platform,
	}
}

func (m *MockEnvironment) Get(key string) string {
	return m.vars[key]
}

func (m *MockEnvironment) Platform() string {
	return m.platform
}

func (m *MockEnvironment) Set(key, value string) {
	m.vars[key] = value
}

// Test cases based on tui-research-analyst recommendations
func TestCapabilitiesDetector_NO_COLOR(t *testing.T) {
	env := NewMockEnvironment("linux")
	env.Set("NO_COLOR", "1")
	env.Set("TERM", "xterm-256color")
	env.Set("COLORTERM", "truecolor")

	detector := service.NewCapabilitiesDetector(env)
	caps := detector.Detect()

	// Business rule: NO_COLOR disables everything
	if caps.SupportsANSI() {
		t.Error("NO_COLOR should disable ANSI")
	}
	if caps.ColorDepth() != value.ColorDepthNone {
		t.Errorf("NO_COLOR should set ColorDepthNone, got %v", caps.ColorDepth())
	}
	if caps.SupportsMouse() {
		t.Error("NO_COLOR should disable mouse")
	}
}

func TestCapabilitiesDetector_FORCE_COLOR(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		wantDepth value.ColorDepth
		wantANSI  bool
	}{
		{"force 0", "0", value.ColorDepthNone, false},
		{"force 1", "1", value.ColorDepth8, true},
		{"force 2", "2", value.ColorDepth256, true},
		{"force 3", "3", value.ColorDepthTrueColor, true},
		{"force true", "true", value.ColorDepthTrueColor, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewMockEnvironment("linux")
			env.Set("FORCE_COLOR", tt.value)
			env.Set("TERM", "dumb") // Would normally disable everything

			detector := service.NewCapabilitiesDetector(env)
			caps := detector.Detect()

			if caps.ColorDepth() != tt.wantDepth {
				t.Errorf("expected depth %v, got %v", tt.wantDepth, caps.ColorDepth())
			}
			if caps.SupportsANSI() != tt.wantANSI {
				t.Errorf("expected ANSI %v, got %v", tt.wantANSI, caps.SupportsANSI())
			}
		})
	}
}

func TestCapabilitiesDetector_COLORTERM(t *testing.T) {
	tests := []struct {
		name      string
		colorterm string
		wantDepth value.ColorDepth
	}{
		{"truecolor", "truecolor", value.ColorDepthTrueColor},
		{"24bit", "24bit", value.ColorDepthTrueColor},
		{"TRUECOLOR uppercase", "TRUECOLOR", value.ColorDepthTrueColor},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewMockEnvironment("linux")
			env.Set("TERM", "xterm")
			env.Set("COLORTERM", tt.colorterm)

			detector := service.NewCapabilitiesDetector(env)
			caps := detector.Detect()

			if caps.ColorDepth() != tt.wantDepth {
				t.Errorf("expected %v, got %v", tt.wantDepth, caps.ColorDepth())
			}
		})
	}
}

func TestCapabilitiesDetector_TERM_PROGRAM(t *testing.T) {
	tests := []struct {
		name        string
		termProgram string
		wantDepth   value.ColorDepth
	}{
		{"iTerm2", "iTerm.app", value.ColorDepthTrueColor},
		{"VS Code", "vscode", value.ColorDepthTrueColor},
		{"Hyper", "Hyper", value.ColorDepthTrueColor},
		{"Warp", "WarpTerminal", value.ColorDepthTrueColor},
		{"Apple Terminal", "Apple_Terminal", value.ColorDepth256},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewMockEnvironment("darwin")
			env.Set("TERM", "xterm")
			env.Set("TERM_PROGRAM", tt.termProgram)

			detector := service.NewCapabilitiesDetector(env)
			caps := detector.Detect()

			if caps.ColorDepth() != tt.wantDepth {
				t.Errorf("expected %v, got %v", tt.wantDepth, caps.ColorDepth())
			}
		})
	}
}

func TestCapabilitiesDetector_TERM_Parsing(t *testing.T) {
	tests := []struct {
		name      string
		term      string
		wantDepth value.ColorDepth
		wantANSI  bool
	}{
		{"xterm-256color", "xterm-256color", value.ColorDepth256, true},
		{"xterm", "xterm", value.ColorDepth8, true},
		{"screen-256color", "screen-256color", value.ColorDepth256, true},
		{"tmux-256color", "tmux-256color", value.ColorDepth256, true},
		{"dumb", "dumb", value.ColorDepthNone, false},
		{"empty", "", value.ColorDepthNone, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewMockEnvironment("linux")
			env.Set("TERM", tt.term)

			detector := service.NewCapabilitiesDetector(env)
			caps := detector.Detect()

			if caps.ColorDepth() != tt.wantDepth {
				t.Errorf("expected depth %v, got %v", tt.wantDepth, caps.ColorDepth())
			}
			if caps.SupportsANSI() != tt.wantANSI {
				t.Errorf("expected ANSI %v, got %v", tt.wantANSI, caps.SupportsANSI())
			}
		})
	}
}

func TestCapabilitiesDetector_Windows(t *testing.T) {
	tests := []struct {
		name      string
		wtSession string
		termProg  string
		wantDepth value.ColorDepth
	}{
		{"Windows Terminal", "guid-123", "", value.ColorDepthTrueColor},
		{"VS Code", "", "vscode", value.ColorDepthTrueColor},
		{"Unknown", "", "", value.ColorDepth8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewMockEnvironment("windows")
			env.Set("TERM", "xterm")
			if tt.wtSession != "" {
				env.Set("WT_SESSION", tt.wtSession)
			}
			if tt.termProg != "" {
				env.Set("TERM_PROGRAM", tt.termProg)
			}

			detector := service.NewCapabilitiesDetector(env)
			caps := detector.Detect()

			if caps.ColorDepth() != tt.wantDepth {
				t.Errorf("expected %v, got %v", tt.wantDepth, caps.ColorDepth())
			}
		})
	}
}

func TestCapabilitiesDetector_Priority(t *testing.T) {
	// Test that NO_COLOR overrides everything
	env := NewMockEnvironment("linux")
	env.Set("NO_COLOR", "1")
	env.Set("FORCE_COLOR", "3")
	env.Set("COLORTERM", "truecolor")
	env.Set("TERM", "xterm-256color")

	detector := service.NewCapabilitiesDetector(env)
	caps := detector.Detect()

	if caps.ColorDepth() != value.ColorDepthNone {
		t.Error("NO_COLOR should override everything")
	}
}

func TestCapabilitiesDetector_Immutability(t *testing.T) {
	// Verify detector is stateless
	env := NewMockEnvironment("linux")
	env.Set("TERM", "xterm-256color")

	detector := service.NewCapabilitiesDetector(env)

	caps1 := detector.Detect()
	caps2 := detector.Detect()

	// Should return same values (stateless)
	if caps1.ColorDepth() != caps2.ColorDepth() {
		t.Error("detector should be stateless")
	}

	// Change environment
	env.Set("COLORTERM", "truecolor")

	caps3 := detector.Detect()

	// Should reflect new environment
	if caps3.ColorDepth() != value.ColorDepthTrueColor {
		t.Error("detector should use current environment")
	}
}
