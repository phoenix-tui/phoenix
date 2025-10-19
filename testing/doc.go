// Package testing provides test helpers and mock implementations for Phoenix TUI framework.
//
// This package solves the common testing problem: production code needs terminal operations,
// but tests don't need (or want) actual terminal interaction.
//
// # NullTerminal - No-Op Implementation
//
// Use NullTerminal when your tests don't care about terminal operations:
//
//	func TestMyModel(t *testing.T) {
//	    m := &MyModel{
//	        terminal: testing.NewNullTerminal(),
//	    }
//	    // All terminal operations are no-ops
//	    m.Render() // Won't panic, won't write to terminal
//	}
//
// # MockTerminal - Recording Implementation
//
// Use MockTerminal when you need to verify terminal operations were called:
//
//	func TestRenderCallsClearLine(t *testing.T) {
//	    mock := testing.NewMockTerminal()
//	    m := &MyModel{terminal: mock}
//
//	    m.Render()
//
//	    // Verify terminal operations
//	    assert.Contains(t, mock.Calls, "ClearLine")
//	    assert.Equal(t, 1, mock.CallCount("ClearLine"))
//	}
//
// # Migration from Defensive Nil Checks
//
// Before (ugly nil checks in production code):
//
//	if m.terminal != nil {
//	    _ = m.terminal.ClearLine()
//	}
//
// After (clean code with NullTerminal in tests):
//
//	// Production code - clean!
//	_ = m.terminal.ClearLine()
//
//	// Test code - use NullTerminal
//	m := &MyModel{
//	    terminal: testing.NewNullTerminal(),
//	}
//
// # Performance
//
// Both NullTerminal and MockTerminal are extremely lightweight:
//   - NullTerminal: Zero overhead (all methods are empty)
//   - MockTerminal: Minimal overhead (string append for call tracking)
//
// Use them liberally in tests without performance concerns.
//
// # Thread Safety
//
// MockTerminal is safe for concurrent use. Multiple goroutines can call
// terminal methods, and all calls will be recorded correctly.
package testing
