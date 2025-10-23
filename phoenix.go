// Package phoenix is the root umbrella module for Phoenix TUI Framework.
//
// Phoenix is a modern, high-performance Terminal User Interface framework for Go,
// built with Domain-Driven Design principles and modern Go 1.25+ patterns.
//
// # Architecture
//
// Phoenix consists of 10 independent libraries that can be used together or separately:
//
//   - github.com/phoenix-tui/phoenix/core       - Terminal primitives & Unicode support
//   - github.com/phoenix-tui/phoenix/style      - CSS-like styling system
//   - github.com/phoenix-tui/phoenix/tea        - Elm Architecture (Model-Update-View)
//   - github.com/phoenix-tui/phoenix/layout     - Flexbox & Box Model layouts
//   - github.com/phoenix-tui/phoenix/render     - High-performance differential renderer
//   - github.com/phoenix-tui/phoenix/components - Rich UI component library
//   - github.com/phoenix-tui/phoenix/mouse      - Mouse input handling
//   - github.com/phoenix-tui/phoenix/clipboard  - Cross-platform clipboard operations
//   - github.com/phoenix-tui/phoenix/terminal   - Terminal detection & capabilities
//   - github.com/phoenix-tui/phoenix/testing    - Testing utilities (Mock/Null terminals)
//
// # Quick Start
//
// Install individual libraries:
//
//	go get github.com/phoenix-tui/phoenix/tea@latest
//	go get github.com/phoenix-tui/phoenix/components@latest
//
// Or install all libraries via the root module:
//
//	go get github.com/phoenix-tui/phoenix@latest
//
// # Example: Hello World
//
//	package main
//
//	import (
//	    "fmt"
//	    "os"
//	    tea "github.com/phoenix-tui/phoenix/tea/api"
//	)
//
//	type model struct{ message string }
//
//	func (m model) Init() tea.Cmd { return nil }
//
//	func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//	    if _, ok := msg.(tea.KeyMsg); ok {
//	        return m, tea.Quit
//	    }
//	    return m, nil
//	}
//
//	func (m model) View() string {
//	    return fmt.Sprintf("Hello, %s!\n\nPress any key to quit.", m.message)
//	}
//
//	func main() {
//	    p := tea.NewProgram(model{message: "World"})
//	    if _, err := p.Run(); err != nil {
//	        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
//	        os.Exit(1)
//	    }
//	}
//
// # Why Phoenix?
//
// Phoenix was created to address critical issues in the Charm ecosystem:
//
//   - Perfect Unicode/Emoji support (no width calculation bugs)
//   - 10x better performance (29,000 FPS renderer vs ~450 FPS)
//   - Type-safe API with Go 1.25+ generics
//   - DDD architecture (testable, maintainable)
//   - 94.5% average test coverage
//   - Zero external TUI dependencies
//
// # Multi-Module Monorepo
//
// This repository uses a multi-module structure where each library is independently versioned.
// The root module serves as an umbrella module for convenient installation and documentation.
//
// For more information, see:
//   - Documentation: https://github.com/phoenix-tui/phoenix/tree/main/docs
//   - Examples: https://github.com/phoenix-tui/phoenix/tree/main/examples
//   - Releases: https://github.com/phoenix-tui/phoenix/releases
//
// # Version
//
// Phoenix version is managed through git tags and Go modules.
// To check your installed version:
//
//	go list -m github.com/phoenix-tui/phoenix
//
// This will show the exact version you're using, including any pre-release tags.
package phoenix

import (
	clipboardapi "github.com/phoenix-tui/phoenix/clipboard/api"
	coreapi "github.com/phoenix-tui/phoenix/core/api"
	styleapi "github.com/phoenix-tui/phoenix/style/api"
	teaapi "github.com/phoenix-tui/phoenix/tea/api"
	terminalapi "github.com/phoenix-tui/phoenix/terminal/api"
	terminalinfra "github.com/phoenix-tui/phoenix/terminal/infrastructure"
)

// ┌─────────────────────────────────────────────────────────────┐
// │ Core - Terminal Primitives                                  │
// └─────────────────────────────────────────────────────────────┘

// AutoDetectTerminal creates a Terminal by auto-detecting the current environment.
// This is the recommended way to create a Terminal for most applications.
//
// Example:
//
//	term := phoenix.AutoDetectTerminal()
//	fmt.Printf("Terminal size: %dx%d\n", term.Size().Width, term.Size().Height)
func AutoDetectTerminal() *coreapi.Terminal {
	return coreapi.AutoDetect()
}

// NewTerminal creates a new Terminal with default auto-detected capabilities.
// Equivalent to AutoDetectTerminal() - provided for API completeness.
//
// Example:
//
//	term := phoenix.NewTerminal()
//	fmt.Printf("Terminal size: %dx%d\n", term.Size().Width, term.Size().Height)
func NewTerminal() *coreapi.Terminal {
	return coreapi.NewTerminal()
}

// NewTerminalWithCapabilities creates a new Terminal with the specified capabilities.
// Use this when you need full control over terminal configuration.
//
// Example:
//
//	caps := phoenix.NewCapabilities(true, phoenix.ColorDepth256, true, true, true)
//	term := phoenix.NewTerminalWithCapabilities(caps)
func NewTerminalWithCapabilities(caps *coreapi.Capabilities) *coreapi.Terminal {
	return coreapi.NewTerminalWithCapabilities(caps)
}

// NewSize creates a new terminal size (width x height in cells).
//
// Example:
//
//	size := phoenix.NewSize(80, 24)  // 80 columns, 24 rows
func NewSize(width, height int) coreapi.Size {
	return coreapi.NewSize(width, height)
}

// NewCapabilities creates a new terminal capabilities configuration.
//
// Example:
//
//	caps := phoenix.NewCapabilities(
//		true,                      // ANSI support
//		phoenix.ColorDepth256,     // 256-color support
//		true,                      // Mouse support
//		true,                      // Alt screen support
//		true,                      // Cursor support
//	)
func NewCapabilities(ansi bool, colorDepth coreapi.ColorDepth, mouse, altScreen, cursor bool) *coreapi.Capabilities {
	return coreapi.NewCapabilities(ansi, colorDepth, mouse, altScreen, cursor)
}

// ColorDepth constants (re-exported from core).
const (
	ColorDepthNone      = coreapi.ColorDepthNone
	ColorDepth8         = coreapi.ColorDepth8
	ColorDepth256       = coreapi.ColorDepth256
	ColorDepthTrueColor = coreapi.ColorDepthTrueColor
)

// ┌─────────────────────────────────────────────────────────────┐
// │ Style - CSS-like Styling                                    │
// └─────────────────────────────────────────────────────────────┘

// NewStyle creates a new Style builder for applying colors, borders, padding, etc.
//
// Example:
//
//	s := phoenix.NewStyle().
//		Foreground("#00FF00").
//		Background("#000000").
//		Bold().
//		Padding(1)
//	fmt.Println(s.Render("Styled text"))
func NewStyle() styleapi.Style {
	return styleapi.New()
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Tea - Elm Architecture (Model-Update-View)                  │
// └─────────────────────────────────────────────────────────────┘

// modelConstraint defines the interface that models must implement for Tea programs.
// This is re-exported from tea/api to make the umbrella API self-contained.
type modelConstraint[T any] interface {
	Init() teaapi.Cmd
	Update(teaapi.Msg) (T, teaapi.Cmd)
	View() string
}

// NewProgram creates a new Tea Program with the given model.
// This is the main entry point for building Phoenix TUI applications.
//
// Example:
//
//	type MyModel struct { count int }
//	// ... implement tea.Model interface ...
//
//	p := phoenix.NewProgram(MyModel{}, phoenix.WithAltScreen[MyModel]())
//	if err := p.Run(); err != nil {
//		log.Fatal(err)
//	}
func NewProgram[T modelConstraint[T]](model T, opts ...teaapi.Option[T]) *teaapi.Program[T] {
	return teaapi.New(model, opts...)
}

// WithAltScreen enables the alternate screen buffer.
// This allows your TUI to take over the full terminal without affecting the scrollback.
//
// Example:
//
//	p := phoenix.NewProgram(model, phoenix.WithAltScreen[MyModel]())
func WithAltScreen[T any]() teaapi.Option[T] {
	return teaapi.WithAltScreen[T]()
}

// WithMouseAllMotion enables mouse support with all motion events.
//
// Example:
//
//	p := phoenix.NewProgram(model, phoenix.WithMouseAllMotion[MyModel]())
func WithMouseAllMotion[T any]() teaapi.Option[T] {
	return teaapi.WithMouseAllMotion[T]()
}

// Quit returns a command that quits the program.
//
// Example:
//
//	func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
//		if msg.(tea.KeyMsg).String() == "q" {
//			return m, phoenix.Quit()
//		}
//		return m, nil
//	}
func Quit() teaapi.Cmd {
	return teaapi.Quit()
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Clipboard - Cross-platform Clipboard Operations             │
// └─────────────────────────────────────────────────────────────┘

// ReadClipboard reads text from the system clipboard.
//
// Example:
//
//	text, err := phoenix.ReadClipboard()
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Clipboard:", text)
func ReadClipboard() (string, error) {
	return clipboardapi.Read()
}

// WriteClipboard writes text to the system clipboard.
//
// Example:
//
//	err := phoenix.WriteClipboard("Hello, clipboard!")
//	if err != nil {
//		log.Fatal(err)
//	}
func WriteClipboard(text string) error {
	return clipboardapi.Write(text)
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Terminal - Platform-optimized Terminal Operations           │
// └─────────────────────────────────────────────────────────────┘

// NewPlatformTerminal creates a new platform-optimized Terminal.
// Automatically selects the best implementation for the current platform:
//   - Windows Console API (fastest on Windows cmd.exe/PowerShell)
//   - ANSI fallback (for Git Bash, MinTTY, Unix)
//
// Note: This is from the terminal/infrastructure package, providing lower-level
// operations compared to core/api Terminal which focuses on capabilities detection.
//
// Example:
//
//	term := phoenix.NewPlatformTerminal()
//	term.HideCursor()
//	defer term.ShowCursor()
func NewPlatformTerminal() terminalapi.Terminal {
	return terminalinfra.NewTerminal()
}

// NewANSITerminal creates a new ANSI-based Terminal.
// Use this when you want to force ANSI escape codes (e.g., for SSH, tmux).
//
// Example:
//
//	term := phoenix.NewANSITerminal()
//	term.Clear()
func NewANSITerminal() terminalapi.Terminal {
	return terminalinfra.NewANSITerminal()
}
