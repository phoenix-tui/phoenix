package value

import "testing"

func TestButton_String(t *testing.T) {
	tests := []struct {
		button   Button
		expected string
	}{
		{ButtonNone, "None"},
		{ButtonLeft, "Left"},
		{ButtonMiddle, "Middle"},
		{ButtonRight, "Right"},
		{ButtonWheelUp, "WheelUp"},
		{ButtonWheelDown, "WheelDown"},
		{Button(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.button.String(); got != tt.expected {
				t.Errorf("Button.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestButton_IsWheel(t *testing.T) {
	tests := []struct {
		button   Button
		expected bool
	}{
		{ButtonNone, false},
		{ButtonLeft, false},
		{ButtonMiddle, false},
		{ButtonRight, false},
		{ButtonWheelUp, true},
		{ButtonWheelDown, true},
	}

	for _, tt := range tests {
		t.Run(tt.button.String(), func(t *testing.T) {
			if got := tt.button.IsWheel(); got != tt.expected {
				t.Errorf("Button.IsWheel() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestButton_IsButton(t *testing.T) {
	tests := []struct {
		button   Button
		expected bool
	}{
		{ButtonNone, false},
		{ButtonLeft, true},
		{ButtonMiddle, true},
		{ButtonRight, true},
		{ButtonWheelUp, false},
		{ButtonWheelDown, false},
	}

	for _, tt := range tests {
		t.Run(tt.button.String(), func(t *testing.T) {
			if got := tt.button.IsButton(); got != tt.expected {
				t.Errorf("Button.IsButton() = %v, want %v", got, tt.expected)
			}
		})
	}
}
