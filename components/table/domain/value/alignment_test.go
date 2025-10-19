package value

import "testing"

func TestAlignment_String(t *testing.T) {
	tests := []struct {
		name      string
		alignment Alignment
		want      string
	}{
		{"Left", AlignmentLeft, "left"},
		{"Center", AlignmentCenter, "center"},
		{"Right", AlignmentRight, "right"},
		{"Unknown", Alignment(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.alignment.String(); got != tt.want {
				t.Errorf("Alignment.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlignment_IsLeft(t *testing.T) {
	tests := []struct {
		name      string
		alignment Alignment
		want      bool
	}{
		{"Left", AlignmentLeft, true},
		{"Center", AlignmentCenter, false},
		{"Right", AlignmentRight, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.alignment.IsLeft(); got != tt.want {
				t.Errorf("Alignment.IsLeft() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlignment_IsCenter(t *testing.T) {
	tests := []struct {
		name      string
		alignment Alignment
		want      bool
	}{
		{"Left", AlignmentLeft, false},
		{"Center", AlignmentCenter, true},
		{"Right", AlignmentRight, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.alignment.IsCenter(); got != tt.want {
				t.Errorf("Alignment.IsCenter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlignment_IsRight(t *testing.T) {
	tests := []struct {
		name      string
		alignment Alignment
		want      bool
	}{
		{"Left", AlignmentLeft, false},
		{"Center", AlignmentCenter, false},
		{"Right", AlignmentRight, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.alignment.IsRight(); got != tt.want {
				t.Errorf("Alignment.IsRight() = %v, want %v", got, tt.want)
			}
		})
	}
}
