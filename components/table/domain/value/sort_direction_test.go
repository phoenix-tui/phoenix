package value

import "testing"

func TestSortDirection_String(t *testing.T) {
	tests := []struct {
		name      string
		direction SortDirection
		want      string
	}{
		{"None", SortDirectionNone, "none"},
		{"Ascending", SortDirectionAsc, "asc"},
		{"Descending", SortDirectionDesc, "desc"},
		{"Unknown", SortDirection(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.direction.String(); got != tt.want {
				t.Errorf("SortDirection.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortDirection_IsAscending(t *testing.T) {
	tests := []struct {
		name      string
		direction SortDirection
		want      bool
	}{
		{"None", SortDirectionNone, false},
		{"Ascending", SortDirectionAsc, true},
		{"Descending", SortDirectionDesc, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.direction.IsAscending(); got != tt.want {
				t.Errorf("SortDirection.IsAscending() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortDirection_IsDescending(t *testing.T) {
	tests := []struct {
		name      string
		direction SortDirection
		want      bool
	}{
		{"None", SortDirectionNone, false},
		{"Ascending", SortDirectionAsc, false},
		{"Descending", SortDirectionDesc, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.direction.IsDescending(); got != tt.want {
				t.Errorf("SortDirection.IsDescending() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortDirection_IsNone(t *testing.T) {
	tests := []struct {
		name      string
		direction SortDirection
		want      bool
	}{
		{"None", SortDirectionNone, true},
		{"Ascending", SortDirectionAsc, false},
		{"Descending", SortDirectionDesc, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.direction.IsNone(); got != tt.want {
				t.Errorf("SortDirection.IsNone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortDirection_Toggle(t *testing.T) {
	tests := []struct {
		name      string
		direction SortDirection
		want      SortDirection
	}{
		{"None to Asc", SortDirectionNone, SortDirectionAsc},
		{"Asc to Desc", SortDirectionAsc, SortDirectionDesc},
		{"Desc to Asc", SortDirectionDesc, SortDirectionAsc},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.direction.Toggle(); got != tt.want {
				t.Errorf("SortDirection.Toggle() = %v, want %v", got, tt.want)
			}
		})
	}
}
