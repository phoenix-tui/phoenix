package tea_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/phoenix-tui/phoenix/tea"
)

// Simple test model
type TestModel struct {
	value int
}

func (m TestModel) Init() tea.Cmd {
	return nil
}

func (m TestModel) Update(msg tea.Msg) (TestModel, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "+" {
			m.value++
		}
		if keyMsg.String() == "q" {
			return m, tea.Quit()
		}
	}
	return m, nil
}

func (m TestModel) View() string {
	return "Value: " + string(rune('0'+m.value))
}

func TestAPI_New(t *testing.T) {
	m := TestModel{value: 0}
	p := tea.New(m)

	if p == nil {
		t.Error("New should return program")
	}
}

func TestAPI_Program_Lifecycle(t *testing.T) {
	var buf bytes.Buffer

	m := TestModel{value: 0}
	p := tea.New(m, tea.WithOutput[TestModel](&buf))

	// Start
	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	if !p.IsRunning() {
		t.Error("should be running")
	}

	// Send quit message instead of Stop() to avoid blocking
	if err := p.Send(tea.QuitMsg{}); err != nil {
		t.Errorf("Send quit failed: %v", err)
	}

	// Wait a bit for goroutine to stop
	time.Sleep(200 * time.Millisecond)

	if p.IsRunning() {
		// Force stop if still running
		p.Stop()
		t.Error("should have stopped after quit message")
	}
}

func TestAPI_Send(t *testing.T) {
	var buf bytes.Buffer

	m := TestModel{value: 0}
	p := tea.New(m, tea.WithOutput[TestModel](&buf))

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}
	defer p.Stop()

	// Send message
	err := p.Send(tea.KeyMsg{Type: tea.KeyRune, Rune: '+'})
	if err != nil {
		t.Errorf("Send failed: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	// Stop program to prevent race condition when reading buffer
	p.Stop()

	// Check output (safe now - no concurrent writes)
	output := buf.String()
	if !strings.Contains(output, "Value: 1") {
		t.Errorf("expected Value: 1 in output, got: %s", output)
	}
}

func TestAPI_Quit(t *testing.T) {
	var buf bytes.Buffer

	m := TestModel{}
	p := tea.New(m, tea.WithOutput[TestModel](&buf))

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	// Send quit
	if err := p.Send(tea.QuitMsg{}); err != nil {
		t.Errorf("Send quit failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if p.IsRunning() {
		t.Error("should have quit")
	}
}

func TestAPI_QuitCommand(t *testing.T) {
	var buf bytes.Buffer

	m := TestModel{}
	p := tea.New(m, tea.WithOutput[TestModel](&buf))

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	// Send 'q' key which triggers Quit() command
	if err := p.Send(tea.KeyMsg{Type: tea.KeyRune, Rune: 'q'}); err != nil {
		t.Errorf("Send 'q' failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if p.IsRunning() {
		t.Error("should have quit via Quit() command")
	}
}

func TestAPI_Batch(t *testing.T) {
	cmd := tea.Batch(
		tea.Println("test1"),
		tea.Println("test2"),
	)

	if cmd == nil {
		t.Error("Batch should return command")
	}
}

func TestAPI_Sequence(t *testing.T) {
	cmd := tea.Sequence(
		tea.Println("test1"),
		tea.Println("test2"),
	)

	if cmd == nil {
		t.Error("Sequence should return command")
	}
}

func TestAPI_Tick(t *testing.T) {
	cmd := tea.Tick(10 * time.Millisecond)

	if cmd == nil {
		t.Error("Tick should return command")
	}
}

func TestAPI_KeyTypes(_ *testing.T) {
	// Just verify constants are accessible
	_ = tea.KeyEnter
	_ = tea.KeyUp
	_ = tea.KeyDown
	_ = tea.KeyLeft
	_ = tea.KeyRight
	_ = tea.KeyF1
	_ = tea.KeyF2
	_ = tea.KeyF3
	_ = tea.KeyF4
	_ = tea.KeyF5
	_ = tea.KeyF6
	_ = tea.KeyF7
	_ = tea.KeyF8
	_ = tea.KeyF9
	_ = tea.KeyF10
	_ = tea.KeyF11
	_ = tea.KeyF12
	_ = tea.KeyCtrlC
	_ = tea.KeyBackspace
	_ = tea.KeyTab
	_ = tea.KeyEsc
	_ = tea.KeySpace
	_ = tea.KeyHome
	_ = tea.KeyEnd
	_ = tea.KeyPgUp
	_ = tea.KeyPgDown
	_ = tea.KeyDelete
	_ = tea.KeyInsert
}

func TestAPI_MouseTypes(_ *testing.T) {
	// Just verify constants are accessible
	_ = tea.MouseButtonNone
	_ = tea.MouseButtonLeft
	_ = tea.MouseButtonMiddle
	_ = tea.MouseButtonRight
	_ = tea.MouseButtonWheelUp
	_ = tea.MouseButtonWheelDown

	_ = tea.MouseActionPress
	_ = tea.MouseActionRelease
	_ = tea.MouseActionMotion
}

func TestAPI_Options(t *testing.T) {
	var buf bytes.Buffer

	m := TestModel{}
	p := tea.New(
		m,
		tea.WithInput[TestModel](strings.NewReader("test")),
		tea.WithOutput[TestModel](&buf),
		tea.WithAltScreen[TestModel](),
		tea.WithMouseAllMotion[TestModel](),
	)

	if p == nil {
		t.Error("should create program with options")
	}
}

func TestAPI_KeyMsg_String(t *testing.T) {
	tests := []struct {
		name     string
		key      tea.KeyMsg
		expected string
	}{
		{
			name:     "rune",
			key:      tea.KeyMsg{Type: tea.KeyRune, Rune: 'a'},
			expected: "a",
		},
		{
			name:     "enter",
			key:      tea.KeyMsg{Type: tea.KeyEnter},
			expected: "enter",
		},
		{
			name:     "up arrow",
			key:      tea.KeyMsg{Type: tea.KeyUp},
			expected: "↑",
		},
		{
			name:     "ctrl+c",
			key:      tea.KeyMsg{Type: tea.KeyCtrlC},
			expected: "ctrl+c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.key.String()
			if got != tt.expected {
				t.Errorf("KeyMsg.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestAPI_MouseMsg(t *testing.T) {
	msg := tea.MouseMsg{
		X:      10,
		Y:      20,
		Button: tea.MouseButtonLeft,
		Action: tea.MouseActionPress,
		Ctrl:   false,
		Alt:    false,
		Shift:  false,
	}

	if msg.X != 10 {
		t.Errorf("MouseMsg.X = %d, want 10", msg.X)
	}
	if msg.Y != 20 {
		t.Errorf("MouseMsg.Y = %d, want 20", msg.Y)
	}
	if msg.Button != tea.MouseButtonLeft {
		t.Errorf("MouseMsg.Button = %v, want MouseButtonLeft", msg.Button)
	}
	if msg.Action != tea.MouseActionPress {
		t.Errorf("MouseMsg.Action = %v, want MouseActionPress", msg.Action)
	}
}

func TestAPI_WindowSizeMsg(t *testing.T) {
	msg := tea.WindowSizeMsg{
		Width:  80,
		Height: 24,
	}

	if msg.Width != 80 {
		t.Errorf("WindowSizeMsg.Width = %d, want 80", msg.Width)
	}
	if msg.Height != 24 {
		t.Errorf("WindowSizeMsg.Height = %d, want 24", msg.Height)
	}
}

// Benchmark API overhead
func BenchmarkAPI_New(b *testing.B) {
	m := TestModel{}
	for i := 0; i < b.N; i++ {
		_ = tea.New(m)
	}
}

func BenchmarkAPI_Send(b *testing.B) {
	var buf bytes.Buffer
	m := TestModel{}
	p := tea.New(m, tea.WithOutput[TestModel](&buf))
	p.Start()
	defer p.Stop()

	msg := tea.KeyMsg{Type: tea.KeyRune, Rune: 'a'}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = p.Send(msg)
	}
}
