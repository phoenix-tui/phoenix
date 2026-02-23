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
	"github.com/phoenix-tui/phoenix/tea/internal/infrastructure/renderer"
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

	// Inline renderer for non-alt-screen mode.
	// Initialized lazily in renderView() on the first render call.
	// Tracks linesRendered so subsequent renders overwrite the previous frame
	// instead of appending below it.
	inlineRenderer *renderer.InlineRenderer

	// Suspend/Resume state (for ExecProcess and public API)
	suspended    bool            // True if TUI is suspended
	suspendState *suspendedState // Saved state when suspended
}

// suspendedState holds terminal state saved during Suspend().
// Used to restore correct state during Resume().
type suspendedState struct {
	wasInRawMode   bool
	wasInAltScreen bool
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

			// Intercept WindowSizeMsg to keep inline renderer dimensions current.
			if sizeMsg, ok := msg.(model2.WindowSizeMsg); ok && !p.altScreen {
				if p.inlineRenderer != nil {
					p.inlineRenderer.Resize(sizeMsg.Width, sizeMsg.Height)
				}
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

				// Intercept WindowSizeMsg to keep inline renderer dimensions current.
				if sizeMsg, ok := msg.(model2.WindowSizeMsg); ok && !p.altScreen {
					if p.inlineRenderer != nil {
						p.inlineRenderer.Resize(sizeMsg.Width, sizeMsg.Height)
					}
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
//
// In inline (non-alt-screen) mode the InlineRenderer is used to overwrite
// the previous frame using ANSI cursor-up sequences, preventing the
// stacked-output bug where each update appends below the last.
//
// In alt-screen mode a plain write is used; the alt-screen renderer will be
// integrated in a future release.
func (p *Program[T]) renderView() {
	view := p.model.View()

	if !p.altScreen {
		// Lazily initialize the inline renderer on the first call.
		// Width/height start at 0 (disables truncation/clipping) until
		// a WindowSizeMsg updates the dimensions.
		if p.inlineRenderer == nil {
			p.inlineRenderer = renderer.NewInlineRenderer(p.output, 0, 0)
		}
		// Render errors are non-fatal (e.g. write to closed pipe during tests).
		_ = p.inlineRenderer.Render(view)
		return
	}

	// Alt-screen mode: plain write (cursor is at absolute position after \x1b[H).
	_, _ = p.output.Write([]byte(view))
}

// startInputReader starts reading input in a goroutine.
// Creates a new goroutine with cancellation support for ExecProcess.
//
// CRITICAL (v0.1.1): Always creates a NEW Reader because CancelableReader
// cannot be reused after Cancel() - it permanently returns EOF.
func (p *Program[T]) startInputReader() {
	// Always create a new Reader (CancelableReader cannot be reused after Cancel)
	// This ensures fresh state after ExecProcess
	p.inputReader = input.NewReader(p.input)

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

			// Only clear state if we're still the current generation
			// (prevents race with restart after stop timeout)
			p.mu.Lock()
			if p.inputReaderGeneration == generation {
				p.inputReaderRunning = false
				p.inputReaderCancel = nil
				p.inputReaderDone = nil
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
//
// CRITICAL FIX (v0.1.1): Now calls inputReader.Cancel() to immediately
// unblock any pending Read() operations. Without this, the goroutine
// would remain blocked in Read() causing race conditions with ExecProcess.
func (p *Program[T]) stopInputReader() {
	p.mu.Lock()
	running := p.inputReaderRunning
	cancel := p.inputReaderCancel
	done := p.inputReaderDone
	inputReader := p.inputReader // Capture for Cancel()
	p.mu.Unlock()

	if !running || cancel == nil {
		// Already stopped or never started
		return
	}

	// STEP 1: Cancel the reader itself (unblocks Read immediately).
	// Uses pipe-based relay: closing pipe writer causes readLoopPipe to return
	// io.EOF instantly. SetReadDeadline + UnblockStdinRead unblock the relay.
	if inputReader != nil {
		inputReader.Cancel()
	}

	// STEP 2: Signal context cancellation
	cancel()

	// STEP 3: Wait for goroutine to exit (now guaranteed to succeed quickly
	// thanks to pipe-based CancelableReader — Cancel closes the pipe)
	select {
	case <-done:
		// Goroutine exited gracefully
	case <-time.After(200 * time.Millisecond):
		// Safety net — should rarely hit with pipe-based CancelableReader.
		// Relay goroutine may still be blocked in stdin Read on platforms
		// where SetReadDeadline is not supported; it will exit on next input.
	}

	// STEP 4: Wait for CancelableReader background goroutines to fully stop
	if inputReader != nil {
		inputReader.WaitForShutdown()
	}

	// Clean up state
	// Increment generation to invalidate old goroutine's defer
	// (prevents race where old goroutine clears flag after restart)
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
// │ Suspend/Resume (Level 1 TTY Control)                            │.
// └─────────────────────────────────────────────────────────────────┘.

// Suspend temporarily suspends the TUI and restores terminal to normal mode.
// This is the first step when running interactive external commands.
//
// Suspend performs the following in order:
//  1. Stops the inputReader goroutine (releases stdin)
//  2. Saves current terminal state (raw mode, alt screen)
//  3. Exits raw mode (restores cooked mode)
//  4. Exits alternate screen (if active)
//  5. Shows cursor
//
// After Suspend, the terminal is in a normal state suitable for:
//   - Running interactive commands (vim, ssh, python REPL)
//   - User interaction outside the TUI
//   - Shell commands that expect cooked mode
//
// Call Resume() to restore the TUI after the external operation.
//
// Example:
//
//	func (m Model) Update(msg Msg) (Model, Cmd) {
//	    switch msg := msg.(type) {
//	    case RunVimMsg:
//	        return m, func() Msg {
//	            p.Suspend()
//	            cmd := exec.Command("vim", "file.txt")
//	            err := cmd.Run()
//	            p.Resume()
//	            return VimFinishedMsg{Err: err}
//	        }
//	    }
//	    return m, nil
//	}
//
// Returns error if suspension fails (e.g., terminal operation error).
// Safe to call multiple times - subsequent calls are no-ops.
func (p *Program[T]) Suspend() error {
	p.mu.Lock()
	if p.suspended {
		p.mu.Unlock()
		return nil // Already suspended
	}

	// Validate terminal exists
	if p.terminal == nil {
		p.terminal = terminal.New()
	}
	p.mu.Unlock()

	// STEP 1: Stop inputReader goroutine (releases stdin)
	// Must be done outside mutex (stopInputReader uses mutex internally)
	p.stopInputReader()

	p.mu.Lock()
	defer p.mu.Unlock()

	// STEP 2: Save current terminal state
	state := &suspendedState{
		wasInRawMode:   p.terminal.IsInRawMode(),
		wasInAltScreen: p.terminal.IsInAltScreen(),
	}

	// STEP 3: Exit raw mode (restore cooked mode for external commands)
	if state.wasInRawMode {
		if err := p.terminal.ExitRawMode(); err != nil {
			// Restart inputReader before returning error
			p.mu.Unlock()
			p.restartInputReader()
			p.mu.Lock()
			return fmt.Errorf("suspend: failed to exit raw mode: %w", err)
		}
	}

	// STEP 4: Exit alternate screen (if active)
	if state.wasInAltScreen {
		if err := p.terminal.ExitAltScreen(); err != nil {
			// Restore raw mode before returning error
			if state.wasInRawMode {
				_ = p.terminal.EnterRawMode() // Best effort
			}
			// Restart inputReader before returning error
			p.mu.Unlock()
			p.restartInputReader()
			p.mu.Lock()
			return fmt.Errorf("suspend: failed to exit alt screen: %w", err)
		}
	}

	// STEP 5: Show cursor (always restore visibility for external commands)
	if err := p.terminal.ShowCursor(); err != nil {
		// Non-fatal - continue anyway
		// Some terminals may not support cursor control
	}

	// Mark as suspended and save state
	p.suspended = true
	p.suspendState = state

	return nil
}

// Resume restores the TUI after a Suspend.
// This is the second step after running interactive external commands.
//
// Resume performs the following in order:
//  1. Hides cursor
//  2. Re-enters alternate screen (if was active before Suspend)
//  3. Re-enters raw mode (if was active before Suspend)
//  4. Restarts the inputReader goroutine
//  5. Forces a full redraw
//
// Example:
//
//	p.Suspend()
//	cmd := exec.Command("vim", "file.txt")
//	cmd.Stdin = os.Stdin
//	cmd.Stdout = os.Stdout
//	cmd.Stderr = os.Stderr
//	err := cmd.Run()
//	p.Resume() // Restore TUI
//
// Returns error if restoration fails (e.g., terminal operation error).
// Safe to call multiple times - subsequent calls are no-ops.
// If Resume fails partway, terminal may be in inconsistent state.
func (p *Program[T]) Resume() error {
	p.mu.Lock()
	if !p.suspended {
		p.mu.Unlock()
		return nil // Not suspended
	}

	state := p.suspendState
	if state == nil {
		// Should not happen, but be defensive
		p.suspended = false
		p.mu.Unlock()
		return nil
	}
	p.mu.Unlock()

	// STEP 1: Hide cursor (restore TUI cursor state)
	p.mu.Lock()
	if err := p.terminal.HideCursor(); err != nil {
		// Non-fatal - continue anyway
	}
	p.mu.Unlock()

	// STEP 2: Re-enter alternate screen (if was active before)
	p.mu.Lock()
	if state.wasInAltScreen {
		if err := p.terminal.EnterAltScreen(); err != nil {
			p.mu.Unlock()
			return fmt.Errorf("resume: failed to restore alt screen: %w", err)
		}
	}
	p.mu.Unlock()

	// STEP 3: Re-enter raw mode (if was active before)
	p.mu.Lock()
	if state.wasInRawMode {
		if err := p.terminal.EnterRawMode(); err != nil {
			p.mu.Unlock()
			return fmt.Errorf("resume: failed to restore raw mode: %w", err)
		}
	}

	// Clear suspended state
	p.suspended = false
	p.suspendState = nil
	p.mu.Unlock()

	// STEP 4: Restart inputReader goroutine
	p.restartInputReader()

	// STEP 5: Force full redraw.
	// Clear the inline renderer diff cache so all lines are repainted.
	// The external command may have written arbitrary content to the terminal,
	// invalidating our previous frame tracking.
	if p.inlineRenderer != nil {
		p.inlineRenderer.Repaint()
	}
	p.renderView()

	return nil
}

// IsSuspended returns true if the TUI is currently suspended.
func (p *Program[T]) IsSuspended() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.suspended
}

// ┌─────────────────────────────────────────────────────────────────┐.
// │ External Process Execution                                      │.
// └─────────────────────────────────────────────────────────────────┘.

// ExecProcess executes an external interactive command with full terminal control.
//
// The program temporarily suspends the TUI (via Suspend) and gives the command
// full control of stdin/stdout/stderr. After the command completes, the TUI is
// restored (via Resume).
//
// This is essential for running interactive commands like:
//   - Text editors (vim, nano, emacs)
//   - Interactive shells (bash, python REPL, claude)
//   - Pagers (less, more)
//   - SSH sessions
//   - Any command requiring TTY control
//
// Internally uses Suspend/Resume pattern:
//  1. Suspend() - stops input, exits raw mode, exits alt screen, shows cursor
//  2. Run command (BLOCKING)
//  3. Resume() - restores raw mode, alt screen, restarts input, redraws
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
// IMPORTANT:
//   - Must be called from a Cmd goroutine (NOT from Update directly)
//   - Blocks until command completes
//   - Requires terminal to be set (auto-created in Run or via WithTerminal)
//   - Uses Suspend/Resume internally for terminal management
//
// Returns error if command execution fails.
func (p *Program[T]) ExecProcess(cmd *exec.Cmd) error {
	// Validate command
	if cmd == nil {
		return fmt.Errorf("exec: cmd is nil")
	}

	// STEP 1: Suspend TUI (stop input, exit raw mode, exit alt screen, show cursor)
	if err := p.Suspend(); err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	// STEP 2: Setup command I/O - give it full terminal control
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// STEP 3: Run command (BLOCKING)
	cmdErr := cmd.Run()

	// STEP 4: Resume TUI (restore raw mode, alt screen, restart input, redraw)
	// ALWAYS resume, even if command failed
	if err := p.Resume(); err != nil {
		if cmdErr != nil {
			return fmt.Errorf("exec: command failed (%v) and %w", cmdErr, err)
		}
		return fmt.Errorf("exec: %w", err)
	}

	// Return original command error (if any)
	return cmdErr
}

// ┌─────────────────────────────────────────────────────────────────┐.
// │ Advanced TTY Control (Level 2)                                  │.
// └─────────────────────────────────────────────────────────────────┘.

// ExecProcessWithTTY executes an external command with advanced TTY control.
//
// This provides Level 2 TTY control with platform-specific enhancements:
//   - Unix/Linux/macOS: Uses tcsetpgrp() for proper foreground process group transfer
//   - Windows: Enhanced console mode management
//
// Advantages over ExecProcess:
//   - Proper job control (Ctrl+Z in child suspends child, not parent)
//   - Child can use its own signal handlers
//   - Better isolation between parent and child processes
//
// Example:
//
//	func (m Model) Update(msg Msg) (Model, Cmd) {
//	    switch msg := msg.(type) {
//	    case RunShellMsg:
//	        return m, func() Msg {
//	            cmd := exec.Command("bash")
//	            opts := TTYOptions{
//	                TransferForeground: true,
//	                CreateProcessGroup: true,
//	            }
//	            err := m.program.ExecProcessWithTTY(cmd, opts)
//	            return ShellExitedMsg{Err: err}
//	        }
//	    }
//	    return m, nil
//	}
//
// IMPORTANT:
//   - Must be called from a Cmd goroutine (NOT from Update directly)
//   - Blocks until command completes
//   - Falls back to ExecProcess if TTY control unavailable
//   - Platform-specific implementation (see tty_control_unix.go, tty_control_windows.go)
//
// Returns error if command execution fails.
func (p *Program[T]) ExecProcessWithTTY(cmd *exec.Cmd, opts TTYOptions) error {
	// Platform-specific implementation in tty_control_unix.go or tty_control_windows.go
	return p.execWithTTYControl(cmd, opts)
}

// execWithTTYControl is implemented in platform-specific files:
//   - tty_control_unix.go (Unix/Linux/macOS) - uses tcsetpgrp()
//   - tty_control_windows.go (Windows) - uses SetConsoleMode()
//
// The build tags ensure only one implementation is compiled per platform.
