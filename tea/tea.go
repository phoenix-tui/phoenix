// Package tea provides Elm Architecture (MVU pattern) event loop for Phoenix TUI framework.
//
// # Overview
//
// Package tea implements the Model-View-Update pattern for building interactive terminal UIs:
//   - Type-safe event loop with Go generics
//   - Command system for side effects (async operations)
//   - Keyboard and mouse event handling
//   - Program lifecycle management (init, update, view, shutdown)
//   - Composable commands (Batch for parallel, Sequence for sequential)
//   - Zero external dependencies (pure Phoenix architecture)
//
// # Features
//
//   - Pure functional Model-Update-View pattern (predictable state management)
//   - Type-safe models with Go 1.25+ generics (compile-time safety)
//   - Command system for async operations (network, timers, file I/O)
//   - Rich event types (KeyMsg, MouseMsg, WindowSizeMsg, custom messages)
//   - Composable commands (Batch, Sequence, Every for timers)
//   - Alternate screen support (full-screen TUI mode)
//   - Raw mode management (character-by-character input)
//   - High test coverage (domain/service 100%, application 80%, infrastructure 82-98%)
//
// # Quick Start
//
// Basic counter application:
//
//	import "github.com/phoenix-tui/phoenix/tea"
//
//	type Model struct {
//		count int
//	}
//
//	func (m Model) Init() tea.Cmd {
//		return nil
//	}
//
//	func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
//		switch msg := msg.(type) {
//		case tea.KeyMsg:
//			if msg.String() == "+" {
//				m.count++
//			} else if msg.String() == "q" {
//				return m, tea.Quit
//			}
//		}
//		return m, nil
//	}
//
//	func (m Model) View() string {
//		return fmt.Sprintf("Count: %d\nPress + to increment, q to quit", m.count)
//	}
//
//	func main() {
//		p := tea.New(Model{})
//		if err := p.Run(); err != nil {
//			log.Fatal(err)
//		}
//	}
//
// Async operations with commands:
//
//	func LoadData() tea.Cmd {
//		return func() tea.Msg {
//			data := fetchFromAPI() // Runs in goroutine
//			return DataLoadedMsg{Data: data}
//		}
//	}
//
//	func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
//		switch msg := msg.(type) {
//		case tea.KeyMsg:
//			if msg.String() == "r" {
//				return m, LoadData() // Trigger async load
//			}
//		case DataLoadedMsg:
//			m.data = msg.Data // Handle result
//		}
//		return m, nil
//	}
//
// Timer-based updates:
//
//	func (m Model) Init() tea.Cmd {
//		return tea.Every(time.Second, func(t time.Time) tea.Msg {
//			return TickMsg{Time: t}
//		})
//	}
//
// # Architecture
//
// Phoenix TEA implements the Elm Architecture pattern:
//
//	┌──────────────────────────────────┐
//	│  Init() → Initial Model + Cmd    │
//	└──────────────┬───────────────────┘
//	               ↓
//	┌──────────────────────────────────┐
//	│  Event Loop (goroutine-based)    │
//	├──────────────────────────────────┤
//	│  1. Wait for Msg (key/mouse/cmd) │
//	│  2. Update(msg) → new Model + Cmd│
//	│  3. View() → string              │
//	│  4. Render to terminal           │
//	│  5. Execute Cmd (async)          │
//	└──────────────┬───────────────────┘
//	               ↓
//	      (repeat until Quit)
//
// DDD structure:
//   - internal/domain/model    - Model, Message, Command domain logic
//   - internal/domain/service  - Event processing, input handling
//   - internal/application     - Program orchestration, lifecycle
//   - tea.go (this file)       - Public API (wrapper types)
//
// # Performance
//
// Event loop is optimized for responsiveness:
//   - Non-blocking event processing (goroutine per command)
//   - Zero allocations on hot paths (message handling)
//   - Efficient terminal I/O (buffered rendering)
//   - Production-ready (100% coverage in critical domain/service layers)
package tea

import (
	"fmt"
	"io"
	"os/exec"
	"time"

	program2 "github.com/phoenix-tui/phoenix/tea/internal/application/program"
	model2 "github.com/phoenix-tui/phoenix/tea/internal/domain/model"
	"github.com/phoenix-tui/phoenix/tea/internal/domain/service"
	"github.com/phoenix-tui/phoenix/terminal"
)

// Msg represents any message that can be sent through the event loop.
// This is a marker interface - any type can be a message.
type Msg interface{}

// Cmd is a function that produces a message asynchronously.
//
// Commands are the way to perform side effects in the Elm Architecture.
// When Update needs to do something async (network call, timer, file I/O),
// it returns a Cmd that will run in a separate goroutine and eventually
// send a message back to Update.
//
// Example:
//
//	func LoadData() Cmd {
//		return func() Msg {
//			data := fetchFromAPI() // This runs in a goroutine
//			return DataLoadedMsg{Data: data}
//		}
//	}
//
// Commands can be combined using Batch (parallel) or Sequence (sequential).
type Cmd func() Msg

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

// KeyType represents the type of key pressed.
type KeyType int

// Key type constants define all supported keyboard events.
const (
	KeyRune KeyType = iota // Regular character key
	KeyEnter
	KeyBackspace
	KeyTab
	KeyEsc
	KeySpace
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	KeyHome
	KeyEnd
	KeyPgUp
	KeyPgDown
	KeyDelete
	KeyInsert
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyCtrlC // Ctrl+C (common, so dedicated type)
)

// KeyMsg represents a keyboard event.
//
// Zero value: KeyMsg with Type=0, Rune=0, all modifiers false is valid (represents no key).
// Construct using field literals for specific keys.
//
//	var k tea.KeyMsg                   // Zero value - valid, no key
//	k2 := tea.KeyMsg{Type: tea.KeyEnter}  // Explicit - Enter key
type KeyMsg struct {
	Type  KeyType // Type of key pressed
	Rune  rune    // The actual rune (for KeyRune type)
	Alt   bool    // Alt modifier
	Ctrl  bool    // Ctrl modifier
	Shift bool    // Shift modifier
}

// String returns a human-readable representation of the key.
//
// Examples:
//   - KeyMsg{Type: KeyRune, Rune: 'a'}                    → "a"
//   - KeyMsg{Type: KeyRune, Rune: 'A', Shift: true}       → "A"
//   - KeyMsg{Type: KeyEnter}                              → "enter"
//   - KeyMsg{Type: KeyRune, Rune: 'c', Ctrl: true}        → "ctrl+c"
//   - KeyMsg{Type: KeyCtrlC}                              → "ctrl+c"
//   - KeyMsg{Type: KeyUp}                                 → "↑"
//   - KeyMsg{Type: KeyF1}                                 → "F1"
func (k KeyMsg) String() string {
	// Convert to internal format and delegate
	internal := model2.KeyMsg{
		Type:  model2.KeyType(k.Type),
		Rune:  k.Rune,
		Alt:   k.Alt,
		Ctrl:  k.Ctrl,
		Shift: k.Shift,
	}
	return internal.String()
}

// MouseButton represents which mouse button was used.
type MouseButton int

// Mouse button constants define all supported mouse buttons.
const (
	MouseButtonNone MouseButton = iota
	MouseButtonLeft
	MouseButtonMiddle
	MouseButtonRight
	MouseButtonWheelUp
	MouseButtonWheelDown
)

// MouseAction represents what the mouse did.
type MouseAction int

// Mouse action constants define all supported mouse actions.
const (
	MouseActionPress   MouseAction = iota // Button pressed
	MouseActionRelease                    // Button released
	MouseActionMotion                     // Mouse moved
)

// MouseMsg represents a mouse event.
//
// Zero value: MouseMsg at (0,0) with no button/action is valid (represents no mouse event).
// Construct using field literals for specific mouse events.
//
//	var m tea.MouseMsg                  // Zero value - valid, no mouse event
//	m2 := tea.MouseMsg{X: 10, Y: 5, Button: tea.MouseButtonLeft, Action: tea.MouseActionPress}
type MouseMsg struct {
	X      int         // Column (0-based)
	Y      int         // Row (0-based)
	Button MouseButton // Which button
	Action MouseAction // What happened
	Alt    bool        // Alt modifier
	Ctrl   bool        // Ctrl modifier
	Shift  bool        // Shift modifier
}

// String returns a human-readable representation of the mouse event.
//
// Examples:
//   - MouseMsg{X: 10, Y: 5, Button: MouseButtonLeft, Action: MouseActionPress}    → "left press at (10, 5)"
//   - MouseMsg{X: 20, Y: 10, Button: MouseButtonRight, Action: MouseActionRelease} → "right release at (20, 10)"
//   - MouseMsg{X: 15, Y: 8, Button: MouseButtonNone, Action: MouseActionMotion}   → "mouse motion at (15, 8)"
//   - MouseMsg{X: 0, Y: 0, Button: MouseButtonWheelUp, Action: MouseActionPress}  → "wheel up at (0, 0)"
func (m MouseMsg) String() string {
	// Convert to internal format and delegate
	internal := model2.MouseMsg{
		X:      m.X,
		Y:      m.Y,
		Button: model2.MouseButton(m.Button),
		Action: model2.MouseAction(m.Action),
		Alt:    m.Alt,
		Ctrl:   m.Ctrl,
		Shift:  m.Shift,
	}
	return internal.String()
}

// WindowSizeMsg represents a terminal resize event.
//
// Zero value: WindowSizeMsg with Width=0, Height=0 is valid but invalid as a window size.
// Use IsValid() to check if dimensions are positive.
//
//	var w tea.WindowSizeMsg             // Zero value - valid struct, but invalid size
//	w2 := tea.WindowSizeMsg{Width: 80, Height: 24}  // Explicit - valid size
type WindowSizeMsg struct {
	Width  int // Terminal width in columns
	Height int // Terminal height in rows
}

// String returns a human-readable representation.
//
// Example:
//   - WindowSizeMsg{Width: 80, Height: 24} → "window resize: 80x24"
func (w WindowSizeMsg) String() string {
	internal := model2.WindowSizeMsg{Width: w.Width, Height: w.Height}
	return internal.String()
}

// IsValid checks if the window size is valid (positive dimensions).
func (w WindowSizeMsg) IsValid() bool {
	internal := model2.WindowSizeMsg{Width: w.Width, Height: w.Height}
	return internal.IsValid()
}

// QuitMsg signals the program to quit.
// This is a message, not a command. The application can choose to ignore it
// or perform cleanup before actually quitting.
type QuitMsg struct{}

// String returns a human-readable representation.
func (q QuitMsg) String() string {
	return "quit"
}

// BatchMsg contains messages from commands executed in parallel via Batch().
//
// The order of messages is undefined since commands run concurrently.
// Your Update function should handle BatchMsg and process each message.
//
// Example:
//
//	func (m AppModel) Update(msg Msg) (Model[AppModel], Cmd) {
//		switch msg := msg.(type) {
//		case BatchMsg:
//			for _, innerMsg := range msg.Messages {
//				// Process each message from parallel execution
//				m, _ = m.Update(innerMsg)
//			}
//			return m, nil
//		}
//		return m, nil
//	}
type BatchMsg struct {
	Messages []Msg
}

// String returns a human-readable representation.
func (b BatchMsg) String() string {
	return fmt.Sprintf("batch (%d messages)", len(b.Messages))
}

// SequenceMsg contains messages from commands executed sequentially via Sequence().
//
// Messages are in the same order as the input commands to Sequence().
// Your Update function should handle SequenceMsg and process messages in order.
//
// Example:
//
//	func (m AppModel) Update(msg Msg) (Model[AppModel], Cmd) {
//		switch msg := msg.(type) {
//		case SequenceMsg:
//			for _, innerMsg := range msg.Messages {
//				// Process each message in sequence
//				m, _ = m.Update(innerMsg)
//			}
//			return m, nil
//		}
//		return m, nil
//	}
type SequenceMsg struct {
	Messages []Msg
}

// String returns a human-readable representation.
func (s SequenceMsg) String() string {
	return fmt.Sprintf("sequence (%d messages)", len(s.Messages))
}

// PrintlnMsg is sent by Println command for debugging output.
type PrintlnMsg struct {
	Message string
}

// TickMsg is sent by Tick command after a duration has elapsed.
type TickMsg struct {
	Time time.Time
}

// Quit returns a command that quits the program.
func Quit() Cmd {
	return func() Msg {
		return QuitMsg{}
	}
}

// Println returns a command that prints a message (for debugging).
func Println(msg string) Cmd {
	return func() Msg {
		return PrintlnMsg{Message: msg}
	}
}

// Tick returns a command that waits for a duration then sends a TickMsg.
func Tick(d time.Duration) Cmd {
	return func() Msg {
		time.Sleep(d)
		return TickMsg{Time: time.Now()}
	}
}

// Batch executes multiple commands concurrently.
//
// Commands run in parallel via goroutines, and messages are collected into
// a BatchMsg. This is useful when you have independent operations that can
// run simultaneously (e.g., multiple API calls).
//
// Optimizations:
//   - Nil commands are filtered out
//   - If no commands remain, returns nil
//   - If only one command remains, returns it directly (no BatchMsg overhead)
//
// Example:
//
//	cmd := Batch(
//		LoadUserData(),
//		LoadSettings(),
//		LoadPreferences(),
//	) // All three run concurrently
//
// The order of messages in BatchMsg is undefined since commands run in parallel.
func Batch(cmds ...Cmd) Cmd {
	// Filter out nil commands
	filtered := make([]Cmd, 0, len(cmds))
	for _, cmd := range cmds {
		if cmd != nil {
			filtered = append(filtered, cmd)
		}
	}

	// Optimization: no commands
	if len(filtered) == 0 {
		return nil
	}

	// Optimization: single command
	if len(filtered) == 1 {
		return filtered[0]
	}

	// Multiple commands: run in parallel
	return func() Msg {
		results := make(chan Msg, len(filtered))

		// Launch all commands in parallel
		for _, cmd := range filtered {
			go func(c Cmd) {
				results <- c()
			}(cmd)
		}

		// Collect all results
		msgs := make([]Msg, 0, len(filtered))
		for i := 0; i < len(filtered); i++ {
			msgs = append(msgs, <-results)
		}

		return BatchMsg{Messages: msgs}
	}
}

// Sequence executes commands sequentially.
//
// Commands run synchronously in order, and messages are collected into
// a SequenceMsg. This is useful when operations must happen in a specific
// order (e.g., login then load data).
//
// Optimizations:
//   - Nil commands are filtered out
//   - If no commands remain, returns nil
//   - If only one command remains, returns it directly (no SequenceMsg overhead)
//
// Example:
//
//	cmd := Sequence(
//		Login(),
//		LoadUserData(),
//		LoadDashboard(),
//	) // Runs in order: login → data → dashboard
//
// The order of messages in SequenceMsg matches the order of input commands.
func Sequence(cmds ...Cmd) Cmd {
	// Filter out nil commands
	filtered := make([]Cmd, 0, len(cmds))
	for _, cmd := range cmds {
		if cmd != nil {
			filtered = append(filtered, cmd)
		}
	}

	// Optimization: no commands
	if len(filtered) == 0 {
		return nil
	}

	// Optimization: single command
	if len(filtered) == 1 {
		return filtered[0]
	}

	// Multiple commands: run sequentially
	return func() Msg {
		msgs := make([]Msg, 0, len(filtered))

		// Execute commands one by one
		for _, cmd := range filtered {
			msg := cmd() // Synchronous execution
			msgs = append(msgs, msg)
		}

		return SequenceMsg{Messages: msgs}
	}
}

// Program orchestrates the Elm Architecture event loop.
//
// Zero value: Program with zero value has nil internal state and will panic if used.
// Always use New() to create a valid Program instance.
//
//	var p tea.Program[Model]      // Zero value - INVALID, will panic
//	p2 := tea.New(Model{})        // Correct - use constructor
//
// Thread safety: Program is NOT safe for concurrent use.
// Program manages event loop state and must be used from a single goroutine.
// Use Send() method to send messages from other goroutines safely.
//
//	// UNSAFE - concurrent Program operations
//	go p.Run()
//	go p.Stop()  // Race condition!
//
//	// SAFE - Send from goroutines, Program runs in main
//	p := tea.New(model)
//	go func() { p.Send(MyMsg{}) }()  // Safe - Send is thread-safe
//	p.Run()  // Event loop in main goroutine
type Program[T any] struct {
	p *program2.Program[T]
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
func New[T modelConstraint[T]](m T, opts ...Option[T]) *Program[T] {
	// Wrap the model to satisfy model.Model[T] interface
	wrapped := modelWrapper[T]{model: m}

	// Convert options
	internalOpts := make([]program2.Option[T], 0, len(opts))
	for _, opt := range opts {
		internalOpts = append(internalOpts, program2.Option[T](opt))
	}

	return &Program[T]{
		p: program2.New(wrapped, internalOpts...),
	}
}

// modelWrapper wraps a user model to satisfy the internal model.Model[T] interface.
type modelWrapper[T modelConstraint[T]] struct {
	model T
}

func (w modelWrapper[T]) Init() model2.Cmd {
	publicCmd := w.model.Init()
	return convertCmdToInternal(publicCmd)
}

func (w modelWrapper[T]) Update(msg model2.Msg) (model2.Model[T], model2.Cmd) {
	publicMsg := convertMsgToPublic(msg)
	updated, cmd := w.model.Update(publicMsg)
	internalCmd := convertCmdToInternal(cmd)
	return modelWrapper[T]{model: updated}, internalCmd
}

func (w modelWrapper[T]) View() string {
	return w.model.View()
}

// convertMsgToPublic converts internal messages to public API messages.
func convertMsgToPublic(msg model2.Msg) Msg {
	switch m := msg.(type) {
	case model2.KeyMsg:
		return KeyMsg{
			Type:  KeyType(m.Type),
			Rune:  m.Rune,
			Alt:   m.Alt,
			Ctrl:  m.Ctrl,
			Shift: m.Shift,
		}
	case model2.MouseMsg:
		return MouseMsg{
			X:      m.X,
			Y:      m.Y,
			Button: MouseButton(m.Button),
			Action: MouseAction(m.Action),
			Alt:    m.Alt,
			Ctrl:   m.Ctrl,
			Shift:  m.Shift,
		}
	case model2.WindowSizeMsg:
		return WindowSizeMsg{
			Width:  m.Width,
			Height: m.Height,
		}
	case model2.QuitMsg:
		return QuitMsg{}
	case model2.BatchMsg:
		publicMsgs := make([]Msg, len(m.Messages))
		for i, msg := range m.Messages {
			publicMsgs[i] = convertMsgToPublic(msg)
		}
		return BatchMsg{Messages: publicMsgs}
	case model2.SequenceMsg:
		publicMsgs := make([]Msg, len(m.Messages))
		for i, msg := range m.Messages {
			publicMsgs[i] = convertMsgToPublic(msg)
		}
		return SequenceMsg{Messages: publicMsgs}
	case service.PrintlnMsg:
		return PrintlnMsg{Message: m.Message}
	case service.TickMsg:
		return TickMsg{Time: m.Time}
	default:
		// Pass through unknown messages as-is
		return m
	}
}

// convertMsgToInternal converts public API messages to internal messages.
func convertMsgToInternal(msg Msg) model2.Msg {
	switch m := msg.(type) {
	case KeyMsg:
		return model2.KeyMsg{
			Type:  model2.KeyType(m.Type),
			Rune:  m.Rune,
			Alt:   m.Alt,
			Ctrl:  m.Ctrl,
			Shift: m.Shift,
		}
	case MouseMsg:
		return model2.MouseMsg{
			X:      m.X,
			Y:      m.Y,
			Button: model2.MouseButton(m.Button),
			Action: model2.MouseAction(m.Action),
			Alt:    m.Alt,
			Ctrl:   m.Ctrl,
			Shift:  m.Shift,
		}
	case WindowSizeMsg:
		return model2.WindowSizeMsg{
			Width:  m.Width,
			Height: m.Height,
		}
	case QuitMsg:
		return model2.QuitMsg{}
	case BatchMsg:
		internalMsgs := make([]model2.Msg, len(m.Messages))
		for i, msg := range m.Messages {
			internalMsgs[i] = convertMsgToInternal(msg)
		}
		return model2.BatchMsg{Messages: internalMsgs}
	case SequenceMsg:
		internalMsgs := make([]model2.Msg, len(m.Messages))
		for i, msg := range m.Messages {
			internalMsgs[i] = convertMsgToInternal(msg)
		}
		return model2.SequenceMsg{Messages: internalMsgs}
	case PrintlnMsg:
		return service.PrintlnMsg{Message: m.Message}
	case TickMsg:
		return service.TickMsg{Time: m.Time}
	default:
		// Pass through unknown messages as-is
		return m
	}
}

// convertCmdToInternal converts public Cmd to internal Cmd.
func convertCmdToInternal(cmd Cmd) model2.Cmd {
	if cmd == nil {
		return nil
	}
	return func() model2.Msg {
		publicMsg := cmd()
		return convertMsgToInternal(publicMsg)
	}
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
	internalMsg := convertMsgToInternal(msg)
	return p.p.Send(internalMsg)
}

// IsRunning returns true if the program is running.
func (p *Program[T]) IsRunning() bool {
	return p.p.IsRunning()
}

// Option configures a Program.
type Option[T any] program2.Option[T]

// WithInput sets a custom input reader.
func WithInput[T any](r io.Reader) Option[T] {
	return Option[T](program2.WithInput[T](r))
}

// WithOutput sets a custom output writer.
func WithOutput[T any](w io.Writer) Option[T] {
	return Option[T](program2.WithOutput[T](w))
}

// WithAltScreen enables alternate screen buffer.
func WithAltScreen[T any]() Option[T] {
	return Option[T](program2.WithAltScreen[T]())
}

// WithMouseAllMotion enables all mouse motion events.
func WithMouseAllMotion[T any]() Option[T] {
	return Option[T](program2.WithMouseAllMotion[T]())
}

// WithTerminal sets a custom terminal instance (for testing).
func WithTerminal[T any](term terminal.Terminal) Option[T] {
	return Option[T](program2.WithTerminal[T](term))
}

// ExecProcess executes an external interactive command with full terminal control.
//
// This method temporarily suspends the TUI, giving the external command full.
// control of stdin/stdout/stderr. When the command exits, the TUI is restored.
//
// Essential for running:.
//   - Text editors: vim, nano, emacs.
//   - Interactive shells: bash, python REPL.
//   - Pagers: less, more.
//   - SSH sessions.
//
// Example:
//
//	func (m Model) Update(msg Msg) (Model, Cmd) {
//	    switch msg := msg.(type) {
//	    case RunVimMsg:
//	        return m, func() Msg {
//	            cmd := exec.Command("vim", "file.txt")
//	            err := m.program.ExecProcess(cmd)
//	            return VimExitedMsg{Err: err}
//	        }
//	    }
//	    return m, nil
//	}
//
// IMPORTANT:.
//   - Must be called from Cmd goroutine (NOT directly from Update).
//   - Blocks until command completes.
//   - Auto-restores TUI state even if command fails.
//
// Returns error if command execution fails.
func (p *Program[T]) ExecProcess(cmd *exec.Cmd) error {
	return p.p.ExecProcess(cmd)
}
