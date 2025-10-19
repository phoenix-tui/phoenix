package api_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/phoenix-tui/phoenix/tea/api"
)

// Simple test model
type TestModel struct {
	value int
}

func (m TestModel) Init() api.Cmd {
	return nil
}

func (m TestModel) Update(msg api.Msg) (TestModel, api.Cmd) {
	if keyMsg, ok := msg.(api.KeyMsg); ok {
		if keyMsg.String() == "+" {
			m.value++
		}
		if keyMsg.String() == "q" {
			return m, api.Quit()
		}
	}
	return m, nil
}

func (m TestModel) View() string {
	return "Value: " + string(rune('0'+m.value))
}

func TestAPI_New(t *testing.T) {
	m := TestModel{value: 0}
	p := api.New(m)

	if p == nil {
		t.Error("New should return program")
	}
}

func TestAPI_Program_Lifecycle(t *testing.T) {
	var buf bytes.Buffer

	m := TestModel{value: 0}
	p := api.New(m, api.WithOutput[TestModel](&buf))

	// Start
	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	if !p.IsRunning() {
		t.Error("should be running")
	}

	// Send quit message instead of Stop() to avoid blocking
	if err := p.Send(api.QuitMsg{}); err != nil {
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
	p := api.New(m, api.WithOutput[TestModel](&buf))

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}
	defer p.Stop()

	// Send message
	err := p.Send(api.KeyMsg{Type: api.KeyRune, Rune: '+'})
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
	p := api.New(m, api.WithOutput[TestModel](&buf))

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	// Send quit
	if err := p.Send(api.QuitMsg{}); err != nil {
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
	p := api.New(m, api.WithOutput[TestModel](&buf))

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	// Send 'q' key which triggers Quit() command
	if err := p.Send(api.KeyMsg{Type: api.KeyRune, Rune: 'q'}); err != nil {
		t.Errorf("Send 'q' failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if p.IsRunning() {
		t.Error("should have quit via Quit() command")
	}
}

func TestAPI_Batch(t *testing.T) {
	cmd := api.Batch(
		api.Println("test1"),
		api.Println("test2"),
	)

	if cmd == nil {
		t.Error("Batch should return command")
	}
}

func TestAPI_Sequence(t *testing.T) {
	cmd := api.Sequence(
		api.Println("test1"),
		api.Println("test2"),
	)

	if cmd == nil {
		t.Error("Sequence should return command")
	}
}

func TestAPI_Tick(t *testing.T) {
	cmd := api.Tick(10 * time.Millisecond)

	if cmd == nil {
		t.Error("Tick should return command")
	}
}

func TestAPI_KeyTypes(_ *testing.T) {
	// Just verify constants are accessible
	_ = api.KeyEnter
	_ = api.KeyUp
	_ = api.KeyDown
	_ = api.KeyLeft
	_ = api.KeyRight
	_ = api.KeyF1
	_ = api.KeyF2
	_ = api.KeyF3
	_ = api.KeyF4
	_ = api.KeyF5
	_ = api.KeyF6
	_ = api.KeyF7
	_ = api.KeyF8
	_ = api.KeyF9
	_ = api.KeyF10
	_ = api.KeyF11
	_ = api.KeyF12
	_ = api.KeyCtrlC
	_ = api.KeyBackspace
	_ = api.KeyTab
	_ = api.KeyEsc
	_ = api.KeySpace
	_ = api.KeyHome
	_ = api.KeyEnd
	_ = api.KeyPgUp
	_ = api.KeyPgDown
	_ = api.KeyDelete
	_ = api.KeyInsert
}

func TestAPI_MouseTypes(_ *testing.T) {
	// Just verify constants are accessible
	_ = api.MouseButtonNone
	_ = api.MouseButtonLeft
	_ = api.MouseButtonMiddle
	_ = api.MouseButtonRight
	_ = api.MouseButtonWheelUp
	_ = api.MouseButtonWheelDown

	_ = api.MouseActionPress
	_ = api.MouseActionRelease
	_ = api.MouseActionMotion
}

func TestAPI_Options(t *testing.T) {
	var buf bytes.Buffer

	m := TestModel{}
	p := api.New(
		m,
		api.WithInput[TestModel](strings.NewReader("test")),
		api.WithOutput[TestModel](&buf),
		api.WithAltScreen[TestModel](),
		api.WithMouseAllMotion[TestModel](),
	)

	if p == nil {
		t.Error("should create program with options")
	}
}

func TestAPI_KeyMsg_String(t *testing.T) {
	tests := []struct {
		name     string
		key      api.KeyMsg
		expected string
	}{
		{
			name:     "rune",
			key:      api.KeyMsg{Type: api.KeyRune, Rune: 'a'},
			expected: "a",
		},
		{
			name:     "enter",
			key:      api.KeyMsg{Type: api.KeyEnter},
			expected: "enter",
		},
		{
			name:     "up arrow",
			key:      api.KeyMsg{Type: api.KeyUp},
			expected: "↑",
		},
		{
			name:     "ctrl+c",
			key:      api.KeyMsg{Type: api.KeyCtrlC},
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
	msg := api.MouseMsg{
		X:      10,
		Y:      20,
		Button: api.MouseButtonLeft,
		Action: api.MouseActionPress,
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
	if msg.Button != api.MouseButtonLeft {
		t.Errorf("MouseMsg.Button = %v, want MouseButtonLeft", msg.Button)
	}
	if msg.Action != api.MouseActionPress {
		t.Errorf("MouseMsg.Action = %v, want MouseActionPress", msg.Action)
	}
}

func TestAPI_WindowSizeMsg(t *testing.T) {
	msg := api.WindowSizeMsg{
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
		_ = api.New(m)
	}
}

func BenchmarkAPI_Send(b *testing.B) {
	var buf bytes.Buffer
	m := TestModel{}
	p := api.New(m, api.WithOutput[TestModel](&buf))
	p.Start()
	defer p.Stop()

	msg := api.KeyMsg{Type: api.KeyRune, Rune: 'a'}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = p.Send(msg)
	}
}
