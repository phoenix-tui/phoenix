package program

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/phoenix-tui/phoenix/tea/domain/model"
)

// TestModel is a simple test model for testing Program.
type TestModel struct {
	value       int
	updateCount int
	initCalled  bool
	lastMsg     string
}

func (m TestModel) Init() model.Cmd {
	// Return command that sends a test message
	return func() model.Msg {
		return testInitMsg{}
	}
}

func (m TestModel) Update(msg model.Msg) (model.Model[TestModel], model.Cmd) {
	m.updateCount++

	switch msg := msg.(type) {
	case testInitMsg:
		m.lastMsg = "init"
		m.initCalled = true
	case model.KeyMsg:
		m.lastMsg = msg.String()
		if msg.String() == "+" {
			m.value++
		}
		if msg.String() == "q" {
			return m, func() model.Msg { return model.QuitMsg{} }
		}
	case model.QuitMsg:
		m.lastMsg = "quit"
	}
	return m, nil
}

func (m TestModel) View() string {
	return fmt.Sprintf("Value: %d, Updates: %d, Last: %s\n", m.value, m.updateCount, m.lastMsg)
}

// Test message types
type testInitMsg struct{}

func (t testInitMsg) String() string { return "init" }

// TestNew verifies constructor creates program correctly.
func TestNew(t *testing.T) {
	m := TestModel{value: 42}
	p := New(m)

	// Verify model set
	if p.model == nil {
		t.Error("model should be set")
	}

	// Verify default input/output (os.Stdin/Stdout)
	if p.input != os.Stdin {
		t.Error("default input should be os.Stdin")
	}
	if p.output != os.Stdout {
		t.Error("default output should be os.Stdout")
	}

	// Verify quitCh created
	if p.quitCh == nil {
		t.Error("quitCh should be created")
	}

	// Verify not running initially
	if p.running {
		t.Error("program should not be running initially")
	}
}

// TestNew_WithOptions verifies options are applied correctly.
func TestNew_WithOptions(t *testing.T) {
	m := TestModel{}

	customInput := strings.NewReader("test input")
	var customOutput bytes.Buffer

	p := New(
		m,
		WithInput[TestModel](customInput),
		WithOutput[TestModel](&customOutput),
		WithAltScreen[TestModel](),
		WithMouseAllMotion[TestModel](),
	)

	// Verify WithInput sets custom reader
	if p.input != customInput {
		t.Error("WithInput should set custom reader")
	}

	// Verify WithOutput sets custom writer
	if p.output != &customOutput {
		t.Error("WithOutput should set custom writer")
	}

	// Verify WithAltScreen sets flag
	if !p.altScreen {
		t.Error("WithAltScreen should set altScreen flag")
	}

	// Verify WithMouseAllMotion sets flag
	if !p.mouseAllMotion {
		t.Error("WithMouseAllMotion should set mouseAllMotion flag")
	}
}

// TestProgram_IsRunning verifies status check works correctly.
func TestProgram_IsRunning(t *testing.T) {
	p := New(TestModel{})

	// Initially false
	if p.IsRunning() {
		t.Error("program should not be running initially")
	}

	// True after Start()
	if err := p.Start(); err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	// Give goroutine time to start
	time.Sleep(50 * time.Millisecond)

	if !p.IsRunning() {
		t.Error("program should be running after Start()")
	}

	// False after Stop()
	p.Stop()

	if p.IsRunning() {
		t.Error("program should not be running after Stop()")
	}
}

// TestProgram_Start_Stop verifies lifecycle works correctly.
func TestProgram_Start_Stop(t *testing.T) {
	p := New(TestModel{})

	// Start() sets running = true
	if err := p.Start(); err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	// Give goroutine time to start
	time.Sleep(50 * time.Millisecond)

	if !p.IsRunning() {
		t.Error("IsRunning() should return true after Start()")
	}

	// Stop() sets running = false
	p.Stop()

	if p.IsRunning() {
		t.Error("IsRunning() should return false after Stop()")
	}

	// Stop() is idempotent (safe to call twice)
	p.Stop() // Should not panic
}

// TestProgram_Start_AlreadyRunning verifies error handling for double start.
func TestProgram_Start_AlreadyRunning(t *testing.T) {
	p := New(TestModel{})

	// First Start() succeeds
	if err := p.Start(); err != nil {
		t.Fatalf("first Start() failed: %v", err)
	}

	// Give goroutine time to start
	time.Sleep(50 * time.Millisecond)

	// Second Start() returns error
	err := p.Start()
	if err == nil {
		t.Error("second Start() should return error")
	}

	// Error message contains "already running"
	if !strings.Contains(err.Error(), "already running") {
		t.Errorf("error message should contain 'already running', got: %v", err)
	}

	p.Stop()
}

// TestProgram_Run verifies blocking execution works correctly.
func TestProgram_Run(t *testing.T) {
	p := New(TestModel{})

	done := make(chan error)
	go func() {
		done <- p.Run() // Run in goroutine
	}()

	// Give it time to start
	time.Sleep(50 * time.Millisecond)

	// Should still be running
	if !p.IsRunning() {
		t.Error("program should be running")
	}

	// Quit
	p.Quit()

	// Should finish
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Run() returned error: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Run() did not finish after Quit()")
	}

	// Should no longer be running
	if p.IsRunning() {
		t.Error("program should not be running after Quit()")
	}
}

// TestProgram_Run_AlreadyRunning verifies error handling for Run after Start.
func TestProgram_Run_AlreadyRunning(t *testing.T) {
	p := New(TestModel{})

	// Start() first
	if err := p.Start(); err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	// Give goroutine time to start
	time.Sleep(50 * time.Millisecond)

	// Run() should return error
	err := p.Run()
	if err == nil {
		t.Error("Run() should return error when already running")
	}

	// Error message contains "already running"
	if !strings.Contains(err.Error(), "already running") {
		t.Errorf("error message should contain 'already running', got: %v", err)
	}

	p.Stop()
}

// TestProgram_Quit verifies quit signal works correctly.
func TestProgram_Quit(t *testing.T) {
	p := New(TestModel{})

	done := make(chan error)
	go func() {
		done <- p.Run()
	}()

	// Give it time to start
	time.Sleep(50 * time.Millisecond)

	// Quit() sends to quitCh
	p.Quit()

	// Should finish quickly
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Run() returned error: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Run() did not finish after Quit()")
	}

	// Multiple Quit() calls don't panic (buffered channel handling)
	p.Quit() // Should not block or panic
	p.Quit() // Should not block or panic
}

// TestProgram_Stop_NotRunning verifies edge case of stopping non-running program.
func TestProgram_Stop_NotRunning(t *testing.T) {
	p := New(TestModel{})

	// Stop() on non-running program doesn't panic
	p.Stop() // Should not panic

	// Stop() is idempotent
	p.Stop() // Should not panic
	p.Stop() // Should not panic
}

// TestProgram_Concurrency verifies thread safety of concurrent access.
func TestProgram_Concurrency(t *testing.T) {
	p := New(TestModel{})

	var wg sync.WaitGroup

	// Start program
	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	// Give goroutine time to start
	time.Sleep(50 * time.Millisecond)

	// 10 goroutines calling IsRunning()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_ = p.IsRunning() // Should not race
			}
		}()
	}

	wg.Wait()
	p.Stop()
}

// TestProgram_MultipleOptions verifies multiple options work together.
func TestProgram_MultipleOptions(t *testing.T) {
	m := TestModel{}

	customInput := strings.NewReader("test")
	var customOutput bytes.Buffer

	p := New(
		m,
		WithInput[TestModel](customInput),
		WithOutput[TestModel](&customOutput),
		WithAltScreen[TestModel](),
		WithMouseAllMotion[TestModel](),
	)

	// All flags/fields should be set correctly
	if p.input != customInput {
		t.Error("input not set correctly")
	}
	if p.output != &customOutput {
		t.Error("output not set correctly")
	}
	if !p.altScreen {
		t.Error("altScreen not set correctly")
	}
	if !p.mouseAllMotion {
		t.Error("mouseAllMotion not set correctly")
	}
}

// TestProgram_RunTwice verifies error when calling Run twice.
func TestProgram_RunTwice(t *testing.T) {
	p := New(TestModel{})

	done1 := make(chan error)
	go func() {
		done1 <- p.Run()
	}()

	// Give first Run time to start
	time.Sleep(50 * time.Millisecond)

	// Second Run should fail immediately
	err := p.Run()
	if err == nil {
		t.Error("second Run() should return error")
	}
	if !strings.Contains(err.Error(), "already running") {
		t.Errorf("error should contain 'already running', got: %v", err)
	}

	// Clean up first Run
	p.Quit()
	<-done1
}

// ========================================
// EVENT LOOP TESTS (Day 4)
// ========================================

// TestProgram_EventLoop_Init verifies Init command is executed.
func TestProgram_EventLoop_Init(t *testing.T) {
	var buf bytes.Buffer

	m := TestModel{}
	p := New(m, WithOutput[TestModel](&buf))

	// Start program
	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	// Give event loop time to process Init
	time.Sleep(150 * time.Millisecond)

	// Send quit
	if err := p.Send(model.QuitMsg{}); err != nil {
		t.Fatal(err)
	}
	p.Stop()

	// Check that Init was called and view was rendered
	output := buf.String()

	// Should have rendered at least once
	if !strings.Contains(output, "Value:") {
		t.Errorf("view should have been rendered, got: %s", output)
	}

	// Should show init message was processed
	if !strings.Contains(output, "Last: init") {
		t.Errorf("init command should have been executed, got: %s", output)
	}
}

// TestProgram_EventLoop_Update verifies Update is called for messages.
func TestProgram_EventLoop_Update(t *testing.T) {
	var buf bytes.Buffer

	m := TestModel{value: 0}
	p := New(m, WithOutput[TestModel](&buf))

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	// Give time for init
	time.Sleep(50 * time.Millisecond)

	// Send messages
	if err := p.Send(model.KeyMsg{Type: model.KeyRune, Rune: '+'}); err != nil {
		t.Fatal(err)
	}
	time.Sleep(50 * time.Millisecond)

	if err := p.Send(model.KeyMsg{Type: model.KeyRune, Rune: '+'}); err != nil {
		t.Fatal(err)
	}
	time.Sleep(50 * time.Millisecond)

	if err := p.Send(model.QuitMsg{}); err != nil {
		t.Fatal(err)
	}
	p.Stop()

	output := buf.String()

	// Should show value incremented twice
	if !strings.Contains(output, "Value: 2") {
		t.Errorf("expected Value: 2 in output, got: %s", output)
	}
}

// TestProgram_EventLoop_View verifies View is rendered.
func TestProgram_EventLoop_View(t *testing.T) {
	var buf bytes.Buffer

	m := TestModel{}
	p := New(m, WithOutput[TestModel](&buf))

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(100 * time.Millisecond)

	if err := p.Send(model.QuitMsg{}); err != nil {
		t.Fatal(err)
	}
	p.Stop()

	output := buf.String()

	// View should have been rendered at least once (initial render)
	if len(output) == 0 {
		t.Error("view should have been rendered")
	}

	// Should contain view format
	if !strings.Contains(output, "Value:") {
		t.Errorf("output should contain view content, got: %s", output)
	}
}

// TestProgram_EventLoop_Quit verifies QuitMsg stops the loop.
func TestProgram_EventLoop_Quit(t *testing.T) {
	var buf bytes.Buffer

	m := TestModel{}
	p := New(m, WithOutput[TestModel](&buf))

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	if !p.IsRunning() {
		t.Error("program should be running")
	}

	// Send QuitMsg
	if err := p.Send(model.QuitMsg{}); err != nil {
		t.Fatal(err)
	}

	// Should stop within reasonable time
	time.Sleep(200 * time.Millisecond)

	if p.IsRunning() {
		t.Error("program should have stopped after QuitMsg")
	}
}

// TestProgram_EventLoop_BatchMsg verifies BatchMsg expansion.
func TestProgram_EventLoop_BatchMsg(t *testing.T) {
	var buf bytes.Buffer

	m := TestModel{value: 0}
	p := New(m, WithOutput[TestModel](&buf))

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(100 * time.Millisecond)

	// Send BatchMsg with 3 increment messages
	batchMsg := model.BatchMsg{
		Messages: []model.Msg{
			model.KeyMsg{Type: model.KeyRune, Rune: '+'},
			model.KeyMsg{Type: model.KeyRune, Rune: '+'},
			model.KeyMsg{Type: model.KeyRune, Rune: '+'},
		},
	}

	if err := p.Send(batchMsg); err != nil {
		t.Fatal(err)
	}
	time.Sleep(150 * time.Millisecond)

	if err := p.Send(model.QuitMsg{}); err != nil {
		t.Fatal(err)
	}
	p.Stop()

	output := buf.String()

	// Should have processed all 3 messages
	if !strings.Contains(output, "Value: 3") {
		t.Errorf("expected Value: 3, got: %s", output)
	}
}

// TestProgram_EventLoop_SequenceMsg verifies SequenceMsg expansion.
func TestProgram_EventLoop_SequenceMsg(t *testing.T) {
	var buf bytes.Buffer

	m := TestModel{value: 0}
	p := New(m, WithOutput[TestModel](&buf))

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(100 * time.Millisecond)

	// Send SequenceMsg
	seqMsg := model.SequenceMsg{
		Messages: []model.Msg{
			model.KeyMsg{Type: model.KeyRune, Rune: '+'},
			model.KeyMsg{Type: model.KeyRune, Rune: '+'},
		},
	}

	if err := p.Send(seqMsg); err != nil {
		t.Fatal(err)
	}
	time.Sleep(150 * time.Millisecond)

	if err := p.Send(model.QuitMsg{}); err != nil {
		t.Fatal(err)
	}
	p.Stop()

	output := buf.String()

	// Should have processed both messages
	if !strings.Contains(output, "Value: 2") {
		t.Errorf("expected Value: 2, got: %s", output)
	}
}

// TestProgram_Send verifies Send method works correctly.
func TestProgram_Send(t *testing.T) {
	m := TestModel{}
	// Use empty input to prevent blocking on os.Stdin
	p := New(m, WithInput[TestModel](bytes.NewReader([]byte{})))

	// Send before start should error
	err := p.Send(model.KeyMsg{Type: model.KeyEnter})
	if err == nil {
		t.Error("Send should error when program not running")
	}

	// Start program
	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	// Give time to start
	time.Sleep(50 * time.Millisecond)

	// Send should work now
	err = p.Send(model.KeyMsg{Type: model.KeyEnter})
	if err != nil {
		t.Errorf("Send should work when running: %v", err)
	}

	p.Stop()

	// Send after stop should error
	err = p.Send(model.KeyMsg{Type: model.KeyEnter})
	if err == nil {
		t.Error("Send should error after stop")
	}
}

// TestProgram_EventLoop_Run verifies Run blocks until quit.
func TestProgram_EventLoop_Run(t *testing.T) {
	var buf bytes.Buffer

	m := TestModel{}
	p := New(m, WithOutput[TestModel](&buf))

	done := make(chan error)
	go func() {
		done <- p.Run() // Run in goroutine
	}()

	// Give it time to start and process init
	time.Sleep(100 * time.Millisecond)

	// Should still be running
	if !p.IsRunning() {
		t.Error("program should be running")
	}

	// Send quit
	if err := p.Send(model.QuitMsg{}); err != nil {
		t.Fatal(err)
	}

	// Should finish
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Run() returned error: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Run() did not finish after QuitMsg")
	}

	// Should no longer be running
	if p.IsRunning() {
		t.Error("program should not be running after quit")
	}

	// Check view was rendered
	output := buf.String()
	if !strings.Contains(output, "Value:") {
		t.Errorf("view should have been rendered, got: %s", output)
	}
}
