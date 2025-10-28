package program

import (
	"io"

	"github.com/phoenix-tui/phoenix/terminal"
)

// Option configures a Program at creation time.
// Uses the functional options pattern for clean, extensible API.
type Option[T any] func(*Program[T])

// WithInput sets a custom input reader (default: os.Stdin).
//
// Example:
//
//	customInput := strings.NewReader("test input")
//	p := program.New(model, program.WithInput(customInput))
func WithInput[T any](r io.Reader) Option[T] {
	return func(p *Program[T]) {
		p.input = r
	}
}

// WithOutput sets a custom output writer (default: os.Stdout).
//
// Example:
//
//	var buf bytes.Buffer
//	p := program.New(model, program.WithOutput(&buf))
func WithOutput[T any](w io.Writer) Option[T] {
	return func(p *Program[T]) {
		p.output = w
	}
}

// WithAltScreen enables alternate screen buffer.
// The TUI will take over the entire terminal screen.
//
// Example:
//
//	p := program.New(model, program.WithAltScreen())
func WithAltScreen[T any]() Option[T] {
	return func(p *Program[T]) {
		p.altScreen = true
	}
}

// WithMouseAllMotion enables all mouse motion events.
// Without this, only click/release events are captured.
//
// Example:
//
//	p := program.New(model, program.WithMouseAllMotion())
func WithMouseAllMotion[T any]() Option[T] {
	return func(p *Program[T]) {
		p.mouseAllMotion = true
	}
}

// WithTerminal sets a custom terminal instance (default: auto-detected).
//
// By default, Program auto-detects the best terminal implementation.
// Use this option for testing with mock terminals.
//
// Example:
//
//	mockTerm := testing.NewMockTerminal()
//	p := program.New(model, program.WithTerminal(mockTerm))
func WithTerminal[T any](term terminal.Terminal) Option[T] {
	return func(p *Program[T]) {
		p.terminal = term
	}
}
