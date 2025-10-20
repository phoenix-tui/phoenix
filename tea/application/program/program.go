// Package program provides the Program type that orchestrates the Elm Architecture event loop.
// It manages application lifecycle, message passing, and rendering coordination.
package program

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/phoenix-tui/phoenix/tea/domain/model"
	"github.com/phoenix-tui/phoenix/tea/infrastructure/input"
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
	model model.Model[T]

	// I/O streams
	input  io.Reader
	output io.Writer

	// Input reader for parsing stdin
	inputReader *input.Reader

	// Configuration flags
	altScreen      bool // Use alternate screen buffer
	mouseAllMotion bool // Enable mouse motion events

	// Lifecycle management
	running bool
	mu      sync.Mutex

	// Event loop channels
	msgCh  chan model.Msg // Incoming messages
	cmdCh  chan model.Cmd // Commands to execute
	viewCh chan string    // View updates for rendering

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
func New[T any](m model.Model[T], opts ...Option[T]) *Program[T] {
	p := &Program[T]{
		model:  m,
		input:  os.Stdin,  // Default
		output: os.Stdout, // Default
		quitCh: make(chan struct{}),
		msgCh:  make(chan model.Msg, 100), // Buffered for performance
		cmdCh:  make(chan model.Cmd, 10),
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
	p.mu.Unlock()

	// Cleanup on exit
	defer func() {
		p.mu.Lock()
		p.running = false
		p.mu.Unlock()
	}()

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
			if _, isQuit := msg.(model.QuitMsg); isQuit {
				return nil // Exit loop
			}

			// Handle BatchMsg - expand to individual messages
			if batchMsg, ok := msg.(model.BatchMsg); ok {
				for _, m := range batchMsg.Messages {
					p.msgCh <- m
				}
				continue
			}

			// Handle SequenceMsg - expand to individual messages (in order)
			if seqMsg, ok := msg.(model.SequenceMsg); ok {
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
	p.mu.Unlock()

	go func() {
		defer func() {
			p.mu.Lock()
			p.running = false
			p.mu.Unlock()
		}()

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
				if _, isQuit := msg.(model.QuitMsg); isQuit {
					return
				}

				// Handle BatchMsg
				if batchMsg, ok := msg.(model.BatchMsg); ok {
					for _, m := range batchMsg.Messages {
						p.msgCh <- m
					}
					continue
				}

				// Handle SequenceMsg
				if seqMsg, ok := msg.(model.SequenceMsg); ok {
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
func (p *Program[T]) Send(msg model.Msg) error {
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
func (p *Program[T]) executeCommand(cmd model.Cmd) {
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
func (p *Program[T]) startInputReader() {
	// Create input reader if not yet created
	if p.inputReader == nil {
		p.inputReader = input.NewReader(p.input)
	}

	go func() {
		for {
			// Read input (blocks until input available)
			msg, err := p.inputReader.Read()
			if err != nil {
				// EOF or error - stop reading
				return
			}

			// Skip nil messages (unknown sequences)
			if msg == nil {
				continue
			}

			// Send to event loop
			select {
			case p.msgCh <- msg:
			case <-p.quitCh:
				return
			}
		}
	}()
}
