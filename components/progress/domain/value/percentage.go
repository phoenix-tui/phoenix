package value

// Percentage represents progress percentage (0-100).
// It provides value object semantics with immutability and clamping.
type Percentage struct {
	value int // Clamped to 0-100
}

// NewPercentage creates a new Percentage value object.
// The value is automatically clamped to the range [0, 100].
func NewPercentage(value int) *Percentage {
	return &Percentage{
		value: clamp(value, 0, 100),
	}
}

// Value returns the current percentage value (0-100).
func (p *Percentage) Value() int {
	return p.value
}

// Add returns a new Percentage with the delta added.
// The result is clamped to [0, 100].
func (p *Percentage) Add(delta int) *Percentage {
	return NewPercentage(p.value + delta)
}

// Subtract returns a new Percentage with the delta subtracted.
// The result is clamped to [0, 100].
func (p *Percentage) Subtract(delta int) *Percentage {
	return NewPercentage(p.value - delta)
}

// IsComplete returns true if the percentage is 100%.
func (p *Percentage) IsComplete() bool {
	return p.value == 100
}

// clamp restricts a value to the range [min, max].
func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
