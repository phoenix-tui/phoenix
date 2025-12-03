// Package value contains value objects for the form component.
package value

import (
	"fmt"
	"regexp"
	"strings"
)

// Validator is a function that validates a value and returns an error if invalid.
type Validator func(value interface{}) error

// Required validates that a value is not empty.
func Required() Validator {
	return func(value interface{}) error {
		if value == nil {
			return fmt.Errorf("required")
		}

		// String validation
		if str, ok := value.(string); ok {
			if strings.TrimSpace(str) == "" {
				return fmt.Errorf("required")
			}
		}

		return nil
	}
}

// MinLength validates that a string has at least n characters.
func MinLength(n int) Validator {
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("invalid type (expected string)")
		}

		if len(str) < n {
			return fmt.Errorf("must be at least %d characters", n)
		}

		return nil
	}
}

// MaxLength validates that a string has at most n characters.
func MaxLength(n int) Validator {
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("invalid type (expected string)")
		}

		if len(str) > n {
			return fmt.Errorf("must be at most %d characters", n)
		}

		return nil
	}
}

// Email validates that a string is a valid email format.
func Email() Validator {
	// Simple email regex (RFC 5322 simplified)
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("invalid type (expected string)")
		}

		if !emailRegex.MatchString(str) {
			return fmt.Errorf("invalid email format")
		}

		return nil
	}
}

// Pattern validates that a string matches a regex pattern.
func Pattern(regex string) Validator {
	compiled := regexp.MustCompile(regex)

	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("invalid type (expected string)")
		}

		if !compiled.MatchString(str) {
			return fmt.Errorf("invalid format")
		}

		return nil
	}
}

// Custom creates a validator from a custom function.
func Custom(fn func(interface{}) error) Validator {
	return fn
}
