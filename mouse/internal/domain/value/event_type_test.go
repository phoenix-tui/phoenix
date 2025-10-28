package value

import "testing"

func TestEventType_String(t *testing.T) {
	tests := []struct {
		eventType EventType
		expected  string
	}{
		{EventPress, "Press"},
		{EventRelease, "Release"},
		{EventClick, "Click"},
		{EventDoubleClick, "DoubleClick"},
		{EventTripleClick, "TripleClick"},
		{EventDrag, "Drag"},
		{EventMotion, "Motion"},
		{EventScroll, "Scroll"},
		{EventType(99), "Unknown"}, // Edge case: unknown event type
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.eventType.String(); got != tt.expected {
				t.Errorf("EventType.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestEventType_IsClick(t *testing.T) {
	tests := []struct {
		eventType EventType
		expected  bool
	}{
		{EventPress, false},
		{EventRelease, false},
		{EventClick, true},
		{EventDoubleClick, true},
		{EventTripleClick, true},
		{EventDrag, false},
		{EventMotion, false},
		{EventScroll, false},
		{EventType(99), false}, // Edge case: unknown event type
	}

	for _, tt := range tests {
		t.Run(tt.eventType.String(), func(t *testing.T) {
			if got := tt.eventType.IsClick(); got != tt.expected {
				t.Errorf("EventType.IsClick() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestEventType_IsDrag(t *testing.T) {
	tests := []struct {
		eventType EventType
		expected  bool
	}{
		{EventPress, false},
		{EventRelease, false},
		{EventClick, false},
		{EventDoubleClick, false},
		{EventTripleClick, false},
		{EventDrag, true},
		{EventMotion, false},
		{EventScroll, false},
		{EventType(99), false}, // Edge case: unknown event type
	}

	for _, tt := range tests {
		t.Run(tt.eventType.String(), func(t *testing.T) {
			if got := tt.eventType.IsDrag(); got != tt.expected {
				t.Errorf("EventType.IsDrag() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestEventType_IsScroll(t *testing.T) {
	tests := []struct {
		eventType EventType
		expected  bool
	}{
		{EventPress, false},
		{EventRelease, false},
		{EventClick, false},
		{EventDoubleClick, false},
		{EventTripleClick, false},
		{EventDrag, false},
		{EventMotion, false},
		{EventScroll, true},
		{EventType(99), false}, // Edge case: unknown event type
	}

	for _, tt := range tests {
		t.Run(tt.eventType.String(), func(t *testing.T) {
			if got := tt.eventType.IsScroll(); got != tt.expected {
				t.Errorf("EventType.IsScroll() = %v, want %v", got, tt.expected)
			}
		})
	}
}
