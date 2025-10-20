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
package phoenix

// Version is the current Phoenix TUI Framework version.
const Version = "v0.1.0-beta.2"
