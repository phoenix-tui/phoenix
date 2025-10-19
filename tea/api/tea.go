// Package api provides a simple and elegant way to build terminal user interfaces.
//
// Phoenix TEA implements the Elm Architecture (Model-View-Update) pattern:
//
//	Init → Update → View
//	  ↑       ↓
//	  └─ Cmd ─┘
//
// Example:
//
//	type MyModel struct {
//		count int
//	}
//
//	func (m MyModel) Init() tea.Cmd {
//		return nil
//	}
//
//	func (m MyModel) Update(msg tea.Msg) (tea.Model[MyModel], tea.Cmd) {
//		switch msg := msg.(type) {
//		case tea.KeyMsg:
//			if msg.String() == "+" {
//				m.count++
//			}
//			if msg.String() == "q" {
//				return m, tea.Quit
//			}
//		}
//		return m, nil
//	}
//
//	func (m MyModel) View() string {
//		return fmt.Sprintf("Count: %d\nPress + to increment, q to quit", m.count)
//	}
//
//	func main() {
//		p := tea.New(MyModel{count: 0})
//		if err := p.Run(); err != nil {
//			log.Fatal(err)
//		}
//	}
package api

import (
	"io"
	"time"

	"github.com/phoenix-tui/phoenix/tea/application/program"
	"github.com/phoenix-tui/phoenix/tea/domain/model"
	"github.com/phoenix-tui/phoenix/tea/domain/service"
)

// Msg represents any message that can be sent through the event loop.
type Msg = model.Msg

// Cmd represents a command (side effect) to execute.
type Cmd = model.Cmd

// Model represents the Elm Architecture contract.
// Your application must implement Init, Update, and View.
//
// The Model interface is satisfied by any type T that implements:
//   - Init() Cmd
//   - Update(Msg) (T, Cmd)
//   - View() string
//
// Note: We don't export this as a named interface because Go generics
// don't allow type aliases for generic interfaces. Instead, your concrete
// types just need to implement these methods.
type modelConstraint[T any] interface {
	Init() Cmd
	Update(Msg) (T, Cmd)
	View() string
}

// KeyMsg represents a keyboard event.
type KeyMsg = model.KeyMsg

// KeyType represents the type of key pressed.
type KeyType = model.KeyType

// Key types.
const (
	KeyRune      = model.KeyRune
	KeyEnter     = model.KeyEnter
	KeyBackspace = model.KeyBackspace
	KeyTab       = model.KeyTab
	KeyEsc       = model.KeyEsc
	KeySpace     = model.KeySpace
	KeyUp        = model.KeyUp
	KeyDown      = model.KeyDown
	KeyLeft      = model.KeyLeft
	KeyRight     = model.KeyRight
	KeyHome      = model.KeyHome
	KeyEnd       = model.KeyEnd
	KeyPgUp      = model.KeyPgUp
	KeyPgDown    = model.KeyPgDown
	KeyDelete    = model.KeyDelete
	KeyInsert    = model.KeyInsert
	KeyF1        = model.KeyF1
	KeyF2        = model.KeyF2
	KeyF3        = model.KeyF3
	KeyF4        = model.KeyF4
	KeyF5        = model.KeyF5
	KeyF6        = model.KeyF6
	KeyF7        = model.KeyF7
	KeyF8        = model.KeyF8
	KeyF9        = model.KeyF9
	KeyF10       = model.KeyF10
	KeyF11       = model.KeyF11
	KeyF12       = model.KeyF12
	KeyCtrlC     = model.KeyCtrlC
)

// MouseMsg represents a mouse event.
type MouseMsg = model.MouseMsg

// MouseButton represents which mouse button was used.
type MouseButton = model.MouseButton

// Mouse buttons.
const (
	MouseButtonNone      = model.MouseButtonNone
	MouseButtonLeft      = model.MouseButtonLeft
	MouseButtonMiddle    = model.MouseButtonMiddle
	MouseButtonRight     = model.MouseButtonRight
	MouseButtonWheelUp   = model.MouseButtonWheelUp
	MouseButtonWheelDown = model.MouseButtonWheelDown
)

// MouseAction represents what the mouse did.
type MouseAction = model.MouseAction

// Mouse actions.
const (
	MouseActionPress   = model.MouseActionPress
	MouseActionRelease = model.MouseActionRelease
	MouseActionMotion  = model.MouseActionMotion
)

// WindowSizeMsg represents a terminal resize event.
type WindowSizeMsg = model.WindowSizeMsg

// QuitMsg signals the program to quit.
type QuitMsg = model.QuitMsg

// BatchMsg contains messages from a Batch command.
type BatchMsg = model.BatchMsg

// SequenceMsg contains messages from a Sequence command.
type SequenceMsg = model.SequenceMsg

// PrintlnMsg is sent by Println command.
type PrintlnMsg = service.PrintlnMsg

// TickMsg is sent by Tick command.
type TickMsg = service.TickMsg

// Quit returns a command that quits the program.
func Quit() Cmd {
	return service.Quit()
}

// Println returns a command that prints a message (for debugging).
func Println(msg string) Cmd {
	return service.Println(msg)
}

// Tick returns a command that waits for a duration then sends a TickMsg.
func Tick(d time.Duration) Cmd {
	return service.Tick(d)
}

// Batch executes multiple commands concurrently.
func Batch(cmds ...Cmd) Cmd {
	return model.Batch(cmds...)
}

// Sequence executes multiple commands sequentially.
func Sequence(cmds ...Cmd) Cmd {
	return model.Sequence(cmds...)
}

// Program orchestrates the Elm Architecture event loop.
type Program[T any] struct {
	p *program.Program[T]
}

// New creates a new program with the given model.
//
// The model must implement:
//   - Init() Cmd
//   - Update(Msg) (T, Cmd)
//   - View() string
//
// Example:
//
//	p := tea.New(MyModel{})
//	p.Run()
func New[T modelConstraint[T]](m T, opts ...ProgramOption[T]) *Program[T] {
	// Wrap the model to satisfy model.Model[T] interface
	wrapped := modelWrapper[T]{model: m}

	// Convert options
	internalOpts := make([]program.ProgramOption[T], 0, len(opts))
	for _, opt := range opts {
		internalOpts = append(internalOpts, program.ProgramOption[T](opt))
	}

	return &Program[T]{
		p: program.New(wrapped, internalOpts...),
	}
}

// modelWrapper wraps a user model to satisfy the internal model.Model[T] interface.
type modelWrapper[T modelConstraint[T]] struct {
	model T
}

func (w modelWrapper[T]) Init() model.Cmd {
	return w.model.Init()
}

func (w modelWrapper[T]) Update(msg model.Msg) (model.Model[T], model.Cmd) {
	updated, cmd := w.model.Update(msg)
	return modelWrapper[T]{model: updated}, cmd
}

func (w modelWrapper[T]) View() string {
	return w.model.View()
}

// Run starts the program and blocks until it quits.
func (p *Program[T]) Run() error {
	return p.p.Run()
}

// Start starts the program in a goroutine.
func (p *Program[T]) Start() error {
	return p.p.Start()
}

// Stop stops a running program gracefully.
func (p *Program[T]) Stop() {
	p.p.Stop()
}

// Send sends a message to the event loop.
func (p *Program[T]) Send(msg Msg) error {
	return p.p.Send(msg)
}

// IsRunning returns true if the program is running.
func (p *Program[T]) IsRunning() bool {
	return p.p.IsRunning()
}

// ProgramOption configures a Program.
type ProgramOption[T any] program.ProgramOption[T]

// WithInput sets a custom input reader.
func WithInput[T any](r io.Reader) ProgramOption[T] {
	return ProgramOption[T](program.WithInput[T](r))
}

// WithOutput sets a custom output writer.
func WithOutput[T any](w io.Writer) ProgramOption[T] {
	return ProgramOption[T](program.WithOutput[T](w))
}

// WithAltScreen enables alternate screen buffer.
func WithAltScreen[T any]() ProgramOption[T] {
	return ProgramOption[T](program.WithAltScreen[T]())
}

// WithMouseAllMotion enables all mouse motion events.
func WithMouseAllMotion[T any]() ProgramOption[T] {
	return ProgramOption[T](program.WithMouseAllMotion[T]())
}
