package value

import (
	"testing"
)

func TestSelectionMode_IsSingle(t *testing.T) {
	tests := []struct {
		name string
		mode SelectionMode
		want bool
	}{
		{
			name: "single mode returns true",
			mode: SelectionModeSingle,
			want: true,
		},
		{
			name: "multi mode returns false",
			mode: SelectionModeMulti,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mode.IsSingle(); got != tt.want {
				t.Errorf("SelectionMode.IsSingle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelectionMode_IsMulti(t *testing.T) {
	tests := []struct {
		name string
		mode SelectionMode
		want bool
	}{
		{
			name: "multi mode returns true",
			mode: SelectionModeMulti,
			want: true,
		},
		{
			name: "single mode returns false",
			mode: SelectionModeSingle,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mode.IsMulti(); got != tt.want {
				t.Errorf("SelectionMode.IsMulti() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelectionMode_String(t *testing.T) {
	tests := []struct {
		name string
		mode SelectionMode
		want string
	}{
		{
			name: "single mode string",
			mode: SelectionModeSingle,
			want: "single",
		},
		{
			name: "multi mode string",
			mode: SelectionModeMulti,
			want: "multi",
		},
		{
			name: "invalid mode string",
			mode: SelectionMode(999),
			want: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mode.String(); got != tt.want {
				t.Errorf("SelectionMode.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
