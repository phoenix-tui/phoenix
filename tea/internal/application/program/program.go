// Package program provides the Program type that orchestrates the Elm Architecture event loop.
// It manages application lifecycle, message passing, and rendering coordination.
package program

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"

	model2 "github.com/phoenix-tui/phoenix/tea/internal/domain/model"
	"github.com/phoenix-tui/phoenix/tea/internal/infrastructure/input"
	"github.com/phoenix-tui/phoenix/terminal"
)

// Program orchestrates the Elm Architecture event loop.
// It manages the application lifecycle, message passing, and rendering.
//
// Type parameter T is the concrete model type.
//
// Example:
//
//	type MyModel struct { count int }
//	// ... implement Model[MyModel] interface ...
//
//	p := program.New(MyModel{count: 0})
//	if err := p.Run(); err != nil {
//		log.Fatal(err)
//	}
type Program[T any] struct {
	// Model instance
	model model2.Model[T]

	// I/O streams
	input  io.Reader
	output io.Writer

	// Terminal for screen management (alternate screen, cursor, etc.)
	// Created automatically in Run() via detect.NewTerminal()
	// Can be set via WithTerminal() option for testing
	terminal terminal.Terminal

	// Input reader for parsing stdin
	inputReader *input.Reader

	// Input reader goroutine lifecycle management
	inputReaderCancel     context.CancelFunc // Cancel function to stop inputReader goroutine
	inputReaderDone       chan struct{}      // Signals when inputReader goroutine has exited
	inputReaderRunning    bool               // True if inputReader goroutine is active
	inputReaderGeneration uint64             // Generation counter to prevent race conditions

	// Configuration flags
	altScreen      bool // Use alternate screen buffer
	mouseAllMotion bool // Enable mouse motion events

	// Lifecycle management
	running bool
	mu      sync.Mutex

	// Event loop channels
	msgCh  chan model2.Msg // Incoming messages
	cmdCh  chan model2.Cmd // Commands to execute
	viewCh chan string     // View updates for rendering

	// Quit channel
	quitCh chan struct{}
}

// New creates a new Program with the given model.
// By default, uses stdin/stdout and no special terminal modes.
//
// Use With* options to customize behavior:
//
//	p := program.New(
//		MyModel{},
//		program.WithInput(customReader),
//		program.WithAltScreen(),
//		program.WithMouseAllMotion(),
//	)
func New[T any](m model2.Model[T], opts ...Option[T]) *Program[T] {
	p := &Program[T]{
		model:  m,
		input:  os.Stdin,  // Default
		output: os.Stdout, // Default
		quitCh: make(chan struct{}),
		msgCh:  make(chan model2.Msg, 100), // Buffered for performance
		cmdCh:  make(chan model2.Cmd, 10),
		viewCh: make(chan string, 10),
	}

	// Apply options
	for _, opt := range opts {
		opt(p)
	}

	return p
}

// Run starts the program and blocks until it quits.
// This is the main entry point for most applications.
//
// The event loop follows the Elm Architecture:
//  1. Call Init() to get initial command
//  2. Render initial view
//  3. Loop:
//     - Wait for message
//     - Check for QuitMsg (exit if found)
//     - Handle BatchMsg/SequenceMsg (expand to individual messages)
//     - Call Update(msg) to get new model and command
//     - Execute command (if any) in goroutine
//     - Render view
//
// Example:
//
//	p := program.New(MyModel{})
//	if err := p.Run(); err != nil {
//		log.Fatal(err)
//	}
//
// Returns error if program is already running or initialization fails.
//
//nolint:gocognit // Event loop orchestration requires sequential logic
func (p *Program[T]) Run() error {
	p.mu.Lock()
	if p.running {
		p.mu.Unlock()
		return fmt.Errorf("program already running")
	}
	p.running = true

	// Initialize terminal if not set (auto-detect best implementation)
	if p.terminal == nil {
		p.terminal = terminal.New()
	}
	p.mu.Unlock()

	// Cleanup on exit
	defer func() {
		p.mu.Lock()

		// Exit raw mode if we're in it
		if p.terminal != nil && p.terminal.IsInRawMode() {
			_ = p.terminal.ExitRawMode() // Best effort cleanup
		}

		// Exit alt screen if we're in it
		if p.terminal != nil && p.altScreen && p.terminal.IsInAltScreen() {
			_ = p.terminal.ExitAltScreen() // Best effort cleanup
		}

		p.running = false
		p.mu.Unlock()
	}()

	// Enter raw mode for TUI (best effort - may fail in test environments)
	p.mu.Lock()
	rawModeErr := p.terminal.EnterRawMode()
	// Note: We continue even if raw mode fails (for tests without actual terminal)

	// Enter alt screen if enabled (only if raw mode succeeded or not required)
	if p.altScreen {
		if err := p.terminal.EnterAltScreen(); err != nil {
			if rawModeErr == nil {
				_ = p.terminal.ExitRawMode() // Cleanup raw mode if we entered it
			}
			p.mu.Unlock()
			return fmt.Errorf("failed to enter alt screen: %w", err)
		}
	}
	p.mu.Unlock()

	// STEP 1: Call Init() to get initial command
	initCmd := p.model.Init()
	if initCmd != nil {
		p.executeCommand(initCmd)
	}
	// Start input reader
	p.startInputReader()

	// STEP 2: Render initial view
	p.renderView()

	// STEP 3: EVENT LOOP - THE HEART OF ELM ARCHITECTURE
	for {
		select {
		case msg := <-p.msgCh:
			// Check for quit
			if _, isQuit := msg.(model2.QuitMsg); isQuit {
				return nil // Exit loop
			}

			// Handle BatchMsg - expand to individual messages
			if batchMsg, ok := msg.(model2.BatchMsg); ok {
				for _, m := range batchMsg.Messages {
					p.msgCh <- m
				}
				continue
			}

			// Handle SequenceMsg - expand to individual messages (in order)
			if seqMsg, ok := msg.(model2.SequenceMsg); ok {
				for _, m := range seqMsg.Messages {
					p.msgCh <- m
				}
				continue
			}

			// Update model
			newModel, cmd := p.model.Update(msg)
			p.model = newModel

			// Execute command (if any)
			if cmd != nil {
				p.executeCommand(cmd)
			}

			// Render view
			p.renderView()

		case <-p.quitCh:
			return nil // External quit signal
		}
	}
}

// Start starts the program in a goroutine and returns immediately.
// Use Stop() to stop the program later.
//
// Runs the same event loop as Run(), but in a background goroutine.
//
// Example:
//
//	p := program.New(MyModel{})
//	if err := p.Start(); err != nil {
//		log.Fatal(err)
//	}
//	// ... do other work ...
//	p.Stop()
//
// Returns error if program is already running.
//
//nolint:gocognit // Event loop orchestration requires sequential logic
func (p *Program[T]) Start() error {
	p.mu.Lock()
	if p.running {
		p.mu.Unlock()
		return fmt.Errorf("program already running")
	}
	p.running = true

	// Initialize terminal if not set (auto-detect best implementation)
	if p.terminal == nil {
		p.terminal = terminal.New()
	}
	p.mu.Unlock()

	go func() {
		defer func() {
			p.mu.Lock()

			// Exit raw mode if we're in it
			if p.terminal != nil && p.terminal.IsInRawMode() {
				_ = p.terminal.ExitRawMode() // Best effort cleanup
			}

			// Exit alt screen if we're in it
			if p.terminal != nil && p.altScreen && p.terminal.IsInAltScreen() {
				_ = p.terminal.ExitAltScreen() // Best effort cleanup
			}

			p.running = false
			p.mu.Unlock()
		}()

		// Enter raw mode for TUI
		p.mu.Lock()
		rawModeErr := p.terminal.EnterRawMode()
		if rawModeErr != nil {
			// Can't return error from goroutine - just don't enter raw mode
			// This allows tests to run without requiring actual terminal
		} else {
			// Enter alt screen if enabled and raw mode succeeded
			if p.altScreen {
				if altScreenErr := p.terminal.EnterAltScreen(); altScreenErr != nil {
					_ = p.terminal.ExitRawMode() // Cleanup raw mode
					p.mu.Unlock()
					return
				}
			}
		}
		p.mu.Unlock()

		// Same event loop as Run(), but in goroutine
		initCmd := p.model.Init()
		if initCmd != nil {
			p.executeCommand(initCmd)
		}
		// Start input reader
		p.startInputReader()

		p.renderView()

		for {
			select {
			case msg := <-p.msgCh:
				if _, isQuit := msg.(model2.QuitMsg); isQuit {
					return
				}

				// Handle BatchMsg
				if batchMsg, ok := msg.(model2.BatchMsg); ok {
					for _, m := range batchMsg.Messages {
						p.msgCh <- m
					}
					continue
				}

				// Handle SequenceMsg
				if seqMsg, ok := msg.(model2.SequenceMsg); ok {
					for _, m := range seqMsg.Messages {
						p.msgCh <- m
					}
					continue
				}

				// Update
				newModel, cmd := p.model.Update(msg)
				p.model = newModel

				if cmd != nil {
					p.executeCommand(cmd)
				}

				p.renderView()

			case <-p.quitCh:
				return
			}
		}
	}()

	return nil
}

// Stop stops a running program gracefully.
// Blocks until the program has fully stopped.
//
// Safe to call multiple times.
func (p *Program[T]) Stop() {
	p.mu.Lock()
	if !p.running {
		p.mu.Unlock()
		return
	}
	p.mu.Unlock()

	// Signal quit
	select {
	case p.quitCh <- struct{}{}:
	default:
		// Already quitting
	}

	// Wait for running to become false (with timeout)
	timeout := time.After(1 * time.Second)
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			// Timeout - force stop
			p.mu.Lock()
			p.running = false
			p.mu.Unlock()
			return
		case <-ticker.C:
			p.mu.Lock()
			running := p.running
			p.mu.Unlock()

			if !running {
				return
			}
		}
	}
}

// Quit signals the program to quit.
// This is typically called from Update when receiving QuitMsg.
//
// Example:
//
//	func (m MyModel) Update(msg model.Msg) (model.Model[MyModel], model.Cmd) {
//		switch msg.(type) {
//		case model.QuitMsg:
//			// Program will call this internally
//			return m, nil
//		}
//		return m, nil
//	}
//
// Internal use - will be called from event loop in Day 4.
func (p *Program[T]) Quit() {
	select {
	case p.quitCh <- struct{}{}:
	default:
		// Already quitting
	}
}

// IsRunning returns true if the program is currently running.
func (p *Program[T]) IsRunning() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.running
}

// Send sends a message to the event loop from external code.
// This is useful for injecting messages from outside the program.
//
// Example:
//
//	p := program.New(model)
//	p.Start()
//
//	// From another goroutine:
//	p.Send(model.KeyMsg{Type: model.KeyEnter})
//
// Returns error if program is not running.
func (p *Program[T]) Send(msg model2.Msg) error {
	p.mu.Lock()
	running := p.running
	p.mu.Unlock()

	if !running {
		return fmt.Errorf("program not running")
	}

	select {
	case p.msgCh <- msg:
		return nil
	case <-time.After(100 * time.Millisecond):
		return fmt.Errorf("timeout sending message")
	}
}

// executeCommand runs a command in a goroutine and sends result to msgCh.
func (p *Program[T]) executeCommand(cmd model2.Cmd) {
	go func() {
		msg := cmd() // Execute command (may block)

		// Send result back to event loop
		select {
		case p.msgCh <- msg:
		case <-p.quitCh:
			// Program quitting, don't send
		}
	}()
}

// renderView renders the current model's view to output.
func (p *Program[T]) renderView() {
	view := p.model.View()

	// Write to output
	// Day 4: simple write. Day 5: will add diff rendering
	_, _ = p.output.Write([]byte(view))
	// Ignore write errors for now (may occur during tests)
}

// startInputReader starts reading input in a goroutine.
// Creates a new goroutine with cancellation support for ExecProcess.
func (p *Program[T]) startInputReader() {
	// Create input reader if not yet created
	if p.inputReader == nil {
		p.inputReader = input.NewReader(p.input)
	}

	// Create cancellation context for this inputReader goroutine
	ctx, cancel := context.WithCancel(context.Background())

	// Set state under mutex protection
	p.mu.Lock()
	p.inputReaderCancel = cancel
	p.inputReaderDone = make(chan struct{})
	p.inputReaderRunning = true
	p.inputReaderGeneration++             // Increment generation for new goroutine
	generation := p.inputReaderGeneration // Capture for this goroutine
	p.mu.Unlock()

	go func() {
		defer func() {
			// Signal that goroutine has exited
			close(p.inputReaderDone)

			// Only clear flag if we're still the current generation
			// (prevents race with restart after stop timeout)
			p.mu.Lock()
			if p.inputReaderGeneration == generation {
				p.inputReaderRunning = false
			}
			p.mu.Unlock()
		}()

		for {
			// Check if context canceled before reading
			select {
			case <-ctx.Done():
				return
			default:
			}

			// Read input (blocks until input available)
			// Note: Read() itself is blocking, but we check ctx.Done() on each iteration
			msg, err := p.inputReader.Read()
			if err != nil {
				// EOF or error - stop reading
				return
			}

			// Skip nil messages (unknown sequences)
			if msg == nil {
				continue
			}

			// Send to event loop (with cancellation check)
			select {
			case p.msgCh <- msg:
			case <-ctx.Done():
				return
			case <-p.quitCh:
				return
			}
		}
	}()
}

// stopInputReader gracefully stops the inputReader goroutine.
// Blocks until the goroutine has fully exited.
// Safe to call even if inputReader is not running.
//
// This MUST be called before ExecProcess to prevent stdin stealing.
func (p *Program[T]) stopInputReader() {
	p.mu.Lock()
	running := p.inputReaderRunning
	cancel := p.inputReaderCancel
	done := p.inputReaderDone
	p.mu.Unlock()

	if !running || cancel == nil {
		// Already stopped or never started
		return
	}

	// Signal cancellation
	cancel()

	// Wait for goroutine to exit (with timeout)
	select {
	case <-done:
		// Goroutine exited gracefully
	case <-time.After(500 * time.Millisecond):
		// Timeout - goroutine may be stuck in Read()
		// This is acceptable - inputReader.Read() will unblock on next stdin activity
		// The goroutine will exit when it checks ctx.Done()
	}

	// Increment generation to invalidate old goroutine's defer
	// (prevents race where old goroutine clears flag after restart)
	// Then ensure flag is cleared for synchronization with restartInputReader
	p.mu.Lock()
	p.inputReaderGeneration++ // Invalidate old goroutine
	p.inputReaderRunning = false
	p.inputReaderCancel = nil
	p.inputReaderDone = nil
	p.mu.Unlock()
}

// restartInputReader restarts the inputReader goroutine after ExecProcess.
// Must be called after stopInputReader to resume normal TUI input handling.
func (p *Program[T]) restartInputReader() {
	p.mu.Lock()
	running := p.inputReaderRunning
	p.mu.Unlock()

	if running {
		// Already running - don't start duplicate
		return
	}

	// Start new inputReader goroutine
	p.startInputReader()
}

// ┌─────────────────────────────────────────────────────────────────┐.
// │ External Process Execution                                      │.
// └─────────────────────────────────────────────────────────────────┘.

// ExecProcess executes an external interactive command with full terminal control.
//
// The program temporarily:.
//  1. Stops inputReader goroutine (prevents stdin stealing).
//  2. Exits raw mode (restores cooked mode for external command).
//  3. Exits alternate screen buffer (if active).
//  4. Shows cursor.
//  5. Gives command full control of stdin/stdout/stderr.
//  6. Waits for command completion (BLOCKING - call from Cmd goroutine!).
//  7. Hides cursor.
//  8. Re-enters alternate screen buffer (if was active).
//  9. Re-enters raw mode (restores TUI input handling).
//
// 10. Restarts inputReader goroutine.
// 11. Forces full TUI refresh.
//
// This is essential for running interactive commands like:.
//   - Text editors (vim, nano, emacs).
//   - Interactive shells (bash, python REPL, claude).
//   - Pagers (less, more).
//   - SSH sessions.
//   - Any command requiring TTY control.
//
// Example:
//
//	func (m Model) Update(msg Msg) (Model, Cmd) {
//	    switch msg := msg.(type) {
//	    case ExecVimMsg:
//	        return m, func() Msg {
//	            cmd := exec.Command("vim", "file.txt")
//	            err := m.program.ExecProcess(cmd)
//	            return VimFinishedMsg{Err: err}
//	        }
//	    }
//	    return m, nil
//	}
//
// IMPORTANT:.
//   - Must be called from a Cmd goroutine (NOT from Update directly).
//   - Blocks until command completes.
//   - Requires terminal to be set (auto-created in Run or via WithTerminal).
//   - inputReader is stopped before command and restarted after.
//
// Returns error if command execution fails.
func (p *Program[T]) ExecProcess(cmd *exec.Cmd) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Validate terminal exists.
	if p.terminal == nil {
		// Auto-create terminal if not set (fallback for edge cases).
		p.terminal = terminal.New()
	}

	// Validate command.
	if cmd == nil {
		return fmt.Errorf("exec: cmd is nil")
	}

	// STEP 1: Stop inputReader goroutine (CRITICAL FIX!)
	// Must release mutex before calling stopInputReader (it needs mutex internally)
	p.mu.Unlock()
	p.stopInputReader()
	p.mu.Lock()

	// STEP 2: Save TUI state.
	// Remember if we were in raw mode and alt screen.
	wasInRawMode := p.terminal.IsInRawMode()
	wasInAltScreen := p.terminal.IsInAltScreen()

	// STEP 3: Exit raw mode (restore cooked mode for external command).
	// CRITICAL: External commands (vim, ssh, python REPL) expect cooked mode!
	if wasInRawMode {
		if err := p.terminal.ExitRawMode(); err != nil {
			// Restore inputReader before returning error
			p.mu.Unlock()
			p.restartInputReader()
			p.mu.Lock()
			return fmt.Errorf("exec: failed to exit raw mode: %w", err)
		}
	}

	// STEP 4: Exit alternate screen (if active).
	if wasInAltScreen {
		if err := p.terminal.ExitAltScreen(); err != nil {
			// Restore raw mode before returning error
			if wasInRawMode {
				_ = p.terminal.EnterRawMode() // Best effort
			}
			// Restore inputReader before returning error
			p.mu.Unlock()
			p.restartInputReader()
			p.mu.Lock()
			return fmt.Errorf("exec: failed to exit alt screen: %w", err)
		}
	}

	// STEP 5: Show cursor (always restore visibility).
	if err := p.terminal.ShowCursor(); err != nil {
		// Non-fatal - continue anyway.
		// Some terminals may not support cursor control.
	}

	// STEP 6: Setup command I/O - give it full terminal control.
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// STEP 7: Run command (BLOCKING).
	// Release mutex while command runs (may take long time)
	p.mu.Unlock()
	cmdErr := cmd.Run()
	p.mu.Lock()

	// STEP 8: ALWAYS restore TUI state (even if command failed).
	// Hide cursor first (restore TUI cursor state).
	if err := p.terminal.HideCursor(); err != nil {
		// Non-fatal - continue anyway.
	}

	// STEP 9: Re-enter alternate screen (if we were in it before).
	if wasInAltScreen {
		if err := p.terminal.EnterAltScreen(); err != nil {
			// CRITICAL: TUI state corrupted.
			// Still restart inputReader before returning
			p.mu.Unlock()
			p.restartInputReader()
			p.mu.Lock()
			if cmdErr != nil {
				return fmt.Errorf("exec: command failed (%v) and failed to restore alt screen: %w", cmdErr, err)
			}
			return fmt.Errorf("exec: failed to restore alt screen: %w", err)
		}
	}

	// STEP 10: Re-enter raw mode (restore TUI input mode).
	// CRITICAL: Must restore raw mode for TUI to receive character input!
	if wasInRawMode {
		if err := p.terminal.EnterRawMode(); err != nil {
			// CRITICAL: TUI can't receive input without raw mode.
			// Restore inputReader before returning
			p.mu.Unlock()
			p.restartInputReader()
			p.mu.Lock()
			if cmdErr != nil {
				return fmt.Errorf("exec: command failed (%v) and failed to re-enter raw mode: %w", cmdErr, err)
			}
			return fmt.Errorf("exec: failed to re-enter raw mode: %w", err)
		}
	}

	// STEP 11: Restart inputReader goroutine (CRITICAL FIX!)
	// Release mutex before restarting
	p.mu.Unlock()
	p.restartInputReader()
	p.mu.Lock()

	// STEP 12: Force full redraw.
	// The TUI needs to repaint after external command.
	p.renderView()

	// Return original command error (if any).
	return cmdErr
}
