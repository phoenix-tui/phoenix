package program

import "io"

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
