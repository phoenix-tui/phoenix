package value

import "testing"

func TestAlignItems_String(t *testing.T) {
	tests := []struct {
		name  string
		align AlignItems
		want  string
	}{
		{
			name:  "Stretch to string",
			align: AlignItemsStretch,
			want:  "stretch",
		},
		{
			name:  "Start to string",
			align: AlignItemsStart,
			want:  "start",
		},
		{
			name:  "End to string",
			align: AlignItemsEnd,
			want:  "end",
		},
		{
			name:  "Center to string",
			align: AlignItemsCenter,
			want:  "center",
		},
		{
			name:  "Invalid value",
			align: AlignItems(99),
			want:  "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.align.String()
			if got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlignItems_Validate(t *testing.T) {
	tests := []struct {
		name  string
		align AlignItems
		want  bool
	}{
		{
			name:  "Stretch is valid",
			align: AlignItemsStretch,
			want:  true,
		},
		{
			name:  "Start is valid",
			align: AlignItemsStart,
			want:  true,
		},
		{
			name:  "End is valid",
			align: AlignItemsEnd,
			want:  true,
		},
		{
			name:  "Center is valid",
			align: AlignItemsCenter,
			want:  true,
		},
		{
			name:  "Invalid value",
			align: AlignItems(99),
			want:  false,
		},
		{
			name:  "Negative value",
			align: AlignItems(-1),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.align.Validate()
			if got != tt.want {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlignItems_IsDefault(t *testing.T) {
	tests := []struct {
		name  string
		align AlignItems
		want  bool
	}{
		{
			name:  "Stretch is default",
			align: AlignItemsStretch,
			want:  true,
		},
		{
			name:  "Start is not default",
			align: AlignItemsStart,
			want:  false,
		},
		{
			name:  "End is not default",
			align: AlignItemsEnd,
			want:  false,
		},
		{
			name:  "Center is not default",
			align: AlignItemsCenter,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.align.IsDefault()
			if got != tt.want {
				t.Errorf("IsDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlignItems_RequiresStretching(t *testing.T) {
	tests := []struct {
		name  string
		align AlignItems
		want  bool
	}{
		{
			name:  "Stretch requires stretching",
			align: AlignItemsStretch,
			want:  true,
		},
		{
			name:  "Start doesn't require stretching",
			align: AlignItemsStart,
			want:  false,
		},
		{
			name:  "End doesn't require stretching",
			align: AlignItemsEnd,
			want:  false,
		},
		{
			name:  "Center doesn't require stretching",
			align: AlignItemsCenter,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.align.RequiresStretching()
			if got != tt.want {
				t.Errorf("RequiresStretching() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlignItems_RequiresAlignment(t *testing.T) {
	tests := []struct {
		name  string
		align AlignItems
		want  bool
	}{
		{
			name:  "Stretch doesn't require alignment",
			align: AlignItemsStretch,
			want:  false,
		},
		{
			name:  "Start requires alignment",
			align: AlignItemsStart,
			want:  true,
		},
		{
			name:  "End requires alignment",
			align: AlignItemsEnd,
			want:  true,
		},
		{
			name:  "Center requires alignment",
			align: AlignItemsCenter,
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.align.RequiresAlignment()
			if got != tt.want {
				t.Errorf("RequiresAlignment() = %v, want %v", got, tt.want)
			}
		})
	}
}
