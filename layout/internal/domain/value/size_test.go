package value

import (
	"testing"
)

func TestNewSize(t *testing.T) {
	tests := []struct {
		name          string
		width         int
		height        int
		minWidth      int
		maxWidth      int
		minHeight     int
		maxHeight     int
		wantWidth     int
		wantHeight    int
		wantMinWidth  int
		wantMaxWidth  int
		wantMinHeight int
		wantMaxHeight int
	}{
		{
			name:  "exact size",
			width: 80, height: 24,
			minWidth: -1, maxWidth: -1, minHeight: -1, maxHeight: -1,
			wantWidth: 80, wantHeight: 24,
			wantMinWidth: -1, wantMaxWidth: -1, wantMinHeight: -1, wantMaxHeight: -1,
		},
		{
			name:  "min/max constraints only",
			width: -1, height: -1,
			minWidth: 40, maxWidth: 120, minHeight: 10, maxHeight: 30,
			wantWidth: -1, wantHeight: -1,
			wantMinWidth: 40, wantMaxWidth: 120, wantMinHeight: 10, wantMaxHeight: 30,
		},
		{
			name:  "unconstrained",
			width: -1, height: -1,
			minWidth: -1, maxWidth: -1, minHeight: -1, maxHeight: -1,
			wantWidth: -1, wantHeight: -1,
			wantMinWidth: -1, wantMaxWidth: -1, wantMinHeight: -1, wantMaxHeight: -1,
		},
		{
			name:  "only min constraints",
			width: -1, height: -1,
			minWidth: 20, maxWidth: -1, minHeight: 5, maxHeight: -1,
			wantWidth: -1, wantHeight: -1,
			wantMinWidth: 20, wantMaxWidth: -1, wantMinHeight: 5, wantMaxHeight: -1,
		},
		{
			name:  "only max constraints",
			width: -1, height: -1,
			minWidth: -1, maxWidth: 100, minHeight: -1, maxHeight: 50,
			wantWidth: -1, wantHeight: -1,
			wantMinWidth: -1, wantMaxWidth: 100, wantMinHeight: -1, wantMaxHeight: 50,
		},
		{
			name:  "zero values normalized to -1",
			width: -5, height: -10,
			minWidth: -3, maxWidth: -7, minHeight: -2, maxHeight: -9,
			wantWidth: -1, wantHeight: -1,
			wantMinWidth: -1, wantMaxWidth: -1, wantMinHeight: -1, wantMaxHeight: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSize(tt.width, tt.height, tt.minWidth, tt.maxWidth, tt.minHeight, tt.maxHeight)

			if s.Width() != tt.wantWidth {
				t.Errorf("Width() = %d, want %d", s.Width(), tt.wantWidth)
			}
			if s.Height() != tt.wantHeight {
				t.Errorf("Height() = %d, want %d", s.Height(), tt.wantHeight)
			}
			if s.MinWidth() != tt.wantMinWidth {
				t.Errorf("MinWidth() = %d, want %d", s.MinWidth(), tt.wantMinWidth)
			}
			if s.MaxWidth() != tt.wantMaxWidth {
				t.Errorf("MaxWidth() = %d, want %d", s.MaxWidth(), tt.wantMaxWidth)
			}
			if s.MinHeight() != tt.wantMinHeight {
				t.Errorf("MinHeight() = %d, want %d", s.MinHeight(), tt.wantMinHeight)
			}
			if s.MaxHeight() != tt.wantMaxHeight {
				t.Errorf("MaxHeight() = %d, want %d", s.MaxHeight(), tt.wantMaxHeight)
			}
		})
	}
}

func TestNewSize_Panics(t *testing.T) {
	tests := []struct {
		name      string
		width     int
		height    int
		minWidth  int
		maxWidth  int
		minHeight int
		maxHeight int
	}{
		{
			name:  "minWidth > maxWidth",
			width: -1, height: -1,
			minWidth: 100, maxWidth: 50, minHeight: -1, maxHeight: -1,
		},
		{
			name:  "minHeight > maxHeight",
			width: -1, height: -1,
			minWidth: -1, maxWidth: -1, minHeight: 50, maxHeight: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("NewSize() did not panic")
				}
			}()
			NewSize(tt.width, tt.height, tt.minWidth, tt.maxWidth, tt.minHeight, tt.maxHeight)
		})
	}
}

func TestNewSizeExact(t *testing.T) {
	s := NewSizeExact(80, 24)

	if !s.HasWidth() || s.Width() != 80 {
		t.Errorf("Width = %d, want 80", s.Width())
	}
	if !s.HasHeight() || s.Height() != 24 {
		t.Errorf("Height = %d, want 24", s.Height())
	}
	if s.HasMinWidth() || s.HasMaxWidth() || s.HasMinHeight() || s.HasMaxHeight() {
		t.Error("Expected no min/max constraints")
	}
}

func TestNewSizeUnconstrained(t *testing.T) {
	s := NewSizeUnconstrained()

	if !s.IsUnconstrained() {
		t.Error("Expected unconstrained size")
	}
	if s.HasWidth() || s.HasHeight() || s.HasMinWidth() || s.HasMaxWidth() || s.HasMinHeight() || s.HasMaxHeight() {
		t.Error("Expected no constraints")
	}
}

func TestSize_Has(t *testing.T) {
	tests := []struct {
		name          string
		size          Size
		wantWidth     bool
		wantHeight    bool
		wantMinWidth  bool
		wantMaxWidth  bool
		wantMinHeight bool
		wantMaxHeight bool
	}{
		{
			name:      "exact size",
			size:      NewSizeExact(80, 24),
			wantWidth: true, wantHeight: true,
			wantMinWidth: false, wantMaxWidth: false,
			wantMinHeight: false, wantMaxHeight: false,
		},
		{
			name:      "min/max only",
			size:      NewSize(-1, -1, 40, 120, 10, 30),
			wantWidth: false, wantHeight: false,
			wantMinWidth: true, wantMaxWidth: true,
			wantMinHeight: true, wantMaxHeight: true,
		},
		{
			name:      "unconstrained",
			size:      NewSizeUnconstrained(),
			wantWidth: false, wantHeight: false,
			wantMinWidth: false, wantMaxWidth: false,
			wantMinHeight: false, wantMaxHeight: false,
		},
		{
			name:      "width with max height",
			size:      NewSize(80, -1, -1, -1, -1, 24),
			wantWidth: true, wantHeight: false,
			wantMinWidth: false, wantMaxWidth: false,
			wantMinHeight: false, wantMaxHeight: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.size.HasWidth() != tt.wantWidth {
				t.Errorf("HasWidth() = %v, want %v", tt.size.HasWidth(), tt.wantWidth)
			}
			if tt.size.HasHeight() != tt.wantHeight {
				t.Errorf("HasHeight() = %v, want %v", tt.size.HasHeight(), tt.wantHeight)
			}
			if tt.size.HasMinWidth() != tt.wantMinWidth {
				t.Errorf("HasMinWidth() = %v, want %v", tt.size.HasMinWidth(), tt.wantMinWidth)
			}
			if tt.size.HasMaxWidth() != tt.wantMaxWidth {
				t.Errorf("HasMaxWidth() = %v, want %v", tt.size.HasMaxWidth(), tt.wantMaxWidth)
			}
			if tt.size.HasMinHeight() != tt.wantMinHeight {
				t.Errorf("HasMinHeight() = %v, want %v", tt.size.HasMinHeight(), tt.wantMinHeight)
			}
			if tt.size.HasMaxHeight() != tt.wantMaxHeight {
				t.Errorf("HasMaxHeight() = %v, want %v", tt.size.HasMaxHeight(), tt.wantMaxHeight)
			}
		})
	}
}

func TestSize_With(t *testing.T) {
	base := NewSizeUnconstrained()

	t.Run("WithWidth", func(t *testing.T) {
		s := base.WithWidth(80)
		if s.Width() != 80 {
			t.Errorf("Width() = %d, want 80", s.Width())
		}
		if base.HasWidth() {
			t.Error("Original size should not be modified (immutable)")
		}
	})

	t.Run("WithHeight", func(t *testing.T) {
		s := base.WithHeight(24)
		if s.Height() != 24 {
			t.Errorf("Height() = %d, want 24", s.Height())
		}
	})

	t.Run("WithMinWidth", func(t *testing.T) {
		s := base.WithMinWidth(40)
		if s.MinWidth() != 40 {
			t.Errorf("MinWidth() = %d, want 40", s.MinWidth())
		}
	})

	t.Run("WithMaxWidth", func(t *testing.T) {
		s := base.WithMaxWidth(120)
		if s.MaxWidth() != 120 {
			t.Errorf("MaxWidth() = %d, want 120", s.MaxWidth())
		}
	})

	t.Run("WithMinHeight", func(t *testing.T) {
		s := base.WithMinHeight(10)
		if s.MinHeight() != 10 {
			t.Errorf("MinHeight() = %d, want 10", s.MinHeight())
		}
	})

	t.Run("WithMaxHeight", func(t *testing.T) {
		s := base.WithMaxHeight(30)
		if s.MaxHeight() != 30 {
			t.Errorf("MaxHeight() = %d, want 30", s.MaxHeight())
		}
	})

	t.Run("chaining", func(t *testing.T) {
		s := base.WithWidth(80).WithHeight(24).WithMinWidth(40)
		if s.Width() != 80 || s.Height() != 24 || s.MinWidth() != 40 {
			t.Error("Chaining did not work correctly")
		}
	})
}

func TestSize_Constrain(t *testing.T) {
	tests := []struct {
		name       string
		size       Size
		inputW     int
		inputH     int
		wantWidth  int
		wantHeight int
	}{
		{
			name:   "exact size overrides input",
			size:   NewSizeExact(80, 24),
			inputW: 100, inputH: 50,
			wantWidth: 80, wantHeight: 24,
		},
		{
			name:   "min width enforced",
			size:   NewSize(-1, -1, 40, -1, -1, -1),
			inputW: 20, inputH: 10,
			wantWidth: 40, wantHeight: 10,
		},
		{
			name:   "max width enforced",
			size:   NewSize(-1, -1, -1, 100, -1, -1),
			inputW: 150, inputH: 10,
			wantWidth: 100, wantHeight: 10,
		},
		{
			name:   "min and max width enforced",
			size:   NewSize(-1, -1, 40, 120, -1, -1),
			inputW: 20, inputH: 10,
			wantWidth: 40, wantHeight: 10,
		},
		{
			name:   "within range - no change",
			size:   NewSize(-1, -1, 40, 120, 10, 30),
			inputW: 80, inputH: 20,
			wantWidth: 80, wantHeight: 20,
		},
		{
			name:   "min height enforced",
			size:   NewSize(-1, -1, -1, -1, 10, -1),
			inputW: 80, inputH: 5,
			wantWidth: 80, wantHeight: 10,
		},
		{
			name:   "max height enforced",
			size:   NewSize(-1, -1, -1, -1, -1, 30),
			inputW: 80, inputH: 50,
			wantWidth: 80, wantHeight: 30,
		},
		{
			name:   "unconstrained - no change",
			size:   NewSizeUnconstrained(),
			inputW: 123, inputH: 456,
			wantWidth: 123, wantHeight: 456,
		},
		{
			name:   "exact width, max height",
			size:   NewSize(80, -1, -1, -1, -1, 30),
			inputW: 100, inputH: 50,
			wantWidth: 80, wantHeight: 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWidth, gotHeight := tt.size.Constrain(tt.inputW, tt.inputH)
			if gotWidth != tt.wantWidth {
				t.Errorf("Constrain() width = %d, want %d", gotWidth, tt.wantWidth)
			}
			if gotHeight != tt.wantHeight {
				t.Errorf("Constrain() height = %d, want %d", gotHeight, tt.wantHeight)
			}
		})
	}
}

func TestSize_String(t *testing.T) {
	tests := []struct {
		name string
		size Size
		want string
	}{
		{
			name: "unconstrained",
			size: NewSizeUnconstrained(),
			want: "Size{unconstrained}",
		},
		{
			name: "exact size",
			size: NewSizeExact(80, 24),
			want: "Size{width=80 height=24}",
		},
		{
			name: "min/max width",
			size: NewSize(-1, -1, 40, 120, -1, -1),
			want: "Size{minW=40,maxW=120}",
		},
		{
			name: "min/max height",
			size: NewSize(-1, -1, -1, -1, 10, 30),
			want: "Size{minH=10,maxH=30}",
		},
		{
			name: "all constraints",
			size: NewSize(-1, -1, 40, 120, 10, 30),
			want: "Size{minW=40,maxW=120 minH=10,maxH=30}",
		},
		{
			name: "only min",
			size: NewSize(-1, -1, 20, -1, 5, -1),
			want: "Size{minW=20 minH=5}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.size.String()
			if got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSize_IsUnconstrained(t *testing.T) {
	tests := []struct {
		name string
		size Size
		want bool
	}{
		{
			name: "unconstrained",
			size: NewSizeUnconstrained(),
			want: true,
		},
		{
			name: "exact size",
			size: NewSizeExact(80, 24),
			want: false,
		},
		{
			name: "only width",
			size: NewSize(80, -1, -1, -1, -1, -1),
			want: false,
		},
		{
			name: "only min width",
			size: NewSize(-1, -1, 40, -1, -1, -1),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.size.IsUnconstrained()
			if got != tt.want {
				t.Errorf("IsUnconstrained() = %v, want %v", got, tt.want)
			}
		})
	}
}
