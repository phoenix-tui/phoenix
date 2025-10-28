package value

import (
	"testing"
)

func TestFlexDirection_IsHorizontal(t *testing.T) {
	tests := []struct {
		name      string
		direction FlexDirection
		want      bool
	}{
		{
			name:      "Row is horizontal",
			direction: FlexDirectionRow,
			want:      true,
		},
		{
			name:      "Column is not horizontal",
			direction: FlexDirectionColumn,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.direction.IsHorizontal()
			if got != tt.want {
				t.Errorf("IsHorizontal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlexDirection_IsVertical(t *testing.T) {
	tests := []struct {
		name      string
		direction FlexDirection
		want      bool
	}{
		{
			name:      "Column is vertical",
			direction: FlexDirectionColumn,
			want:      true,
		},
		{
			name:      "Row is not vertical",
			direction: FlexDirectionRow,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.direction.IsVertical()
			if got != tt.want {
				t.Errorf("IsVertical() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlexDirection_String(t *testing.T) {
	tests := []struct {
		name      string
		direction FlexDirection
		want      string
	}{
		{
			name:      "Row to string",
			direction: FlexDirectionRow,
			want:      "row",
		},
		{
			name:      "Column to string",
			direction: FlexDirectionColumn,
			want:      "column",
		},
		{
			name:      "Invalid direction",
			direction: FlexDirection(99),
			want:      "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.direction.String()
			if got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlexDirection_Validate(t *testing.T) {
	tests := []struct {
		name      string
		direction FlexDirection
		want      bool
	}{
		{
			name:      "Row is valid",
			direction: FlexDirectionRow,
			want:      true,
		},
		{
			name:      "Column is valid",
			direction: FlexDirectionColumn,
			want:      true,
		},
		{
			name:      "Invalid value",
			direction: FlexDirection(99),
			want:      false,
		},
		{
			name:      "Negative value",
			direction: FlexDirection(-1),
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.direction.Validate()
			if got != tt.want {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
