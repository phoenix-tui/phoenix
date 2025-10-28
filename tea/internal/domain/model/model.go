package model

// Model represents the Elm Architecture contract for Phoenix TUI applications.
//
// The Model interface defines the three core functions of the MVU (Model-View-Update)
// pattern: initialization, message handling, and rendering.
//
// Type parameter T must be the concrete model type implementing this interface.
// This enables type-safe model transformations through the Update cycle.
//
// Example:
//
//	type CounterModel struct {
//		count int
//	}
//
//	func (m CounterModel) Init() Cmd {
//		return nil // No initial command
//	}
//
//	func (m CounterModel) Update(msg Msg) (Model[CounterModel], Cmd) {
//		switch msg := msg.(type) {
//		case KeyMsg:
//			if msg.String() == "+" {
//				m.count++
//			} else if msg.String() == "-" {
//				m.count--
//			}
//		}
//		return m, nil
//	}
//
//	func (m CounterModel) View() string {
//		return fmt.Sprintf("Count: %d\nPress + or - to change", m.count)
//	}
//
// Design Decisions:
//   - Generic type parameter T ensures Update returns the correct model type
//   - Init/Update return optional Cmd for asynchronous operations
//   - View returns string (renderer handles ANSI conversion)
//   - Model should be immutable - Update returns new instance
type Model[T any] interface {
	// Init returns an optional initial command to run when the model is created.
	//
	// This is called once at startup. Return nil if no initialization is needed,
	// or return a command to perform async operations (e.g., loading data).
	//
	// Example:
	//   func (m AppModel) Init() Cmd {
	//       return Tick(time.Second) // Start a timer
	//   }
	Init() Cmd

	// Update handles incoming messages and returns an updated model and optional command.
	//
	// This is the heart of the Elm Architecture. When a message arrives:
	//   1. Inspect the message type (use type switch)
	//   2. Update model state based on message
	//   3. Return new model instance (immutability!)
	//   4. Optionally return a command for async operations
	//
	// Example:
	//   func (m CounterModel) Update(msg Msg) (Model[CounterModel], Cmd) {
	//       switch msg := msg.(type) {
	//       case KeyMsg:
	//           if msg.String() == "q" {
	//               return m, Quit()
	//           }
	//       case TickMsg:
	//           m.count++
	//           return m, Tick(time.Second)
	//       }
	//       return m, nil
	//   }
	Update(msg Msg) (Model[T], Cmd)

	// View renders the current model state to a string.
	//
	// This should be a pure function - no side effects, no I/O.
	// The string can include ANSI escape codes for styling.
	//
	// Example:
	//   func (m CounterModel) View() string {
	//       return fmt.Sprintf("Count: %d\n\nPress q to quit", m.count)
	//   }
	View() string
}
