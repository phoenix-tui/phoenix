package model

import "fmt"

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

// Batch executes multiple commands concurrently and collects their results.
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

// Sequence executes commands sequentially (one after another) and collects their results.
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
