// Package core provides terminal primitives and Unicode/ANSI handling for Phoenix TUI framework.
//
// # Overview
//
// Package core is the foundation layer of Phoenix, providing:
//   - Terminal capability detection (ANSI, color depth, mouse support)
//   - Correct Unicode/Emoji width calculation (fixes Lipgloss #562)
//   - Terminal primitives (Size, Position, Cell)
//   - Raw mode management
//   - Zero external dependencies (stdlib only)
//
// # Features
//
//   - Auto-detection of terminal capabilities from environment
//   - Correct grapheme cluster width (emoji, CJK characters)
//   - 46x faster Unicode processing than alternatives
//   - Fluent, immutable API with method chaining
//   - Type-safe operations with compile-time guarantees
//   - Cross-platform support (Unix, Windows)
//
// # Quick Start
//
// Basic terminal detection:
//
//	import "github.com/phoenix-tui/phoenix/core"
//
//	term := core.AutoDetect()
//	caps := term.Capabilities()
//
//	if caps.SupportsTrueColor() {
//	    fmt.Println("24-bit color supported!")
//	}
//
// Unicode width calculation:
//
//	width := core.StringWidth("Hello 👋 World 🌍")  // Returns 17 (correct!)
//	// Compare: Lipgloss returns 19 (wrong)
//
// # Architecture
//
// This package follows Domain-Driven Design (DDD):
//   - internal/domain/model    - Core business logic (Terminal, Cell)
//   - internal/domain/value    - Value objects (Capabilities, Size)
//   - internal/domain/service  - Domain services (Unicode processing)
//   - internal/infrastructure  - Platform-specific code (platform detection)
//   - core.go (this file)      - Public API (wrapper types)
//
// The public API provides clean wrapper types that delegate to internal domain models,
// following the Relica pattern for better pkg.go.dev visibility.
//
// # Performance
//
// Unicode processing is optimized for speed:
//   - ASCII strings: <50 ns/op
//   - Emoji strings: <200 ns/op (46x faster than alternatives)
//   - CJK strings: <200 ns/op
//   - Zero allocations in hot paths
package core

import (
	model2 "github.com/phoenix-tui/phoenix/core/internal/domain/model"
	service2 "github.com/phoenix-tui/phoenix/core/internal/domain/service"
	value2 "github.com/phoenix-tui/phoenix/core/internal/domain/value"
	"github.com/phoenix-tui/phoenix/core/internal/infrastructure/platform"
)

// Terminal represents the public API for terminal operations.
// This is the main entry point for phoenix/core library.
//
// Design Philosophy:
//   - Fluent API with method chaining
//   - Immutable operations (returns new instances)
//   - Type-safe with compile-time guarantees
//   - Zero external dependencies (stdlib only)
//
// Zero value: A Terminal with zero value has nil internal state and will panic if used.
// Always use NewTerminal() or AutoDetect() to create a valid Terminal instance.
//
//	var t core.Terminal      // Zero value - INVALID, will panic
//	t2 := core.NewTerminal() // Correct - use constructor
//	t3 := core.AutoDetect()  // Recommended - auto-detect capabilities
//
// Thread safety: Terminal is safe for concurrent reads (immutable).
// Concurrent writes are not applicable since methods return new instances.
//
// Example Usage:
//
//	term := core.NewTerminal()
//	caps := term.Capabilities()
//
//	if caps.SupportsColor() {
//	    fmt.Println("Terminal supports color!")
//	}
//
//	// Detect capabilities automatically
//	term = core.AutoDetect()
//	fmt.Printf("Color depth: %d\n", term.Capabilities().ColorDepth())
type Terminal struct {
	domain *model2.Terminal
}

// NewTerminal creates a new Terminal with default capabilities.
// Default assumes no ANSI support, no colors, VT100 size (80x24).
//
// For automatic detection, use AutoDetect() instead.
func NewTerminal() *Terminal {
	caps := value2.NewCapabilities(false, value2.ColorDepthNone, false, false, false)
	return &Terminal{domain: model2.NewTerminal(caps)}
}

// AutoDetect creates a Terminal with automatically detected capabilities
// based on environment variables (TERM, COLORTERM, NO_COLOR, etc.).
//
// This is the recommended way to create a Terminal for most use cases.
//
// Example:
//
//	term := core.AutoDetect()
//	if term.Capabilities().SupportsTrueColor() {
//	    // Use 24-bit colors
//	}
func AutoDetect() *Terminal {
	env := platform.OsEnvironmentProvider{}
	detector := service2.NewCapabilitiesDetector(env)
	caps := detector.Detect()

	return &Terminal{domain: model2.NewTerminal(caps)}
}

// NewTerminalWithCapabilities creates a Terminal with specific capabilities.
// Useful for testing or when you know exact terminal capabilities.
func NewTerminalWithCapabilities(caps *Capabilities) *Terminal {
	if caps == nil {
		return NewTerminal()
	}
	return &Terminal{domain: model2.NewTerminal(caps.domain)}
}

// Capabilities returns the terminal's capabilities.
// The returned Capabilities object is immutable.
func (t *Terminal) Capabilities() *Capabilities {
	return &Capabilities{domain: t.domain.Capabilities()}
}

// Size returns the current terminal size.
func (t *Terminal) Size() Size {
	domainSize := t.domain.Size()
	return Size{Width: domainSize.Width, Height: domainSize.Height}
}

// WithSize returns a new Terminal with the specified size.
// Original Terminal is not modified (immutable operation).
//
// Example:
//
//	term := core.NewTerminal()
//	resized := term.WithSize(core.NewSize(120, 40))
func (t *Terminal) WithSize(size Size) *Terminal {
	domainSize := value2.NewSize(size.Width, size.Height)
	newDomain := t.domain.WithSize(domainSize)
	return &Terminal{domain: newDomain}
}

// IsRawModeEnabled returns whether raw mode is currently enabled.
func (t *Terminal) IsRawModeEnabled() bool {
	return t.domain.IsRawMode()
}

// WithRawMode returns a new Terminal with the specified raw mode.
// Original Terminal is not modified (immutable operation).
//
// Note: This is a low-level API. Most users should use EnableRawMode/DisableRawMode
// from the platform-specific packages instead.
func (t *Terminal) WithRawMode(rm *RawMode) *Terminal {
	if rm == nil {
		return t
	}
	newDomain := t.domain.WithRawMode(rm.domain)
	return &Terminal{domain: newDomain}
}

// Size represents terminal dimensions.
//
// Zero value: Size{Width: 0, Height: 0} is valid but represents an empty terminal.
// Use NewSize() for validated dimensions (minimum 1x1).
//
//	var s core.Size          // Zero value - (0, 0), may cause issues
//	s2 := core.NewSize(80, 24)  // Validated - ensures minimum 1x1
type Size struct {
	Width  int
	Height int
}

// NewSize creates a new Size with validation (minimum 1x1).
func NewSize(width, height int) Size {
	s := value2.NewSize(width, height)
	return Size{Width: s.Width, Height: s.Height}
}

// Position represents a position in the terminal (0-based).
//
// Zero value: Position{Row: 0, Col: 0} is valid and represents top-left corner.
// This is a valid and commonly used position.
//
//	var p core.Position      // Zero value - (0, 0) top-left, valid
//	p2 := core.NewPosition(5, 10)  // Explicit position
type Position struct {
	Row int
	Col int
}

// NewPosition creates a new Position with validation (non-negative).
func NewPosition(row, col int) Position {
	p := value2.NewPosition(row, col)
	return Position{Row: p.Row, Col: p.Col}
}

// Add returns a new Position offset by delta row/column.
func (p Position) Add(deltaRow, deltaCol int) Position {
	domainPos := value2.NewPosition(p.Row, p.Col)
	result := domainPos.Add(deltaRow, deltaCol)
	return Position{Row: result.Row, Col: result.Col}
}

// Cell represents a terminal cell with grapheme cluster support.
//
// Zero value: Cell{Content: "", Width: 0} is valid and represents an empty cell.
// Use NewCellAuto() for automatic Unicode width calculation (recommended).
//
//	var c core.Cell              // Zero value - empty cell, valid
//	c2 := core.NewCellAuto("A")  // Recommended - auto width
//	c3 := core.NewCell("A", 1)   // Manual width control
type Cell struct {
	Content string
	Width   int
}

// NewCell creates a new Cell with the given content and manual width.
// This is useful when you need explicit control over width (advanced use cases).
//
// For automatic Unicode-aware width calculation (recommended), use NewCellAuto().
//
// Example:
//
//	cell := core.NewCell("A", 1)     // Manual width control
//	cell := core.NewCell("👋", 3)    // Wrong width, but user controls
func NewCell(content string, width int) Cell {
	c := value2.NewCell(content, width)
	return Cell{Content: c.Content(), Width: c.Width()}
}

// NewCellAuto creates a Cell with automatic Unicode width calculation.
// This is the recommended way to create cells as it handles all Unicode correctly:
//   - Emoji: "👋" -> width 2
//   - CJK: "中" -> width 2
//   - ASCII: "A" -> width 1
//   - Combining: "é" -> width 1
//   - Zero-width: correctly handled
//
// This fixes Charm's lipgloss#562 bug with incorrect emoji/Unicode rendering.
//
// Example:
//
//	cell := core.NewCellAuto("👋")   // content "👋", width 2 (auto)
//	cell := core.NewCellAuto("中文")  // content "中文", width 4 (auto)
//	cell := core.NewCellAuto("Hello") // content "Hello", width 5 (auto)
//	cell := core.NewCellAuto("Café")  // content "Café", width 4 (auto)
//
// For manual width control (rare cases), use NewCell(content, width).
func NewCellAuto(content string) Cell {
	// Use UnicodeService to calculate correct width
	unicodeService := service2.NewUnicodeService()
	width := unicodeService.StringWidth(content)

	// Create Cell with calculated width
	domainCell := value2.NewCell(content, width)
	return Cell{
		Content: domainCell.Content(),
		Width:   domainCell.Width(),
	}
}

// RawMode represents raw mode state with original terminal state preservation.
// This is a low-level API for platform-specific implementations.
//
// Zero value: RawMode with zero value has nil internal state and will panic if used.
// Always use NewRawMode() to create a valid RawMode instance.
//
//	var rm core.RawMode           // Zero value - INVALID, will panic
//	rm2, err := core.NewRawMode(originalState)  // Correct - use constructor
type RawMode struct {
	domain *model2.RawMode
}

// NewRawMode creates a new RawMode with the original terminal state.
// The originalState should be platform-specific state (e.g., syscall.Termios on Unix).
func NewRawMode(originalState interface{}) (*RawMode, error) {
	domain, err := model2.NewRawMode(originalState)
	if err != nil {
		return nil, err
	}
	return &RawMode{domain: domain}, nil
}

// IsEnabled returns whether raw mode is currently enabled.
func (r *RawMode) IsEnabled() bool {
	return r.domain.IsEnabled()
}

// OriginalState returns the original terminal state for restoration.
func (r *RawMode) OriginalState() interface{} {
	return r.domain.OriginalState()
}

// Capabilities wraps domain capabilities with a public API.
//
// Zero value: Capabilities with zero value has nil internal state and will panic if used.
// Always use NewCapabilities() or AutoDetect() to create valid instances.
//
//	var c core.Capabilities           // Zero value - INVALID, will panic
//	c2 := core.NewCapabilities(...)   // Correct - use constructor
//
// Thread safety: Capabilities is immutable and safe for concurrent reads.
type Capabilities struct {
	domain *value2.Capabilities
}

// NewCapabilities creates capabilities with specific features.
func NewCapabilities(ansi bool, colorDepth ColorDepth, mouse, altScreen, cursor bool) *Capabilities {
	domainColorDepth := value2.ColorDepth(colorDepth)
	domain := value2.NewCapabilities(ansi, domainColorDepth, mouse, altScreen, cursor)
	return &Capabilities{domain: domain}
}

// SupportsANSI returns whether terminal supports ANSI escape sequences.
func (c *Capabilities) SupportsANSI() bool {
	return c.domain.SupportsANSI()
}

// ColorDepth returns the terminal's color depth.
func (c *Capabilities) ColorDepth() ColorDepth {
	return ColorDepth(c.domain.ColorDepth())
}

// SupportsColor returns whether terminal supports any color.
func (c *Capabilities) SupportsColor() bool {
	return c.domain.ColorDepth() > value2.ColorDepthNone
}

// SupportsTrueColor returns whether terminal supports 24-bit true color.
func (c *Capabilities) SupportsTrueColor() bool {
	return c.domain.ColorDepth() == value2.ColorDepthTrueColor
}

// SupportsMouse returns whether terminal supports mouse events.
func (c *Capabilities) SupportsMouse() bool {
	return c.domain.SupportsMouse()
}

// SupportsAltScreen returns whether terminal supports alternate screen buffer.
func (c *Capabilities) SupportsAltScreen() bool {
	return c.domain.SupportsAltScreen()
}

// SupportsCursorControl returns whether terminal supports cursor control.
func (c *Capabilities) SupportsCursorControl() bool {
	return c.domain.SupportsCursorControl()
}

// ColorDepth represents terminal color support levels.
type ColorDepth int

const (
	// ColorDepthNone - no color support (monochrome).
	ColorDepthNone ColorDepth = 0

	// ColorDepth8 - 8 colors (3-bit: black, red, green, yellow, blue, magenta, cyan, white).
	ColorDepth8 ColorDepth = 8

	// ColorDepth256 - 256 colors (8-bit: 216 colors + 16 system + 24 grayscale).
	ColorDepth256 ColorDepth = 256

	// ColorDepthTrueColor - 16.7 million colors (24-bit RGB).
	ColorDepthTrueColor ColorDepth = 16777216
)

// String returns a human-readable color depth description.
func (cd ColorDepth) String() string {
	switch cd {
	case ColorDepthNone:
		return "None"
	case ColorDepth8:
		return "8 colors"
	case ColorDepth256:
		return "256 colors"
	case ColorDepthTrueColor:
		return "TrueColor (16.7M colors)"
	default:
		return "Unknown"
	}
}
