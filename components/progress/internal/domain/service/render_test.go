package service

import (
	"strings"
	"testing"

	"github.com/phoenix-tui/phoenix/components/progress/internal/domain/model"
)

func TestNewRenderService(t *testing.T) {
	service := NewRenderService()
	if service == nil {
		t.Fatal("NewRenderService() returned nil")
	}
}

func TestCalculateFilledWidth(t *testing.T) {
	service := NewRenderService()

	tests := []struct {
		name       string
		barWidth   int
		percentage int
		expected   int
	}{
		// Normal cases.
		{"0% of 40", 40, 0, 0},
		{"25% of 40", 40, 25, 10},
		{"50% of 40", 40, 50, 20},
		{"75% of 40", 40, 75, 30},
		{"100% of 40", 40, 100, 40},

		// Edge cases.
		{"0% of 100", 100, 0, 0},
		{"100% of 100", 100, 100, 100},
		{"50% of 1", 1, 50, 0}, // Rounds down
		{"100% of 1", 1, 100, 1},

		// Rounding.
		{"33% of 100", 100, 33, 33},
		{"66% of 100", 100, 66, 66},
		{"10% of 50", 50, 10, 5},
		{"20% of 50", 50, 20, 10},

		// Zero width.
		{"Any % of 0", 0, 50, 0},

		// Negative (invalid but handled)
		{"Negative width", -10, 50, 0},
		{"Negative percentage", 40, -10, 0},

		// Over 100%.
		{"Over 100%", 40, 150, 40},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.CalculateFilledWidth(tt.barWidth, tt.percentage)
			if result != tt.expected {
				t.Errorf("CalculateFilledWidth(%d, %d) = %d, expected %d",
					tt.barWidth, tt.percentage, result, tt.expected)
			}
		})
	}
}

func TestRenderBar(t *testing.T) {
	service := NewRenderService()

	tests := []struct {
		name     string
		bar      *model.Bar
		expected string
	}{
		{
			name:     "Simple bar 0%",
			bar:      model.NewBar(10),
			expected: "░░░░░░░░░░",
		},
		{
			name:     "Simple bar 50%",
			bar:      model.NewBarWithPercentage(10, 50),
			expected: "█████░░░░░",
		},
		{
			name:     "Simple bar 100%",
			bar:      model.NewBarWithPercentage(10, 100),
			expected: "██████████",
		},
		{
			name: "Bar with percentage display",
			bar: func() *model.Bar {
				b := model.NewBarWithPercentage(10, 50).WithShowPercent(true)
				return &b
			}(),
			expected: "█████░░░░░ 050%",
		},
		{
			name: "Bar with label",
			bar: func() *model.Bar {
				b := model.NewBar(10).WithLabel("Loading")
				return &b
			}(),
			expected: "Loading ░░░░░░░░░░",
		},
		{
			name: "Bar with label and percentage",
			bar: func() *model.Bar {
				b := model.NewBarWithPercentage(10, 75).
					WithLabel("Progress").
					WithShowPercent(true)
				return &b
			}(),
			expected: "Progress ███████░░░ 075%",
		},
		{
			name: "Bar with custom chars",
			bar: func() *model.Bar {
				b := model.NewBarWithPercentage(10, 50).
					WithFillChar('▓').
					WithEmptyChar('▒')
				return &b
			}(),
			expected: "▓▓▓▓▓▒▒▒▒▒",
		},
		{
			name:     "Nil bar",
			bar:      nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.RenderBar(tt.bar)
			if result != tt.expected {
				t.Errorf("RenderBar() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestRenderBarStructure(t *testing.T) {
	service := NewRenderService()

	bar := func() *model.Bar {
		b := model.NewBarWithPercentage(20, 50).
			WithLabel("Test").
			WithShowPercent(true)
		return &b
	}()

	result := service.RenderBar(bar)

	// Check for label.
	if !strings.HasPrefix(result, "Test ") {
		t.Errorf("Result should start with label: %q", result)
	}

	// Check for percentage.
	if !strings.HasSuffix(result, " 050%") {
		t.Errorf("Result should end with percentage: %q", result)
	}

	// Check that result contains expected parts.
	// (Don't check byte length as Unicode chars take 3 bytes each)
	parts := strings.Split(result, " ")
	if len(parts) != 3 { // "Test" + bar + "050%"
		t.Errorf("Result should have 3 parts: %q", result)
	}
}

func TestRenderBarEmptyLabel(t *testing.T) {
	service := NewRenderService()

	bar := func() *model.Bar {
		b := model.NewBarWithPercentage(10, 50).
			WithLabel("") // Empty label
		return &b
	}()

	result := service.RenderBar(bar)

	// Should not have leading space.
	if strings.HasPrefix(result, " ") {
		t.Errorf("Result should not start with space: %q", result)
	}

	// Should be just the bar.
	if result != "█████░░░░░" {
		t.Errorf("Result = %q, expected bar only", result)
	}
}

func TestRenderBarDifferentWidths(t *testing.T) {
	service := NewRenderService()

	tests := []struct {
		width      int
		percentage int
	}{
		{1, 0},
		{1, 50},
		{1, 100},
		{5, 50},
		{20, 50},
		{40, 50},
		{100, 50},
	}

	for _, tt := range tests {
		bar := model.NewBarWithPercentage(tt.width, tt.percentage)
		result := service.RenderBar(bar)

		// Count characters (not bytes) - Unicode chars are 3 bytes each.
		totalChars := strings.Count(result, "█") + strings.Count(result, "░")
		if totalChars != tt.width {
			t.Errorf("Width %d, percentage %d: chars = %d, expected %d: %q",
				tt.width, tt.percentage, totalChars, tt.width, result)
		}
	}
}

func TestFormatPercentage(t *testing.T) {
	service := NewRenderService()

	tests := []struct {
		percentage int
		expected   string
	}{
		{0, "000%"},
		{1, "001%"},
		{10, "010%"},
		{50, "050%"},
		{99, "099%"},
		{100, "100%"},
	}

	for _, tt := range tests {
		result := service.formatPercentage(tt.percentage)
		if result != tt.expected {
			t.Errorf("formatPercentage(%d) = %q, expected %q",
				tt.percentage, result, tt.expected)
		}
	}
}

func TestRenderBarAllPercentages(t *testing.T) {
	service := NewRenderService()

	// Test all percentages 0-100.
	for pct := 0; pct <= 100; pct++ {
		bar := model.NewBarWithPercentage(100, pct)
		result := service.RenderBar(bar)

		// Count filled chars.
		filled := strings.Count(result, "█")
		empty := strings.Count(result, "░")

		if filled+empty != 100 {
			t.Errorf("Percentage %d: total chars = %d, expected 100", pct, filled+empty)
		}

		// Verify filled is correct.
		expectedFilled := service.CalculateFilledWidth(100, pct)
		if filled != expectedFilled {
			t.Errorf("Percentage %d: filled = %d, expected %d", pct, filled, expectedFilled)
		}
	}
}
