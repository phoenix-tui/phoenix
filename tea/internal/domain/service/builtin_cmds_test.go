package service

import (
	"strings"
	"testing"
	"time"

	model2 "github.com/phoenix-tui/phoenix/tea/internal/domain/model"
)

// TestQuit verifies Quit command sends QuitMsg.
func TestQuit(t *testing.T) {
	cmd := Quit()

	if cmd == nil {
		t.Fatal("Quit() returned nil")
	}

	msg := cmd()

	quitMsg, ok := msg.(model2.QuitMsg)
	if !ok {
		t.Fatalf("Quit command sent %T, expected model.QuitMsg", msg)
	}

	// Verify QuitMsg string representation
	expected := "quit"
	if quitMsg.String() != expected {
		t.Errorf("QuitMsg.String() = %q, want %q", quitMsg.String(), expected)
	}
}

// TestPrintln verifies Println command sends PrintlnMsg.
func TestPrintln(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "simple message",
			message: "test",
		},
		{
			name:    "empty message",
			message: "",
		},
		{
			name:    "multiline message",
			message: "line1\nline2\nline3",
		},
		{
			name:    "message with special characters",
			message: "Hello, 世界! 🚀",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := Println(tt.message)

			if cmd == nil {
				t.Fatal("Println() returned nil")
			}

			msg := cmd()

			printlnMsg, ok := msg.(PrintlnMsg)
			if !ok {
				t.Fatalf("Println command sent %T, expected PrintlnMsg", msg)
			}

			if printlnMsg.Message != tt.message {
				t.Errorf("PrintlnMsg.Message = %q, want %q", printlnMsg.Message, tt.message)
			}
		})
	}
}

// TestPrintlnMsg_String verifies PrintlnMsg.String() format.
func TestPrintlnMsg_String(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "simple message",
			message:  "hello",
			expected: "println: hello",
		},
		{
			name:     "empty message",
			message:  "",
			expected: "println: ",
		},
		{
			name:     "message with spaces",
			message:  "hello world",
			expected: "println: hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := PrintlnMsg{Message: tt.message}

			str := msg.String()

			if str != tt.expected {
				t.Errorf("String() = %q, want %q", str, tt.expected)
			}
		})
	}
}

// TestTick verifies Tick command waits and sends TickMsg.
func TestTick(t *testing.T) {
	duration := 50 * time.Millisecond

	start := time.Now()
	cmd := Tick(duration)

	if cmd == nil {
		t.Fatal("Tick() returned nil")
	}

	msg := cmd()
	elapsed := time.Since(start)

	// Verify it waited approximately the right duration
	minExpected := duration
	maxExpected := duration + 30*time.Millisecond // Allow margin

	if elapsed < minExpected {
		t.Errorf("Tick waited %v, expected at least %v", elapsed, minExpected)
	}
	if elapsed > maxExpected {
		t.Logf("Warning: Tick waited %v (expected ~%v), may be slow system", elapsed, duration)
	}

	// Verify TickMsg
	tickMsg, ok := msg.(TickMsg)
	if !ok {
		t.Fatalf("Tick command sent %T, expected TickMsg", msg)
	}

	// Verify Time is recent (within last second)
	now := time.Now()
	if tickMsg.Time.After(now) {
		t.Errorf("TickMsg.Time is in the future: %v > %v", tickMsg.Time, now)
	}
	if now.Sub(tickMsg.Time) > time.Second {
		t.Errorf("TickMsg.Time is too old: %v (now: %v)", tickMsg.Time, now)
	}
}

// TestTick_ZeroDuration verifies Tick with zero duration.
func TestTick_ZeroDuration(t *testing.T) {
	start := time.Now()
	cmd := Tick(0)

	if cmd == nil {
		t.Fatal("Tick(0) returned nil")
	}

	msg := cmd()
	elapsed := time.Since(start)

	// Should return almost immediately
	if elapsed > 10*time.Millisecond {
		t.Errorf("Tick(0) took %v, expected near-instant", elapsed)
	}

	_, ok := msg.(TickMsg)
	if !ok {
		t.Fatalf("Tick command sent %T, expected TickMsg", msg)
	}
}

// TestTick_MultipleTicks verifies multiple Tick commands work independently.
func TestTick_MultipleTicks(t *testing.T) {
	cmd1 := Tick(20 * time.Millisecond)
	cmd2 := Tick(40 * time.Millisecond)

	// Execute cmd1
	start1 := time.Now()
	msg1 := cmd1()
	elapsed1 := time.Since(start1)

	// Execute cmd2
	start2 := time.Now()
	msg2 := cmd2()
	elapsed2 := time.Since(start2)

	// Verify cmd1 timing
	if elapsed1 < 20*time.Millisecond || elapsed1 > 50*time.Millisecond {
		t.Errorf("cmd1 took %v, expected ~20ms", elapsed1)
	}

	// Verify cmd2 timing
	if elapsed2 < 40*time.Millisecond || elapsed2 > 70*time.Millisecond {
		t.Errorf("cmd2 took %v, expected ~40ms", elapsed2)
	}

	// Verify both sent TickMsg
	tickMsg1, ok := msg1.(TickMsg)
	if !ok {
		t.Errorf("cmd1 sent %T, expected TickMsg", msg1)
	}
	tickMsg2, ok := msg2.(TickMsg)
	if !ok {
		t.Errorf("cmd2 sent %T, expected TickMsg", msg2)
	}

	// cmd2 should have later time than cmd1
	if !tickMsg2.Time.After(tickMsg1.Time) {
		t.Errorf("tickMsg2.Time (%v) should be after tickMsg1.Time (%v)", tickMsg2.Time, tickMsg1.Time)
	}
}

// TestTickMsg_String verifies TickMsg.String() format.
func TestTickMsg_String(t *testing.T) {
	testTime := time.Date(2025, 10, 16, 12, 30, 45, 0, time.UTC)
	msg := TickMsg{Time: testTime}

	str := msg.String()

	// Should contain "tick at" and RFC3339 formatted time
	if !strings.HasPrefix(str, "tick at ") {
		t.Errorf("String() = %q, expected to start with 'tick at '", str)
	}

	expectedTime := testTime.Format(time.RFC3339)
	if !strings.Contains(str, expectedTime) {
		t.Errorf("String() = %q, expected to contain %q", str, expectedTime)
	}

	expected := "tick at " + expectedTime
	if str != expected {
		t.Errorf("String() = %q, want %q", str, expected)
	}
}

// TestTick_Concurrent verifies Tick commands can run concurrently.
func TestTick_Concurrent(t *testing.T) {
	duration := 50 * time.Millisecond

	cmd1 := Tick(duration)
	cmd2 := Tick(duration)
	cmd3 := Tick(duration)

	// Run all three concurrently
	start := time.Now()

	results := make(chan time.Time, 3)
	go func() {
		msg := cmd1()
		tickMsg := msg.(TickMsg)
		results <- tickMsg.Time
	}()
	go func() {
		msg := cmd2()
		tickMsg := msg.(TickMsg)
		results <- tickMsg.Time
	}()
	go func() {
		msg := cmd3()
		tickMsg := msg.(TickMsg)
		results <- tickMsg.Time
	}()

	// Collect results
	var times []time.Time
	for i := 0; i < 3; i++ {
		times = append(times, <-results)
	}

	elapsed := time.Since(start)

	// All three should complete in ~duration (not 3x duration)
	maxExpected := duration + 30*time.Millisecond
	if elapsed > maxExpected {
		t.Errorf("Concurrent ticks took %v, expected ~%v", elapsed, duration)
	}

	// All times should be within a few milliseconds of each other
	// Note: Allow larger margin on Windows due to timer resolution
	maxDiff := 30 * time.Millisecond
	for i := 1; i < len(times); i++ {
		diff := times[i].Sub(times[0])
		if diff < 0 {
			diff = -diff
		}
		if diff > maxDiff {
			t.Errorf("Tick times differ by %v, expected near-simultaneous (max %v)", diff, maxDiff)
		}
	}
}

// TestBuiltinCmds_Integration verifies commands work together.
func TestBuiltinCmds_Integration(t *testing.T) {
	// Simulate a simple Update flow using builtin commands

	type testModel struct {
		count int
		done  bool
	}

	m := testModel{count: 0}

	// 1. Start with Tick
	tickCmd := Tick(10 * time.Millisecond)
	tickMsg := tickCmd()

	if _, ok := tickMsg.(TickMsg); !ok {
		t.Fatalf("expected TickMsg, got %T", tickMsg)
	}

	// 2. Increment counter
	m.count++

	// 3. Log the increment
	printCmd := Println("incremented")
	printMsg := printCmd()

	printlnMsg, ok := printMsg.(PrintlnMsg)
	if !ok {
		t.Fatalf("expected PrintlnMsg, got %T", printMsg)
	}
	if printlnMsg.Message != "incremented" {
		t.Errorf("message = %q, want %q", printlnMsg.Message, "incremented")
	}

	// 4. Quit after 3 ticks
	if m.count >= 3 {
		quitCmd := Quit()
		quitMsg := quitCmd()

		if _, ok := quitMsg.(model2.QuitMsg); !ok {
			t.Fatalf("expected QuitMsg, got %T", quitMsg)
		}
		m.done = true
	}

	if m.count != 1 {
		t.Errorf("count = %d, want 1", m.count)
	}
	if m.done {
		t.Error("done = true, want false (count < 3)")
	}
}

// TestBuiltinCmds_WithBatch verifies builtin commands work with Batch.
func TestBuiltinCmds_WithBatch(t *testing.T) {
	cmd := model2.Batch(
		Println("first"),
		Println("second"),
		Println("third"),
	)

	if cmd == nil {
		t.Fatal("Batch returned nil")
	}

	msg := cmd()

	batchMsg, ok := msg.(model2.BatchMsg)
	if !ok {
		t.Fatalf("expected BatchMsg, got %T", msg)
	}

	if len(batchMsg.Messages) != 3 {
		t.Fatalf("expected 3 messages, got %d", len(batchMsg.Messages))
	}

	// All should be PrintlnMsg
	messages := make(map[string]bool)
	for _, m := range batchMsg.Messages {
		printMsg, ok := m.(PrintlnMsg)
		if !ok {
			t.Errorf("expected PrintlnMsg, got %T", m)
			continue
		}
		messages[printMsg.Message] = true
	}

	// Check all messages present (order undefined for Batch)
	for _, expected := range []string{"first", "second", "third"} {
		if !messages[expected] {
			t.Errorf("message %q not found", expected)
		}
	}
}

// TestBuiltinCmds_WithSequence verifies builtin commands work with Sequence.
func TestBuiltinCmds_WithSequence(t *testing.T) {
	cmd := model2.Sequence(
		Println("first"),
		Tick(10*time.Millisecond),
		Println("second"),
	)

	if cmd == nil {
		t.Fatal("Sequence returned nil")
	}

	start := time.Now()
	msg := cmd()
	elapsed := time.Since(start)

	// Should wait at least 10ms for Tick
	if elapsed < 10*time.Millisecond {
		t.Errorf("Sequence took %v, expected at least 10ms", elapsed)
	}

	seqMsg, ok := msg.(model2.SequenceMsg)
	if !ok {
		t.Fatalf("expected SequenceMsg, got %T", msg)
	}

	if len(seqMsg.Messages) != 3 {
		t.Fatalf("expected 3 messages, got %d", len(seqMsg.Messages))
	}

	// Verify order
	if _, ok := seqMsg.Messages[0].(PrintlnMsg); !ok {
		t.Errorf("message[0] is %T, expected PrintlnMsg", seqMsg.Messages[0])
	}
	if _, ok := seqMsg.Messages[1].(TickMsg); !ok {
		t.Errorf("message[1] is %T, expected TickMsg", seqMsg.Messages[1])
	}
	if _, ok := seqMsg.Messages[2].(PrintlnMsg); !ok {
		t.Errorf("message[2] is %T, expected PrintlnMsg", seqMsg.Messages[2])
	}

	// Verify content
	firstMsg := seqMsg.Messages[0].(PrintlnMsg)
	if firstMsg.Message != "first" {
		t.Errorf("message[0] = %q, want %q", firstMsg.Message, "first")
	}

	secondMsg := seqMsg.Messages[2].(PrintlnMsg)
	if secondMsg.Message != "second" {
		t.Errorf("message[2] = %q, want %q", secondMsg.Message, "second")
	}
}
