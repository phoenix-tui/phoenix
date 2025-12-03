package value

import (
	"testing"
)

func TestNewButton(t *testing.T) {
	tests := []struct {
		name      string
		label     string
		wantLabel string
		wantKey   rune
	}{
		{
			name:      "Yes button",
			label:     "Yes",
			wantLabel: "Yes",
			wantKey:   'Y',
		},
		{
			name:      "No button",
			label:     "No",
			wantLabel: "No",
			wantKey:   'N',
		},
		{
			name:      "Cancel button",
			label:     "Cancel",
			wantLabel: "Cancel",
			wantKey:   'C',
		},
		{
			name:      "Empty label",
			label:     "",
			wantLabel: "",
			wantKey:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewButton(tt.label)
			if b.Label() != tt.wantLabel {
				t.Errorf("Label() = %v, want %v", b.Label(), tt.wantLabel)
			}
			if b.Key() != tt.wantKey {
				t.Errorf("Key() = %v, want %v", b.Key(), tt.wantKey)
			}
		})
	}
}

func TestButton_MatchesKey(t *testing.T) {
	tests := []struct {
		name   string
		button *Button
		key    rune
		want   bool
	}{
		{
			name:   "Exact match uppercase",
			button: NewButton("Yes"),
			key:    'Y',
			want:   true,
		},
		{
			name:   "Exact match lowercase",
			button: NewButton("Yes"),
			key:    'y',
			want:   true,
		},
		{
			name:   "No match different letter",
			button: NewButton("Yes"),
			key:    'N',
			want:   false,
		},
		{
			name:   "No match different letter lowercase",
			button: NewButton("Yes"),
			key:    'n',
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.button.MatchesKey(tt.key); got != tt.want {
				t.Errorf("MatchesKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
