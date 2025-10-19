// Package service provides domain services for terminal capabilities detection.
package service

import (
	"strings"

	"github.com/phoenix-tui/phoenix/core/domain/value"
)

// CapabilitiesDetector is a domain service that detects terminal capabilities.
// Based on tui-research-analyst recommendations.
//
// Detection Priority:
//  1. NO_COLOR → Disable all
//  2. FORCE_COLOR → User override
//  3. Platform-specific detection
//  4. COLORTERM → Explicit color
//  5. TERM_PROGRAM → Known terminals
//  6. TERM → Parsing
//  7. Conservative defaults
type CapabilitiesDetector struct {
	env EnvironmentProvider
}

// NewCapabilitiesDetector creates detector with environment provider.
func NewCapabilitiesDetector(env EnvironmentProvider) *CapabilitiesDetector {
	return &CapabilitiesDetector{env: env}
}

// Detect analyzes environment and returns capabilities.
func (cd *CapabilitiesDetector) Detect() *value.Capabilities {
	// Priority 1: NO_COLOR
	if cd.env.Get("NO_COLOR") != "" {
		return value.NewCapabilities(false, value.ColorDepthNone, false, false, false)
	}

	// Priority 2: FORCE_COLOR
	if fc := cd.env.Get("FORCE_COLOR"); fc != "" {
		return cd.parseForceColor(fc)
	}

	// Priority 3: Platform-specific
	switch cd.env.Platform() {
	case "windows":
		return cd.detectWindows()
	case "darwin":
		return cd.detectMacOS()
	default:
		return cd.detectUnix()
	}
}

func (cd *CapabilitiesDetector) detectColorDepth() value.ColorDepth {
	// COLORTERM explicit
	ct := strings.ToLower(cd.env.Get("COLORTERM"))
	if ct == "truecolor" || ct == "24bit" {
		return value.ColorDepthTrueColor
	}

	// TERM_PROGRAM known terminals
	switch cd.env.Get("TERM_PROGRAM") {
	case "iTerm.app", "vscode", "Hyper", "WarpTerminal":
		return value.ColorDepthTrueColor
	case "Apple_Terminal":
		return value.ColorDepth256
	}

	// TERM parsing
	term := cd.env.Get("TERM")
	if strings.Contains(term, "256color") {
		return value.ColorDepth256
	}
	if strings.Contains(term, "color") {
		return value.ColorDepth8
	}
	if term == "dumb" || term == "" {
		return value.ColorDepthNone
	}

	return value.ColorDepth8 // Conservative default
}

func (cd *CapabilitiesDetector) detectUnix() *value.Capabilities {
	term := cd.env.Get("TERM")
	if term == "dumb" || term == "" {
		return value.NewCapabilities(false, value.ColorDepthNone, false, false, false)
	}

	colorDepth := cd.detectColorDepth()
	return value.NewCapabilities(true, colorDepth, true, true, true)
}

func (cd *CapabilitiesDetector) detectMacOS() *value.Capabilities {
	colorDepth := cd.detectColorDepth()

	// Apple Terminal.app conservative
	if cd.env.Get("TERM_PROGRAM") == "Apple_Terminal" && colorDepth == value.ColorDepthTrueColor {
		colorDepth = value.ColorDepth256
	}

	return value.NewCapabilities(true, colorDepth, true, true, true)
}

func (cd *CapabilitiesDetector) detectWindows() *value.Capabilities {
	// Windows Terminal
	if cd.env.Get("WT_SESSION") != "" {
		return value.NewCapabilities(true, value.ColorDepthTrueColor, true, true, true)
	}

	// VS Code
	if cd.env.Get("TERM_PROGRAM") == "vscode" {
		return value.NewCapabilities(true, value.ColorDepthTrueColor, true, true, true)
	}

	// Conservative default
	return value.NewCapabilities(true, value.ColorDepth8, false, true, true)
}

func (cd *CapabilitiesDetector) parseForceColor(fc string) *value.Capabilities {
	switch fc {
	case "0":
		return value.NewCapabilities(false, value.ColorDepthNone, false, false, false)
	case "1":
		return value.NewCapabilities(true, value.ColorDepth8, true, true, true)
	case "2":
		return value.NewCapabilities(true, value.ColorDepth256, true, true, true)
	case "3", "true":
		return value.NewCapabilities(true, value.ColorDepthTrueColor, true, true, true)
	default:
		return value.NewCapabilities(true, value.ColorDepthTrueColor, true, true, true)
	}
}
