package value

import "testing"

func TestJustifyContent_String(t *testing.T) {
	tests := []struct {
		name    string
		justify JustifyContent
		want    string
	}{
		{
			name:    "Start to string",
			justify: JustifyContentStart,
			want:    "start",
		},
		{
			name:    "End to string",
			justify: JustifyContentEnd,
			want:    "end",
		},
		{
			name:    "Center to string",
			justify: JustifyContentCenter,
			want:    "center",
		},
		{
			name:    "SpaceBetween to string",
			justify: JustifyContentSpaceBetween,
			want:    "space-between",
		},
		{
			name:    "Invalid value",
			justify: JustifyContent(99),
			want:    "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.justify.String()
			if got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJustifyContent_Validate(t *testing.T) {
	tests := []struct {
		name    string
		justify JustifyContent
		want    bool
	}{
		{
			name:    "Start is valid",
			justify: JustifyContentStart,
			want:    true,
		},
		{
			name:    "End is valid",
			justify: JustifyContentEnd,
			want:    true,
		},
		{
			name:    "Center is valid",
			justify: JustifyContentCenter,
			want:    true,
		},
		{
			name:    "SpaceBetween is valid",
			justify: JustifyContentSpaceBetween,
			want:    true,
		},
		{
			name:    "Invalid value",
			justify: JustifyContent(99),
			want:    false,
		},
		{
			name:    "Negative value",
			justify: JustifyContent(-1),
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.justify.Validate()
			if got != tt.want {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJustifyContent_IsDefault(t *testing.T) {
	tests := []struct {
		name    string
		justify JustifyContent
		want    bool
	}{
		{
			name:    "Start is default",
			justify: JustifyContentStart,
			want:    true,
		},
		{
			name:    "End is not default",
			justify: JustifyContentEnd,
			want:    false,
		},
		{
			name:    "Center is not default",
			justify: JustifyContentCenter,
			want:    false,
		},
		{
			name:    "SpaceBetween is not default",
			justify: JustifyContentSpaceBetween,
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.justify.IsDefault()
			if got != tt.want {
				t.Errorf("IsDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJustifyContent_NeedsDistribution(t *testing.T) {
	tests := []struct {
		name    string
		justify JustifyContent
		want    bool
	}{
		{
			name:    "Start doesn't need distribution",
			justify: JustifyContentStart,
			want:    false,
		},
		{
			name:    "End doesn't need distribution",
			justify: JustifyContentEnd,
			want:    false,
		},
		{
			name:    "Center needs distribution",
			justify: JustifyContentCenter,
			want:    true,
		},
		{
			name:    "SpaceBetween needs distribution",
			justify: JustifyContentSpaceBetween,
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.justify.NeedsDistribution()
			if got != tt.want {
				t.Errorf("NeedsDistribution() = %v, want %v", got, tt.want)
			}
		})
	}
}
