package value_test

import (
	"testing"

	"github.com/phoenix-tui/phoenix/components/form/internal/domain/value"
)

func TestRequired(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{"nil value", nil, true},
		{"empty string", "", true},
		{"whitespace string", "   ", true},
		{"valid string", "hello", false},
		{"zero int", 0, false}, // Not a string, should pass
	}

	validator := value.Required()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Required() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMinLength(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		minLen  int
		wantErr bool
	}{
		{"too short", "ab", 3, true},
		{"exact length", "abc", 3, false},
		{"longer than min", "abcdef", 3, false},
		{"non-string", 123, 3, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := value.MinLength(tt.minLen)
			err := validator(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("MinLength(%d) error = %v, wantErr %v", tt.minLen, err, tt.wantErr)
			}
		})
	}
}

func TestMaxLength(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		maxLen  int
		wantErr bool
	}{
		{"shorter than max", "ab", 5, false},
		{"exact length", "abc", 3, false},
		{"longer than max", "abcdef", 3, true},
		{"non-string", 123, 3, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := value.MaxLength(tt.maxLen)
			err := validator(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("MaxLength(%d) error = %v, wantErr %v", tt.maxLen, err, tt.wantErr)
			}
		})
	}
}

func TestEmail(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{"valid email", "user@example.com", false},
		{"valid email with subdomain", "user@mail.example.com", false},
		{"valid email with plus", "user+tag@example.com", false},
		{"missing @", "userexample.com", true},
		{"missing domain", "user@", true},
		{"missing TLD", "user@example", true},
		{"spaces", "user @example.com", true},
		{"non-string", 123, true},
	}

	validator := value.Email()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Email() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPattern(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		input   interface{}
		wantErr bool
	}{
		{"matches digits", `^\d{3}$`, "123", false},
		{"doesn't match digits", `^\d{3}$`, "12a", true},
		{"matches phone", `^\d{3}-\d{3}-\d{4}$`, "555-123-4567", false},
		{"doesn't match phone", `^\d{3}-\d{3}-\d{4}$`, "555-1234567", true},
		{"non-string", `^\d+$`, 123, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := value.Pattern(tt.pattern)
			err := validator(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pattern(%q) error = %v, wantErr %v", tt.pattern, err, tt.wantErr)
			}
		})
	}
}

func TestCustom(t *testing.T) {
	// Custom validator that checks if string contains "test"
	customFn := func(val interface{}) error {
		str, ok := val.(string)
		if !ok {
			return nil // Skip non-strings
		}
		if len(str) > 0 && str[0] == 't' {
			return nil
		}
		return &validationError{msg: "must start with 't'"}
	}

	validator := value.Custom(customFn)

	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{"starts with t", "test", false},
		{"doesn't start with t", "hello", true},
		{"non-string", 123, false}, // Custom validator skips non-strings
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Custom() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// validationError is a helper type for custom error messages.
type validationError struct {
	msg string
}

func (e *validationError) Error() string {
	return e.msg
}
