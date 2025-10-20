package model

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

// testCounterMsg increments a counter (for testing command execution).
type testCounterMsg struct {
	id    int
	value int
}

func (t testCounterMsg) String() string {
	return fmt.Sprintf("counter(%d)=%d", t.id, t.value)
}

// testDelayMsg is sent after a delay (for testing timing).
//
//nolint:unused // Used in deferred test scenarios
type testDelayMsg struct {
	duration time.Duration
}

//nolint:unused // Used in deferred test scenarios
func (t testDelayMsg) String() string {
	return fmt.Sprintf("delay(%v)", t.duration)
}

// TestCmd_Execution verifies a simple command executes correctly.
func TestCmd_Execution(t *testing.T) {
	expectedValue := 42

	cmd := func() Msg {
		return testCounterMsg{id: 1, value: expectedValue}
	}

	msg := cmd()

	counterMsg, ok := msg.(testCounterMsg)
	if !ok {
		t.Fatalf("cmd sent %T, expected testCounterMsg", msg)
	}
	if counterMsg.value != expectedValue {
		t.Errorf("value = %d, want %d", counterMsg.value, expectedValue)
	}
}

// TestBatch_Nil verifies Batch(nil) returns nil.
func TestBatch_Nil(t *testing.T) {
	cmd := Batch(nil)

	if cmd != nil {
		t.Errorf("Batch(nil) = %v, want nil", cmd)
	}
}

// TestBatch_Empty verifies Batch() with no args returns nil.
func TestBatch_Empty(t *testing.T) {
	cmd := Batch()

	if cmd != nil {
		t.Errorf("Batch() = %v, want nil", cmd)
	}
}

// TestBatch_NilFiltering verifies Batch filters out nil commands.
func TestBatch_NilFiltering(t *testing.T) {
	cmd1 := func() Msg { return testCounterMsg{id: 1, value: 1} }
	cmd2 := func() Msg { return testCounterMsg{id: 2, value: 2} }

	// Mix nil commands with real ones
	cmd := Batch(nil, cmd1, nil, cmd2, nil)

	if cmd == nil {
		t.Fatal("Batch(nil, cmd1, nil, cmd2, nil) returned nil")
	}

	msg := cmd()

	batchMsg, ok := msg.(BatchMsg)
	if !ok {
		t.Fatalf("Batch sent %T, expected BatchMsg", msg)
	}

	if len(batchMsg.Messages) != 2 {
		t.Errorf("BatchMsg has %d messages, want 2", len(batchMsg.Messages))
	}
}

// TestBatch_SingleCommand verifies Batch optimizes single command.
func TestBatch_SingleCommand(t *testing.T) {
	expectedValue := 99
	originalCmd := func() Msg {
		return testCounterMsg{id: 1, value: expectedValue}
	}

	// Batch should return the command directly (optimization)
	cmd := Batch(originalCmd)

	if cmd == nil {
		t.Fatal("Batch(singleCmd) returned nil")
	}

	msg := cmd()

	// Should receive the message directly, NOT BatchMsg
	counterMsg, ok := msg.(testCounterMsg)
	if !ok {
		t.Fatalf("Batch(singleCmd) sent %T, expected testCounterMsg (not BatchMsg)", msg)
	}
	if counterMsg.value != expectedValue {
		t.Errorf("value = %d, want %d", counterMsg.value, expectedValue)
	}
}

// TestBatch_MultipleCommands verifies Batch executes commands in parallel.
func TestBatch_MultipleCommands(t *testing.T) {
	cmd1 := func() Msg { return testCounterMsg{id: 1, value: 10} }
	cmd2 := func() Msg { return testCounterMsg{id: 2, value: 20} }
	cmd3 := func() Msg { return testCounterMsg{id: 3, value: 30} }

	cmd := Batch(cmd1, cmd2, cmd3)

	if cmd == nil {
		t.Fatal("Batch(cmd1, cmd2, cmd3) returned nil")
	}

	msg := cmd()

	batchMsg, ok := msg.(BatchMsg)
	if !ok {
		t.Fatalf("Batch sent %T, expected BatchMsg", msg)
	}

	if len(batchMsg.Messages) != 3 {
		t.Fatalf("BatchMsg has %d messages, want 3", len(batchMsg.Messages))
	}

	// Verify all messages are present (order is undefined for parallel execution)
	ids := make(map[int]bool)
	values := make(map[int]int)

	for _, m := range batchMsg.Messages {
		counterMsg, ok := m.(testCounterMsg)
		if !ok {
			t.Errorf("message is %T, expected testCounterMsg", m)
			continue
		}
		ids[counterMsg.id] = true
		values[counterMsg.id] = counterMsg.value
	}

	// Check all IDs present
	for i := 1; i <= 3; i++ {
		if !ids[i] {
			t.Errorf("message with id=%d not found", i)
		}
	}

	// Check values
	expectedValues := map[int]int{1: 10, 2: 20, 3: 30}
	for id, expectedValue := range expectedValues {
		if values[id] != expectedValue {
			t.Errorf("id=%d: value=%d, want %d", id, values[id], expectedValue)
		}
	}
}

// TestBatch_Parallelism verifies Batch actually runs commands in parallel.
func TestBatch_Parallelism(t *testing.T) {
	delayDuration := 50 * time.Millisecond

	// Create 3 commands that each sleep for 50ms
	cmd1 := func() Msg {
		time.Sleep(delayDuration)
		return testCounterMsg{id: 1, value: 1}
	}
	cmd2 := func() Msg {
		time.Sleep(delayDuration)
		return testCounterMsg{id: 2, value: 2}
	}
	cmd3 := func() Msg {
		time.Sleep(delayDuration)
		return testCounterMsg{id: 3, value: 3}
	}

	start := time.Now()
	cmd := Batch(cmd1, cmd2, cmd3)
	msg := cmd()
	elapsed := time.Since(start)

	// If parallel: ~50ms (all run at once)
	// If sequential: ~150ms (50ms * 3)
	// Allow some margin for goroutine scheduling
	maxExpected := delayDuration + 30*time.Millisecond

	if elapsed > maxExpected {
		t.Errorf("Batch took %v, expected ~%v (parallel execution)", elapsed, delayDuration)
		t.Errorf("Commands may have run sequentially instead of in parallel")
	}

	batchMsg, ok := msg.(BatchMsg)
	if !ok {
		t.Fatalf("Batch sent %T, expected BatchMsg", msg)
	}
	if len(batchMsg.Messages) != 3 {
		t.Errorf("BatchMsg has %d messages, want 3", len(batchMsg.Messages))
	}
}

// TestBatchMsg_String verifies BatchMsg.String() format.
func TestBatchMsg_String(t *testing.T) {
	tests := []struct {
		name     string
		messages []Msg
		expected string
	}{
		{
			name:     "empty",
			messages: []Msg{},
			expected: "batch (0 messages)",
		},
		{
			name:     "one message",
			messages: []Msg{testCounterMsg{id: 1, value: 10}},
			expected: "batch (1 messages)",
		},
		{
			name: "three messages",
			messages: []Msg{
				testCounterMsg{id: 1, value: 10},
				testCounterMsg{id: 2, value: 20},
				testCounterMsg{id: 3, value: 30},
			},
			expected: "batch (3 messages)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := BatchMsg{Messages: tt.messages}

			str := msg.String()

			if str != tt.expected {
				t.Errorf("String() = %q, want %q", str, tt.expected)
			}
		})
	}
}

// TestSequence_Nil verifies Sequence(nil) returns nil.
func TestSequence_Nil(t *testing.T) {
	cmd := Sequence(nil)

	if cmd != nil {
		t.Errorf("Sequence(nil) = %v, want nil", cmd)
	}
}

// TestSequence_Empty verifies Sequence() with no args returns nil.
func TestSequence_Empty(t *testing.T) {
	cmd := Sequence()

	if cmd != nil {
		t.Errorf("Sequence() = %v, want nil", cmd)
	}
}

// TestSequence_NilFiltering verifies Sequence filters out nil commands.
func TestSequence_NilFiltering(t *testing.T) {
	cmd1 := func() Msg { return testCounterMsg{id: 1, value: 1} }
	cmd2 := func() Msg { return testCounterMsg{id: 2, value: 2} }

	cmd := Sequence(nil, cmd1, nil, cmd2, nil)

	if cmd == nil {
		t.Fatal("Sequence(nil, cmd1, nil, cmd2, nil) returned nil")
	}

	msg := cmd()

	seqMsg, ok := msg.(SequenceMsg)
	if !ok {
		t.Fatalf("Sequence sent %T, expected SequenceMsg", msg)
	}

	if len(seqMsg.Messages) != 2 {
		t.Errorf("SequenceMsg has %d messages, want 2", len(seqMsg.Messages))
	}
}

// TestSequence_SingleCommand verifies Sequence optimizes single command.
func TestSequence_SingleCommand(t *testing.T) {
	expectedValue := 77
	originalCmd := func() Msg {
		return testCounterMsg{id: 1, value: expectedValue}
	}

	// Sequence should return the command directly (optimization)
	cmd := Sequence(originalCmd)

	if cmd == nil {
		t.Fatal("Sequence(singleCmd) returned nil")
	}

	msg := cmd()

	// Should receive the message directly, NOT SequenceMsg
	counterMsg, ok := msg.(testCounterMsg)
	if !ok {
		t.Fatalf("Sequence(singleCmd) sent %T, expected testCounterMsg (not SequenceMsg)", msg)
	}
	if counterMsg.value != expectedValue {
		t.Errorf("value = %d, want %d", counterMsg.value, expectedValue)
	}
}

// TestSequence_Order verifies Sequence executes commands in order.
func TestSequence_Order(t *testing.T) {
	// Use atomic counter to track execution order
	var executionOrder int32

	cmd1 := func() Msg {
		order := atomic.AddInt32(&executionOrder, 1)
		return testCounterMsg{id: 1, value: int(order)}
	}
	cmd2 := func() Msg {
		order := atomic.AddInt32(&executionOrder, 1)
		return testCounterMsg{id: 2, value: int(order)}
	}
	cmd3 := func() Msg {
		order := atomic.AddInt32(&executionOrder, 1)
		return testCounterMsg{id: 3, value: int(order)}
	}

	cmd := Sequence(cmd1, cmd2, cmd3)

	if cmd == nil {
		t.Fatal("Sequence(cmd1, cmd2, cmd3) returned nil")
	}

	msg := cmd()

	seqMsg, ok := msg.(SequenceMsg)
	if !ok {
		t.Fatalf("Sequence sent %T, expected SequenceMsg", msg)
	}

	if len(seqMsg.Messages) != 3 {
		t.Fatalf("SequenceMsg has %d messages, want 3", len(seqMsg.Messages))
	}

	// Verify execution order (should be 1, 2, 3)
	for i, m := range seqMsg.Messages {
		counterMsg, ok := m.(testCounterMsg)
		if !ok {
			t.Errorf("message[%d] is %T, expected testCounterMsg", i, m)
			continue
		}

		expectedID := i + 1
		expectedOrder := i + 1

		if counterMsg.id != expectedID {
			t.Errorf("message[%d]: id=%d, want %d", i, counterMsg.id, expectedID)
		}
		if counterMsg.value != expectedOrder {
			t.Errorf("message[%d]: execution order=%d, want %d", i, counterMsg.value, expectedOrder)
		}
	}
}

// TestSequence_Sequential verifies Sequence runs commands sequentially (not parallel).
func TestSequence_Sequential(t *testing.T) {
	delayDuration := 30 * time.Millisecond

	// Create 3 commands that each sleep for 30ms
	cmd1 := func() Msg {
		time.Sleep(delayDuration)
		return testCounterMsg{id: 1, value: 1}
	}
	cmd2 := func() Msg {
		time.Sleep(delayDuration)
		return testCounterMsg{id: 2, value: 2}
	}
	cmd3 := func() Msg {
		time.Sleep(delayDuration)
		return testCounterMsg{id: 3, value: 3}
	}

	start := time.Now()
	cmd := Sequence(cmd1, cmd2, cmd3)
	msg := cmd()
	elapsed := time.Since(start)

	// If sequential: ~90ms (30ms * 3)
	// If parallel: ~30ms (all run at once)
	minExpected := delayDuration * 3
	maxExpected := minExpected + 30*time.Millisecond // Allow margin

	if elapsed < minExpected {
		t.Errorf("Sequence took %v, expected at least %v (sequential execution)", elapsed, minExpected)
		t.Errorf("Commands may have run in parallel instead of sequentially")
	}
	if elapsed > maxExpected {
		t.Logf("Warning: Sequence took %v (expected ~%v), may be slow system", elapsed, minExpected)
	}

	seqMsg, ok := msg.(SequenceMsg)
	if !ok {
		t.Fatalf("Sequence sent %T, expected SequenceMsg", msg)
	}
	if len(seqMsg.Messages) != 3 {
		t.Errorf("SequenceMsg has %d messages, want 3", len(seqMsg.Messages))
	}
}

// TestSequence_OrderWithDelays verifies order is preserved even with different delays.
func TestSequence_OrderWithDelays(t *testing.T) {
	// cmd2 takes longer than cmd3, but should still execute second
	cmd1 := func() Msg {
		time.Sleep(10 * time.Millisecond)
		return testCounterMsg{id: 1, value: 1}
	}
	cmd2 := func() Msg {
		time.Sleep(50 * time.Millisecond) // Longest
		return testCounterMsg{id: 2, value: 2}
	}
	cmd3 := func() Msg {
		time.Sleep(10 * time.Millisecond)
		return testCounterMsg{id: 3, value: 3}
	}

	cmd := Sequence(cmd1, cmd2, cmd3)
	msg := cmd()

	seqMsg, ok := msg.(SequenceMsg)
	if !ok {
		t.Fatalf("Sequence sent %T, expected SequenceMsg", msg)
	}

	// Verify order is 1, 2, 3 (not based on duration)
	expectedOrder := []int{1, 2, 3}
	for i, m := range seqMsg.Messages {
		counterMsg, ok := m.(testCounterMsg)
		if !ok {
			t.Errorf("message[%d] is %T, expected testCounterMsg", i, m)
			continue
		}
		if counterMsg.id != expectedOrder[i] {
			t.Errorf("message[%d]: id=%d, want %d", i, counterMsg.id, expectedOrder[i])
		}
	}
}

// TestSequenceMsg_String verifies SequenceMsg.String() format.
func TestSequenceMsg_String(t *testing.T) {
	tests := []struct {
		name     string
		messages []Msg
		expected string
	}{
		{
			name:     "empty",
			messages: []Msg{},
			expected: "sequence (0 messages)",
		},
		{
			name:     "one message",
			messages: []Msg{testCounterMsg{id: 1, value: 10}},
			expected: "sequence (1 messages)",
		},
		{
			name: "five messages",
			messages: []Msg{
				testCounterMsg{id: 1, value: 10},
				testCounterMsg{id: 2, value: 20},
				testCounterMsg{id: 3, value: 30},
				testCounterMsg{id: 4, value: 40},
				testCounterMsg{id: 5, value: 50},
			},
			expected: "sequence (5 messages)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := SequenceMsg{Messages: tt.messages}

			str := msg.String()

			if str != tt.expected {
				t.Errorf("String() = %q, want %q", str, tt.expected)
			}
		})
	}
}

// TestBatch_vs_Sequence_Comparison demonstrates the difference.
func TestBatch_vs_Sequence_Comparison(t *testing.T) {
	delayDuration := 40 * time.Millisecond

	cmd1 := func() Msg {
		time.Sleep(delayDuration)
		return testCounterMsg{id: 1, value: 1}
	}
	cmd2 := func() Msg {
		time.Sleep(delayDuration)
		return testCounterMsg{id: 2, value: 2}
	}

	// Test Batch (parallel)
	startBatch := time.Now()
	batchCmd := Batch(cmd1, cmd2)
	_ = batchCmd()
	batchElapsed := time.Since(startBatch)

	// Test Sequence (sequential)
	startSeq := time.Now()
	seqCmd := Sequence(cmd1, cmd2)
	_ = seqCmd()
	seqElapsed := time.Since(startSeq)

	t.Logf("Batch (parallel) took: %v", batchElapsed)
	t.Logf("Sequence (sequential) took: %v", seqElapsed)

	// Sequence should take approximately 2x as long as Batch
	// (with some margin for goroutine overhead)
	if seqElapsed < batchElapsed {
		t.Error("Sequence should take longer than Batch")
	}

	// Batch should be close to single delay (parallel)
	if batchElapsed > delayDuration+30*time.Millisecond {
		t.Errorf("Batch took %v, expected ~%v (parallel)", batchElapsed, delayDuration)
	}

	// Sequence should be close to double delay (sequential)
	expectedSeq := delayDuration * 2
	if seqElapsed < expectedSeq-10*time.Millisecond || seqElapsed > expectedSeq+40*time.Millisecond {
		t.Logf("Warning: Sequence took %v, expected ~%v", seqElapsed, expectedSeq)
	}
}
